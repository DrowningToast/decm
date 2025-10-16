package event

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"decm-database/go/generated"
)

type UpdateEventIssuerParams struct {
	IsSigned  int32
	Signature pgtype.Text
}

func (u *EventUsecase) UpdateEventIssuer(ctx context.Context, id uuid.UUID, params UpdateEventIssuerParams) (*generated.EventIssuer, error) {
	updateParams := generated.UpdateEventIssuerParams{
		ID:        id,
		IsSigned:  params.IsSigned,
		Signature: params.Signature,
	}

	return u.EventIssuerDataGateway.UpdateEventIssuer(ctx, updateParams)
}
