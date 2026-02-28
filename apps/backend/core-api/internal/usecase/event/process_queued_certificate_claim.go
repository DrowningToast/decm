package event

import (
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"time"

	"github.com/cockroachdb/errors"
)

// ProcessQueuedCertificateClaim is the exported bridge called by the
// certificate claim worker. It recovers the participant's public key from the
// stored signature and delegates to the existing private claimCertificate
// logic, then marks the signature as broadcasted on success.
func (uc *EventUsecase) ProcessQueuedCertificateClaim(ctx context.Context, pending entity.PendingCertificateClaim) error {
	// 1. Fetch the full certificate record.
	certificate, err := uc.EventCertificateDataGateway.GetEventCertificateByID(ctx, pending.CertificateId)
	if err != nil {
		return errors.Wrap(err, "failed to get event certificate")
	}

	// 2. Idempotency guard: already minted by a previous attempt.
	if certificate.CertificateTokenId != nil {
		now := time.Now()
		_, _ = uc.UserSignatureDg.UpdateUserSignatureBroadcastedAt(ctx, pending.SignatureId, &now)
		return nil
	}

	// 3. Decode the hex-encoded signature back to raw bytes.
	sigBytes, err := hex.DecodeString(pending.Signature)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex signature")
	}

	// 4. Recover the participant's public key from the signature.
	messageHash := cyptoutils.HashEthereumMessage(pending.SignMessage)
	participantPublicKey, err := cyptoutils.RecoverPublicKeyFromSignature(messageHash, sigBytes)
	if err != nil {
		return errors.Wrap(err, "failed to recover public key from signature")
	}

	// 5. Derive the participant's Ethereum address.
	participantAddress := cyptoutils.PublicKeyToAddress(participantPublicKey)

	// 6. Build minimal JWT claims so claimCertificate can look up the attendee
	//    record using the stored credential ID.
	claims := &auth.JwtClaims{UserId: pending.AuthenticationCredentialId}

	// 7. Perform the on-chain mint.
	_, err = uc.claimCertificate(ctx, claims, certificate, sigBytes, pending.SignMessage, &participantAddress, participantPublicKey)
	if err != nil {
		return errors.Wrap(err, "failed to claim certificate")
	}

	// 8. Mark the signature as processed so the worker does not retry it.
	now := time.Now()
	if _, err = uc.UserSignatureDg.UpdateUserSignatureBroadcastedAt(ctx, pending.SignatureId, &now); err != nil {
		return errors.Wrap(err, "failed to update signature broadcasted_at")
	}

	return nil
}
