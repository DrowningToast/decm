package postgres

import (
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"context"

	"github.com/google/uuid"

	"decm-database/go/generated"
)

var _ datagateway.EventContractDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventContract(ctx context.Context, params generated.CreateEventContractParams) (*generated.EventContract, error) {
	result, err := r.queries.CreateEventContract(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventContract, error) {
	result, err := r.queries.GetEventContractByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) UpdateEventContract(ctx context.Context, params generated.UpdateEventContractParams) (*generated.EventContract, error) {
	result, err := r.queries.UpdateEventContract(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) DeleteEventContract(ctx context.Context, eventID uuid.UUID) error {
	return r.queries.DeleteEventContract(ctx, eventID)
}
