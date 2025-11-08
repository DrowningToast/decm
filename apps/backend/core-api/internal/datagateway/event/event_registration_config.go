package datagateway

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type EventRegistrationConfigDataGateway interface {
	CreateEventRegistrationConfig(ctx context.Context, params generated.CreateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error)
	GetEventRegistrationConfigByEventID(ctx context.Context, eventID uuid.UUID) (*entity.EventRegistrationConfig, error)
	GetEventRegistrationConfigPasswordByEventID(ctx context.Context, eventID uuid.UUID) (*string, error)
	UpdateEventRegistrationConfig(ctx context.Context, params generated.UpdateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error)
	DeleteEventRegistrationConfig(ctx context.Context, eventID uuid.UUID) error
}
