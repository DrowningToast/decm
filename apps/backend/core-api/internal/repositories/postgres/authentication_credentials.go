package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

var _ datagateway.AuthenticationCredentialDataGateway = (*Repository)(nil)

func (r *Repository) GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.GetAuthenticationCredentialByID(ctx, generated.GetAuthenticationCredentialByIDParams{
		EncryptionKey: r.piiEncryptionKey,
		ID:            id,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	model := generated.AuthenticationCredential{
		ID:                  query.ID,
		SolutionStatus:      query.SolutionStatus,
		HashedPassword:      query.HashedPassword,
		EncryptedPrivateKey: query.EncryptedPrivateKey,
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  query.GoogleConnectorRef,
		GithubConnectorRef:  query.GithubConnectorRef,
		IsVerifiedOrganizer: query.IsVerifiedOrganizer,
		IsVerifiedStudent:   query.IsVerifiedStudent,
		CreatedAt:           query.CreatedAt,
		UpdatedAt:           query.UpdatedAt,
	}

	entity := entity.MapAuthenticationCredentialToEntity(model)
	return &entity, nil
}

func (r *Repository) GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.GetAuthenticationCredentialByWalletAddress(ctx, generated.GetAuthenticationCredentialByWalletAddressParams{
		EncryptionKey: r.piiEncryptionKey,
		WalletAddress: walletAddress,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	model := generated.AuthenticationCredential{
		ID:                  query.ID,
		SolutionStatus:      query.SolutionStatus,
		HashedPassword:      query.HashedPassword,
		EncryptedPrivateKey: query.EncryptedPrivateKey,
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  query.GoogleConnectorRef,
		GithubConnectorRef:  query.GithubConnectorRef,
		IsVerifiedOrganizer: query.IsVerifiedOrganizer,
		IsVerifiedStudent:   query.IsVerifiedStudent,
		CreatedAt:           query.CreatedAt,
		UpdatedAt:           query.UpdatedAt,
	}

	entity := entity.MapAuthenticationCredentialToEntity(model)
	return &entity, nil
}

func (r *Repository) GetAuthenticationCredentialByGoogleConnectorRef(ctx context.Context, googleConnectorRef string) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.GetAuthenticationCredentialByGoogleConnectorRef(ctx, generated.GetAuthenticationCredentialByGoogleConnectorRefParams{
		EncryptionKey:      r.piiEncryptionKey,
		GoogleConnectorRef: googleConnectorRef,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	model := generated.AuthenticationCredential{
		ID:                  query.ID,
		SolutionStatus:      query.SolutionStatus,
		HashedPassword:      query.HashedPassword,
		EncryptedPrivateKey: query.EncryptedPrivateKey,
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  query.GoogleConnectorRef,
		GithubConnectorRef:  query.GithubConnectorRef,
		IsVerifiedOrganizer: query.IsVerifiedOrganizer,
		IsVerifiedStudent:   query.IsVerifiedStudent,
		CreatedAt:           query.CreatedAt,
		UpdatedAt:           query.UpdatedAt,
	}

	entity := entity.MapAuthenticationCredentialToEntity(model)
	return &entity, nil
}

func (r *Repository) CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.CreateAuthenticationCredential(ctx, generated.CreateAuthenticationCredentialParams{
		SolutionStatus:      int32(credential.SolutionStatus),
		WalletAddress:       credential.WalletAddress,
		HashedPassword:      pgmapper.StringPtrToPgText(credential.HashedPassword),
		EncryptedPrivateKey: pgmapper.StringPtrToPgText(credential.EncryptedPrivateKey),
		GoogleConnectorRef:  pgmapper.StringPtrToPgText(credential.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.StringPtrToPgText(credential.GithubConnectorRef),
		IsVerifiedOrganizer: pgmapper.BoolToInt32(credential.IsVerifiedOrganizer),
		IsVerifiedStudent:   pgmapper.BoolToInt32(credential.IsVerifiedStudent),
		EncryptionKey:       r.piiEncryptionKey,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	model := generated.AuthenticationCredential{
		ID:                  query.ID,
		SolutionStatus:      query.SolutionStatus,
		HashedPassword:      query.HashedPassword,
		EncryptedPrivateKey: query.EncryptedPrivateKey,
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  query.GoogleConnectorRef,
		GithubConnectorRef:  query.GithubConnectorRef,
		IsVerifiedOrganizer: query.IsVerifiedOrganizer,
		IsVerifiedStudent:   query.IsVerifiedStudent,
		CreatedAt:           query.CreatedAt,
		UpdatedAt:           query.UpdatedAt,
	}

	entity := entity.MapAuthenticationCredentialToEntity(model)
	return &entity, nil
}

func (r *Repository) UpdateAuthenticationCredential(ctx context.Context, id uuid.UUID, params datagateway.UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.UpdateAuthenticationCredential(ctx, generated.UpdateAuthenticationCredentialParams{
		ID:                  id,
		SolutionStatus:      pgmapper.Int32ToPgInt4(int32(params.SolutionStatus)),
		HashedPassword:      pgmapper.StringPtrToPgText(params.HashedPassword),
		EncryptedPrivateKey: pgmapper.StringPtrToPgText(params.EncryptedPrivateKey),
		GoogleConnectorRef:  pgmapper.StringPtrToPgText(params.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.StringPtrToPgText(params.GithubConnectorRef),
		IsVerifiedOrganizer: pgmapper.BoolToPgInt4(params.IsVerifiedOrganizer),
		IsVerifiedStudent:   pgmapper.BoolToPgInt4(params.IsVerifiedStudent),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	model := generated.AuthenticationCredential{
		ID:                  query.ID,
		SolutionStatus:      query.SolutionStatus,
		HashedPassword:      query.HashedPassword,
		EncryptedPrivateKey: query.EncryptedPrivateKey,
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  query.GoogleConnectorRef,
		GithubConnectorRef:  query.GithubConnectorRef,
		IsVerifiedOrganizer: query.IsVerifiedOrganizer,
		IsVerifiedStudent:   query.IsVerifiedStudent,
		CreatedAt:           query.CreatedAt,
		UpdatedAt:           query.UpdatedAt,
	}

	entity := entity.MapAuthenticationCredentialToEntity(model)
	return &entity, nil
}

func (r *Repository) DeleteAuthenticationCredential(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteAuthenticationCredential(ctx, id)
	if err != nil {
		return pgerrutils.ParsePgError(err)
	}
	return nil
}
