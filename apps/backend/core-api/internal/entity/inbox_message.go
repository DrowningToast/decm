package entity

import (
	"time"

	"github.com/google/uuid"
)

type InboxMessageType int

const (
	InboxMessageTypeGeneral                     InboxMessageType = 0
	InboxMessageTypeEventRegistrationInvitation InboxMessageType = 1
	InboxMessageTypeEventCertificateInvitation  InboxMessageType = 2
)

func (t InboxMessageType) String() string {
	return []string{"event_registration_invitation", "event_certificate_invitation"}[t]
}

type InboxMessage struct {
	ID                     uuid.UUID  `json:"id"`
	SenderCredentialID     *uuid.UUID `json:"sender_credential_id"`
	ReceiverCredentialID   *uuid.UUID `json:"receiver_credential_id"`
	ReceiverWalletAddress  *string    `json:"receiver_wallet_address"`
	ReceiverEmail          *string    `json:"receiver_email"`
	MessageType            int        `json:"message_type"`
	MessageContent         string     `json:"message_content"`
	FallbackMessageContent *string    `json:"fallback_message_content"`
	IsRead                 int        `json:"is_read"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	HiddenAt               *time.Time `json:"hidden_at"`
	DeletedAt              *time.Time `json:"deleted_at"`
}
