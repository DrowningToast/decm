package usecase

import "apps/backend/core-api/internal/datagateway"

type Usecase struct {
	AuthenticationCredentialDg datagateway.AuthenticationCredentialDataGateway
}

func New(authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway) *Usecase {
	return &Usecase{
		AuthenticationCredentialDg: authenticationCredentialDg,
	}
}
