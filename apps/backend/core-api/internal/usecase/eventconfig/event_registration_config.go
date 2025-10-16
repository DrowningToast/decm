package eventconfig

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"decm-database/go/generated"
)

type CreateEventRegistrationConfigParams struct {
	FirstNameRequirementStatus           pgtype.Int4
	LastNameRequirementStatus            pgtype.Int4
	EmailRequirementStatus               pgtype.Int4
	BioRequirementStatus                 pgtype.Int4
	PhoneNumberRequirementStatus         pgtype.Int4
	AddressRequirementStatus             pgtype.Int4
	AcademicInstitutionRequirementStatus pgtype.Int4
	AcademicEmailRequirementStatus       pgtype.Int4
}

func (uc *EventConfigUsecase) CreateEventRegistrationConfig(ctx context.Context, eventID uuid.UUID, params CreateEventRegistrationConfigParams) (*generated.EventRegistrationConfig, error) {
	// Check if config already exists for this event
	existingConfig, err := uc.EventRegistrationDg.GetEventRegistrationConfigByEventID(ctx, eventID)
	if err == nil && existingConfig != nil {
		return nil, fmt.Errorf("event registration config already exists for event ID: %s", eventID.String())
	}

	// Create new config
	createParams := generated.CreateEventRegistrationConfigParams{
		EventID:                              eventID,
		FirstNameRequirementStatus:           params.FirstNameRequirementStatus,
		LastNameRequirementStatus:            params.LastNameRequirementStatus,
		EmailRequirementStatus:               params.EmailRequirementStatus,
		BioRequirementStatus:                 params.BioRequirementStatus,
		PhoneNumberRequirementStatus:         params.PhoneNumberRequirementStatus,
		AddressRequirementStatus:             params.AddressRequirementStatus,
		AcademicInstitutionRequirementStatus: params.AcademicInstitutionRequirementStatus,
		AcademicEmailRequirementStatus:       params.AcademicEmailRequirementStatus,
	}

	return uc.EventRegistrationDg.CreateEventRegistrationConfig(ctx, createParams)
}

func (uc *EventConfigUsecase) GetEventRegistrationConfigByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventRegistrationConfig, error) {
	return uc.EventRegistrationDg.GetEventRegistrationConfigByEventID(ctx, eventID)
}

type UpdateEventRegistrationConfigParams struct {
	FirstNameRequirementStatus           pgtype.Int4
	LastNameRequirementStatus            pgtype.Int4
	EmailRequirementStatus               pgtype.Int4
	BioRequirementStatus                 pgtype.Int4
	PhoneNumberRequirementStatus         pgtype.Int4
	AddressRequirementStatus             pgtype.Int4
	AcademicInstitutionRequirementStatus pgtype.Int4
	AcademicEmailRequirementStatus       pgtype.Int4
}

func (uc *EventConfigUsecase) UpdateEventRegistrationConfig(ctx context.Context, eventID uuid.UUID, params UpdateEventRegistrationConfigParams) (*generated.EventRegistrationConfig, error) {
	updateParams := generated.UpdateEventRegistrationConfigParams{
		EventID:                              eventID,
		FirstNameRequirementStatus:           params.FirstNameRequirementStatus,
		LastNameRequirementStatus:            params.LastNameRequirementStatus,
		EmailRequirementStatus:               params.EmailRequirementStatus,
		BioRequirementStatus:                 params.BioRequirementStatus,
		PhoneNumberRequirementStatus:         params.PhoneNumberRequirementStatus,
		AddressRequirementStatus:             params.AddressRequirementStatus,
		AcademicInstitutionRequirementStatus: params.AcademicInstitutionRequirementStatus,
		AcademicEmailRequirementStatus:       params.AcademicEmailRequirementStatus,
	}

	return uc.EventRegistrationDg.UpdateEventRegistrationConfig(ctx, updateParams)
}

func (uc *EventConfigUsecase) DeleteEventRegistrationConfig(ctx context.Context, eventID uuid.UUID) error {
	return uc.EventRegistrationDg.DeleteEventRegistrationConfig(ctx, eventID)
}
