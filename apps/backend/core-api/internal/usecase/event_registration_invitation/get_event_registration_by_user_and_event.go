package event_registration_invitation

import (
	"context"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (uc *EventRegistrationInvitationUsecase) GetEventRegistrationByUserAndEvent(ctx context.Context, eventId uuid.UUID, currentUser *auth.JwtClaims) (*entity.EventRegistrationInvitation, *entity.InboxMessage, error) {
	// Get event registration invitation by user and event
	email := currentUser.Email
	walletAddress := currentUser.WalletAddress
	invitation, inbox, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByEventIDAndCredential(ctx, eventId, currentUser.UserId, email, &walletAddress)
	if err != nil {
		var customError *customerror.Err
		if errors.As(err, &customError) {
			if *customError.Code != customerror.ErrNotFound.Code {
				return nil, nil, errors.Wrap(err, "failed to get event registration invitation by event id and credential")
			}
		}
		// If not found, that's fine
	}
	return invitation, inbox, nil
}
