package datagateway

import (
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

type EventContractDataGateway interface {
	CreateEventContract(ctx context.Context, params generated.CreateEventContractParams) (*generated.EventContract, error)
	GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventContract, error)
	UpdateEventContract(ctx context.Context, params generated.UpdateEventContractParams) (*generated.EventContract, error)
	DeleteEventContract(ctx context.Context, eventID uuid.UUID) error
}
