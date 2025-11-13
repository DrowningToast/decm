package inbox

import (
	"context"
	"time"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/datagateway"
	eventdatagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
)

type InboxUsecase struct {
	AuthenticationCredentialDg    datagateway.AuthenticationCredentialDataGateway
	InboxMessageDg                datagateway.InboxMessageDataGateway
	EventRegistrationInvitationDg datagateway.EventRegistrationInvitationDataGateway
	EventDg                       eventdatagateway.EventDataGateway
}

func NewInboxUsecase(authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway, inboxMessageDg datagateway.InboxMessageDataGateway) *InboxUsecase {
	return &InboxUsecase{
		AuthenticationCredentialDg: authenticationCredentialDg,
		InboxMessageDg:             inboxMessageDg,
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

	// Check if receiver email matches
	if message.ReceiverEmail != nil && user.Email != nil && *message.ReceiverEmail == *user.Email {
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
	IsRead                        int                     `json:"is_read"`
	CreatedAt                     time.Time               `json:"created_at"`
	UpdatedAt                     time.Time               `json:"updated_at"`
	HiddenAt                      *time.Time              `json:"hidden_at,omitempty"`
	DeletedAt                     *time.Time              `json:"deleted_at,omitempty"`
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
		ID:                    inboxMessage.Id,
		ReceiverWalletAddress: inboxMessage.ReceiverWalletAddress,
		ReceiverEmail:         inboxMessage.ReceiverEmail,
		MessageType:           entity.InboxMessageType(inboxMessage.MessageType),
		MessageContent:        inboxMessage.MessageContent,
		IsRead:                inboxMessage.IsRead,
		CreatedAt:             inboxMessage.CreatedAt,
		UpdatedAt:             inboxMessage.UpdatedAt,
		HiddenAt:              inboxMessage.HiddenAt,
		DeletedAt:             inboxMessage.DeletedAt,
		SenderCredentialEmail: sender.GoogleConnectorRef,
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

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
}

func (uc *InboxUsecase) ToWithEventRegistrationInvitationViewModel(ctx context.Context, inboxMessage entity.InboxMessage, eventRegistrationInvitation entity.EventRegistrationInvitation, event entity.Event) (*InboxMessagesEventRegistrationInvitationViewModel, error) {
	inboxMessageViewModel, err := uc.ToViewModel(ctx, inboxMessage)
	if err != nil {
		return nil, err
	}
	return &InboxMessagesEventRegistrationInvitationViewModel{
		InboxMessagesViewModel: *inboxMessageViewModel,
		EventId:                event.ID,
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
	}, nil
}
