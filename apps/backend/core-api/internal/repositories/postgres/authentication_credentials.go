package postgres

import (
	"context"
	"decm-database/go/generated"
	"errors"

	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

	"github.com/jackc/pgx/v5"
)

var _ datagateway.AuthenticationCredentialDataGateway = (*Repository)(nil)

func (r *Repository) GetAuthenticationCredentialById(ctx context.Context, id int32) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.GetAuthenticationCredentialByID(ctx, id)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity.MapAuthenticationCredentialsToEntities([]generated.AuthenticationCredential{query})[0], nil
}

func (r *Repository) GetAuthenticationCredentialByPublicKey(ctx context.Context, publicKey string) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.GetAuthenticationCredentialByPublicKey(ctx, publicKey)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity.MapAuthenticationCredentialsToEntities([]generated.AuthenticationCredential{query})[0], nil
}

func (r *Repository) CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.CreateAuthenticationCredential(ctx, generated.CreateAuthenticationCredentialParams{
		SolutionStatus:      int32(credential.SolutionStatus),
		PublicKey:           credential.PublicKey,
		Password:            pgmapper.StringPtrToPgText(credential.Password),
		PrivateKey:          pgmapper.StringPtrToPgText(credential.PrivateKey),
		GoogleConnectorRef:  pgmapper.StringPtrToPgText(credential.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.StringPtrToPgText(credential.GithubConnectorRef),
		IsVerifiedOrganizer: pgmapper.BoolToInt(credential.IsVerifiedOrganizer),
		IsVerifiedStudent:   pgmapper.BoolToInt(credential.IsVerifiedStudent),
	})
	if err != nil {
		return nil, err
	}

	return &entity.MapAuthenticationCredentialsToEntities([]generated.AuthenticationCredential{query})[0], nil
}

func (r *Repository) UpdateAuthenticationCredential(ctx context.Context, id int32, params datagateway.UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error) {
	query, err := r.queries.UpdateAuthenticationCredential(ctx, generated.UpdateAuthenticationCredentialParams{
		ID:                  id,
		SolutionStatus:      int32(params.SolutionStatus),
		Password:            pgmapper.StringPtrToPgText(params.Password),
		PrivateKey:          pgmapper.StringPtrToPgText(params.PrivateKey),
		GoogleConnectorRef:  pgmapper.StringPtrToPgText(params.GoogleConnectorRef),
		GithubConnectorRef:  pgmapper.StringPtrToPgText(params.GithubConnectorRef),
		IsVerifiedOrganizer: pgmapper.BoolToInt(params.IsVerifiedOrganizer),
		IsVerifiedStudent:   pgmapper.BoolToInt(params.IsVerifiedStudent),
	})
	if err != nil {
		return nil, err
	}

	return &entity.MapAuthenticationCredentialsToEntities([]generated.AuthenticationCredential{query})[0], nil
}
