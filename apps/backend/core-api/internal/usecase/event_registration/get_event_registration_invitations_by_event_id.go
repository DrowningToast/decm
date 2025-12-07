package event_registration

import (
	"context"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
)

type GetEventRegistrationInvitationsByEventIDParameters struct {
	EventId uuid.UUID
}

func (uc *EventRegistrationUsecase) GetEventRegistrationInvitationsByEventId(ctx context.Context, params GetEventRegistrationInvitationsByEventIDParameters, currentUser *auth.JwtClaims) ([]*entity.EventRegistrationInvitation, error) {
	// Get all invitations for the specified event
	invitations, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationsByEventID(ctx, params.EventId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return invitations, nil
}
