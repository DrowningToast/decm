package event_registration

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"

	"github.com/google/uuid"
)

func (uc *EventRegistrationUsecase) ImportEventParticipants(ctx context.Context, params ImportEventParticipantsParameters, currentUser *auth.JwtClaims) ([]*entity.EventRegistrationInvitation, error) {
	// Check if user is authorized to import participants for this event
	// This would typically involve checking if the user is the event owner or has admin privileges
	// For now, we'll assume the user is authorized if they have a valid JWT token

	dbEvent, err := uc.EventDg.GetEventById(ctx, params.EventID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	if dbEvent.OwnerCredentialId != currentUser.UserId {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not the event owner"))
	}

	if dbEvent.EventStatus != entity.EventStatusActive {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("event is not active"))
	}

	if len(params.Participants) == 0 {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("participants array is empty"))
	}

	invitations := make([]*entity.EventRegistrationInvitation, 0, len(params.Participants))

	// Process each participant
	for _, participant := range params.Participants {
		// Validate participant data
		if participant.Email == "" {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("email is required"))
		}

		// Create message content with translations
		messageContent := map[string]string{
			"th": fmt.Sprintf("คุณได้รับเชิญให้เข้าร่วมงาน: %s", dbEvent.Title),
			"en": fmt.Sprintf("You have been invited to join the event: %s", dbEvent.Title),
		}
		messageContentJSON, err := json.Marshal(messageContent)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}

		// Create inbox message
		inboxMessageParams := offchain_datagateway.CreateInboxMessageParameters{
			SenderCredentialID:     &currentUser.UserId,
			ReceiverCredentialID:   nil, // Empty as specified
			ReceiverEmail:          participant.Email,
			MessageType:            1, // event_registration_invitation
			MessageContent:         string(messageContentJSON),
			FallbackMessageContent: stringPtr("You have been invited to join the event"),
			IsRead:                 0, // Unread
		}

		inboxMessage, err := uc.InboxMessageDg.CreateInboxMessage(ctx, inboxMessageParams)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}

		// Generate 8 random characters code
		code := uuid.New().String()[:8]

		// Create event registration invitation
		invitationParams := event_datagateway.CreateEventRegistrationInvitationParameters{
			EventID:             params.EventID,
			InboxMessageID:      inboxMessage.Id,
			ValidUntil:          nil, // No expiration by default
			Code:                &code,
			FirstName:           stringPtr(participant.FirstName),
			LastName:            stringPtr(participant.LastName),
			Email:               stringPtr(participant.Email),
			PhoneNumber:         participant.PhoneNumber,
			AcademicInstitution: participant.AcademicInstitution,
		}

		invitation, err := uc.EventRegistrationInvitationDg.CreateEventRegistrationInvitation(ctx, invitationParams)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}

		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func stringPtr(s string) *string {
	return &s
}
