package event_registration_invitation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"apps/backend/common/customerror"
	datagateway "apps/backend/core-api/internal/datagateway"
	eventdatagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
)

type Participant struct {
	Name                string `json:"name"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Email               string `json:"email"`
	PhoneNumber         string `json:"phone_number"`
	AcademicInstitution string `json:"academic_institution"`
}

type ImportEventParticipantsParameters struct {
	EventID      uuid.UUID
	Participants []Participant
}

type EventRegistrationInvitationUsecase struct {
	InboxMessageDg                datagateway.InboxMessageDataGateway
	EventRegistrationInvitationDg datagateway.EventRegistrationInvitationDataGateway
	EventDg                       eventdatagateway.EventDataGateway
}

func NewEventRegistrationInvitationUsecase(
	inboxMessageDg datagateway.InboxMessageDataGateway,
	eventRegistrationInvitationDg datagateway.EventRegistrationInvitationDataGateway,
	eventDg eventdatagateway.EventDataGateway,
) *EventRegistrationInvitationUsecase {
	return &EventRegistrationInvitationUsecase{
		InboxMessageDg:                inboxMessageDg,
		EventRegistrationInvitationDg: eventRegistrationInvitationDg,
		EventDg:                       eventDg,
	}
}

func (uc *EventRegistrationInvitationUsecase) ImportEventParticipants(ctx context.Context, params ImportEventParticipantsParameters, currentUser *auth.JwtClaims) ([]*entity.EventRegistrationInvitation, error) {
	// Check if user is authorized to import participants for this event
	// This would typically involve checking if the user is the event owner or has admin privileges
	// For now, we'll assume the user is authorized if they have a valid JWT token

	dbEvent, err := uc.EventDg.GetEventById(ctx, params.EventID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	if dbEvent.OwnerCredentialID != currentUser.UserId {
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
		inboxMessageParams := datagateway.CreateInboxMessageParameters{
			SenderCredentialID:     &currentUser.UserId,
			ReceiverCredentialID:   nil, // Empty as specified
			ReceiverEmail:          participant.Email,
			MessageType:            1, // event_registration_invitation
			MessageContent:         string(messageContentJSON),
			FallbackMessageContent: stringPtr(fmt.Sprintf("You have been invited to join the event")),
			IsRead:                 0, // Unread
		}

		inboxMessage, err := uc.InboxMessageDg.CreateInboxMessage(ctx, inboxMessageParams)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}

		// Generate 8 random characters code
		code := uuid.New().String()[:8]

		// Create event registration invitation
		invitationParams := datagateway.CreateEventRegistrationInvitationParameters{
			EventID:             params.EventID,
			InboxMessageID:      inboxMessage.ID,
			ValidUntil:          nil, // No expiration by default
			Code:                &code,
			FirstName:           stringPtr(participant.FirstName),
			LastName:            stringPtr(participant.LastName),
			Email:               stringPtr(participant.Email),
			PhoneNumber:         stringPtr(participant.PhoneNumber),
			AcademicInstitution: stringPtr(participant.AcademicInstitution),
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
