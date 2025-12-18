package event

import (
	"context"
	"decm-database/go/generated"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (u *EventUsecase) GetEventIssuersByEventID(ctx context.Context, eventID uuid.UUID) ([]generated.EventIssuer, error) {
	return u.EventIssuerDataGateway.GetEventIssuersByEventID(ctx, eventID)
}

func (u *EventUsecase) GetEventIssuerByID(ctx context.Context, id uuid.UUID) (generated.EventIssuer, error) {
	return u.EventIssuerDataGateway.GetEventIssuerByID(ctx, id)
}

// IsUserIssuerForEvent checks if a user is assigned as an issuer for a specific event.
// Returns true if the user is an issuer, false if not found, or an error for database issues.
func (u *EventUsecase) IsUserIssuerForEvent(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (bool, error) {
	_, err := u.EventIssuerDataGateway.GetEventIssuerByEventIDAndIssuerCredentialID(
		ctx,
		eventID,
		userID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// User is not an issuer for this event
			u.logger.InfoContext(ctx, "No row found in event_issuers table", "event_id", eventID, "user_id", userID)
			return false, nil
		}
		// Database or other error
		return false, errors.Wrap(err, "failed to check if user is an issuer for this event")
	}
	// User is an issuer for this event
	return true, nil
}
