package event

import (
	"apps/backend/common/customerror"
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
)

func (uc *EventUsecase) queueCertificateClaim(ctx context.Context, currentUser *auth.JwtClaims, certificate *entity.EventCertificate, signature []byte, signMessage string, participantAddress *common.Address, participantPublicKey *ecdsa.PublicKey) (*entity.EventCertificate, *entity.UserSignature, error) {
	if currentUser == nil {
		return nil, nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}
	if participantAddress == nil {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("participant address is required"))
	}
	if participantPublicKey == nil {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("participant public key is required for data encryption"))
	}

	// Check if certificate contract is deployed
	if certificate.EventCertificateAddress == nil {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate contract not deployed yet"))
	}

	// Check if already minted (has token ID)
	if certificate.CertificateTokenId != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate already minted"))
	}

	// Check if signature already stored or not
	if certificate.UserClaimSignatureId != nil {
		// Check if signature is broadcasted or not
		userSignature, err := uc.UserSignatureDg.GetUserSignatureByID(ctx, *certificate.UserClaimSignatureId)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get user signature by id")
		}
		// Internal error, should not happen
		if userSignature == nil {
			return nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("user signature not found"))
		}
		if userSignature.BroadcastedAt != nil {
			// 400 Bad request, already broadcasted
			return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("signature already broadcasted"))
		}
		// Allow re-queue when the signature is provably expired OR aborted:
		//   • the worker explicitly marked it as expired (MarkAsExpiredAt set), OR
		//   • the estimated wall-clock deadline has passed, OR
		//   • the worker aborted it (e.g. participant not joined on-chain at the
		//     time the worker ran — the user may have since re-joined).
		// If none of these conditions are true, block so the pending signature
		// is not silently overwritten.
		isExpired := userSignature.MarkAsExpiredAt != nil ||
			(userSignature.EstimatedDeadline != nil && time.Now().After(*userSignature.EstimatedDeadline))
		isAborted := userSignature.AbortedAt != nil
		if !isExpired && !isAborted {
			return nil, nil, customerror.Parse(&customerror.ErrAlreadyPerformed, errors.New("signature not expired, please try again"))
		}
	}

	signatureAsString := hex.EncodeToString(signature)

	participantSignMessageStr := signMessage // The message the participant signed
	deadlineBlockFromMessage, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(participantSignMessageStr)
	// Deadline
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "failed to extract deadline block from participant sign message"))
	}
	deadlineBlock := int32(*deadlineBlockFromMessage)
	estimatedDeadline, err := uc.BlockchainClientDg.EstimateDeadlineTime(ctx, *deadlineBlockFromMessage)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to estimate deadline time"))
	}

	userSignature, err := uc.UserSignatureDg.CreateUserSignature(ctx, offchain_datagateway.CreateUserSignatureParameters{
		Signature:                  signatureAsString,
		SignMessage:                signMessage,
		AuthenticationCredentialId: currentUser.UserId,
		DeadlineBlock:              &deadlineBlock,
		EstimatedDeadline:          estimatedDeadline,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create user signature")
	}

	// Update certificate to be linked to the user signature.
	// NOTE: UpdateEventCertificate is a full-row UPDATE so every column must be
	// supplied; passing nil for a NOT NULL column (e.g. event_contract_address)
	// would cause a constraint violation.  Carry all existing values forward and
	// only change UserClaimSignatureId.
	certificate, err = uc.EventCertificateDataGateway.UpdateEventCertificate(ctx, certificate.Id, event_datagateway.UpdateEventCertificateParameters{
		ReceiverCredentialID:    certificate.ReceiverCredentialId,
		ReceiverEmail:           certificate.ReceiverEmail,
		Name:                    certificate.Name,
		AcademicInstitution:     certificate.AcademicInstitution,
		CertificateTitle:        certificate.CertificateTitle,
		CertificateSubtitle:     certificate.CertificateSubtitle,
		EventContractAddress:    &certificate.EventContractAddress,
		EventCertificateAddress: certificate.EventCertificateAddress,
		CertificateTokenID:      certificate.CertificateTokenId,
		RevokedAt:               certificate.RevokedAt,
		UserClaimSignatureId:    &userSignature.Id,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to update certificate to be linked to the user signature")
	}

	return certificate, userSignature, nil
}
