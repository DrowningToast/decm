package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/pgerrutils"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ datagateway.EventCertificateConfigDataGateway = (*Repository)(nil)

// mapEventCertificateConfigToEntity maps generated.EventCertificateConfig to entity.EventCertificateConfig
func mapEventCertificateConfigToEntity(gen generated.EventCertificateConfig) *entity.EventCertificateConfig {
	return &entity.EventCertificateConfig{
		ID:                              gen.ID,
		EventID:                         gen.EventID,
		BaseCertificateStorageKey:       gen.BaseCertificateStorageKey,
		EventNamePosX:                   pgFloat8ToFloat64Ptr(gen.EventNamePosX),
		EventNamePosY:                   pgFloat8ToFloat64Ptr(gen.EventNamePosY),
		NamePosX:                        gen.NamePosX,
		NamePosY:                        gen.NamePosY,
		AcademicInstitutionPosX:         pgFloat8ToFloat64Ptr(gen.AcademicInstitutionPosX),
		AcademicInstitutionPosY:         pgFloat8ToFloat64Ptr(gen.AcademicInstitutionPosY),
		CertificateTitlePosX:            pgFloat8ToFloat64Ptr(gen.CertificateTitlePosX),
		CertificateTitlePosY:            pgFloat8ToFloat64Ptr(gen.CertificateTitlePosY),
		CertificateSubtitlePosX:         pgFloat8ToFloat64Ptr(gen.CertificateSubtitlePosX),
		CertificateSubtitlePosY:         pgFloat8ToFloat64Ptr(gen.CertificateSubtitlePosY),
		IsPublished:                     gen.IsPublished,
		EventNameFontFamilyID:           pgInt4ToInt32Ptr(gen.EventNameFontFamilyID),
		EventNameFontWeight:             pgInt4ToInt32Ptr(gen.EventNameFontWeight),
		NameFontFamilyID:                pgInt4ToInt32Ptr(gen.NameFontFamilyID),
		NameFontWeight:                  pgInt4ToInt32Ptr(gen.NameFontWeight),
		AcademicInstitutionFontFamilyID: pgInt4ToInt32Ptr(gen.AcademicInstitutionFontFamilyID),
		AcademicInstitutionFontWeight:   pgInt4ToInt32Ptr(gen.AcademicInstitutionFontWeight),
		CertificateTitleFontFamilyID:    pgInt4ToInt32Ptr(gen.CertificateTitleFontFamilyID),
		CertificateTitleFontWeight:      pgInt4ToInt32Ptr(gen.CertificateTitleFontWeight),
		CertificateSubtitleFontFamilyID: pgInt4ToInt32Ptr(gen.CertificateSubtitleFontFamilyID),
		CertificateSubtitleFontWeight:   pgInt4ToInt32Ptr(gen.CertificateSubtitleFontWeight),
		CreatedAt:                       gen.CreatedAt.Time,
		UpdatedAt:                       gen.UpdatedAt.Time,
	}
}

// pgFloat8ToFloat64Ptr converts pgtype.Float8 to *float64
func pgFloat8ToFloat64Ptr(f8 pgtype.Float8) *float64 {
	if !f8.Valid {
		return nil
	}
	return &f8.Float64
}

// pgInt4ToInt32Ptr converts pgtype.Int4 to *int32
func pgInt4ToInt32Ptr(i4 pgtype.Int4) *int32 {
	if !i4.Valid {
		return nil
	}
	return &i4.Int32
}

func (r *Repository) CreateEventCertificateConfig(ctx context.Context, params generated.CreateEventCertificateConfigParams) (*entity.EventCertificateConfig, error) {
	result, err := r.queries.CreateEventCertificateConfig(ctx, params)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return mapEventCertificateConfigToEntity(result), nil
}

func (r *Repository) GetEventCertificateConfigByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificateConfig, error) {
	result, err := r.queries.GetEventCertificateConfigByID(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return mapEventCertificateConfigToEntity(result), nil
}

func (r *Repository) GetEventCertificateConfigByEventID(ctx context.Context, eventID uuid.UUID) (*entity.EventCertificateConfig, error) {
	result, err := r.queries.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return mapEventCertificateConfigToEntity(result), nil
}

func (r *Repository) UpdateEventCertificateConfig(ctx context.Context, params generated.UpdateEventCertificateConfigParams) (*entity.EventCertificateConfig, error) {
	result, err := r.queries.UpdateEventCertificateConfig(ctx, params)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return mapEventCertificateConfigToEntity(result), nil
}

func (r *Repository) UpdateEventCertificateTextConfig(ctx context.Context, params generated.UpdateEventCertificateTextConfigParams) (*entity.EventCertificateConfig, error) {
	result, err := r.queries.UpdateEventCertificateTextConfig(ctx, params)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return mapEventCertificateConfigToEntity(result), nil
}

func (r *Repository) DeleteEventCertificateConfig(ctx context.Context, eventID uuid.UUID) error {
	return r.queries.DeleteEventCertificateConfig(ctx, eventID)
}

func (r *Repository) ToggleEventCertificateConfigPublished(ctx context.Context, params generated.ToggleEventCertificateConfigPublishedParams) (*entity.EventCertificateConfig, error) {
	result, err := r.queries.ToggleEventCertificateConfigPublished(ctx, params)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return mapEventCertificateConfigToEntity(result), nil
}
