package event

import (
	"context"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
)

type GetEventRegistrationInvitationsByEventIDParameters struct {
	EventID uuid.UUID
}

func (uc *EventUsecase) GetEventRegistrationInvitationsByEventID(ctx context.Context, params GetEventRegistrationInvitationsByEventIDParameters, currentUser *auth.JwtClaims) ([]*entity.EventRegistrationInvitation, error) {
	// Get all invitations for the specified event
	invitations, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationsByEventID(ctx, params.EventID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return invitations, nil
}
