package entity

import (
	"time"

	"github.com/google/uuid"
)

// EventRegistrationConfig represents the registration configuration for an event
type EventRegistrationConfig struct {
	ID                                   uuid.UUID  `json:"id"`
	EventID                              uuid.UUID  `json:"event_id"`
	FinalCallForRegistration             *time.Time `json:"final_call_for_registration,omitempty"`
	RegistrationPassword                 *string    `json:"registration_password,omitempty"`
	IsIdentityVerificationRequired       bool       `json:"is_identity_verification_required"`
	FirstNameRequirementStatus           int        `json:"first_name_requirement_status"`
	LastNameRequirementStatus            int        `json:"last_name_requirement_status"`
	EmailRequirementStatus               int        `json:"email_requirement_status"`
	BioRequirementStatus                 int        `json:"bio_requirement_status"`
	PhoneNumberRequirementStatus         int        `json:"phone_number_requirement_status"`
	AddressRequirementStatus             int        `json:"address_requirement_status"`
	AcademicInstitutionRequirementStatus int        `json:"academic_institution_requirement_status"`
	AcademicEmailRequirementStatus       int        `json:"academic_email_requirement_status"`
	CreatedAt                            time.Time  `json:"created_at"`
	UpdatedAt                            time.Time  `json:"updated_at"`
}
