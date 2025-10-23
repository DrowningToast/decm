package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common"
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
	query, err := r.queries.GetAuthenticationCredentialById(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields
	googleConnectorRef, err := pgmapper.DecryptPgTextToStringPtr(query.GoogleConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	githubConnectorRef, err := pgmapper.DecryptPgTextToStringPtr(query.GithubConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.AuthenticationCredential{
		Id:                  query.ID,
		SolutionStatus:      common.SolutionStatus(query.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(query.HashedPassword),
		EncryptedPrivateKey: nil, // EncryptedPrivateKey is []byte, not decrypted - kept as nil for security
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  googleConnectorRef,
		GithubConnectorRef:  githubConnectorRef,
		IsVerifiedOrganizer: query.IsVerifiedOrganizer == 1,
		IsVerifiedStudent:   query.IsVerifiedStudent == 1,
		CreatedAt:           query.CreatedAt.Time,
		UpdatedAt:           query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.GetAuthenticationCredentialByWalletAddress(ctx, walletAddress)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields
	googleConnectorRef, err := pgmapper.DecryptPgTextToStringPtr(query.GoogleConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	githubConnectorRef, err := pgmapper.DecryptPgTextToStringPtr(query.GithubConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.AuthenticationCredential{
		Id:                  query.ID,
		SolutionStatus:      common.SolutionStatus(query.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(query.HashedPassword),
		EncryptedPrivateKey: nil, // EncryptedPrivateKey is []byte, not decrypted - kept as nil for security
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  googleConnectorRef,
		GithubConnectorRef:  githubConnectorRef,
		IsVerifiedOrganizer: query.IsVerifiedOrganizer == 1,
		IsVerifiedStudent:   query.IsVerifiedStudent == 1,
		CreatedAt:           query.CreatedAt.Time,
		UpdatedAt:           query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetAuthenticationCredentialByGoogleConnectorRef(ctx context.Context, googleConnectorRef string) (*entity.AuthenticationCredential, error) {
	// Encrypt the connector ref for searching (database stores encrypted values)
	encryptedRef, err := pgmapper.EncryptPII(googleConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	query, err := r.queries.GetAuthenticationCredentialByGoogleConnectorRef(ctx, pgtype.Text{
		String: encryptedRef,
		Valid:  true,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields
	decryptedGoogleRef, err := pgmapper.DecryptPgTextToStringPtr(query.GoogleConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	githubConnectorRefDecrypted, err := pgmapper.DecryptPgTextToStringPtr(query.GithubConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.AuthenticationCredential{
		Id:                  query.ID,
		SolutionStatus:      common.SolutionStatus(query.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(query.HashedPassword),
		EncryptedPrivateKey: nil, // EncryptedPrivateKey is []byte, not decrypted - kept as nil for security
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  decryptedGoogleRef,
		GithubConnectorRef:  githubConnectorRefDecrypted,
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

	// Encrypt PII fields
	googleConnectorRefEncrypted, err := pgmapper.EncryptStringPtrToPgText(credential.GoogleConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	githubConnectorRefEncrypted, err := pgmapper.EncryptStringPtrToPgText(credential.GithubConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	query, err := r.queries.CreateAuthenticationCredential(ctx, generated.CreateAuthenticationCredentialParams{
		SolutionStatus:      int32(credential.SolutionStatus),
		WalletAddress:       credential.WalletAddress,
		HashedPassword:      pgmapper.StringPtrToPgText(credential.HashedPassword),
		EncryptedPrivateKey: encryptedPrivateKey,
		GoogleConnectorRef:  googleConnectorRefEncrypted,
		GithubConnectorRef:  githubConnectorRefEncrypted,
		IsVerifiedOrganizer: pgmapper.BoolToInt32(credential.IsVerifiedOrganizer),
		IsVerifiedStudent:   pgmapper.BoolToInt32(credential.IsVerifiedStudent),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields for return
	googleConnectorRefDecrypted, err := pgmapper.DecryptPgTextToStringPtr(query.GoogleConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	githubConnectorRefDecrypted, err := pgmapper.DecryptPgTextToStringPtr(query.GithubConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.AuthenticationCredential{
		Id:                  query.ID,
		SolutionStatus:      common.SolutionStatus(query.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(query.HashedPassword),
		EncryptedPrivateKey: nil, // Don't return encrypted private key for security
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  googleConnectorRefDecrypted,
		GithubConnectorRef:  githubConnectorRefDecrypted,
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

	// Encrypt PII fields
	googleConnectorRefEncrypted, err := pgmapper.EncryptStringPtrToPgText(params.GoogleConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	githubConnectorRefEncrypted, err := pgmapper.EncryptStringPtrToPgText(params.GithubConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	query, err := r.queries.UpdateAuthenticationCredential(ctx, generated.UpdateAuthenticationCredentialParams{
		ID:                  id,
		SolutionStatus:      pgmapper.Int32ToPgInt4(int32(params.SolutionStatus)),
		HashedPassword:      pgmapper.StringPtrToPgText(params.HashedPassword),
		EncryptedPrivateKey: encryptedPrivateKey,
		WalletAddress:       pgtype.Text{}, // WalletAddress is not updateable
		GoogleConnectorRef:  googleConnectorRefEncrypted,
		GithubConnectorRef:  githubConnectorRefEncrypted,
		IsVerifiedOrganizer: pgmapper.BoolToPgInt4(params.IsVerifiedOrganizer),
		IsVerifiedStudent:   pgmapper.BoolToPgInt4(params.IsVerifiedStudent),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields for return
	googleConnectorRefDecrypted, err := pgmapper.DecryptPgTextToStringPtr(query.GoogleConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	githubConnectorRefDecrypted, err := pgmapper.DecryptPgTextToStringPtr(query.GithubConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.AuthenticationCredential{
		Id:                  query.ID,
		SolutionStatus:      common.SolutionStatus(query.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(query.HashedPassword),
		EncryptedPrivateKey: nil, // Don't return encrypted private key for security
		WalletAddress:       query.WalletAddress,
		GoogleConnectorRef:  googleConnectorRefDecrypted,
		GithubConnectorRef:  githubConnectorRefDecrypted,
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
