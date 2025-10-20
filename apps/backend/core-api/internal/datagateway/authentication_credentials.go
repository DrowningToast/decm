package datagateway

import (
	"context"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type UpdateAuthenticationCredentialParameters struct {
	SolutionStatus      entity.SolutionStatus
	HashedPassword      *string
	EncryptedPrivateKey *string

	GoogleConnectorRef *string
	GithubConnectorRef *string

	IsVerifiedOrganizer bool
	IsVerifiedStudent   bool
	IsVerifiedIssuer    bool
}

type AuthenticationCredentialDataGateway interface {
	GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error)
	GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, error)
	GetAuthenticationCredentialByGoogleConnectorRef(ctx context.Context, googleConnectorRef string) (*entity.AuthenticationCredential, error)
	CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error)
	UpdateAuthenticationCredential(ctx context.Context, id uuid.UUID, params UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error)
	DeleteAuthenticationCredential(ctx context.Context, id uuid.UUID) error
}
