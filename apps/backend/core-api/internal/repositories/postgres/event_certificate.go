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

var _ datagateway.EventCertificateDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventCertificate(ctx context.Context, params datagateway.CreateEventCertificateParameters) (*entity.EventCertificate, error) {
	// Encrypt PII fields
	nameEnc, err := pgmapper.EncryptStringPtrToPgText(params.Name, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	academicInstitutionEnc, err := pgmapper.EncryptStringPtrToPgText(params.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.CreateEventCertificate(ctx, generated.CreateEventCertificateParams{
		EventID:                 params.EventID,
		ReceiverCredentialID:    pgmapper.UUIDToPgUUID(params.ReceiverCredentialID),
		ReceiverEmail:           pgmapper.StringPtrToPgText(params.ReceiverEmail),
		Name:                    nameEnc,
		AcademicInstitution:     academicInstitutionEnc,
		CertificateTitle:        pgmapper.StringPtrToPgText(params.CertificateTitle),
		CertificateSubtitle:     pgmapper.StringPtrToPgText(params.CertificateSubtitle),
		EventContractAddress:    pgmapper.StringPtrToPgText(&params.EventContractAddress),
		EventCertificateAddress: pgmapper.StringPtrToPgText(params.EventCertificateAddress),
		CertificateTokenID:      pgmapper.StringPtrToPgText(params.CertificateTokenID),
		CertificateDigest:       pgmapper.StringPtrToPgText(params.Digest),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields for return
	nameDec, err := pgmapper.DecryptPgTextToStringPtr(result.Name, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(result.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventCertificate{
		Id:                      result.ID,
		EventId:                 result.EventID,
		ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:           pgmapper.PgTextToStringPtr(result.ReceiverEmail),
		Name:                    nameDec,
		AcademicInstitution:     academicInstitutionDec,
		CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
		CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
		EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
		EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
		CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
		CreatedAt:               result.CreatedAt.Time,
		RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
	}, nil
}

func (r *Repository) GetEventCertificateByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificate, error) {
	result, err := r.queries.GetEventCertificateByID(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields
	name, err := pgmapper.DecryptPgTextToStringPtr(result.Name, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	academicInstitution, err := pgmapper.DecryptPgTextToStringPtr(result.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventCertificate{
		Id:                      result.ID,
		EventId:                 result.EventID,
		ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:           pgmapper.PgTextToStringPtr(result.ReceiverEmail),
		Name:                    name,
		AcademicInstitution:     academicInstitution,
		CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
		CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
		EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
		EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
		CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
		CreatedAt:               result.CreatedAt.Time,
		RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
	}, nil
}

func (r *Repository) GetEventCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	results, err := r.queries.GetEventCertificatesByEventID(ctx, eventID)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	certificates := make([]*entity.EventCertificate, len(results))
	for i, result := range results {
		// Decrypt PII fields
		name, err := pgmapper.DecryptPgTextToStringPtr(result.Name, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		academicInstitution, err := pgmapper.DecryptPgTextToStringPtr(result.AcademicInstitution, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		certificates[i] = &entity.EventCertificate{
			Id:                      result.ID,
			EventId:                 result.EventID,
			ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:           pgmapper.PgTextToStringPtr(result.ReceiverEmail),
			Name:                    name,
			AcademicInstitution:     academicInstitution,
			CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
			CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
			EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
			EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
			CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
			CreatedAt:               result.CreatedAt.Time,
			RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
		}
	}

	return certificates, nil
}

func (r *Repository) UpdateEventCertificate(ctx context.Context, id uuid.UUID, params datagateway.UpdateEventCertificateParameters) (*entity.EventCertificate, error) {
	// Encrypt PII fields
	nameEnc, err := pgmapper.EncryptStringPtrToPgText(params.Name, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	academicInstitutionEnc, err := pgmapper.EncryptStringPtrToPgText(params.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.UpdateEventCertificate(ctx, generated.UpdateEventCertificateParams{
		ID:                      id,
		ReceiverCredentialID:    pgmapper.UUIDToPgUUID(params.ReceiverCredentialID),
		ReceiverEmail:           pgmapper.StringPtrToPgText(params.ReceiverEmail),
		Name:                    nameEnc,
		AcademicInstitution:     academicInstitutionEnc,
		CertificateTitle:        pgmapper.StringPtrToPgText(params.CertificateTitle),
		CertificateSubtitle:     pgmapper.StringPtrToPgText(params.CertificateSubtitle),
		EventContractAddress:    pgmapper.StringPtrToPgText(params.EventContractAddress),
		EventCertificateAddress: pgmapper.StringPtrToPgText(params.EventCertificateAddress),
		CertificateTokenID:      pgmapper.StringPtrToPgText(params.CertificateTokenID),
		RevokedAt:               pgmapper.TimePtrToPgTimestampz(params.RevokedAt),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields for return
	nameDec, err := pgmapper.DecryptPgTextToStringPtr(result.Name, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(result.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventCertificate{
		Id:                      result.ID,
		EventId:                 result.EventID,
		ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:           pgmapper.PgTextToStringPtr(result.ReceiverEmail),
		Name:                    nameDec,
		AcademicInstitution:     academicInstitutionDec,
		CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
		CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
		EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
		EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
		CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
		CreatedAt:               result.CreatedAt.Time,
		RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
	}, nil
}

func (r *Repository) DeleteEventCertificate(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteEventCertificate(ctx, id)
	if err != nil {
		return pgerrutils.ParsePgError(err)
	}
	return nil
}
