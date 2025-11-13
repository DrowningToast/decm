package entity

import (
	"time"

	"github.com/google/uuid"
)

type EventRegistrationInvitation struct {
	Id                  uuid.UUID  `json:"id"`
	EventId             uuid.UUID  `json:"event_id"`
	InboxMessageId      uuid.UUID  `json:"inbox_message_id"`
	ValidUntil          *time.Time `json:"valid_until,omitempty"`
	Code                *string    `json:"code,omitempty"`
	FirstName           *string    `json:"first_name,omitempty"`
	LastName            *string    `json:"last_name,omitempty"`
	Email               *string    `json:"email,omitempty"`
	PhoneNumber         *string    `json:"phone_number,omitempty"`
	AcademicInstitution *string    `json:"academic_institution,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	CancelledAt         *time.Time `json:"cancelled_at,omitempty"`
}
