package event

import (
	"apps/backend/common/customerror"
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

type CreateEventIssuerParams struct {
	EventID            uuid.UUID
	IssuerCredentialID uuid.UUID
	IsSigned           int32
}

func (u *EventUsecase) CreateEventIssuer(ctx context.Context, params CreateEventIssuerParams) (*generated.EventIssuer, error) {
	createParams := generated.CreateEventIssuerParams{
		EventID:            params.EventID,
		IssuerCredentialID: params.IssuerCredentialID,
		IsSigned:           params.IsSigned,
	}

	issuer, err := u.EventIssuerDataGateway.CreateEventIssuer(ctx, createParams)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return issuer, nil
}
