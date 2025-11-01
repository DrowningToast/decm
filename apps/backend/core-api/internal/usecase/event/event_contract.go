package event

import (
	"context"

	"github.com/google/uuid"

	"decm-database/go/generated"
)

func (u *EventUsecase) GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventContract, error) {
	return u.EventContractDataGateway.GetEventContractByEventID(ctx, eventID)
}

func (u *EventUsecase) DeleteEventContract(ctx context.Context, eventID uuid.UUID) error {
	return u.EventContractDataGateway.DeleteEventContract(ctx, eventID)
}
