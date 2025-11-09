package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

var _ datagateway.EventCertificateSignatureDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventCertificateSignature(ctx context.Context, params datagateway.CreateEventCertificateSignatureParameters) (*entity.EventCertificateSignature, error) {
	result, err := r.queries.CreateEventCertificateSignature(ctx, generated.CreateEventCertificateSignatureParams{
		EventCertificateID: params.EventCertificateID,
		IssuerCredentialID: params.IssuerCredentialID,
		IssuerSignature:    pgmapper.StringPtrToPgText(params.IssuerSignature),
		HostSignature:      pgmapper.StringPtrToPgText(&params.HostSignature),
		SignMessage:        pgmapper.StringPtrToPgText(params.SignMessage),
		SignMessageDigest:  pgmapper.StringPtrToPgText(params.SignMessageDigest),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.EventCertificateSignature{
		ID:                 result.ID,
		EventCertificateID: result.EventCertificateID,
		IssuerCredentialID: result.IssuerCredentialID,
		IssuerSignature:    pgmapper.PgTextToStringPtr(result.IssuerSignature),
		HostSignature:      *pgmapper.PgTextToStringPtr(result.HostSignature),
		SignMessage:        pgmapper.PgTextToStringPtr(result.SignMessage),
		SignMessageDigest:  pgmapper.PgTextToStringPtr(result.SignMessageDigest),
	}, nil
}

func (r *Repository) GetEventCertificateSignatureByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificateSignature, error) {
	result, err := r.queries.GetEventCertificateSignatureByID(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.EventCertificateSignature{
		ID:                 result.ID,
		EventCertificateID: result.EventCertificateID,
		IssuerCredentialID: result.IssuerCredentialID,
		IssuerSignature:    pgmapper.PgTextToStringPtr(result.IssuerSignature),
		HostSignature:      *pgmapper.PgTextToStringPtr(result.HostSignature),
		SignMessage:        pgmapper.PgTextToStringPtr(result.SignMessage),
		SignMessageDigest:  pgmapper.PgTextToStringPtr(result.SignMessageDigest),
	}, nil
}

func (r *Repository) GetEventCertificateSignaturesByEventCertificateID(ctx context.Context, eventCertificateID uuid.UUID) ([]*entity.EventCertificateSignature, error) {
	results, err := r.queries.GetEventCertificateSignaturesByEventCertificateID(ctx, eventCertificateID)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]*entity.EventCertificateSignature, len(results))
	for i, result := range results {
		signatures[i] = &entity.EventCertificateSignature{
			ID:                 result.ID,
			EventCertificateID: result.EventCertificateID,
			IssuerCredentialID: result.IssuerCredentialID,
			IssuerSignature:    pgmapper.PgTextToStringPtr(result.IssuerSignature),
			HostSignature:      *pgmapper.PgTextToStringPtr(result.HostSignature),
			SignMessage:        pgmapper.PgTextToStringPtr(result.SignMessage),
			SignMessageDigest:  pgmapper.PgTextToStringPtr(result.SignMessageDigest),
		}
	}

	return signatures, nil
}

func (r *Repository) UpdateEventCertificateSignature(ctx context.Context, id uuid.UUID, params datagateway.UpdateEventCertificateSignatureParameters) (*entity.EventCertificateSignature, error) {
	result, err := r.queries.UpdateEventCertificateSignature(ctx, generated.UpdateEventCertificateSignatureParams{
		ID:                id,
		IssuerSignature:   pgmapper.StringPtrToPgText(params.IssuerSignature),
		HostSignature:     pgmapper.StringPtrToPgText(&params.HostSignature),
		SignMessage:       pgmapper.StringPtrToPgText(params.SignMessage),
		SignMessageDigest: pgmapper.StringPtrToPgText(params.SignMessageDigest),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.EventCertificateSignature{
		ID:                 result.ID,
		EventCertificateID: result.EventCertificateID,
		IssuerCredentialID: result.IssuerCredentialID,
		IssuerSignature:    pgmapper.PgTextToStringPtr(result.IssuerSignature),
		HostSignature:      *pgmapper.PgTextToStringPtr(result.HostSignature),
		SignMessage:        pgmapper.PgTextToStringPtr(result.SignMessage),
		SignMessageDigest:  pgmapper.PgTextToStringPtr(result.SignMessageDigest),
	}, nil
}

func (r *Repository) UpdateEventCertificateIssuerSignature(ctx context.Context, id uuid.UUID, issuerSignature *string) (*entity.EventCertificateSignature, error) {
	result, err := r.queries.UpdateEventCertificateIssuerSignature(ctx, generated.UpdateEventCertificateIssuerSignatureParams{
		ID:              id,
		IssuerSignature: pgmapper.StringPtrToPgText(issuerSignature),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.EventCertificateSignature{
		ID:                 result.ID,
		EventCertificateID: result.EventCertificateID,
		IssuerCredentialID: result.IssuerCredentialID,
		IssuerSignature:    pgmapper.PgTextToStringPtr(result.IssuerSignature),
		HostSignature:      *pgmapper.PgTextToStringPtr(result.HostSignature),
		SignMessage:        pgmapper.PgTextToStringPtr(result.SignMessage),
		SignMessageDigest:  pgmapper.PgTextToStringPtr(result.SignMessageDigest),
	}, nil
}

func (r *Repository) DeleteEventCertificateSignature(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteEventCertificateSignature(ctx, id)
	if err != nil {
		return pgerrutils.ParsePgError(err)
	}
	return nil
}
