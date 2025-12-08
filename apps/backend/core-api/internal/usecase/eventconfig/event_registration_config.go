package eventconfig

import (
	"context"
	"decm-database/go/generated"
	"errors"
	"fmt"

	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	eventDg "apps/backend/core-api/internal/datagateway/event"
)

type CreateEventRegistrationConfigParams struct {
	FinalCallForRegistration             pgtype.Timestamptz
	RegistrationPassword                 pgtype.Text
	FirstNameRequirementStatus           pgtype.Int4
	LastNameRequirementStatus            pgtype.Int4
	EmailRequirementStatus               pgtype.Int4
	BioRequirementStatus                 pgtype.Int4
	PhoneNumberRequirementStatus         pgtype.Int4
	AddressRequirementStatus             pgtype.Int4
	AcademicInstitutionRequirementStatus pgtype.Int4
	AcademicEmailRequirementStatus       pgtype.Int4
}

func (uc *EventConfigUsecase) CreateEventRegistrationConfig(ctx context.Context, eventID uuid.UUID, params CreateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error) {
	// Check if config already exists for this event
	existingConfig, err := uc.EventRegistrationDg.GetEventRegistrationConfigByEventId(ctx, eventID)
	if err == nil && existingConfig != nil {
		return nil, fmt.Errorf("event registration config already exists for event ID: %s", eventID.String())
	}

	// Hash registration password if provided
	var registrationPassword pgtype.Text
	if params.RegistrationPassword.Valid && params.RegistrationPassword.String != "" {
		hashPwd, err := hashutils.HashPassword(params.RegistrationPassword.String)
		if err != nil {
			return nil, err
		}
		registrationPassword = pgmapper.StringPtrToPgText(&hashPwd)
	}

	// Create new config
	createParams := generated.CreateEventRegistrationConfigParams{
		EventID:                              eventID,
		FinalCallForRegistration:             params.FinalCallForRegistration,
		RegistrationPassword:                 registrationPassword,
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

func (uc *EventConfigUsecase) GetEventRegistrationConfigByEventID(ctx context.Context, eventID uuid.UUID) (*entity.EventRegistrationConfig, error) {
	return uc.EventRegistrationDg.GetEventRegistrationConfigByEventId(ctx, eventID)
}

type UpdateEventRegistrationConfigParams struct {
	FinalCallForRegistration             pgtype.Timestamptz
	RegistrationPassword                 pgtype.Text
	FirstNameRequirementStatus           pgtype.Int4
	LastNameRequirementStatus            pgtype.Int4
	EmailRequirementStatus               pgtype.Int4
	BioRequirementStatus                 pgtype.Int4
	PhoneNumberRequirementStatus         pgtype.Int4
	AddressRequirementStatus             pgtype.Int4
	AcademicInstitutionRequirementStatus pgtype.Int4
	AcademicEmailRequirementStatus       pgtype.Int4
	IsBookingRequestRequired             *bool
	IsTicketTransferable                 *bool
	EventType                            *entity.EventType
}

func (uc *EventConfigUsecase) UpdateEventRegistrationConfig(ctx context.Context, eventID uuid.UUID, params UpdateEventRegistrationConfigParams, currentUser *auth.JwtClaims) (*entity.EventRegistrationConfig, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user not authenticated"))
	}

	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	event, err := uc.EventDg.GetEventById(ctx, eventID)
	if err != nil {
		return nil, err
	}

	eventRegistrationConfig, err := uc.EventRegistrationDg.GetEventRegistrationConfigByEventId(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event.OwnerCredentialId != credential.Id {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, err)
	}

	updateParams := generated.UpdateEventRegistrationConfigParams{
		EventID:                              eventID,
		FinalCallForRegistration:             params.FinalCallForRegistration,
		FirstNameRequirementStatus:           params.FirstNameRequirementStatus,
		LastNameRequirementStatus:            params.LastNameRequirementStatus,
		EmailRequirementStatus:               params.EmailRequirementStatus,
		BioRequirementStatus:                 params.BioRequirementStatus,
		PhoneNumberRequirementStatus:         params.PhoneNumberRequirementStatus,
		AddressRequirementStatus:             params.AddressRequirementStatus,
		AcademicInstitutionRequirementStatus: params.AcademicInstitutionRequirementStatus,
		AcademicEmailRequirementStatus:       params.AcademicEmailRequirementStatus,
	}

	if params.RegistrationPassword.Valid && (eventRegistrationConfig.RegistrationPassword == nil || params.RegistrationPassword.String != *eventRegistrationConfig.RegistrationPassword) {
		// If registration password is provided, hash it
		hashPwd, err := hashutils.HashPassword(params.RegistrationPassword.String)
		if err != nil {
			return nil, err
		}
		updateParams.RegistrationPassword = pgmapper.StringPtrToPgText(&hashPwd)
	} else {
		// If registration password is not provided, use the existing registration password
		updateParams.RegistrationPassword = pgmapper.StringPtrToPgText(eventRegistrationConfig.RegistrationPassword)
	}

	_, err = uc.EventRegistrationDg.UpdateEventRegistrationConfig(ctx, updateParams)
	if err != nil {
		return nil, err
	}

	// Only update event if EventType is provided
	if params.EventType != nil {
		updateEventParams := eventDg.UpdateEventParameters{
			EventType:                params.EventType,
			IsBookingRequestRequired: params.IsBookingRequestRequired,
			IsTicketTransferable:     params.IsTicketTransferable,
		}

		_, err = uc.EventDg.UpdateEvent(ctx, eventID, updateEventParams)
		if err != nil {
			return nil, err
		}
	}

	return uc.EventRegistrationDg.GetEventRegistrationConfigByEventId(ctx, eventID)
}

func (uc *EventConfigUsecase) DeleteEventRegistrationConfig(ctx context.Context, eventID uuid.UUID) error {
	return uc.EventRegistrationDg.DeleteEventRegistrationConfig(ctx, eventID)
}
