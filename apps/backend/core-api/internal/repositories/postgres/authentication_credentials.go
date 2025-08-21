package postgres

import (
	"context"
	"decm-database/go/generated"
	"errors"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ datagateway.AuthenticationCredentialDataGateway = (*Repository)(nil)

func (r *Repository) GetAuthenticationCredentialById(ctx context.Context, id int32) (*entity.AuthenticationCredential, *customerror.Err) {
	query, err := r.queries.GetAuthenticationCredentialByID(ctx, generated.GetAuthenticationCredentialByIDParams{
		EncryptionKey: r.piiEncryptionKey,
		ID:            id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customerror.TryParseAsCustomErr(&customerror.ErrNotFound, err)
		}
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	model := generated.AuthenticationCredential{
		ID:                  query.ID,
		SolutionStatus:      query.SolutionStatus,
		HashedPassword:      query.HashedPassword,
		EncryptedPrivateKey: query.EncryptedPrivateKey,
		PublicKey:           query.PublicKey,
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

func (r *Repository) GetAuthenticationCredentialByPublicKey(ctx context.Context, publicKey string) (*entity.AuthenticationCredential, *customerror.Err) {
	query, err := r.queries.GetAuthenticationCredentialByPublicKey(ctx, generated.GetAuthenticationCredentialByPublicKeyParams{
		EncryptionKey: r.piiEncryptionKey,
		PublicKey:     publicKey,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customerror.TryParseAsCustomErr(&customerror.ErrNotFound, err)
		}
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	model := generated.AuthenticationCredential{
		ID:                  query.ID,
		SolutionStatus:      query.SolutionStatus,
		HashedPassword:      query.HashedPassword,
		EncryptedPrivateKey: query.EncryptedPrivateKey,
		PublicKey:           query.PublicKey,
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

func (r *Repository) CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, *customerror.Err) {
	query, err := r.queries.CreateAuthenticationCredential(ctx, generated.CreateAuthenticationCredentialParams{
		SolutionStatus:      int32(credential.SolutionStatus),
		PublicKey:           credential.PublicKey,
		HashedPassword:      pgmapper.StringPtrToPgText(credential.HashedPassword),
		EncryptedPrivateKey: pgmapper.StringPtrToPgText(credential.EncryptedPrivateKey),
		GoogleConnectorRef:  pgmapper.StringPtrToPgText(credential.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.StringPtrToPgText(credential.GithubConnectorRef),
		IsVerifiedOrganizer: pgmapper.BoolToInt32(credential.IsVerifiedOrganizer),
		IsVerifiedStudent:   pgmapper.BoolToInt32(credential.IsVerifiedStudent),
		EncryptionKey:       r.piiEncryptionKey,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, customerror.TryParseAsCustomErr(&customerror.ErrDuplicateEntry, err)
			}
		}
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	model := generated.AuthenticationCredential{
		ID:                  query.ID,
		SolutionStatus:      query.SolutionStatus,
		HashedPassword:      query.HashedPassword,
		EncryptedPrivateKey: query.EncryptedPrivateKey,
		PublicKey:           query.PublicKey,
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

func (r *Repository) UpdateAuthenticationCredential(ctx context.Context, id int32, params datagateway.UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, *customerror.Err) {
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, customerror.TryParseAsCustomErr(&customerror.ErrDuplicateEntry, err)
			}
		}
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	model := generated.AuthenticationCredential{
		ID:                  query.ID,
		SolutionStatus:      query.SolutionStatus,
		HashedPassword:      query.HashedPassword,
		EncryptedPrivateKey: query.EncryptedPrivateKey,
		PublicKey:           query.PublicKey,
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

func (r *Repository) DeleteAuthenticationCredential(ctx context.Context, id int32) error {
	err := r.queries.DeleteAuthenticationCredential(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return customerror.TryParseAsCustomErr(&customerror.ErrNotFound, err)
		}
		return customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}
	return nil
}
