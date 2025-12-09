package event

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"apps/backend/common/customerror"
	eventdatagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"

	"github.com/cockroachdb/errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"

	certificateContract "apps/backend/contracts/certificate"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type ClaimCertificateUserError string

const (
	ClaimCertificateUserErrorCertificateNotFound       ClaimCertificateUserError = "certificate_not_found"
	ClaimCertificateUserErrorCertificateAlreadyClaimed ClaimCertificateUserError = "certificate_already_claimed"
	ClaimCertificateUserErrorCertificateRevoked        ClaimCertificateUserError = "certificate_revoked"
	ClaimCertificateUserErrorNotEligible               ClaimCertificateUserError = "not_eligible"
)

type ClaimCertificatePayload struct {
	CertificateID uuid.UUID `json:"certificate_id"`
}

// returns raw string, then message hash
// CRITICAL: walletAddress parameter should be the address derived from the user's private key,
// NOT from the JWT claims, to ensure cryptographic consistency
func (uc *EventUsecase) GetClaimCertificateSignMessage(ctx context.Context, client *ethclient.Client, walletAddress common.Address, currentUser auth.JwtClaims, certificateContractAddress common.Address, deadlineBlock *uint64) (*string, *ethcommon.Hash, error) {
	// Validation
	if deadlineBlock == nil {
		calculatedDeadlineBlock, err := cyptoutils.GetCalculatedDeadlineBlock(client)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get calculated deadline block")
		}
		deadlineBlock = &calculatedDeadlineBlock
	}

	// Use the provided walletAddress parameter (which should be derived from private key)
	// instead of currentUser.WalletAddress (which comes from JWT claims)
	signMessage, err := cyptoutils.GetSignMessage(walletAddress, certificateContractAddress, *deadlineBlock)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get sign message")
	}

	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	return &signMessage, &messageHash, nil
}

// Checks if the user is eligible to claim the certificate
// Eligibility is proven by:
// 1. Matching authentication credential ID OR matching email
// 2. Certificate is not revoked
// 3. Certificate config is published (EventCertificateAddress is set)
// 4. Certificate has not been claimed before (CertificateTokenId is NULL)
func (uc *EventUsecase) CheckClaimEligibility(ctx context.Context, certificate *entity.EventCertificate, currentUser *auth.JwtClaims) error {
	if currentUser == nil {
		return customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	// Check if certificate is revoked
	if certificate.RevokedAt != nil {
		return customerror.NewWithPreset(&customerror.ErrInvalidArgument, errors.New(string(ClaimCertificateUserErrorCertificateRevoked)))
	}

	// Check if certificate config is published
	if certificate.EventCertificateAddress == nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate is not published yet"))
	}

	// Check if certificate has already been claimed
	if certificate.CertificateTokenId != nil {
		return customerror.NewWithPreset(&customerror.ErrInvalidArgument, errors.New(string(ClaimCertificateUserErrorCertificateAlreadyClaimed)))
	}

	// Check if user is the intended receiver by credential ID
	if certificate.ReceiverCredentialId != nil {
		if *certificate.ReceiverCredentialId != currentUser.UserId {
			return customerror.NewWithPreset(&customerror.ErrUnauthorized, errors.New(string(ClaimCertificateUserErrorNotEligible)))
		}
		return nil
	}

	// Check if user is the intended receiver by email
	if certificate.ReceiverEmail != nil {
		if currentUser.Email == nil || *currentUser.Email != *certificate.ReceiverEmail {
			return customerror.NewWithPreset(&customerror.ErrUnauthorized, errors.New(string(ClaimCertificateUserErrorNotEligible)))
		}
		return nil
	}

	// If no receiver credential or email is set, anyone can claim (open certificate)
	return nil
}

