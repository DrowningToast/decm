package entity

import (
	"time"

	"github.com/google/uuid"
)

type EventAttendee struct {
	Id                   uuid.UUID `json:"id"`
	EventId              uuid.UUID `json:"event_id"`
	AttendeeCredentialId uuid.UUID `json:"attendee_credential_id"`
	WalletAddress        string    `json:"wallet_address"`
	ContractAddress      string    `json:"contract_address"`
	IsAttendeeAccepted   bool      `json:"is_attendee_accepted"`
	FirstName            *string   `json:"first_name,omitempty"`
	LastName             *string   `json:"last_name,omitempty"`
	Email                *string   `json:"email,omitempty"`
	Bio                  *string   `json:"bio,omitempty"`
	PhoneNumber          *string   `json:"phone_number,omitempty"`
	Address              *string   `json:"address,omitempty"`
	AcademicInstitution  *string   `json:"academic_institution,omitempty"`
	AcademicEmail        *string   `json:"academic_email,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
