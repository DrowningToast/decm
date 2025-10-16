package datagateway

import (
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

type EventContractDataGateway interface {
	CreateEventContract(ctx context.Context, dbtx generated.DBTX, params generated.CreateEventContractParams) (*generated.EventContract, error)
	GetEventContractByEventID(ctx context.Context, dbtx generated.DBTX, eventID uuid.UUID) (*generated.EventContract, error)
	UpdateEventContract(ctx context.Context, dbtx generated.DBTX, params generated.UpdateEventContractParams) (*generated.EventContract, error)
	DeleteEventContract(ctx context.Context, dbtx generated.DBTX, eventID uuid.UUID) error
}
