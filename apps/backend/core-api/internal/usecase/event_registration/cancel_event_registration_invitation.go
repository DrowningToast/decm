package event_registration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (uc *EventRegistrationUsecase) CancelEventRegistrationInvitation(ctx context.Context, params CancelEventRegistrationInvitationParameters, currentUser *auth.JwtClaims) (*entity.EventRegistrationInvitation, error) {
	// Check if invitation exists
	invitation, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByID(ctx, params.EventRegistrationInvitationID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, err)
	}

	// Check if invitation is already cancelled
	if invitation.CancelledAt != nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invitation is already cancelled"))
	}

	// Get the original inbox message to retrieve receiver information
	originalInboxMessage, err := uc.InboxMessageDg.GetInboxMessageByID(ctx, invitation.InboxMessageId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("failed to get original inbox message"))
	}

	// Get event to retrieve event title
	event, err := uc.EventDg.GetEventById(ctx, invitation.EventId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("failed to get event"))
	}
	if event == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event not found"))
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

	// Create message content with translations for revocation notification
	messageContent := map[string]string{
		"th": fmt.Sprintf("คำเชิญเข้าร่วมงาน %s ของคุณถูกยกเลิก", event.Title),
		"en": fmt.Sprintf("Your %s invitation has been revoked", event.Title),
	}
	messageContentJSON, err := json.Marshal(messageContent)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("failed to marshal message content"))
	}

	// Get receiver email from invitation or original inbox message
	receiverEmail := ""
	if invitation.Email != nil {
		receiverEmail = *invitation.Email
	} else if originalInboxMessage.ReceiverEmail != nil {
		receiverEmail = *originalInboxMessage.ReceiverEmail
	}

	if receiverEmail == "" {
		// If we can't determine the receiver, continue without sending the message
		// This shouldn't happen in normal flow, but we don't want to fail the cancellation
		return updatedInvitation, nil
	}

	// Create inbox message to notify about revocation
	inboxMessageParams := datagateway.CreateInboxMessageParameters{
		SenderCredentialID:     &currentUser.UserId,
		ReceiverCredentialID:   originalInboxMessage.ReceiverCredentialId,
		ReceiverEmail:          receiverEmail,
		MessageType:            int(entity.InboxMessageTypeGeneral), // General type
		MessageContent:         string(messageContentJSON),
		FallbackMessageContent: stringPtr(fmt.Sprintf("Your %s invitation has been revoked", event.Title)),
		IsRead:                 0, // Unread
	}

	_, err = uc.InboxMessageDg.CreateInboxMessage(ctx, inboxMessageParams)
	if err != nil {
		// Log error but don't fail the cancellation - the invitation was already cancelled
		// Returning an error here would be confusing since the cancellation succeeded
		// In production, you might want to log this error for monitoring
	}

	return updatedInvitation, nil
}
