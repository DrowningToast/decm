package eventconfig

import (
	"github.com/google/uuid"
)

type EventRegistrationConfigResponse struct {
	ID                                   uuid.UUID `json:"id"`
	EventID                              uuid.UUID `json:"event_id"`
	FirstNameRequirementStatus           int32     `json:"first_name_requirement_status"`
	LastNameRequirementStatus            int32     `json:"last_name_requirement_status"`
	EmailRequirementStatus               int32     `json:"email_requirement_status"`
	BioRequirementStatus                 int32     `json:"bio_requirement_status"`
	PhoneNumberRequirementStatus         int32     `json:"phone_number_requirement_status"`
	AddressRequirementStatus             int32     `json:"address_requirement_status"`
	AcademicInstitutionRequirementStatus int32     `json:"academic_institution_requirement_status"`
	AcademicEmailRequirementStatus       int32     `json:"academic_email_requirement_status"`
	CreatedAt                            string    `json:"created_at"`
	UpdatedAt                            string    `json:"updated_at"`
}
