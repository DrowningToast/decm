package datagateway

import (
	"context"

	"apps/backend/core-api/internal/entity"
)

type UpdateAuthenticationCredentialParameters struct {
	SolutionStatus entity.SolutionStatus
	Password       *string
	PrivateKey     *string

	GoogleConnectorRef *string
	GithubConnectorRef *string

	IsVerifiedOrganizer bool
	IsVerifiedStudent   bool
}

type AuthenticationCredentialDataGateway interface {
	GetAuthenticationCredentialById(ctx context.Context, id int32) (*entity.AuthenticationCredential, error)
	GetAuthenticationCredentialByPublicKey(ctx context.Context, publicKey string) (*entity.AuthenticationCredential, error)
	CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error)
	UpdateAuthenticationCredential(ctx context.Context, id int32, params UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error)
	DeleteAuthenticationCredential(ctx context.Context, id int32) error
}
