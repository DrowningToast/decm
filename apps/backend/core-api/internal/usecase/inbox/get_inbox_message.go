package inbox

import (
	"context"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type GetInboxMessageResult struct {
	InboxMessage                entity.InboxMessage
	EventRegistrationInvitation *entity.EventRegistrationInvitation
	Event                       *entity.Event
}

func (uc *InboxUsecase) GetInboxMessage(ctx context.Context, user auth.JwtClaims, messageID uuid.UUID) (*GetInboxMessageResult, error) {
	message, err := uc.InboxMessageDg.GetInboxMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("inbox message not found"))
	}

	// Check if user is authorized to read this message
	if !uc.isAuthorizedToReadMessage(message, user) {
		return nil, errors.Wrap(customerror.Parse(&customerror.ErrForbidden, errors.New("you are not allowed to read this message")), "user is not authorized to read this message")
	}

	switch entity.InboxMessageType(message.MessageType) {
	case entity.InboxMessageTypeEventRegistrationInvitation:
		eventRegistrationInvitation, event, err := uc.GetRelatedEventRegistrationInvitation(ctx, message.Id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get related event registration invitation")
		}
		return &GetInboxMessageResult{
			InboxMessage:                *message,
			EventRegistrationInvitation: eventRegistrationInvitation,
			Event:                       event,
		}, nil
	}

	return &GetInboxMessageResult{
		InboxMessage: *message,
	}, nil
}

func (uc *InboxUsecase) GetRelatedEventRegistrationInvitation(ctx context.Context, inboxId uuid.UUID) (*entity.EventRegistrationInvitation, *entity.Event, error) {
	eventRegistrationInvitation, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByInboxMessageID(ctx, inboxId)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get event registration invitation by inbox message id")
	}
	if eventRegistrationInvitation == nil {
		return nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event registration invitation not found"))
	}
	event, err := uc.EventDg.GetEventById(ctx, eventRegistrationInvitation.EventId)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get event by id")
	}
	if event == nil {
		return nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event not found"))
	}
	return eventRegistrationInvitation, event, nil
}