func (uc *EventUsecase) ClaimCertificateWithPin(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, certificateId uuid.UUID, password string) (*entity.EventCertificate, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	// Get certificate
	certificate, err := uc.EventCertificateDataGateway.GetEventCertificateByID(ctx, certificateId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if certificate == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New(string(ClaimCertificateUserErrorCertificateNotFound)))
	}

	// Check eligibility (includes published check, not claimed check, revocation check, and user match)
	if err := uc.CheckClaimEligibility(ctx, certificate, currentUser); err != nil {
		return nil, err
	}

	// CRITICAL FIX: Get the actual wallet address derived from the private key
	// instead of using currentUser.WalletAddress from JWT claims
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if credential == nil || credential.EncryptedPrivateKey == nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("encrypted private key not found"))
	}

	// Decrypt to get the address that cryptographically matches the private key
	privateKey, participantAddress, err := cyptoutils.DecryptPrivateKey(*credential.EncryptedPrivateKey, password)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("invalid password or failed to decrypt private key"))
	}

	deadlineBlock, err := cyptoutils.GetCalculatedDeadlineBlock(client)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Use the derived address from the private key, NOT the JWT address
	signMessage, err := cyptoutils.GetSignMessage(*participantAddress, common.HexToAddress(*certificate.EventCertificateAddress), deadlineBlock)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Sign the message directly with the decrypted private key
	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	signature, err := cyptoutils.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Derive public key from private key for ECIES encryption
	participantPublicKey, err := cyptoutils.GetPublicKeyFromPrivateKey(privateKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to derive public key"))
	}

	return uc.claimCertificate(ctx, client, currentUser, certificate, signature, signMessage, participantAddress, participantPublicKey)
}

func (uc *EventUsecase) ClaimCertificateWithSignature(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, certificateId uuid.UUID, signature []byte, signMessage string) (*entity.EventCertificate, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	// Get certificate
	certificate, err := uc.EventCertificateDataGateway.GetEventCertificateByID(ctx, certificateId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if certificate == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New(string(ClaimCertificateUserErrorCertificateNotFound)))
	}

	// CRITICAL FIX: Get the actual wallet address derived from the private key
	// The signature verification must use the address that actually signed the message
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if credential == nil || credential.WalletAddress == "" {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("wallet address not found"))
	}

	// Use the wallet address from the credential (which is derived from the private key)
	participantAddress := common.HexToAddress(credential.WalletAddress)

	// validate original sign message
	deadlineBlock, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(signMessage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract deadline block from sign message")
	}
	isValid, err := cyptoutils.ValidateSignMessage(signMessage, participantAddress, common.HexToAddress(*certificate.EventCertificateAddress), deadlineBlock)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate sign message")
	}
	if !isValid {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("sign message is not valid"))
	}
	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	// check if the signature matches the sign message or not
	isValidHash, err := cyptoutils.VerifySignatureByDigest(participantAddress, messageHash, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify signature by digest")
	}
	if !isValidHash {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("signature does not match the sign message"))
	}

	// Check eligibility (includes published check, not claimed check, revocation check, and user match)
	if err := uc.CheckClaimEligibility(ctx, certificate, currentUser); err != nil {
		return nil, err
	}

	// WALLET EXTENSION FLOW:
	// User signs a message with their wallet extension (no PII sent)
	// Backend:
	// 1. Verifies the signature (already done above)
	// 2. Recovers public key from the signature
	// 3. Encrypts PII CSV with user's public key
	// 4. User can decrypt later with their wallet's private key

	// Recover public key from signature
	messageHashForRecovery := cyptoutils.HashEthereumMessage(signMessage)
	participantPublicKey, err := cyptoutils.RecoverPublicKeyFromSignature(messageHashForRecovery, signature)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to recover public key from signature"))
	}

	// Verify recovered address matches the participant address
	recoveredAddress := cyptoutils.PublicKeyToAddress(participantPublicKey)
	if recoveredAddress.Hex() != participantAddress.Hex() {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.Errorf("recovered address (%s) does not match participant address (%s)", recoveredAddress.Hex(), participantAddress.Hex()))
	}

	// Proceed with claiming using recovered public key
	return uc.claimCertificate(ctx, client, currentUser, certificate, signature, signMessage, &participantAddress, participantPublicKey)
}

