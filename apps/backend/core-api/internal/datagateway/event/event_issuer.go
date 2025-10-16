package datagateway

import (
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

type EventIssuerDataGateway interface {
	CreateEventIssuer(ctx context.Context, params generated.CreateEventIssuerParams) (*generated.EventIssuer, error)
	GetEventIssuerByID(ctx context.Context, id uuid.UUID) (generated.EventIssuer, error)
	GetEventIssuersByEventID(ctx context.Context, eventID uuid.UUID) ([]generated.EventIssuer, error)
	UpdateEventIssuer(ctx context.Context, params generated.UpdateEventIssuerParams) (*generated.EventIssuer, error)
	DeleteEventIssuer(ctx context.Context, eventID uuid.UUID) error
}
