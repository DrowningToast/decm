package entity

import (
	"time"

	"github.com/google/uuid"
)

type EventAttendee struct {
	ID                   uuid.UUID `json:"id"`
	EventID              uuid.UUID `json:"event_id"`
	AttendeeCredentialID uuid.UUID `json:"attendee_credential_id"`
	ContractAddress      string    `json:"contract_address"`
	IsAttendeeAccepted   bool      `json:"is_attendee_accepted"`
	FirstName            *string   `json:"first_name"`
	LastName             *string   `json:"last_name"`
	Email                *string   `json:"email"`
	Bio                  *string   `json:"bio"`
	PhoneNumber          *string   `json:"phone_number"`
	Address              *string   `json:"address"`
	AcademicInstitution  *string   `json:"academic_institution"`
	AcademicEmail        *string   `json:"academic_email"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
