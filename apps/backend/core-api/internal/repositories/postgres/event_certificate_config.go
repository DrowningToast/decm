package postgres

import (
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"context"

	"github.com/google/uuid"

	"decm-database/go/generated"
)

var _ datagateway.EventCertificateConfigDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventCertificateConfig(ctx context.Context, params generated.CreateEventCertificateConfigParams) (*generated.EventCertificateConfig, error) {
	result, err := r.queries.CreateEventCertificateConfig(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) GetEventCertificateConfigByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventCertificateConfig, error) {
	result, err := r.queries.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) UpdateEventCertificateConfig(ctx context.Context, params generated.UpdateEventCertificateConfigParams) (*generated.EventCertificateConfig, error) {
	result, err := r.queries.UpdateEventCertificateConfig(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) DeleteEventCertificateConfig(ctx context.Context, eventID uuid.UUID) error {
	return r.queries.DeleteEventCertificateConfig(ctx, eventID)
}
