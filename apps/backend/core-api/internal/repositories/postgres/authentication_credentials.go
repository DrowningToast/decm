package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/log"
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ datagateway.AuthenticationCredentialDataGateway = (*Repository)(nil)

func (r *Repository) GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.GetAuthenticationCredentialById(ctx, generated.GetAuthenticationCredentialByIdParams{
		EncryptionKey: r.piiEncryptionKey,
		ID:            id,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.AuthenticationCredential{
		Id:                  query.ID,
		SolutionStatus:      entity.SolutionStatus(query.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(query.HashedPassword),
		EncryptedPrivateKey: nil, // EncryptedPrivateKey is []byte, not decrypted - kept as nil for security
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  pgmapper.PgTextToStringPtr(query.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.PgTextToStringPtr(query.GithubConnectorRef),
		IsVerifiedOrganizer: query.IsVerifiedOrganizer == 1,
		IsVerifiedStudent:   query.IsVerifiedStudent == 1,
		CreatedAt:           query.CreatedAt.Time,
		UpdatedAt:           query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.GetAuthenticationCredentialByWalletAddress(ctx, generated.GetAuthenticationCredentialByWalletAddressParams{
		EncryptionKey: r.piiEncryptionKey,
		WalletAddress: walletAddress,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.AuthenticationCredential{
		Id:                  query.ID,
		SolutionStatus:      entity.SolutionStatus(query.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(query.HashedPassword),
		EncryptedPrivateKey: nil, // EncryptedPrivateKey is []byte, not decrypted - kept as nil for security
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  pgmapper.PgTextToStringPtr(query.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.PgTextToStringPtr(query.GithubConnectorRef),
		IsVerifiedOrganizer: query.IsVerifiedOrganizer == 1,
		IsVerifiedStudent:   query.IsVerifiedStudent == 1,
		CreatedAt:           query.CreatedAt.Time,
		UpdatedAt:           query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetAuthenticationCredentialByGoogleConnectorRef(ctx context.Context, googleConnectorRef string) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.GetAuthenticationCredentialByGoogleConnectorRef(ctx, generated.GetAuthenticationCredentialByGoogleConnectorRefParams{
		EncryptionKey: pgmapper.StringPtrToPgText(&r.piiEncryptionKey),
		GoogleEmail:   pgmapper.StringPtrToPgText(&googleConnectorRef),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.AuthenticationCredential{
		Id:                  query.ID,
		SolutionStatus:      entity.SolutionStatus(query.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(query.HashedPassword),
		EncryptedPrivateKey: nil, // EncryptedPrivateKey is []byte, not decrypted - kept as nil for security
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  pgmapper.PgTextToStringPtr(query.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.PgTextToStringPtr(query.GithubConnectorRef),
		IsVerifiedOrganizer: query.IsVerifiedOrganizer == 1,
		IsVerifiedStudent:   query.IsVerifiedStudent == 1,
		CreatedAt:           query.CreatedAt.Time,
		UpdatedAt:           query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error) {
	// Convert EncryptedPrivateKey from *string to []byte if present
	var encryptedPrivateKey []byte
	if credential.EncryptedPrivateKey != nil {
		encryptedPrivateKey = []byte(*credential.EncryptedPrivateKey)
	}

	logger := log.LoadLogger()
	logger.Info("Creating authentication credential", "credential", credential)
	logger.Info("Google Connector Ref", "google_connector_ref", *credential.GoogleConnectorRef)

	query, err := r.queries.CreateAuthenticationCredential(ctx, generated.CreateAuthenticationCredentialParams{
		SolutionStatus:      int32(credential.SolutionStatus),
		WalletAddress:       credential.WalletAddress,
		HashedPassword:      pgmapper.StringPtrToPgText(credential.HashedPassword),
		EncryptedPrivateKey: encryptedPrivateKey,
		GoogleConnectorRef:  pgmapper.StringPtrToPgText(credential.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.StringPtrToPgText(credential.GithubConnectorRef),
		IsVerifiedOrganizer: pgmapper.BoolToInt32(credential.IsVerifiedOrganizer),
		IsVerifiedStudent:   pgmapper.BoolToInt32(credential.IsVerifiedStudent),
		EncryptionKey:       pgmapper.StringPtrToPgText(&r.piiEncryptionKey),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.AuthenticationCredential{
		Id:                  query.ID,
		SolutionStatus:      entity.SolutionStatus(query.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(query.HashedPassword),
		EncryptedPrivateKey: nil, // Don't return encrypted private key for security
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  pgmapper.PgTextToStringPtr(query.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.PgTextToStringPtr(query.GithubConnectorRef),
		IsVerifiedOrganizer: query.IsVerifiedOrganizer == 1,
		IsVerifiedStudent:   query.IsVerifiedStudent == 1,
		CreatedAt:           query.CreatedAt.Time,
		UpdatedAt:           query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) UpdateAuthenticationCredential(ctx context.Context, id uuid.UUID, params datagateway.UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error) {
	// Convert EncryptedPrivateKey from *string to []byte if present
	var encryptedPrivateKey []byte
	if params.EncryptedPrivateKey != nil {
		encryptedPrivateKey = []byte(*params.EncryptedPrivateKey)
	}

	query, err := r.queries.UpdateAuthenticationCredential(ctx, generated.UpdateAuthenticationCredentialParams{
		ID:                  id,
		SolutionStatus:      pgmapper.Int32ToPgInt4(int32(params.SolutionStatus)),
		HashedPassword:      pgmapper.StringPtrToPgText(params.HashedPassword),
		EncryptedPrivateKey: encryptedPrivateKey,
		WalletAddress:       pgtype.Text{}, // WalletAddress is not updateable
		GoogleConnectorRef:  pgmapper.StringPtrToPgText(params.GoogleConnectorRef),
		EncryptionKey:       r.piiEncryptionKey,
		GithubConnectorRef:  pgmapper.StringPtrToPgText(params.GithubConnectorRef),
		IsVerifiedOrganizer: pgmapper.BoolToPgInt4(params.IsVerifiedOrganizer),
		IsVerifiedStudent:   pgmapper.BoolToPgInt4(params.IsVerifiedStudent),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.AuthenticationCredential{
		Id:                  query.ID,
		SolutionStatus:      entity.SolutionStatus(query.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(query.HashedPassword),
		EncryptedPrivateKey: nil, // Don't return encrypted private key for security
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  pgmapper.PgTextToStringPtr(query.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.PgTextToStringPtr(query.GithubConnectorRef),
		IsVerifiedOrganizer: query.IsVerifiedOrganizer == 1,
		IsVerifiedStudent:   query.IsVerifiedStudent == 1,
		CreatedAt:           query.CreatedAt.Time,
		UpdatedAt:           query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) DeleteAuthenticationCredential(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteAuthenticationCredential(ctx, id)
	if err != nil {
		return pgerrutils.ParsePgError(err)
	}
	return nil
}
