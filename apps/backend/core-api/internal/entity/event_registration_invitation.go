package entity

import (
	"time"

	"github.com/google/uuid"
)

type EventRegistrationInvitation struct {
	ID             uuid.UUID  `json:"id"`
	EventID        uuid.UUID  `json:"event_id"`
	InboxMessageID uuid.UUID  `json:"inbox_message_id"`
	ValidUntil     *time.Time `json:"valid_until"`
	Code           *string    `json:"code"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CancelledAt    *time.Time `json:"cancelled_at"`
}
