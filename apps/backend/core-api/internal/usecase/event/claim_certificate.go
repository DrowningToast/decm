package event

import (
	"context"
	"math/big"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

type CheckClaimEligibilityParams struct {
	CertificatePassword *string
}

// Performs checks if the user is eligible to claim the certificate
func (uc *EventUsecase) CheckClaimEligibility(ctx context.Context, certificate *entity.EventCertificate, currentUser *auth.JwtClaims, params CheckClaimEligibilityParams) (bool, error) {
	if currentUser == nil {
		return false, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	// Check if certificate is revoked
	if certificate.RevokedAt != nil {
		return false, customerror.NewWithPreset(&customerror.ErrInvalidArgument, errors.New(string(ClaimCertificateUserErrorCertificateRevoked)))
	}

	// Check if user is the intended receiver by credential ID
	if certificate.ReceiverCredentialId != nil {
		if *certificate.ReceiverCredentialId != currentUser.UserId {
			return false, customerror.NewWithPreset(&customerror.ErrUnauthorized, errors.New(string(ClaimCertificateUserErrorNotEligible)))
		}
		return true, nil
	}

	// Check if user is the intended receiver by email
	if certificate.ReceiverEmail != nil {
		if currentUser.Email == nil || *currentUser.Email != *certificate.ReceiverEmail {
			return false, customerror.NewWithPreset(&customerror.ErrUnauthorized, errors.New(string(ClaimCertificateUserErrorNotEligible)))
		}
		return true, nil
	}

	// If no receiver credential or email is set, anyone can claim (open certificate)
	// But if a password is required, validate it
	if params.CertificatePassword != nil {
		// TODO: Implement password validation logic if needed
		// This would require a certificate_password field in the database
		return true, nil
	}

	return true, nil
}

func (uc *EventUsecase) ClaimCertificateWithPin(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, certificateId uuid.UUID, eligibilityProof CheckClaimEligibilityParams, password string) (*entity.EventCertificate, error) {
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

	// Check if certificate contract address is set (certificate is published)
	if certificate.EventCertificateAddress == nil || certificate.CertificateTokenId == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate is not published yet"))
	}

	// Check eligibility
	isEligible, err := uc.CheckClaimEligibility(ctx, certificate, currentUser, eligibilityProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check claim eligibility")
	}
	if !isEligible {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New(string(ClaimCertificateUserErrorNotEligible)))
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

	return uc.claimCertificate(ctx, client, currentUser, certificate, signature, signMessage, participantAddress)
}

func (uc *EventUsecase) ClaimCertificateWithSignature(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, certificateId uuid.UUID, eligibilityProof CheckClaimEligibilityParams, signature []byte, signMessage string) (*entity.EventCertificate, error) {
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

	// Check if certificate contract address is set (certificate is published)
	if certificate.EventCertificateAddress == nil || certificate.CertificateTokenId == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate is not published yet"))
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

	// check if the user is eligible to claim the certificate
	isEligible, err := uc.CheckClaimEligibility(ctx, certificate, currentUser, eligibilityProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check claim eligibility")
	}
	if !isEligible {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New(string(ClaimCertificateUserErrorNotEligible)))
	}

	return uc.claimCertificate(ctx, client, currentUser, certificate, signature, signMessage, &participantAddress)
}

func (uc *EventUsecase) claimCertificate(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, certificate *entity.EventCertificate, signature []byte, signMessage string, participantAddress *common.Address) (*entity.EventCertificate, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}
	if participantAddress == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("participant address is required"))
	}

	if certificate.EventCertificateAddress == nil || certificate.CertificateTokenId == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate is not published yet"))
	}

	certificateContractInstance, err := certificateContract.NewEventCertificate(common.HexToAddress(*certificate.EventCertificateAddress), client)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Parse token ID
	tokenId := new(big.Int)
	tokenId, ok := tokenId.SetString(*certificate.CertificateTokenId, 10)
	if !ok {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("failed to parse certificate token id"))
	}

	// Check if the certificate NFT exists and get its owner
	owner, err := certificateContractInstance.OwnerOf(&bind.CallOpts{Context: ctx}, tokenId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get certificate owner"))
	}

	// Verify that the participant address matches the NFT owner
	if owner.Cmp(*participantAddress) != 0 {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("wallet address does not match certificate owner"))
	}

	// Check if already claimed on-chain
	// Note: In the contract, when participantSignedCertificate is called, an event is emitted
	// We can check the event logs or maintain a database flag for this
	// For now, we'll attempt the transaction and handle errors

	transactor, err := cyptoutils.GetKeyedTransactor()
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Call participantSignedCertificate on the contract
	tx, err := certificateContractInstance.ParticipantSignedCertificate(transactor, tokenId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "failed to claim certificate on blockchain: wallet=%s, contract=%s, tokenId=%s", participantAddress.Hex(), *certificate.EventCertificateAddress, *certificate.CertificateTokenId))
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "transaction mining failed: tx=%s", tx.Hash().Hex()))
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Errorf("transaction reverted (tx=%s, gas=%d)", tx.Hash().Hex(), receipt.GasUsed))
	}

	// TODO: Update database to mark certificate as claimed
	// This might involve updating the event_certificates table with a claimed_at timestamp
	// or updating the receiver_credential_id if it was previously null

	return certificate, nil
}

