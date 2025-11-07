package event_registration_invitation

import (
	"context"
	"errors"
	"time"

	"apps/backend/common/customerror"
	datagateway "apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
)

type CancelEventRegistrationInvitationParameters struct {
	EventRegistrationInvitationID uuid.UUID
}

func (uc *EventRegistrationInvitationUsecase) CancelEventRegistrationInvitation(ctx context.Context, params CancelEventRegistrationInvitationParameters, currentUser *auth.JwtClaims) (*entity.EventRegistrationInvitation, error) {
	// Check if invitation exists
	invitation, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByID(ctx, params.EventRegistrationInvitationID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, err)
	}

	// Check if invitation is already cancelled
	if invitation.CancelledAt != nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invitation is already cancelled"))
	}

	// Cancel the invitation by setting cancelled_at timestamp
	now := time.Now()
	updateParams := datagateway.UpdateEventRegistrationInvitationParameters{
		CancelledAt: &now,
	}

	updatedInvitation, err := uc.EventRegistrationInvitationDg.UpdateEventRegistrationInvitation(ctx, params.EventRegistrationInvitationID, updateParams)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return updatedInvitation, nil
}
