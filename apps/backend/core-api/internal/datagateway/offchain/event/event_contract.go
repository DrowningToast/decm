package offchain_datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

type EventContractDataGateway interface {
	CreateEventContract(ctx context.Context, params generated.CreateEventContractParams) (*entity.EventContract, error)
	GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*entity.EventContract, error)
	UpdateEventContract(ctx context.Context, params generated.UpdateEventContractParams) (*entity.EventContract, error)
	DeleteEventContract(ctx context.Context, eventID uuid.UUID) error
}
