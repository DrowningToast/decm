package event

import (
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

func (u *EventUsecase) GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventContract, error) {
	return u.EventContractDataGateway.GetEventContractByEventID(ctx, eventID)
}

func (u *EventUsecase) DeleteEventContract(ctx context.Context, eventID uuid.UUID) error {
	return u.EventContractDataGateway.DeleteEventContract(ctx, eventID)
}
