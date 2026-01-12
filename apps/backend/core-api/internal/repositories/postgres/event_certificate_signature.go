package postgres

import (
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"

	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"

	"github.com/google/uuid"
)

var _ eventdatagateway.EventCertificateSignatureDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventCertificateSignature(ctx context.Context, params eventdatagateway.CreateEventCertificateSignatureParameters) (*entity.EventCertificateSignature, error) {
	result, err := r.queries.CreateEventCertificateSignature(ctx, generated.CreateEventCertificateSignatureParams{
		EventCertificateConfigID: params.EventCertificateConfigID,
		IssuerCredentialID:       params.IssuerCredentialID,
		IssuerSignature:          pgmapper.StringPtrToPgText(params.IssuerSignature),
		HostSignature:            pgmapper.StringPtrToPgText(&params.HostSignature),
		SignMessage:              pgmapper.StringPtrToPgText(params.SignMessage),
		SignMessageDigest:        pgmapper.StringPtrToPgText(params.SignMessageDigest),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.EventCertificateSignature{
		Id:                       result.ID,
		EventCertificateConfigId: result.EventCertificateConfigID,
		IssuerCredentialId:       result.IssuerCredentialID,
		IssuerSignature:          pgmapper.PgTextToStringPtr(result.IssuerSignature),
		HostSignature:            *pgmapper.PgTextToStringPtr(result.HostSignature),
		SignMessage:              pgmapper.PgTextToStringPtr(result.SignMessage),
		SignMessageDigest:        pgmapper.PgTextToStringPtr(result.SignMessageDigest),
	}, nil
}

func (r *Repository) GetEventCertificateSignatureByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificateSignature, error) {
	result, err := r.queries.GetEventCertificateSignatureByID(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.EventCertificateSignature{
		Id:                       result.ID,
		EventCertificateConfigId: result.EventCertificateConfigID,
		IssuerCredentialId:       result.IssuerCredentialID,
		IssuerSignature:          pgmapper.PgTextToStringPtr(result.IssuerSignature),
		HostSignature:            *pgmapper.PgTextToStringPtr(result.HostSignature),
		SignMessage:              pgmapper.PgTextToStringPtr(result.SignMessage),
		SignMessageDigest:        pgmapper.PgTextToStringPtr(result.SignMessageDigest),
	}, nil
}

func (r *Repository) GetEventCertificateSignaturesByEventCertificateConfigID(ctx context.Context, eventCertificateConfigID uuid.UUID) ([]*entity.EventCertificateSignature, error) {
	results, err := r.queries.GetEventCertificateSignaturesByEventCertificateConfigID(ctx, eventCertificateConfigID)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]*entity.EventCertificateSignature, len(results))
	for i, result := range results {
		signatures[i] = &entity.EventCertificateSignature{
			Id:                       result.ID,
			EventCertificateConfigId: result.EventCertificateConfigID,
			IssuerCredentialId:       result.IssuerCredentialID,
			IssuerSignature:          pgmapper.PgTextToStringPtr(result.IssuerSignature),
			HostSignature:            *pgmapper.PgTextToStringPtr(result.HostSignature),
			SignMessage:              pgmapper.PgTextToStringPtr(result.SignMessage),
			SignMessageDigest:        pgmapper.PgTextToStringPtr(result.SignMessageDigest),
		}
	}

	return signatures, nil
}

func (r *Repository) UpdateEventCertificateSignature(ctx context.Context, id uuid.UUID, params eventdatagateway.UpdateEventCertificateSignatureParameters) (*entity.EventCertificateSignature, error) {
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
		Id:                       result.ID,
		EventCertificateConfigId: result.EventCertificateConfigID,
		IssuerCredentialId:       result.IssuerCredentialID,
		IssuerSignature:          pgmapper.PgTextToStringPtr(result.IssuerSignature),
		HostSignature:            *pgmapper.PgTextToStringPtr(result.HostSignature),
		SignMessage:              pgmapper.PgTextToStringPtr(result.SignMessage),
		SignMessageDigest:        pgmapper.PgTextToStringPtr(result.SignMessageDigest),
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
		Id:                       result.ID,
		EventCertificateConfigId: result.EventCertificateConfigID,
		IssuerCredentialId:       result.IssuerCredentialID,
		IssuerSignature:          pgmapper.PgTextToStringPtr(result.IssuerSignature),
		HostSignature:            *pgmapper.PgTextToStringPtr(result.HostSignature),
		SignMessage:              pgmapper.PgTextToStringPtr(result.SignMessage),
		SignMessageDigest:        pgmapper.PgTextToStringPtr(result.SignMessageDigest),
	}, nil
}

func (r *Repository) DeleteEventCertificateSignature(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteEventCertificateSignature(ctx, id)
	if err != nil {
		return pgerrutils.ParsePgError(err)
	}
	return nil
}
