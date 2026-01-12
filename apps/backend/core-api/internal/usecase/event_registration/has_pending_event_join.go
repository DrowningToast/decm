package event_registration

import (
	"context"

	"github.com/google/uuid"
)

// HasPendingEventJoin returns true when there is an unbroadcasted, unexpired,
// unaborted event-join signature in the queue for the given event and
// authentication credential.  It delegates directly to the EventAttendee data
// gateway which owns the join-queue view.
func (uc *EventRegistrationUsecase) HasPendingEventJoin(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (bool, error) {
	return uc.EventAttendeeDg.HasPendingEventJoinByEventAndCredential(ctx, eventId, credentialId)
}
