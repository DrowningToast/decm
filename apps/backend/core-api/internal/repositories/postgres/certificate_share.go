package postgres

import (
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"
	"errors"

	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ event_datagateway.CertificateShareDataGateway = (*Repository)(nil)

func (r *Repository) CreateCertificateShare(ctx context.Context, params event_datagateway.CreateCertificateShareParameters) (*entity.CertificateShare, error) {
	result, err := r.queries.CreateCertificateShare(ctx, generated.CreateCertificateShareParams{
		EventCertificateID: params.EventCertificateId,
		Active:             params.Active,
		Handle:             pgmapper.StringPtrToPgText(&params.Handle),
		Password:           pgmapper.StringPtrToPgText(params.Password),
	})
	if err != nil {
		return nil, err
	}

	return &entity.CertificateShare{
		Id:                 result.ID,
		EventCertificateId: result.EventCertificateID,
		Active:             result.Active,
		Handle:             result.Handle.String,
		Password:           pgmapper.PgTextToStringPtr(result.Password),
		CreatedAt:          result.CreatedAt.Time,
		UpdatedAt:          result.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetCertificateShareByHandle(ctx context.Context, handle string) (*entity.CertificateShare, error) {
	result, err := r.queries.GetCertificateShareByHandle(ctx, pgmapper.StringPtrToPgText(&handle))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity.CertificateShare{
		Id:                 result.ID,
		EventCertificateId: result.EventCertificateID,
		Active:             result.Active,
		Handle:             result.Handle.String,
		Password:           pgmapper.PgTextToStringPtr(result.Password),
		CreatedAt:          result.CreatedAt.Time,
		UpdatedAt:          result.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetCertificateShareByID(ctx context.Context, id uuid.UUID) (*entity.CertificateShare, error) {
	result, err := r.queries.GetCertificateShareByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity.CertificateShare{
		Id:                 result.ID,
		EventCertificateId: result.EventCertificateID,
		Active:             result.Active,
		Handle:             result.Handle.String,
		Password:           pgmapper.PgTextToStringPtr(result.Password),
		CreatedAt:          result.CreatedAt.Time,
		UpdatedAt:          result.UpdatedAt.Time,
	}, nil
}

func (r *Repository) UpdateCertificateShare(ctx context.Context, id uuid.UUID, params event_datagateway.UpdateCertificateShareParameters) (*entity.CertificateShare, error) {
	result, err := r.queries.UpdateCertificateShare(ctx, generated.UpdateCertificateShareParams{
		ID:       id,
		Password: pgmapper.StringPtrToPgText(params.Password),
	})
	if err != nil {
		return nil, err
	}

	return &entity.CertificateShare{
		Id:                 result.ID,
		EventCertificateId: result.EventCertificateID,
		Active:             result.Active,
		Handle:             result.Handle.String,
		Password:           pgmapper.PgTextToStringPtr(result.Password),
		CreatedAt:          result.CreatedAt.Time,
		UpdatedAt:          result.UpdatedAt.Time,
	}, nil
}
