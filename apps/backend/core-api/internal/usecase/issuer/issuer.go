package issuer

import (
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
