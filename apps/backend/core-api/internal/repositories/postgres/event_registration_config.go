package postgres

import (
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"
	"time"

	datagateway "apps/backend/core-api/internal/datagateway/event"

	"github.com/google/uuid"
)

var _ datagateway.EventRegistrationConfigDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventRegistrationConfig(ctx context.Context, params generated.CreateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error) {
	result, err := r.queries.CreateEventRegistrationConfig(ctx, params)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return mapGeneratedToEntityEventRegistrationConfig(&result), nil
}

func (r *Repository) GetEventRegistrationConfigByEventId(ctx context.Context, eventID uuid.UUID) (*entity.EventRegistrationConfig, error) {
	result, err := r.queries.GetEventRegistrationConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return mapGeneratedToEntityEventRegistrationConfig(&result), nil
}

func (r *Repository) UpdateEventRegistrationConfig(ctx context.Context, params generated.UpdateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error) {
	result, err := r.queries.UpdateEventRegistrationConfig(ctx, params)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return mapGeneratedToEntityEventRegistrationConfig(&result), nil
}

func (r *Repository) DeleteEventRegistrationConfig(ctx context.Context, eventID uuid.UUID) error {
	return r.queries.DeleteEventRegistrationConfig(ctx, eventID)
}

func (r *Repository) GetEventRegistrationConfigPasswordByEventId(ctx context.Context, eventID uuid.UUID) (*entity.EventRegistrationConfig, error) {
	result, err := r.queries.GetEventRegistrationConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return mapGeneratedToEntityEventRegistrationConfig(&result), nil
}

// mapGeneratedToEntityEventRegistrationConfig converts generated.EventRegistrationConfig to entity.EventRegistrationConfig
func mapGeneratedToEntityEventRegistrationConfig(config *generated.EventRegistrationConfig) *entity.EventRegistrationConfig {
	if config == nil {
		return nil
	}

	var finalCallForRegistration *time.Time
	if config.FinalCallForRegistration.Valid {
		finalCallForRegistration = &config.FinalCallForRegistration.Time
	}

	return &entity.EventRegistrationConfig{
		ID:                                   config.ID,
		EventID:                              config.EventID,
		FinalCallForRegistration:             finalCallForRegistration,
		RegistrationPassword:                 pgmapper.PgTextToStringPtr(config.RegistrationPassword),
		IsIdentityVerificationRequired:       config.IsIdentityVerificationRequired.Int32 != 0,
		FirstNameRequirementStatus:           int(config.FirstNameRequirementStatus.Int32),
		LastNameRequirementStatus:            int(config.LastNameRequirementStatus.Int32),
		EmailRequirementStatus:               int(config.EmailRequirementStatus.Int32),
		BioRequirementStatus:                 int(config.BioRequirementStatus.Int32),
		PhoneNumberRequirementStatus:         int(config.PhoneNumberRequirementStatus.Int32),
		AddressRequirementStatus:             int(config.AddressRequirementStatus.Int32),
		AcademicInstitutionRequirementStatus: int(config.AcademicInstitutionRequirementStatus.Int32),
		AcademicEmailRequirementStatus:       int(config.AcademicEmailRequirementStatus.Int32),
		CreatedAt:                            config.CreatedAt.Time,
		UpdatedAt:                            config.UpdatedAt.Time,
	}
}
