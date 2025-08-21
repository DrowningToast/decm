package datagateway

import (
	"context"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
)

type UpdateAuthenticationCredentialParameters struct {
	SolutionStatus      entity.SolutionStatus
	HashedPassword      *string
	EncryptedPrivateKey *string

	GoogleConnectorRef *string
	GithubConnectorRef *string

	IsVerifiedOrganizer bool
	IsVerifiedStudent   bool
}

type AuthenticationCredentialDataGateway interface {
	GetAuthenticationCredentialById(ctx context.Context, id int32) (*entity.AuthenticationCredential, *customerror.Err)
	GetAuthenticationCredentialByPublicKey(ctx context.Context, publicKey string) (*entity.AuthenticationCredential, *customerror.Err)
	CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, *customerror.Err)
	UpdateAuthenticationCredential(ctx context.Context, id int32, params UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, *customerror.Err)
	DeleteAuthenticationCredential(ctx context.Context, id int32) error
}
