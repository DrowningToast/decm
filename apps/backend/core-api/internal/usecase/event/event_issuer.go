package event

import (
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

func (u *EventUsecase) GetEventIssuersByEventID(ctx context.Context, eventID uuid.UUID) ([]generated.EventIssuer, error) {
	return u.EventIssuerDataGateway.GetEventIssuersByEventID(ctx, eventID)
}

func (u *EventUsecase) GetEventIssuerByID(ctx context.Context, id uuid.UUID) (generated.EventIssuer, error) {
	return u.EventIssuerDataGateway.GetEventIssuerByID(ctx, id)
}
