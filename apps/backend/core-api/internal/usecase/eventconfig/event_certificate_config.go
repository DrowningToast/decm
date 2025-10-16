package eventconfig

import (
	"context"
	"decm-database/go/generated"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateEventCertificateConfigParams struct {
	BaseCertificateStorageKey string
	EventNamePosX             int32
	EventNamePosY             int32
	NamePosX                  int32
	NamePosY                  int32
	AcademicInstitutionPosX   pgtype.Int4
	AcademicInstitutionPosY   pgtype.Int4
}

func (uc *EventConfigUsecase) CreateEventCertificateConfig(ctx context.Context, eventID uuid.UUID, params CreateEventCertificateConfigParams) (*generated.EventCertificateConfig, error) {
	// Check if config already exists for this event
	existingConfig, err := uc.EventCertificateDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err == nil && existingConfig != nil {
		return nil, fmt.Errorf("event certificate config already exists for event ID: %s", eventID.String())
	}

	// Create new config
	createParams := generated.CreateEventCertificateConfigParams{
		EventID:                   eventID,
		BaseCertificateStorageKey: params.BaseCertificateStorageKey,
		EventNamePosX:             params.EventNamePosX,
		EventNamePosY:             params.EventNamePosY,
		NamePosX:                  params.NamePosX,
		NamePosY:                  params.NamePosY,
		AcademicInstitutionPosX:   params.AcademicInstitutionPosX,
		AcademicInstitutionPosY:   params.AcademicInstitutionPosY,
	}

	return uc.EventCertificateDg.CreateEventCertificateConfig(ctx, createParams)
}

type UpdateEventCertificateConfigParams struct {
	BaseCertificateStorageKey string
	EventNamePosX             int32
	EventNamePosY             int32
	NamePosX                  int32
	NamePosY                  int32
	AcademicInstitutionPosX   pgtype.Int4
	AcademicInstitutionPosY   pgtype.Int4
}

func (uc *EventConfigUsecase) UpdateEventCertificateConfig(ctx context.Context, eventID uuid.UUID, params UpdateEventCertificateConfigParams) (*generated.EventCertificateConfig, error) {
	updateParams := generated.UpdateEventCertificateConfigParams{
		EventID:                   eventID,
		BaseCertificateStorageKey: params.BaseCertificateStorageKey,
		EventNamePosX:             params.EventNamePosX,
		EventNamePosY:             params.EventNamePosY,
		NamePosX:                  params.NamePosX,
		NamePosY:                  params.NamePosY,
		AcademicInstitutionPosX:   params.AcademicInstitutionPosX,
		AcademicInstitutionPosY:   params.AcademicInstitutionPosY,
	}

	return uc.EventCertificateDg.UpdateEventCertificateConfig(ctx, updateParams)
}

func (uc *EventConfigUsecase) GetEventCertificateConfigByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventCertificateConfig, error) {
	return uc.EventCertificateDg.GetEventCertificateConfigByEventID(ctx, eventID)
}

func (uc *EventConfigUsecase) DeleteEventCertificateConfig(ctx context.Context, eventID uuid.UUID) error {
	return uc.EventCertificateDg.DeleteEventCertificateConfig(ctx, eventID)
}
