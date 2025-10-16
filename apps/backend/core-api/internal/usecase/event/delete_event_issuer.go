package event

import (
	"context"

	"github.com/google/uuid"
)

func (u *EventUsecase) DeleteEventIssuer(ctx context.Context, id uuid.UUID) error {
	return u.EventIssuerDataGateway.DeleteEventIssuer(ctx, id)
}
