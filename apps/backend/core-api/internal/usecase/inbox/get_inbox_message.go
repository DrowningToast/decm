package inbox

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type GetInboxMessageResult struct {
	InboxMessage                entity.InboxMessage
	EventRegistrationInvitation *entity.EventRegistrationInvitation
	EventAttendee               *entity.EventAttendee

	EventCertificate *entity.EventCertificate

	Event *entity.Event
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
		eventRegistrationInvitation, event, eventAttendee, err := uc.GetRelatedEventRegistrationInvitation(ctx, message.Id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get related event registration invitation")
		}
		return &GetInboxMessageResult{
			InboxMessage:                *message,
			EventRegistrationInvitation: eventRegistrationInvitation,
			Event:                       event,
			EventAttendee:               eventAttendee,
		}, nil
	case entity.InboxMessageTypeEventCertificateInvitation:
		eventCertificate, event, err := uc.GetRelatedEventCertificateInvitation(ctx, message.Id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get related event certificate invitation")
		}
		return &GetInboxMessageResult{
			InboxMessage:     *message,
			EventCertificate: eventCertificate,
			Event:            event,
		}, nil

	}

	return &GetInboxMessageResult{
		InboxMessage: *message,
	}, nil
}

func (uc *InboxUsecase) GetRelatedEventRegistrationInvitation(ctx context.Context, inboxId uuid.UUID) (*entity.EventRegistrationInvitation, *entity.Event, *entity.EventAttendee, error) {
	eventRegistrationInvitation, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByInboxMessageID(ctx, inboxId)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to get event registration invitation by inbox message id")
	}
	if eventRegistrationInvitation == nil {
		return nil, nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event registration invitation not found"))
	}
	event, err := uc.EventDg.GetEventById(ctx, eventRegistrationInvitation.EventId)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to get event by id")
	}
	if event == nil {
		return nil, nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event not found"))
	}
	// check if the invited is onboarded or not
	targetUser, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddress(ctx, datagateway.GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddressParameters{
		GoogleConnectorRef: eventRegistrationInvitation.Email,
		// TODO: Allow wallet address to be used as well
		WalletAddress: nil,
	})
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to get authentication credential by email and wallet address")
	}
	if targetUser == nil {
		return nil, nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("target user not found"))
	}
	eventAttendee, err := uc.EventAttendeeDg.GetEventAttendeeByEventIdAndCredentialId(ctx, event.Id, targetUser.Id)
	if err != nil {
		var customError *customerror.Err
		if errors.As(err, &customError) {
			if *customError.Code != customerror.ErrNotFound.Code {
				return nil, nil, nil, errors.Wrap(err, "failed to get event attendee by event id and credential id")
			}
		} else {
			return nil, nil, nil, errors.Wrap(err, "failed to get event attendee by event id and credential id")
		}
	}
	return eventRegistrationInvitation, event, eventAttendee, nil
}

func (uc *InboxUsecase) GetRelatedEventCertificateInvitation(ctx context.Context, inboxId uuid.UUID) (*entity.EventCertificate, *entity.Event, error) {
	eventCertificate, err := uc.EventCertificateDg.GetEventCertificateByInboxMessageID(ctx, inboxId)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get event certificate by inbox message id")
	}
	if eventCertificate == nil {
		return nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event certificate not found"))
	}
	event, err := uc.EventDg.GetEventById(ctx, eventCertificate.EventId)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get event by id")
	}
	if event == nil {
		return nil, nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event not found"))
	}
	return eventCertificate, event, nil
}
