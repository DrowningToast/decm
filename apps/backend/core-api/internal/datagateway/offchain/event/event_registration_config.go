package offchain_datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

type EventRegistrationConfigDataGateway interface {
	CreateEventRegistrationConfig(ctx context.Context, params generated.CreateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error)
	GetEventRegistrationConfigByEventId(ctx context.Context, eventId uuid.UUID) (*entity.EventRegistrationConfig, error)
	GetEventRegistrationConfigPasswordByEventId(ctx context.Context, eventId uuid.UUID) (*entity.EventRegistrationConfig, error)
	UpdateEventRegistrationConfig(ctx context.Context, params generated.UpdateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error)
	DeleteEventRegistrationConfig(ctx context.Context, eventID uuid.UUID) error
}
