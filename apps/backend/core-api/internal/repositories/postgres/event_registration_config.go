package postgres

import (
	"context"
	"decm-database/go/generated"

	datagateway "apps/backend/core-api/internal/datagateway/event"

	"github.com/google/uuid"
)

var _ datagateway.EventRegistrationConfigDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventRegistrationConfig(ctx context.Context, params generated.CreateEventRegistrationConfigParams) (*generated.EventRegistrationConfig, error) {
	result, err := r.queries.CreateEventRegistrationConfig(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) GetEventRegistrationConfigByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventRegistrationConfig, error) {
	result, err := r.queries.GetEventRegistrationConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) UpdateEventRegistrationConfig(ctx context.Context, params generated.UpdateEventRegistrationConfigParams) (*generated.EventRegistrationConfig, error) {
	result, err := r.queries.UpdateEventRegistrationConfig(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) DeleteEventRegistrationConfig(ctx context.Context, eventID uuid.UUID) error {
	return r.queries.DeleteEventRegistrationConfig(ctx, eventID)
}
