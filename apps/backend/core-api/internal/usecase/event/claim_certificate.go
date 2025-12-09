package event

import (
	"context"
	"crypto/ecdsa"
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

	// DEFENSIVE: Make a copy before verification (participant signature not used in contract,
	// but good practice to avoid mutation side effects)
	signatureCopy := make([]byte, len(signature))
	copy(signatureCopy, signature)

	// check if the signature matches the sign message or not (using copy)
	isValidHash, err := cyptoutils.VerifySignatureByDigest(participantAddress, messageHash, signatureCopy)
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

	// Parameter 8, 9, 10, 12: Use PARTICIPANT's signature for claim transaction
	// CRITICAL: The signature must come from the PARTICIPANT who is claiming the certificate,
	// NOT from the host/system. The participant proves authorization by signing the message.
	//
	// The participant's signature and signMessage were already created in ClaimCertificateWithPin
	// and passed to this function. We use those for parameters 8 and 9.

	// First, retrieve the stored signMessage from database for parameter 12 (metadata storage)
	// This is the original JSON body from the certificate signature table
	firstSignature := certificateSignatures[0]
	if firstSignature.SignMessage == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate signature signMessage is missing"))
	}
	if firstSignature.HostSignature == "" {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate signature hostSignature is missing"))
	}
	storedSignMessageStr := *firstSignature.SignMessage // This is the raw JSON from the database for parameter 12

	// Use the participant's signMessage and signature that were passed to this function
	// These were created in ClaimCertificateWithPin using the participant's private key
	participantSignMessageStr := signMessage // The message the participant signed
	participantSignatureBytes := signature   // The signature bytes from the participant

	// Parameter 10: Host's signature of the raw JSON (storedSignMessageStr)
	// This is the host's signature of the original JSON message stored in the database
	hostSignatureStr := firstSignature.HostSignature
	// Ensure it has 0x prefix if not already present
	if !strings.HasPrefix(hostSignatureStr, "0x") {
		hostSignatureStr = "0x" + hostSignatureStr
	}

	// Convert participant signature to hex string for logging/debugging (not used in contract call)
	participantSignatureStr := fmt.Sprintf("0x%x", participantSignatureBytes)

	slog.Info("📝 Using signatures for claim transaction",
		"participant_sign_message", participantSignMessageStr,
		"participant_sign_message_length", len(participantSignMessageStr),
		"participant_address", participantAddress.Hex(),
		"stored_sign_message", storedSignMessageStr,
		"stored_sign_message_length", len(storedSignMessageStr),
		"host_signature_hex", hostSignatureStr,
		"participant_signature_bytes_length", len(participantSignatureBytes),
		"note", "Participant signature used for param 8&9 (verification), host signature used for param 10 (metadata), stored message used for param 12 (metadata)",
	)

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
	// CREATE FRESH SIGNATURE FOR CLAIM TRANSACTION
	// ============================================
	// Get system transactor first (needed for signature creation and transaction)
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
		"participant_sign_message", participantSignMessageStr,
		"stored_sign_message", storedSignMessageStr,
		"participant_signature_length", len(participantSignatureBytes),
	)

	// Pre-flight check: Verify signature hasn't been used (replay prevention)
	// The contract's recoverSigner marks signatures as used, so if this was used before, it will revert
	var isSignatureUsed bool
	isSignatureUsed, err = certificateContractInstance.UsedSignatures(nil, participantSignatureBytes)
	if err != nil {
		slog.Warn("⚠️ Could not check if signature was already used", "error", err.Error())
		isSignatureUsed = false // Default to false if check fails (will be caught by contract if actually used)
	} else if isSignatureUsed {
		logger.Error("❌ CRITICAL: Participant signature has already been used! This will cause contract revert.",
			"participant_signature", participantSignatureStr,
			"note", "Each signature can only be used once. Participant must create a new signature to claim.",
		)
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.New("participant signature has already been used in a previous transaction"))
	} else {
		logger.Info("✅ Signature replay check passed", "signature_not_used", true)
	}

	// Note: Unlike join_event's addParticipant (which has no access control),
	// mintNft DOES check access control via requireHostOrAdmin(signer, msg.sender).
	// The signer must be the participant claiming the certificate, and the system transactor
	// must be authorized to send the transaction.
	logger.Info("🔐 Access control will be verified by contract",
		"participant_address", participantAddress.Hex(),
		"system_transactor", transactor.From.Hex(),
		"note", "Contract will verify: signer is participant AND transactor is allowed sender",
	)

	// Pre-flight check: Verify receiver address is not zero
	if receiverAddress == (common.Address{}) {
		logger.Error("❌ CRITICAL: Receiver address is zero address! Contract will revert.")
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.New("receiver address cannot be zero address"))
	}
	logger.Info("✅ Receiver address check passed", "receiver_address", receiverAddress.Hex())

	// DIAGNOSTIC: Check if receiver is a contract (must implement ERC721Receiver)
	code, err := client.CodeAt(ctx, receiverAddress, nil)
	if err == nil && len(code) > 0 {
		logger.Warn("⚠️ Receiver is a smart contract - must implement ERC721Receiver interface",
			"receiver", receiverAddress.Hex(),
			"code_length", len(code),
		)
	}

	// DIAGNOSTIC: Verify signature locally before sending to contract
	// This mimics exactly what the contract will do in ThemisUtils.recoverSigner
	slog.Info("🔬 LOCAL SIGNATURE VERIFICATION (mimics contract)",
		"participant_sign_message", participantSignMessageStr,
		"participant_sign_message_length", len(participantSignMessageStr),
		"participant_signature_hex", participantSignatureStr,
		"participant_signature_bytes_length", len(participantSignatureBytes),
		"signature_v_value", participantSignatureBytes[64],
	)

	// Verify signature V value
	if len(participantSignatureBytes) != 65 {
		logger.Error("❌ CRITICAL: Invalid signature length!", "length", len(participantSignatureBytes), "expected", 65)
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.Errorf("invalid signature length: got %d, expected 65", len(participantSignatureBytes)))
	}
	if participantSignatureBytes[64] != 27 && participantSignatureBytes[64] != 28 {
		logger.Error("❌ CRITICAL: Invalid signature V value!",
			"v", participantSignatureBytes[64],
			"expected", "27 or 28",
			"note", "Signature was mutated or created incorrectly",
		)
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.Errorf("invalid signature v value: got %d, expected 27 or 28", participantSignatureBytes[64]))
	}

	// Simulate contract's signature recovery (using participant's message for verification)
	contractMessageHash := cyptoutils.HashEthereumMessage(participantSignMessageStr)

	// Make a copy before verification to avoid mutation
	participantSignatureBytesCopy := make([]byte, len(participantSignatureBytes))
	copy(participantSignatureBytesCopy, participantSignatureBytes)

	localRecoveredPubKey, err := cyptoutils.RecoverPublicKeyFromSignature(contractMessageHash, participantSignatureBytesCopy)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.Wrapf(err, "local signature recovery failed - contract will reject this signature"))
	}

	localRecoveredAddress := cyptoutils.PublicKeyToAddress(localRecoveredPubKey)
	expectedParticipantAddress := *participantAddress

	// Verify that the signature was created with the participant's private key
	if localRecoveredAddress.Hex() != expectedParticipantAddress.Hex() {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.Errorf("signature verification failed: recovered %s, expected participant address %s",
				localRecoveredAddress.Hex(), expectedParticipantAddress.Hex()))
	}

	// The contract's recoverSigner expects the original message string (not the digest)
	// because it applies MessageHashUtils.toEthSignedMessageHash() to the input
	// We sign with HashEthereumMessage which matches what the contract will compute
	//
	// CRITICAL: The signature for parameters 8 and 9 MUST come from the PARTICIPANT
	// who is claiming the certificate, NOT from the host/system. The participant proves
	// authorization by signing the message with their private key.
	//
	// Parameter 10: Host's signature of the raw JSON (storedSignMessageStr) - this is the
	// host's signature of the original JSON message that was stored when certificates were imported.
	//
	// Comparison with update_event.go (working):
	// - update_event.go: Host signs message using GetSignMessage, hashes it, signs hash, passes original message to contract
	// - claim_certificate.go: PARTICIPANT signs message using GetSignMessage (params 8&9), HOST signature used for param 10 (metadata)
	participantSignMsgPreview := participantSignMessageStr
	if len(participantSignMsgPreview) > 100 {
		participantSignMsgPreview = participantSignMsgPreview[:100] + "..."
	}
	storedSignMsgPreview := storedSignMessageStr
	if len(storedSignMsgPreview) > 100 {
		storedSignMsgPreview = storedSignMsgPreview[:100] + "..."
	}
	tx, err := certificateContractInstance.MintNft(
		transactor,
		receiverAddress,
		userId,
		certificateId,
		issuerId,
		encryptedUserData,        // Param 5: Encrypted user data (attendee profile JSON) using user's public key
		backendEncryptedUserData, // Param 6: Encrypted user data (attendee profile JSON) using backend public key
		issuerAddresses,
		participantSignMessageStr, // Param 8: Participant's signed message for signature verification (contract will hash it with Ethereum prefix)
		participantSignatureBytes, // Param 9: Participant's signature bytes (signed with participant's private key)
		hostSignatureStr,          // Param 10: Host's signature hex string of the raw JSON (for metadata)
		hostPublicKey,             // Param 11: Host public key
		storedSignMessageStr,      // Param 12: Stored JSON message from database (original from certificate signature table)
		userEncryptedProof,        // Param 13: Encrypted user proof (ceretificate as .csv encrypted with user public key)
		backendEncryptedProof,     // Param 14: Encrypted user proof (ceretificate as .csv encrypted with backend public key)
		certificateTitle,
		certificateSubtitle,
		userDataHashStr, // Hash of the CSV data for blockchain verification
		issuerProofs,
	)
	if err != nil {
		logger.Error("❌ Failed to mint certificate NFT", "error", err.Error())

		// Try to extract detailed revert reason from error message
		revertReason := extractRevertReasonFromError(err)

		// Enhanced error diagnostics for transaction revert
		if strings.Contains(err.Error(), "execution reverted") || strings.Contains(err.Error(), "revert") {
			logger.Error("🔍 Transaction reverted - Detailed analysis:",
				"recovered_signer", hostCredential.WalletAddress,
				"system_transactor", transactor.From.Hex(),
				"receiver_address", receiverAddress.Hex(),
				"participant_signature", participantSignatureStr,
				"host_signature", hostSignatureStr,
				"participant_sign_message", participantSignMessageStr,
				"stored_sign_message", storedSignMessageStr,
				"participant_address", participantAddress.Hex(),
				"revert_reason_extracted", revertReason,
				"possible_issues", []string{
					"1. Signature replay: Participant signature already used (check signature usage)",
					"2. Invalid signature: Participant signature format/validity issue",
					"3. Access control: Participant not authorized OR system transactor not allowed",
					"4. Zero address: Receiver address is zero (unlikely - we check this)",
					"5. ERC721 hook failure: Receiver contract's onERC721Received failed",
					"6. Out of gas: Parameters too large (check string/array lengths)",
					"7. Contract mismatch: Sign message contract address doesn't match",
				},
				"preflight_checks", map[string]interface{}{
					"signature_used": isSignatureUsed,
					"receiver_zero":  receiverAddress == (common.Address{}),
				},
				"detected_errors", func() []string {
					errStr := err.Error()
					errors := []string{}
					if strings.Contains(errStr, "Themis__InvalidSignature") || strings.Contains(errStr, "InvalidSignature") {
						errors = append(errors, "Themis__InvalidSignature: Signature recovery failed")
					}
					if strings.Contains(errStr, "Themis__SignatureAlreadyUsed") || strings.Contains(errStr, "SignatureAlreadyUsed") {
						errors = append(errors, "Themis__SignatureAlreadyUsed: Signature replay")
					}
					if strings.Contains(errStr, "Not host or admin") || strings.Contains(errStr, "allowed msg sender") {
						errors = append(errors, "Access Control: Neither signer is host/admin NOR transactor is allowed")
					}
					return errors
				}(),
				"full_error_message", err.Error(),
			)
		}
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

// extractRevertReasonFromError attempts to extract the revert reason from an error
func extractRevertReasonFromError(err error) string {
	if err == nil {
		return ""
	}

	errMsg := err.Error()

	// Log full error for debugging
	slog.Debug("Full error message for revert reason extraction", "error", errMsg)

	// Try to extract revert reason from common patterns (check most specific first)
	patterns := []struct {
		prefix string
		suffix string
	}{
		{"execution reverted: ", "\n"},
		{"execution reverted:", "\n"},
		{"revert ", "\n"},
		{"revert: ", "\n"},
		{"VM execution error.\n\nrevert ", "\n"},
		{"revert ", ""},
		{"execution reverted: ", ""},
	}

	for _, pattern := range patterns {
		if idx := strings.Index(errMsg, pattern.prefix); idx != -1 {
			reason := strings.TrimSpace(errMsg[idx+len(pattern.prefix):])
			if pattern.suffix != "" {
				if suffixIdx := strings.Index(reason, pattern.suffix); suffixIdx != -1 {
					reason = reason[:suffixIdx]
				}
			}
			// Also trim common trailing parts
			reason = strings.TrimSpace(reason)
			if len(reason) > 0 && reason != "" {
				return reason
			}
		}
	}

	// Check for wrapped errors that might contain revert data
	if strings.Contains(errMsg, "revert") || strings.Contains(errMsg, "reverted") {
		// Try to find any text after revert/reverted
		for _, keyword := range []string{"revert ", "reverted ", "reverted: ", "reverted:"} {
			if idx := strings.Index(errMsg, keyword); idx != -1 {
				remaining := errMsg[idx+len(keyword):]
				// Take up to 200 chars or first newline
				if newlineIdx := strings.Index(remaining, "\n"); newlineIdx != -1 && newlineIdx < 200 {
					return strings.TrimSpace(remaining[:newlineIdx])
				}
				if len(remaining) > 200 {
					return strings.TrimSpace(remaining[:200]) + "..."
				}
				return strings.TrimSpace(remaining)
			}
		}
	}

	maxLen := 200
	if len(errMsg) < maxLen {
		maxLen = len(errMsg)
	}
	return "Could not extract revert reason - error message: " + errMsg[:maxLen]
}
