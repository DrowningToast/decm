package datagateway

import (
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

type EventIssuerDataGateway interface {
	CreateEventIssuer(ctx context.Context, dbtx generated.DBTX, params generated.CreateEventIssuerParams) (*generated.EventIssuer, error)
	GetEventIssuersByEventID(ctx context.Context, dbtx generated.DBTX, eventID uuid.UUID) ([]generated.EventIssuer, error)
	UpdateEventIssuer(ctx context.Context, dbtx generated.DBTX, params generated.UpdateEventIssuerParams) (*generated.EventIssuer, error)
	DeleteEventIssuer(ctx context.Context, dbtx generated.DBTX, eventID uuid.UUID) error
}
