package event

import (
	"apps/backend/core-api/internal/entity"
	"context"

	"github.com/google/uuid"
)

func (u *EventUsecase) GetEventContractByEventId(ctx context.Context, eventID uuid.UUID) (*entity.EventContract, error) {
	return u.EventContractDataGateway.GetEventContractByEventID(ctx, eventID)
}

func (u *EventUsecase) DeleteEventContract(ctx context.Context, eventID uuid.UUID) error {
	return u.EventContractDataGateway.DeleteEventContract(ctx, eventID)
}
