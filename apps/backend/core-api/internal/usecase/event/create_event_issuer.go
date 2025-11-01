package event

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/common/customerror"
	"decm-database/go/generated"
)

type CreateEventIssuerParams struct {
	EventID            uuid.UUID
	IssuerCredentialID uuid.UUID
	IsSigned           int32
	Signature          pgtype.Text
	SignMessage        pgtype.Text
}

func (u *EventUsecase) CreateEventIssuer(ctx context.Context, params CreateEventIssuerParams) (*generated.EventIssuer, error) {
	createParams := generated.CreateEventIssuerParams{
		EventID:            params.EventID,
		IssuerCredentialID: params.IssuerCredentialID,
		IsSigned:           params.IsSigned,
		Signature:          params.Signature,
		SignMessage:        params.SignMessage,
	}

	issuer, err := u.EventIssuerDataGateway.CreateEventIssuer(ctx, createParams)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return issuer, nil
}
