package event_registration

import (
	"apps/backend/common/customerror"
	"apps/backend/services/auth"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// AcceptEventRegistrationInvitation marks an event registration invitation as accepted
func (uc *EventRegistrationUsecase) AcceptEventRegistrationInvitation(ctx context.Context, invitationID uuid.UUID, currentUser *auth.JwtClaims) error {
	if currentUser == nil {
		return customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	// Get the invitation to verify it exists
	invitation, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByID(ctx, invitationID)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if invitation == nil {
		return customerror.Parse(&customerror.ErrNotFound, errors.New("invitation not found"))
	}

	// Verify the invitation belongs to the current user by checking if we can retrieve it
	// using the user's credentials (this verifies ownership through inbox message)
	email := currentUser.Email
	walletAddress := currentUser.WalletAddress
	userInvitation, _, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByEventIDAndCredential(ctx, invitation.EventId, currentUser.UserId, email, &walletAddress)
	if err != nil {
		return customerror.Parse(&customerror.ErrUnauthorized, errors.New("invitation does not belong to the current user"))
	}
	if userInvitation == nil || userInvitation.Id != invitationID {
		return customerror.Parse(&customerror.ErrUnauthorized, errors.New("invitation does not belong to the current user"))
	}

	// Mark invitation as accepted with current timestamp
	now := time.Now()
	_, err = uc.EventRegistrationInvitationDg.UpdateEventRegistrationInvitationAcceptedStatus(ctx, invitationID, &now)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return nil
}
