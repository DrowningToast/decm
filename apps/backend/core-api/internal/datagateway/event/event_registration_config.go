package datagateway

import (
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

type EventRegistrationConfigDataGateway interface {
	CreateEventRegistrationConfig(ctx context.Context, params generated.CreateEventRegistrationConfigParams) (*generated.EventRegistrationConfig, error)
	GetEventRegistrationConfigByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventRegistrationConfig, error)
	UpdateEventRegistrationConfig(ctx context.Context, params generated.UpdateEventRegistrationConfigParams) (*generated.EventRegistrationConfig, error)
	DeleteEventRegistrationConfig(ctx context.Context, eventID uuid.UUID) error
}
