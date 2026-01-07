package inbox

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"strings"
	"time"

	eventdatagateway "apps/backend/core-api/internal/datagateway/event"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type InboxUsecase struct {
	AuthenticationCredentialDg    datagateway.AuthenticationCredentialDataGateway
	InboxMessageDg                datagateway.InboxMessageDataGateway
	EventRegistrationInvitationDg datagateway.EventRegistrationInvitationDataGateway
	EventDg                       eventdatagateway.EventDataGateway
	EventCertificateDg            eventdatagateway.EventCertificateDataGateway
	EventAttendeeDg               datagateway.EventAttendeeDataGateway
}

func NewInboxUsecase(
	authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway,
	inboxMessageDg datagateway.InboxMessageDataGateway,
	eventRegistrationInvitationDg datagateway.EventRegistrationInvitationDataGateway,
	eventDg eventdatagateway.EventDataGateway,
	eventCertificateDg eventdatagateway.EventCertificateDataGateway,
	eventAttendeeDg datagateway.EventAttendeeDataGateway,
) *InboxUsecase {
	return &InboxUsecase{
		AuthenticationCredentialDg:    authenticationCredentialDg,
		InboxMessageDg:                inboxMessageDg,
		EventRegistrationInvitationDg: eventRegistrationInvitationDg,
		EventDg:                       eventDg,
		EventCertificateDg:            eventCertificateDg,
		EventAttendeeDg:               eventAttendeeDg,
	}
}

// isAuthorizedToReadMessage checks if the user is authorized to read the message
// User is authorized if any of the following conditions are met:
// - ReceiverCredentialID matches user's ID
// - ReceiverEmail matches user's email
// - ReceiverWalletAddress matches user's wallet address
func (uc *InboxUsecase) isAuthorizedToReadMessage(message *entity.InboxMessage, user auth.JwtClaims) bool {
	// Check if receiver credential ID matches
	if message.ReceiverCredentialId != nil && *message.ReceiverCredentialId == user.UserId {
		return true
	}

	// Check if receiver email matches (case-insensitive)
	if message.ReceiverEmail != nil && user.Email != nil && strings.EqualFold(*message.ReceiverEmail, *user.Email) {
		return true
	}

	// Check if receiver wallet address matches
	if message.ReceiverWalletAddress != nil && *message.ReceiverWalletAddress == user.WalletAddress {
		return true
	}

	return false
}

type InboxMessagesViewModel struct {
	ID                            uuid.UUID               `json:"id"`
	SenderCredentialEmail         *string                 `json:"sender_credential_email,omitempty"`
	SenderCredentialWalletAddress *string                 `json:"sender_credential_wallet_address,omitempty"`
	ReceiverWalletAddress         *string                 `json:"receiver_wallet_address,omitempty"`
	ReceiverEmail                 *string                 `json:"receiver_email,omitempty"`
	MessageType                   entity.InboxMessageType `json:"message_type"`
	MessageContent                string                  `json:"message_content"`
	MessageContentFallback        *string                 `json:"message_content_fallback,omitempty"`
	IsRead                        int                     `json:"is_read"`
	CreatedAt                     time.Time               `json:"created_at"`
	UpdatedAt                     time.Time               `json:"updated_at"`
	HiddenAt                      *time.Time              `json:"hidden_at,omitempty"`
	DeletedAt                     *time.Time              `json:"deleted_at,omitempty"`

	// Certificate-specific fields (only populated for certificate invitation messages)
	EventId                   *uuid.UUID `json:"event_id,omitempty"`
	CertificateId             *uuid.UUID `json:"certificate_id,omitempty"`
	CertificateTitle          *string    `json:"certificate_title,omitempty"`
	TokenId                   *string    `json:"token_id,omitempty"`
	HasParticipantJoinedEvent *bool      `json:"has_participant_joined_event,omitempty"`
	EventName                 *string    `json:"event_name,omitempty"`
}

func (uc *InboxUsecase) ToViewModel(ctx context.Context, inboxMessage entity.InboxMessage) (*InboxMessagesViewModel, error) {
	sender, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, *inboxMessage.SenderCredentialId)
	if err != nil {
		return nil, err
	}
	if sender == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, err)
	}
	return &InboxMessagesViewModel{
		ID:                     inboxMessage.Id,
		ReceiverWalletAddress:  inboxMessage.ReceiverWalletAddress,
		ReceiverEmail:          inboxMessage.ReceiverEmail,
		MessageType:            entity.InboxMessageType(inboxMessage.MessageType),
		MessageContent:         inboxMessage.MessageContent,
		MessageContentFallback: inboxMessage.FallbackMessageContent,
		IsRead:                 inboxMessage.IsRead,
		CreatedAt:              inboxMessage.CreatedAt,
		UpdatedAt:              inboxMessage.UpdatedAt,
		HiddenAt:               inboxMessage.HiddenAt,
		DeletedAt:              inboxMessage.DeletedAt,
		SenderCredentialEmail:  sender.GoogleConnectorRef,
	}, nil
}

