package event

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"time"

	"github.com/google/uuid"
)

// AbortCertificateClaim marks the pending signature as aborted with the given
// reason code so the blockchain submission worker will never retry it.  It is
// called when the worker determines the claim can never succeed — for example
// when the claimant has not joined the event on-chain and has no pending join
// in the queue.
func (uc *EventUsecase) AbortCertificateClaim(ctx context.Context, signatureId uuid.UUID, reason entity.UserSignatureAbortReason) error {
	_, err := uc.UserSignatureDg.UpdateUserSignatureAbortedAt(ctx, signatureId, time.Now(), reason)
	return err
}
