package eventconfig

import (
	"errors"
	"time"

	"apps/backend/common/customerror"

	"github.com/google/uuid"
)

type EventRegistrationConfigRequirementStatus int32

const (
	EventRegistrationConfigRequirementStatusNotRequired EventRegistrationConfigRequirementStatus = 0
	EventRegistrationConfigRequirementStatusRequired    EventRegistrationConfigRequirementStatus = 1
	EventRegistrationConfigRequirementStatusOptional    EventRegistrationConfigRequirementStatus = 2
)

func (s EventRegistrationConfigRequirementStatus) String() (*string, error) {
	switch s {
	case EventRegistrationConfigRequirementStatusNotRequired:
		str := "not_required"
		return &str, nil
	case EventRegistrationConfigRequirementStatusRequired:
		str := "required"
		return &str, nil
	case EventRegistrationConfigRequirementStatusOptional:
		str := "optional"
		return &str, nil
	default:
		return nil, customerror.NewWithPreset(&customerror.ErrInvalidArgument, errors.New("invalid event registration config requirement status"))
	}
}

type EventRegistrationConfigResponse struct {
	ID                                   uuid.UUID                                `json:"id"`
	EventID                              uuid.UUID                                `json:"event_id"`
	FinalCallForRegistration             *time.Time                               `json:"final_call_for_registration,omitempty"`
	RegistrationPassword                 *string                                  `json:"registration_password,omitempty"`
	FirstNameRequirementStatus           EventRegistrationConfigRequirementStatus `json:"first_name_requirement_status,omitempty"`
	LastNameRequirementStatus            EventRegistrationConfigRequirementStatus `json:"last_name_requirement_status"`
	EmailRequirementStatus               EventRegistrationConfigRequirementStatus `json:"email_requirement_status"`
	BioRequirementStatus                 EventRegistrationConfigRequirementStatus `json:"bio_requirement_status"`
	PhoneNumberRequirementStatus         EventRegistrationConfigRequirementStatus `json:"phone_number_requirement_status"`
	AddressRequirementStatus             EventRegistrationConfigRequirementStatus `json:"address_requirement_status"`
	AcademicInstitutionRequirementStatus EventRegistrationConfigRequirementStatus `json:"academic_institution_requirement_status"`
	AcademicEmailRequirementStatus       EventRegistrationConfigRequirementStatus `json:"academic_email_requirement_status"`
	CreatedAt                            string                                   `json:"created_at"`
	UpdatedAt                            string                                   `json:"updated_at"`
}