type InboxMessagesEventRegistrationInvitationViewModel struct {
	InboxMessagesViewModel

	EventId             uuid.UUID  `json:"event_id"`
	ValidUntil          *time.Time `json:"valid_until,omitempty"`
	Code                *string    `json:"code,omitempty"`
	FirstName           *string    `json:"first_name,omitempty"`
	LastName            *string    `json:"last_name,omitempty"`
	Email               *string    `json:"email,omitempty"`
	PhoneNumber         *string    `json:"phone_number,omitempty"`
	AcademicInstitution *string    `json:"academic_institution,omitempty"`

	AcceptedAt *time.Time `json:"accepted_at,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
}

func (uc *InboxUsecase) ToWithEventRegistrationInvitationViewModel(ctx context.Context, inboxMessage entity.InboxMessage, eventRegistrationInvitation entity.EventRegistrationInvitation, event entity.Event, eventAttendee *entity.EventAttendee) (*InboxMessagesEventRegistrationInvitationViewModel, error) {
	inboxMessageViewModel, err := uc.ToViewModel(ctx, inboxMessage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert inbox message to view model")
	}
	var acceptedAt *time.Time = nil
	if eventAttendee != nil {
		acceptedAt = &eventAttendee.CreatedAt
	}
	return &InboxMessagesEventRegistrationInvitationViewModel{
		InboxMessagesViewModel: *inboxMessageViewModel,
		EventId:                event.Id,
		ValidUntil:             eventRegistrationInvitation.ValidUntil,
		Code:                   eventRegistrationInvitation.Code,
		FirstName:              eventRegistrationInvitation.FirstName,
		LastName:               eventRegistrationInvitation.LastName,
		Email:                  eventRegistrationInvitation.Email,
		PhoneNumber:            eventRegistrationInvitation.PhoneNumber,
		AcademicInstitution:    eventRegistrationInvitation.AcademicInstitution,
		CreatedAt:              eventRegistrationInvitation.CreatedAt,
		UpdatedAt:              eventRegistrationInvitation.UpdatedAt,
		CancelledAt:            eventRegistrationInvitation.CancelledAt,
		AcceptedAt:             acceptedAt,
	}, nil
}

type InboxMessageCertificateInvitationViewModel struct {
	InboxMessagesViewModel

	EventName        string    `json:"event_name"`
	EventId          uuid.UUID `json:"event_id"`
	CertificateId    uuid.UUID `json:"certificate_id"`
	CertificateTitle *string   `json:"certificate_title,omitempty"`
	TokenId          *string   `json:"token_id,omitempty"`

	HasParticipantJoinedEvent bool `json:"has_participant_joined_event"`

	CreatedAt time.Time `json:"created_at"`

	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

func (uc *InboxUsecase) ToWithCertificateInvitationViewModel(ctx context.Context, inboxMessage entity.InboxMessage, eventCertificate entity.EventCertificate, user auth.JwtClaims) (*InboxMessageCertificateInvitationViewModel, error) {
	// verify event type
	if inboxMessage.MessageType != entity.InboxMessageTypeEventCertificateInvitation {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid message type"))
	}
	inboxMessageViewModel, err := uc.ToViewModel(ctx, inboxMessage)
	if err != nil {
		return nil, err
	}
	// check if the user has joined or not
	eventAttendee, err := uc.EventAttendeeDg.GetEventAttendeeByEventIdAndCredentialId(ctx, eventCertificate.EventId, user.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get event attendee by event id and credential id")
	}
	return &InboxMessageCertificateInvitationViewModel{
		InboxMessagesViewModel:    *inboxMessageViewModel,
		EventName:                 *eventCertificate.Name,
		CertificateId:             eventCertificate.Id,
		CertificateTitle:          eventCertificate.CertificateTitle,
		EventId:                   eventCertificate.EventId,
		CreatedAt:                 eventCertificate.CreatedAt,
		RevokedAt:                 eventCertificate.RevokedAt,
		TokenId:                   eventCertificate.CertificateTokenId,
		HasParticipantJoinedEvent: eventAttendee != nil,
	}, nil
}

// ToEnrichedCertificateViewModel enriches a basic view model with certificate-specific fields for list view
func (uc *InboxUsecase) ToEnrichedCertificateViewModel(ctx context.Context, inboxMessage entity.InboxMessage, user auth.JwtClaims) (*InboxMessagesViewModel, error) {
	// verify event type
	if inboxMessage.MessageType != entity.InboxMessageTypeEventCertificateInvitation {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid message type"))
	}

	inboxMessageViewModel, err := uc.ToViewModel(ctx, inboxMessage)
	if err != nil {
		return nil, err
	}

	// Get the certificate associated with this inbox message
	eventCertificate, err := uc.EventCertificateDg.GetEventCertificateByInboxMessageID(ctx, inboxMessage.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get event certificate by inbox message id")
	}
	if eventCertificate == nil {
		// If certificate not found, return basic view model without enrichment
		return inboxMessageViewModel, nil
	}

	// Check if the user has joined the event
	eventAttendee, err := uc.EventAttendeeDg.GetEventAttendeeByEventIdAndCredentialId(ctx, eventCertificate.EventId, user.UserId)
	if err != nil {
		// If error getting attendee, assume not joined
		hasJoined := false
		inboxMessageViewModel.HasParticipantJoinedEvent = &hasJoined
	} else {
		hasJoined := eventAttendee != nil
		inboxMessageViewModel.HasParticipantJoinedEvent = &hasJoined
	}

	// Enrich with certificate-specific fields
	eventName := ""
	if eventCertificate.Name != nil {
		eventName = *eventCertificate.Name
	}
	inboxMessageViewModel.EventId = &eventCertificate.EventId
	inboxMessageViewModel.CertificateId = &eventCertificate.Id
	inboxMessageViewModel.CertificateTitle = eventCertificate.CertificateTitle
	inboxMessageViewModel.TokenId = eventCertificate.CertificateTokenId
	inboxMessageViewModel.EventName = &eventName

	return inboxMessageViewModel, nil
}
