package postgres

import (
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"context"

	"github.com/google/uuid"
)

var _ datagateway.EventDataGateway = (*Repository)(nil)

func (r *Repository) CreateEvent(ctx context.Context, event datagateway.CreateEventParameters) (*entity.Event, error) {
	return &entity.Event{}, nil
}

func (r *Repository) GetEventById(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	return &entity.Event{}, nil
}
