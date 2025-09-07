package datagateway

import (
	"context"

	"apps/backend/common/customerror"
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
}

type AuthenticationCredentialDataGateway interface {
	GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, *customerror.Err)
	GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, *customerror.Err)
	CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, *customerror.Err)
	UpdateAuthenticationCredential(ctx context.Context, id uuid.UUID, params UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, *customerror.Err)
	DeleteAuthenticationCredential(ctx context.Context, id uuid.UUID) error
}
