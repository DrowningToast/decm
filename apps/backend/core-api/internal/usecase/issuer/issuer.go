package issuer

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/core-api/internal/datagateway"
)

type IssuerUsecase struct {
	IssuerDg datagateway.IssuerDataGateway
}

func NewIssuerUsecase(issuerDg datagateway.IssuerDataGateway) *IssuerUsecase {
	return &IssuerUsecase{
		IssuerDg: issuerDg,
	}
}

func (u *IssuerUsecase) GetEventsByIssuerCredentialID(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetEventIssuersByCredentialIDRow, error) {
	events, err := u.IssuerDg.GetEventsByIssuerCredentialID(ctx, issuerCredentialID, limitCount, offsetCount)
	if err != nil {
		return nil, err
	}

	return events, nil
}
