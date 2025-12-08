package datagateway

import (
	"context"

	"apps/backend/common"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type UpdateAuthenticationCredentialParameters struct {
	SolutionStatus      common.SolutionStatus
	HashedPassword      *string
	EncryptedPrivateKey *string

	GoogleConnectorRef *string
	GithubConnectorRef *string

	IsVerifiedOrganizer bool
	IsVerifiedStudent   bool
	IsVerifiedIssuer    bool
}

// GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddressParameters are parameters for searching by google connector ref or wallet address
type GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddressParameters struct {
	// GoogleConnectorRef is the Google OAuth email (optional, at least one must be provided)
	GoogleConnectorRef *string
	// WalletAddress is the wallet address (optional, at least one must be provided)
	WalletAddress *string
}

type AuthenticationCredentialDataGateway interface {
	GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error)
	GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error)
	GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, error)
	GetAuthenticationCredentialByGoogleConnectorRef(ctx context.Context, googleConnectorRef string) (*entity.AuthenticationCredential, error)
	// GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddress fetches a credential by Google OAuth email OR wallet address.
	// At least one parameter must be provided. Returns the first matching credential.
	GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddress(ctx context.Context, params GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddressParameters) (*entity.AuthenticationCredential, error)
	CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error)
	UpdateAuthenticationCredential(ctx context.Context, id uuid.UUID, params UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error)
	DeleteAuthenticationCredential(ctx context.Context, id uuid.UUID) error
}
