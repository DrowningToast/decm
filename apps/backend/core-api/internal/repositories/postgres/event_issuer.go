package postgres

import (
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"context"

	"github.com/google/uuid"

	"decm-database/go/generated"
)

var _ datagateway.EventIssuerDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventIssuer(ctx context.Context, dbtx generated.DBTX, params generated.CreateEventIssuerParams) (*generated.EventIssuer, error) {
	queries := generated.New(dbtx)
	result, err := queries.CreateEventIssuer(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) GetEventIssuersByEventID(ctx context.Context, dbtx generated.DBTX, eventID uuid.UUID) ([]generated.EventIssuer, error) {
	queries := generated.New(dbtx)
	return queries.GetEventIssuersByEventID(ctx, eventID)
}

func (r *Repository) GetEventIssuerByID(ctx context.Context, dbtx generated.DBTX, id uuid.UUID) (*generated.EventIssuer, error) {
	queries := generated.New(dbtx)
	result, err := queries.GetEventIssuerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) UpdateEventIssuer(ctx context.Context, dbtx generated.DBTX, params generated.UpdateEventIssuerParams) (*generated.EventIssuer, error) {
	queries := generated.New(dbtx)
	result, err := queries.UpdateEventIssuer(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) DeleteEventIssuer(ctx context.Context, dbtx generated.DBTX, id uuid.UUID) error {
	queries := generated.New(dbtx)
	return queries.DeleteEventIssuer(ctx, id)
}
