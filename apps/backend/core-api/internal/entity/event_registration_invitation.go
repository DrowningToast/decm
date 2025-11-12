package entity

import (
	"time"

	"github.com/google/uuid"
)

type EventRegistrationInvitation struct {
	Id                  uuid.UUID  `json:"id"`
	EventId             uuid.UUID  `json:"event_id"`
	InboxMessageId      uuid.UUID  `json:"inbox_message_id"`
	ValidUntil          *time.Time `json:"valid_until"`
	Code                *string    `json:"code"`
	FirstName           *string    `json:"first_name"`
	LastName            *string    `json:"last_name"`
	Email               *string    `json:"email"`
	PhoneNumber         *string    `json:"phone_number"`
	AcademicInstitution *string    `json:"academic_institution"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	CancelledAt         *time.Time `json:"cancelled_at"`
}