func (uc *EventUsecase) claimCertificate(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, certificate *entity.EventCertificate, signature []byte, signMessage string, participantAddress *common.Address, participantPublicKey *ecdsa.PublicKey) (*entity.EventCertificate, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}
	if participantAddress == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("participant address is required"))
	}
	if participantPublicKey == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("participant public key is required for data encryption"))
	}

	// Check if certificate contract is deployed
	if certificate.EventCertificateAddress == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate contract not deployed yet"))
	}

	// Check if already minted (has token ID)
	if certificate.CertificateTokenId != nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate already minted"))
	}

	// ============================================
	// PREPARE DATA FOR MINTING
	// ============================================

	// 1. Get all event issuers
	issuers, err := uc.EventIssuerDataGateway.GetEventIssuersByEventID(ctx, certificate.EventId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get event issuers"))
	}
	if len(issuers) == 0 {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("no issuers found for this event"))
	}

	// 2. Get all certificate signatures (one per issuer)
	certificateSignatures, err := uc.EventCertificateSignatureDataGateway.GetEventCertificateSignaturesByEventCertificateID(ctx, certificate.Id)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get certificate signatures"))
	}
	if len(certificateSignatures) != len(issuers) {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("not all issuers have signed the certificate"))
	}

	// 3. Get event details
	event, err := uc.EventDataGateway.GetEventById(ctx, certificate.EventId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get event"))
	}

	// 4. Get host credentials (for public key)
	hostCredential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, event.OwnerCredentialId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get host credentials"))
	}

	// ============================================
	// BUILD MINT PARAMETERS
	// ============================================

	// Parameter 1: receiverAddress
	receiverAddress := *participantAddress

	// Parameter 2: userId (receiver's credential ID or generate one)
	userId := ""
	if certificate.ReceiverCredentialId != nil {
		userId = certificate.ReceiverCredentialId.String()
	} else {
		userId = currentUser.UserId.String()
	}

	// Parameter 3: certificateId
	certificateId := certificate.Id.String()

	// Parameter 4: issuerId (first issuer's credential ID)
	issuerId := issuers[0].IssuerCredentialID.String()

	// Parameter 5 & 6: encryptedUserData and backendEncryptedUserData
	// Build UserData CSV structure (PII - Raw Certificate Data)
	// UserData format: name,academic_institution,certificate_title,certificate_subtitle
	// This matches the pattern used in import_certificate_receivers.go
	name := ""
	if certificate.Name != nil {
		name = *certificate.Name
	}
	academicInstitution := ""
	if certificate.AcademicInstitution != nil {
		academicInstitution = *certificate.AcademicInstitution
	}
	certTitle := ""
	if certificate.CertificateTitle != nil {
		certTitle = *certificate.CertificateTitle
	}
	certSubtitle := ""
	if certificate.CertificateSubtitle != nil {
		certSubtitle = *certificate.CertificateSubtitle
	}

	// ============================================
	// DATA SECTION: Participant Profile from event_attendees
	// ============================================

	// 1. Get attendee data from event_attendees table
	// MUST exist - if not, user hasn't joined the event yet
	// Use currentUser.UserId (authenticated user's credential ID) to get their attendee record
	attendee, err := uc.EventAttendeeDg.GetEventAttendeeByEventIdAndCredentialId(ctx, certificate.EventId, currentUser.UserId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.Wrap(err, "attendee record not found - user must join event first"))
	}

	// 2. Build attendee profile JSON with all 8 PII fields
	// Keep null values as null before stringifying
	type AttendeeProfileData struct {
		FirstName           *string `json:"first_name"`
		LastName            *string `json:"last_name"`
		Email               *string `json:"email"`
		Bio                 *string `json:"bio"`
		PhoneNumber         *string `json:"phone_number"`
		Address             *string `json:"address"`
		AcademicInstitution *string `json:"academic_institution"`
		AcademicEmail       *string `json:"academic_email"`
	}

	attendeeProfile := AttendeeProfileData{
		FirstName:           attendee.FirstName,
		LastName:            attendee.LastName,
		Email:               attendee.Email,
		Bio:                 attendee.Bio,
		PhoneNumber:         attendee.PhoneNumber,
		Address:             attendee.Address,
		AcademicInstitution: attendee.AcademicInstitution,
		AcademicEmail:       attendee.AcademicEmail,
	}

	// 3. Convert to JSON string
	attendeeProfileJSON, err := json.Marshal(attendeeProfile)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to marshal attendee profile"))
	}
	attendeeProfileStr := string(attendeeProfileJSON)

	// Derive backend public key from blockchain private key for encryption
	backendPrivateKey, err := crypto.HexToECDSA(uc.cfg.Blockchain.PrivateKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to parse blockchain private key"))
	}
	backendPublicKey, err := cyptoutils.GetPublicKeyFromPrivateKey(backendPrivateKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to derive backend public key from blockchain private key"))
	}

	// 4. DUAL ENCRYPTION for DATA section (Participant Profile)
	// Parameter 5: encryptedUserData - Encrypted with user's wallet public key (ECIES)
	encryptedUserData, err := cyptoutils.EncryptWithPublicKeyBytes(attendeeProfileStr, participantPublicKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to encrypt attendee profile with public key"))
	}

	// Parameter 6: backendEncryptedUserData - Encrypted with backend public key (ECIES, derived from BLOCKCHAIN_PRIVATE_KEY)
	backendEncryptedUserData, err := cyptoutils.EncryptWithPublicKeyBytes(attendeeProfileStr, backendPublicKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to encrypt attendee profile with backend public key"))
	}

	// ============================================
	// PROOF SECTION: Certificate PII CSV
	// ============================================

	// 5. Create certificate PII CSV: name,academic_institution,certificate_title,certificate_subtitle
	// DO NOT mutate or merge fields - keep as is
	certificatePIIcsv := fmt.Sprintf("%s,%s,%s,%s", name, academicInstitution, certTitle, certSubtitle)

	// 6. DUAL ENCRYPTION for PROOF section (Certificate PII)
	// Parameter 13: userEncryptedProof - Encrypted with user's wallet public key (ECIES)
	userEncryptedProof, err := cyptoutils.EncryptWithPublicKeyBytes(certificatePIIcsv, participantPublicKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to encrypt certificate PII with public key"))
	}

	// Parameter 14: backendEncryptedProof - Encrypted with backend public key (ECIES, derived from BLOCKCHAIN_PRIVATE_KEY)
	backendEncryptedProof, err := cyptoutils.EncryptWithPublicKeyBytes(certificatePIIcsv, backendPublicKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to encrypt certificate PII with backend public key"))
	}

	// 7. Get certificate digest from database (pre-computed hash)
	// This should match the hash stored when the certificate was created
	if certificate.CertificateDigest == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate digest is missing"))
	}
	userDataHashStr := *certificate.CertificateDigest

	// Parameter 7: issuerAddresses (array of issuer wallet addresses)
	issuerAddresses := make([]common.Address, len(issuers))
	for i, issuer := range issuers {
		// Get issuer's wallet address from their credentials
		issuerCred, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, issuer.IssuerCredentialID)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "failed to get issuer %d credentials", i))
		}
		issuerAddresses[i] = common.HexToAddress(issuerCred.WalletAddress)
	}

	// Parameter 8, 9, 10, 12: signature data from first certificate signature
	firstSignature := certificateSignatures[0]
	if firstSignature.SignMessageDigest == nil || firstSignature.SignMessage == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate signature data is incomplete"))
	}

	signedMessageDigest := *firstSignature.SignMessageDigest
	hostSignatureStr := firstSignature.HostSignature
	signMessageStr := *firstSignature.SignMessage

	slog.Info("📝 Certificate signature data for contract verification",
		"signed_message_digest", signedMessageDigest,
		"host_signature", hostSignatureStr,
		"host_signature_length", len(hostSignatureStr),
		"sign_message", signMessageStr,
		"host_wallet_address", hostCredential.WalletAddress,
	)

	// Decode host signature from hex string to bytes
	hostSignatureBytes, err := hex.DecodeString(strings.TrimPrefix(hostSignatureStr, "0x"))
	if err != nil {
		slog.Error("❌ Failed to decode host signature", "error", err.Error())
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to decode host signature"))
	}

	slog.Info("🔐 Decoded host signature",
		"signature_bytes_length", len(hostSignatureBytes),
		"signature_hex", fmt.Sprintf("0x%x", hostSignatureBytes),
	)

	// Try to recover the signer to verify signature validity before sending to contract
	// Use the original message string since signatures are now created with HashEthereumMessage
	recoveredSigner, err := cyptoutils.GetAddressFromSignature(signMessageStr, hostSignatureStr)
	if err != nil {
		slog.Warn("⚠️ Could not recover signer from signature (this may be expected)", "error", err.Error())
	} else {
		signerMatch := strings.EqualFold(recoveredSigner.Hex(), hostCredential.WalletAddress)
		slog.Info("✅ Signature recovery test",
			"recovered_signer", recoveredSigner.Hex(),
			"expected_host_address", hostCredential.WalletAddress,
			"addresses_match", signerMatch,
		)
		if !signerMatch {
			slog.Error("❌ CRITICAL: Recovered signer does not match host wallet address! Contract will likely revert.")
		}
	}

	// Parameter 11: hostPublicKey
	hostPublicKey := hostCredential.WalletAddress // Using wallet address as public key identifier

	// Parameter 15 & 16: certificateTitle and certificateSubtitle
	certificateTitle := certTitle
	certificateSubtitle := certSubtitle

	// Parameter 17: issuerProofs (array of {issuerSignature, issuerPublicKey})
	type IssuerProof struct {
		IssuerSignature string
		IssuerPublicKey string
	}
	issuerProofs := make([]certificateContract.CertificateVCStructsIssuerProof, len(certificateSignatures))
	for i, sig := range certificateSignatures {
		issuerSig := ""
		if sig.IssuerSignature != nil {
			issuerSig = *sig.IssuerSignature
		}

		// Get issuer's public key (wallet address)
		issuerCred, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, sig.IssuerCredentialId)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "failed to get issuer credentials for proof %d", i))
		}

		issuerProofs[i] = certificateContract.CertificateVCStructsIssuerProof{
			IssuerSignature: issuerSig,
			IssuerPublicKey: issuerCred.WalletAddress,
		}
	}

	// ============================================
	// IDEMPOTENCY CHECK: Handle different states
	// ============================================

	// Check current state: Is NFT minted on-chain? Is DB updated?
	certificateContractInstance, err := certificateContract.NewEventCertificate(common.HexToAddress(*certificate.EventCertificateAddress), client)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to create contract instance"))
	}

	// Try to get token data from contract using certificateId as lookup
	// If it exists, NFT is already minted
	isNftMinted := false
	var onChainTokenId *uint64

	slog.Info("🔒 Checking idempotency state",
		"certificate_id", certificate.Id.String(),
		"db_token_id", certificate.CertificateTokenId,
		"contract_address", certificate.EventCertificateAddress,
	)

	// Check if token counter exists and try to find minted certificate
	// Note: This is a simplified check - in production you might want to query events
	// or use a mapping in the contract for certificateId -> tokenId lookup
	if certificate.CertificateTokenId != nil {
		// DB says it's minted, verify on-chain
		tokenId := new(big.Int)
		tokenId.SetString(*certificate.CertificateTokenId, 10)

		slog.Info("📋 Database indicates certificate already has token ID, verifying on-chain",
			"token_id", *certificate.CertificateTokenId,
		)

		// Try to get token data - if it succeeds, NFT exists
		_, err := certificateContractInstance.GetTokenData(nil, tokenId)
		if err == nil {
			isNftMinted = true
			tokenIdUint64 := tokenId.Uint64()
			onChainTokenId = &tokenIdUint64
			slog.Info("✅ NFT verified on-chain", "token_id", tokenIdUint64)
		} else {
			slog.Warn("⚠️ Database has token ID but NFT not found on-chain", "error", err.Error())
		}
	}

	// Decision matrix based on state
	isDbUpdated := certificate.CertificateTokenId != nil

	slog.Info("🔍 Idempotency check result",
		"is_nft_minted", isNftMinted,
		"is_db_updated", isDbUpdated,
	)

	// Case 1: NFT minted AND DB updated → Already claimed
	if isNftMinted && isDbUpdated {
		slog.Info("🚫 Certificate already claimed (Case 1: NFT minted + DB updated)")
		return nil, customerror.NewWithPreset(
			&customerror.ErrInvalidArgument,
			errors.New(string(ClaimCertificateUserErrorCertificateAlreadyClaimed)),
		)
	}

	// Case 2: NFT minted BUT DB not updated → Only update DB
	if isNftMinted && !isDbUpdated {
		slog.Info("🔄 NFT already minted but DB not updated (Case 2: Recovery mode)")
		if onChainTokenId == nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("token ID not found despite NFT being minted"))
		}

		tokenIdStr := fmt.Sprintf("%d", *onChainTokenId)

		// Update database only
		updatedCert, err := uc.EventCertificateDataGateway.UpdateEventCertificate(ctx, certificate.Id, eventdatagateway.UpdateEventCertificateParameters{
			CertificateTokenID: &tokenIdStr,
		})
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to update certificate with existing token ID"))
		}

		slog.Info("✅ Database updated with existing token ID", "token_id", tokenIdStr)
		return updatedCert, nil
	}

	// Case 3: NFT NOT minted (regardless of DB state) → Mint and update DB
	slog.Info("🎯 Proceeding to mint NFT (Case 3: NFT not minted)")
	// ============================================
	// MINT NFT ON BLOCKCHAIN
	// ============================================

	transactor, err := cyptoutils.GetKeyedTransactor()
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get system transactor"))
	}

	// ============================================
	// LOG MINTING PARAMETERS FOR DEBUGGING
	// ============================================
	logger := slog.With(
		"certificate_id", certificateId,
		"receiver_address", receiverAddress.Hex(),
		"contract_address", certificate.EventCertificateAddress,
		"transactor_address", transactor.From.Hex(),
		"user_id", userId,
		"issuer_id", issuerId,
		"num_issuer_addresses", len(issuerAddresses),
		"num_issuer_proofs", len(issuerProofs),
		"certificate_title", certificateTitle,
		"certificate_subtitle", certificateSubtitle,
		"data_hash", userDataHashStr,
	)

	logger.Info("🚀 Attempting to mint certificate NFT")
	logger.Debug("Mint parameters",
		"encrypted_user_data_length", len(encryptedUserData),
		"backend_encrypted_user_data_length", len(backendEncryptedUserData),
		"user_encrypted_proof_length", len(userEncryptedProof),
		"backend_encrypted_proof_length", len(backendEncryptedProof),
		"sign_message_digest", signedMessageDigest,
		"host_signature_length", len(hostSignatureBytes),
	)

	// The contract's recoverSigner expects the original message string (not the digest)
	// because it applies MessageHashUtils.toEthSignedMessageHash() to the input
	// We sign with HashEthereumMessage which matches what the contract will compute
	tx, err := certificateContractInstance.MintNft(
		transactor,
		receiverAddress,
		userId,
		certificateId,
		issuerId,
		encryptedUserData,
		backendEncryptedUserData,
		issuerAddresses,
		signMessageStr, // Pass original message, contract will hash it with Ethereum prefix
		hostSignatureBytes,
		hostSignatureStr,
		hostPublicKey,
		signMessageStr,
		userEncryptedProof,
		backendEncryptedProof,
		certificateTitle,
		certificateSubtitle,
		userDataHashStr, // Hash of the CSV data for blockchain verification
		issuerProofs,
	)
	if err != nil {
		logger.Error("❌ Failed to mint certificate NFT", "error", err.Error())
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to mint certificate NFT"))
	}

	logger.Info("✅ Transaction submitted successfully", "tx_hash", tx.Hash().Hex())

	logger.Info("⏳ Waiting for transaction to be mined...")
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		logger.Error("❌ Transaction mining failed", "tx_hash", tx.Hash().Hex(), "error", err.Error())
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "transaction mining failed: tx=%s", tx.Hash().Hex()))
	}

	logger.Info("📦 Transaction mined", "tx_hash", tx.Hash().Hex(), "gas_used", receipt.GasUsed, "status", receipt.Status)

	if receipt.Status != types.ReceiptStatusSuccessful {
		logger.Error("❌ Transaction reverted on-chain",
			"tx_hash", tx.Hash().Hex(),
			"gas_used", receipt.GasUsed,
			"block_number", receipt.BlockNumber,
		)
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Errorf("transaction reverted (tx=%s, gas=%d)", tx.Hash().Hex(), receipt.GasUsed))
	}

	// ============================================
	// EXTRACT TOKEN ID FROM EVENT LOGS
	// ============================================

	logger.Info("🔍 Parsing transaction logs for CertificateMinted event", "num_logs", len(receipt.Logs))

	// Parse CertificateMinted event to get tokenId
	var mintedTokenId *big.Int
	for i, log := range receipt.Logs {
		event, err := certificateContractInstance.ParseCertificateMinted(*log)
		if err != nil {
			// Not the event we're looking for, continue
			logger.Debug("Skipping log entry", "log_index", i, "topics_count", len(log.Topics))
			continue
		}
		// Found the CertificateMinted event
		mintedTokenId = event.TokenId
		logger.Info("✅ Found CertificateMinted event", "token_id", mintedTokenId.String(), "log_index", i)
		break
	}

	if mintedTokenId == nil {
		logger.Error("❌ Failed to find CertificateMinted event in transaction logs")
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("failed to extract token ID from minting event"))
	}

	// ============================================
	// UPDATE DATABASE WITH TOKEN ID
	// ============================================

	tokenIdStr := mintedTokenId.String()

	updatedCertificate, err := uc.EventCertificateDataGateway.UpdateEventCertificate(ctx, certificate.Id, eventdatagateway.UpdateEventCertificateParameters{
		CertificateTokenID: &tokenIdStr,
	})
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "NFT minted (tokenId=%s, tx=%s) but failed to update database", tokenIdStr, tx.Hash().Hex()))
	}

	return updatedCertificate, nil
}
