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
func (uc *EventUsecase) GetClaimCertificateSignMessage(ctx context.Context, client *ethclient.Client, walletAddress common.Address, currentUser auth.JwtClaims, certificateContractAddress common.Address, deadlineBlock *uint64) (*string, *common.Hash, error) {
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

	// 2. Get certificate config to get config ID
	certificateConfig, err := uc.EventCertificateConfigDg.GetEventCertificateConfigByEventID(ctx, certificate.EventId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get certificate config"))
	}

	// 3. Get all certificate signatures (one per issuer)
	certificateSignatures, err := uc.EventCertificateSignatureDataGateway.GetEventCertificateSignaturesByEventCertificateConfigID(ctx, certificateConfig.ID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get certificate signatures"))
	}
	if len(certificateSignatures) != len(issuers) {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("not all issuers have signed the certificate"))
	}

	// 4. Get event details
	event, err := uc.EventDataGateway.GetEventById(ctx, certificate.EventId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get event"))
	}

	// 5. Get host credentials (for public key)
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

	// CRITICAL: Make a defensive copy of the signature to ensure it's never mutated
	// This is the signature that will be passed to the contract - it must remain unchanged
	participantSignatureBytes := make([]byte, len(signature))
	copy(participantSignatureBytes, signature)

	// Parameter 10: Host's signature of the raw JSON (storedSignMessageStr)
	// This is the host's signature of the original JSON message stored in the database
	hostSignatureStr := firstSignature.HostSignature
	// Ensure it has 0x prefix if not already present
	if !strings.HasPrefix(hostSignatureStr, "0x") {
		hostSignatureStr = "0x" + hostSignatureStr
	}

	// Parameter 11: hostPublicKey
	hostPublicKey := hostCredential.WalletAddress // Using wallet address as public key identifier

	// Parameter 15 & 16: certificateTitle and certificateSubtitle
	certificateTitle := certTitle
	certificateSubtitle := certSubtitle

	// Parameter 17: issuerProofs (array of {issuerSignature, issuerPublicKey})
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
	//
	// This section implements a three-tier idempotency system to handle:
	// 1. Normal flow: NFT not minted → Mint and update DB
	// 2. Recovery flow: NFT minted but DB not updated → Query events and update DB only
	// 3. Already claimed: NFT minted AND DB updated → Return error
	//
	// RECOVERY EFFICIENCY:
	// - Uses indexed event parameters (receiverAddress) for O(1) blockchain filtering
	// - CertificateMinted event has tokenId and receiverAddress indexed
	// - Filters by receiverAddress first, then matches certificateId (non-indexed string)
	// - This is much more efficient than iterating through all tokens or scanning all blocks
	//
	// Check current state: Is NFT minted on-chain? Is DB updated?
	certificateContractInstance, err := certificateContract.NewEventCertificate(common.HexToAddress(*certificate.EventCertificateAddress), client)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to create contract instance"))
	}
	// Try to get token data from contract using certificateId as lookup
	// If it exists, NFT is already minted
	isNftMinted := false
	var onChainTokenId *uint64

	// Check if token counter exists and try to find minted certificate
	// Strategy 1: If DB has token ID, verify it on-chain
	// Strategy 2: If DB doesn't have token ID, query blockchain events for this certificateId
	// This uses indexed event parameters (receiverAddress) for efficient O(1) filtering
	if certificate.CertificateTokenId != nil {
		// DB says it's minted, verify on-chain
		tokenId := new(big.Int)
		tokenId.SetString(*certificate.CertificateTokenId, 10)

		// Try to get token data - if it succeeds, NFT exists
		_, err := certificateContractInstance.GetTokenData(nil, tokenId)
		if err == nil {
			isNftMinted = true
			tokenIdUint64 := tokenId.Uint64()
			onChainTokenId = &tokenIdUint64
			slog.Info("Verified NFT exists on-chain from DB token ID", "token_id", tokenIdUint64, "certificate_id", certificate.Id)
		} else {
			slog.Warn("Database has token ID but NFT not found on-chain - possible token burn or invalid state",
				"token_id", *certificate.CertificateTokenId,
				"certificate_id", certificate.Id,
				"error", err.Error())
		}
	} else {
		// DB doesn't have token ID - query blockchain for CertificateMinted events
		// This handles the recovery case where minting succeeded but DB update failed
		// EFFICIENT: Uses indexed receiverAddress parameter for O(1) event filtering
		currentBlock, err := client.BlockNumber(ctx)
		if err != nil {
			slog.Warn("Failed to get current block number for recovery - skipping event query",
				"certificate_id", certificate.Id,
				"error", err.Error())
		} else {
			// Look back a reasonable number of blocks (2000 blocks = ~7 hours on most chains)
			// Most recovery scenarios happen within hours of minting
			// If needed, we can extend this or make it configurable
			fromBlock := uint64(0)
			lookbackBlocks := uint64(2000)
			if currentBlock > lookbackBlocks {
				fromBlock = currentBlock - lookbackBlocks
			}

			// Chunk block range queries to comply with RPC provider limits
			// Free tier providers (e.g., Alchemy, Infura) typically limit to 10 blocks per request
			// We query in chunks and aggregate results
			const maxBlockRange = uint64(10) // Free tier limit for eth_getLogs
			targetCertId := certificate.Id.String()
			receiverAddresses := []common.Address{*participantAddress}
			var matchingEvents []*certificateContract.EventCertificateCertificateMinted

			// Query in chunks of maxBlockRange blocks
			chunkStart := fromBlock
			for chunkStart <= currentBlock {
				chunkEnd := chunkStart + maxBlockRange - 1
				if chunkEnd > currentBlock {
					chunkEnd = currentBlock
				}

				filterOpts := &bind.FilterOpts{
					Start:   chunkStart,
					End:     &chunkEnd,
					Context: ctx,
				}

				// Filter by indexed receiverAddress for efficient O(1) lookup
				// Then iterate through filtered results to match certificateId (non-indexed string)
				iter, err := certificateContractInstance.FilterCertificateMinted(filterOpts, nil, receiverAddresses)
				if err != nil {
					slog.Warn("Failed to filter CertificateMinted events for recovery chunk",
						"certificate_id", certificate.Id,
						"receiver_address", participantAddress.Hex(),
						"from_block", chunkStart,
						"to_block", chunkEnd,
						"error", err.Error())
					// Continue to next chunk even if this one fails
					chunkStart = chunkEnd + 1
					continue
				}

				// Collect matching events from this chunk
				for iter.Next() {
					event := iter.Event
					if event.CertificateId == targetCertId {
						matchingEvents = append(matchingEvents, event)
					}
				}

				if iter.Error() != nil {
					slog.Warn("Error iterating CertificateMinted events in chunk",
						"certificate_id", certificate.Id,
						"from_block", chunkStart,
						"to_block", chunkEnd,
						"error", iter.Error())
				}

				_ = iter.Close()

				// If we found a match, we can stop searching (optional optimization)
				// But continue searching to detect duplicates
				chunkStart = chunkEnd + 1
			}

			// Handle matching events
			if len(matchingEvents) > 0 {
				if len(matchingEvents) > 1 {
					slog.Warn("Multiple CertificateMinted events found for certificate - using first match",
						"certificate_id", certificate.Id,
						"count", len(matchingEvents))
				}

				// Use the first (and typically only) matching event
				event := matchingEvents[0]
				tokenIdUint64 := event.TokenId.Uint64()

				// CRITICAL: Verify the token still exists on-chain before marking as minted
				// This handles edge cases like token burns or contract state changes
				tokenIdBigInt := new(big.Int).SetUint64(tokenIdUint64)
				_, err := certificateContractInstance.GetTokenData(nil, tokenIdBigInt)
				if err != nil {
					slog.Warn("Found CertificateMinted event but token no longer exists on-chain - possible burn",
						"certificate_id", certificate.Id,
						"token_id", tokenIdUint64,
						"tx_hash", event.Raw.TxHash.Hex(),
						"error", err.Error())
					// Don't mark as minted if token doesn't exist
				} else {
					isNftMinted = true
					onChainTokenId = &tokenIdUint64
					slog.Info("Recovery: Found already-minted NFT via event query",
						"certificate_id", certificate.Id,
						"token_id", tokenIdUint64,
						"tx_hash", event.Raw.TxHash.Hex(),
						"block_number", event.Raw.BlockNumber,
						"receiver_address", event.ReceiverAddress.Hex())
				}
			} else {
				slog.Debug("No CertificateMinted events found for certificate in lookback window",
					"certificate_id", certificate.Id,
					"receiver_address", participantAddress.Hex(),
					"from_block", fromBlock,
					"to_block", currentBlock)
			}
		}
	}

	// Decision matrix based on state
	isDbUpdated := certificate.CertificateTokenId != nil

	// Case 1: NFT minted AND DB updated → Already claimed
	if isNftMinted && isDbUpdated {
		return nil, customerror.NewWithPreset(
			&customerror.ErrInvalidArgument,
			errors.New(string(ClaimCertificateUserErrorCertificateAlreadyClaimed)),
		)
	}

	// Case 2: NFT minted BUT DB not updated → Only update DB (Recovery Flow)
	// This handles the scenario where:
	// - Minting transaction succeeded on-chain
	// - Database update failed or was interrupted
	// - User retries claiming, and we recover the token ID from blockchain events
	if isNftMinted && !isDbUpdated {
		slog.Info("Recovery mode: NFT minted on-chain but DB not updated - syncing database",
			"certificate_id", certificate.Id,
			"token_id", onChainTokenId)

		if onChainTokenId == nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer,
				errors.New("token ID not found despite NFT being minted - this should not happen"))
		}

		// Final verification: Double-check token exists before updating DB
		tokenIdBigInt := new(big.Int).SetUint64(*onChainTokenId)
		_, err := certificateContractInstance.GetTokenData(nil, tokenIdBigInt)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer,
				errors.Wrapf(err, "token ID %d found in events but token does not exist on-chain - possible burn or invalid state", *onChainTokenId))
		}

		tokenIdStr := fmt.Sprintf("%d", *onChainTokenId)

		// Update database with recovered token ID
		// CRITICAL: Preserve all existing fields to avoid null constraint violations
		updatedCert, err := uc.EventCertificateDataGateway.UpdateEventCertificate(ctx, certificate.Id, eventdatagateway.UpdateEventCertificateParameters{
			ReceiverCredentialID:    certificate.ReceiverCredentialId,
			ReceiverEmail:           certificate.ReceiverEmail,
			Name:                    certificate.Name,
			AcademicInstitution:     certificate.AcademicInstitution,
			CertificateTitle:        certificate.CertificateTitle,
			CertificateSubtitle:     certificate.CertificateSubtitle,
			EventContractAddress:    &certificate.EventContractAddress,
			EventCertificateAddress: certificate.EventCertificateAddress,
			CertificateTokenID:      &tokenIdStr,
			RevokedAt:               certificate.RevokedAt,
			Digest:                  certificate.CertificateDigest,
		})
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer,
				errors.Wrapf(err, "failed to update certificate with recovered token ID %s", tokenIdStr))
		}

		slog.Info("Recovery successful: Database synced with on-chain token ID",
			"certificate_id", certificate.Id,
			"token_id", tokenIdStr)

		return updatedCert, nil
	}

	// Case 3: NFT NOT minted (regardless of DB state) → Mint and update DB
	// ============================================
	// VALIDATE PARTICIPANT SIGN MESSAGE MATCHES CONTRACT
	// ============================================
	// CRITICAL: Verify that the participant's sign message uses the correct contract address
	// This must match the certificate contract address used in the contract call
	certificateContractAddr := common.HexToAddress(*certificate.EventCertificateAddress)
	deadlineBlockFromMessage, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(participantSignMessageStr)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "failed to extract deadline block from participant sign message"))
	}
	isValidMessage, err := cyptoutils.ValidateSignMessage(participantSignMessageStr, *participantAddress, certificateContractAddr, deadlineBlockFromMessage)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "failed to validate participant sign message"))
	}
	if !isValidMessage {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("participant sign message does not match certificate contract address or participant address"))
	}

	transactor, err := cyptoutils.GetKeyedTransactor()
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get system transactor"))
	}

	// Pre-flight checks
	var isSignatureUsed bool
	isSignatureUsed, err = certificateContractInstance.UsedSignatures(nil, participantSignatureBytes)
	if err != nil {
		slog.Warn("Could not check if signature was already used", "error", err.Error())
	} else if isSignatureUsed {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.New("participant signature has already been used in a previous transaction"))
	}

	if receiverAddress == (common.Address{}) {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.New("receiver address cannot be zero address"))
	}

	// Verify signature format
	if len(participantSignatureBytes) != 65 {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.Errorf("invalid signature length: got %d, expected 65", len(participantSignatureBytes)))
	}
	if participantSignatureBytes[64] != 27 && participantSignatureBytes[64] != 28 {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.Errorf("invalid signature v value: got %d, expected 27 or 28", participantSignatureBytes[64]))
	}

	// Verify signature locally before sending to contract
	contractMessageHash := cyptoutils.HashEthereumMessage(participantSignMessageStr)
	participantSignatureBytesCopy := make([]byte, len(participantSignatureBytes))
	copy(participantSignatureBytesCopy, participantSignatureBytes)

	localRecoveredPubKey, err := cyptoutils.RecoverPublicKeyFromSignature(contractMessageHash, participantSignatureBytesCopy)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.Wrapf(err, "local signature recovery failed"))
	}

	localRecoveredAddress := cyptoutils.PublicKeyToAddress(localRecoveredPubKey)
	if localRecoveredAddress.Hex() != participantAddress.Hex() {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument,
			errors.Errorf("signature verification failed: recovered %s, expected %s",
				localRecoveredAddress.Hex(), participantAddress.Hex()))
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
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to mint certificate NFT"))
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "transaction mining failed: tx=%s", tx.Hash().Hex()))
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Errorf("transaction reverted (tx=%s, gas=%d)", tx.Hash().Hex(), receipt.GasUsed))
	}

	// Extract token ID from event logs
	var mintedTokenId *big.Int
	for _, log := range receipt.Logs {
		event, err := certificateContractInstance.ParseCertificateMinted(*log)
		if err != nil {
			continue
		}
		mintedTokenId = event.TokenId
		break
	}

	if mintedTokenId == nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("failed to extract token ID from minting event"))
	}

	// ============================================
	// UPDATE DATABASE WITH TOKEN ID
	// ============================================

	tokenIdStr := mintedTokenId.String()

	updatedCertificate, err := uc.EventCertificateDataGateway.UpdateEventCertificate(ctx, certificate.Id, eventdatagateway.UpdateEventCertificateParameters{
		ReceiverCredentialID:    certificate.ReceiverCredentialId,
		ReceiverEmail:           certificate.ReceiverEmail,
		Name:                    certificate.Name,
		AcademicInstitution:     certificate.AcademicInstitution,
		CertificateTitle:        certificate.CertificateTitle,
		CertificateSubtitle:     certificate.CertificateSubtitle,
		EventContractAddress:    &certificate.EventContractAddress,
		EventCertificateAddress: certificate.EventCertificateAddress,
		CertificateTokenID:      &tokenIdStr,
		RevokedAt:               certificate.RevokedAt,
		Digest:                  certificate.CertificateDigest,
	})
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "NFT minted (tokenId=%s, tx=%s) but failed to update database", tokenIdStr, tx.Hash().Hex()))
	}

	return updatedCertificate, nil
}
