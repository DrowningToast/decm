package postgres

import (
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"context"

	"github.com/google/uuid"

	"decm-database/go/generated"
)

var _ datagateway.EventContractDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventContract(ctx context.Context, dbtx generated.DBTX, params generated.CreateEventContractParams) (*generated.EventContract, error) {
	queries := generated.New(dbtx)
	result, err := queries.CreateEventContract(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) GetEventContractByEventID(ctx context.Context, dbtx generated.DBTX, eventID uuid.UUID) (*generated.EventContract, error) {
	queries := generated.New(dbtx)
	result, err := queries.GetEventContractByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) UpdateEventContract(ctx context.Context, dbtx generated.DBTX, params generated.UpdateEventContractParams) (*generated.EventContract, error) {
	queries := generated.New(dbtx)
	result, err := queries.UpdateEventContract(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) DeleteEventContract(ctx context.Context, dbtx generated.DBTX, eventID uuid.UUID) error {
	queries := generated.New(dbtx)
	return queries.DeleteEventContract(ctx, eventID)
}
