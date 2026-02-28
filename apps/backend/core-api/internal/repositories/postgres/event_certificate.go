package postgres

import (
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"
	"strings"

	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"

	"github.com/google/uuid"
)

var _ event_datagateway.EventCertificateDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventCertificate(ctx context.Context, params event_datagateway.CreateEventCertificateParameters) (*entity.EventCertificate, error) {
	// Encrypt PII fields
	nameEnc, err := pgmapper.EncryptStringPtrToPgText(params.Name, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	academicInstitutionEnc, err := pgmapper.EncryptStringPtrToPgText(params.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	// Normalize email to lowercase for case-insensitive comparison
	normalizedEmail := params.ReceiverEmail
	if params.ReceiverEmail != nil && *params.ReceiverEmail != "" {
		lowercaseEmail := strings.ToLower(*params.ReceiverEmail)
		normalizedEmail = &lowercaseEmail
	}

	receiverEmailEnc, err := pgmapper.EncryptStringPtrToPgText(normalizedEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.CreateEventCertificate(ctx, generated.CreateEventCertificateParams{
		EventID:                 params.EventID,
		ReceiverCredentialID:    pgmapper.UUIDPtrToPgUUID(params.ReceiverCredentialID),
		ReceiverEmail:           receiverEmailEnc,
		Name:                    nameEnc,
		AcademicInstitution:     academicInstitutionEnc,
		CertificateTitle:        pgmapper.StringPtrToPgText(params.CertificateTitle),
		CertificateSubtitle:     pgmapper.StringPtrToPgText(params.CertificateSubtitle),
		EventContractAddress:    pgmapper.StringPtrToPgText(&params.EventContractAddress),
		EventCertificateAddress: pgmapper.StringPtrToPgText(params.EventCertificateAddress),
		CertificateTokenID:      pgmapper.StringPtrToPgText(params.CertificateTokenID),
		CertificateDigest:       pgmapper.StringPtrToPgText(params.Digest),
		UserClaimSignatureID:    pgmapper.UUIDPtrToPgUUID(params.UserClaimSignatureId),
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

	receiverEmailDec, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventCertificate{
		Id:                      result.ID,
		EventId:                 result.EventID,
		ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:           receiverEmailDec,
		Name:                    nameDec,
		AcademicInstitution:     academicInstitutionDec,
		CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
		CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
		EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
		EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
		CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
		CertificateDigest:       pgmapper.PgTextToStringPtr(result.CertificateDigest),
		InboxMessageId:          pgmapper.PgUUIDToUUIDPtr(result.InboxMessageID),
		UserClaimSignatureId:    pgmapper.PgUUIDToUUIDPtr(result.UserClaimSignatureID),
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

	receiverEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventCertificate{
		Id:                      result.ID,
		EventId:                 result.EventID,
		ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:           receiverEmail,
		Name:                    name,
		AcademicInstitution:     academicInstitution,
		CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
		CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
		EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
		EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
		CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
		CertificateDigest:       pgmapper.PgTextToStringPtr(result.CertificateDigest),
		InboxMessageId:          pgmapper.PgUUIDToUUIDPtr(result.InboxMessageID),
		UserClaimSignatureId:    pgmapper.PgUUIDToUUIDPtr(result.UserClaimSignatureID),
		CreatedAt:               result.CreatedAt.Time,
		RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
	}, nil
}

func (r *Repository) GetEventCertificateWithSignature(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*event_datagateway.EventCertificateWithSignature, error) {
	row, err := r.queries.GetEventCertificateWithSignature(ctx, generated.GetEventCertificateWithSignatureParams{
		EventID:              eventID,
		ReceiverCredentialID: pgmapper.UUIDToPgUUID(credentialID),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &event_datagateway.EventCertificateWithSignature{
		Id:                 row.ID,
		EventId:            row.EventID,
		CredentialId:       pgmapper.PgUUIDToUUIDPtr(row.ReceiverCredentialID),
		TokenId:            pgmapper.PgTextToStringPtr(row.CertificateTokenID),
		JoinedAt:           row.CreatedAt.Time,
		IsBroadcasted:      row.BroadcastedAt.Valid,
		BroadcastedAt:      pgmapper.PgTimestamptzToTimePtr(row.BroadcastedAt),
		SignatureCreatedAt: pgmapper.PgTimestamptzToTimePtr(row.SignatureCreatedAt),
		EstimatedDeadline:  pgmapper.PgTimestamptzToTimePtr(row.EstimatedDeadline),
		Signature:          row.Signature.String,
		SignMessage:        row.SignMessage.String,
	}, nil
}

func (r *Repository) GetEventCertificateByInboxMessageID(ctx context.Context, inboxMessageID uuid.UUID) (*entity.EventCertificate, error) {
	result, err := r.queries.GetEventCertificateByInboxMessageID(ctx, pgmapper.UUIDToPgUUID(inboxMessageID))
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

	receiverEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventCertificate{
		Id:                      result.ID,
		EventId:                 result.EventID,
		ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:           receiverEmail,
		Name:                    name,
		AcademicInstitution:     academicInstitution,
		CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
		CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
		EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
		EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
		CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
		CertificateDigest:       pgmapper.PgTextToStringPtr(result.CertificateDigest),
		InboxMessageId:          pgmapper.PgUUIDToUUIDPtr(result.InboxMessageID),
		UserClaimSignatureId:    pgmapper.PgUUIDToUUIDPtr(result.UserClaimSignatureID),
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

		receiverEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		certificates[i] = &entity.EventCertificate{
			Id:                      result.ID,
			EventId:                 result.EventID,
			ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:           receiverEmail,
			Name:                    name,
			AcademicInstitution:     academicInstitution,
			CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
			CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
			EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
			EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
			CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
			InboxMessageId:          pgmapper.PgUUIDToUUIDPtr(result.InboxMessageID),
			UserClaimSignatureId:    pgmapper.PgUUIDToUUIDPtr(result.UserClaimSignatureID),
			CreatedAt:               result.CreatedAt.Time,
			RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
		}
	}

	return certificates, nil
}

func (r *Repository) GetAllEventCertificateIDsByEventID(ctx context.Context, eventID uuid.UUID) ([]uuid.UUID, error) {
	results, err := r.queries.GetAllEventCertificateIDsByEventID(ctx, eventID)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// The query already returns UUIDs directly, so we can return them as-is
	return results, nil
}

func (r *Repository) UpdateEventCertificate(ctx context.Context, id uuid.UUID, params event_datagateway.UpdateEventCertificateParameters) (*entity.EventCertificate, error) {
	// Encrypt PII fields
	nameEnc, err := pgmapper.EncryptStringPtrToPgText(params.Name, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	academicInstitutionEnc, err := pgmapper.EncryptStringPtrToPgText(params.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	receiverEmailEnc, err := pgmapper.EncryptStringPtrToPgText(params.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.UpdateEventCertificate(ctx, generated.UpdateEventCertificateParams{
		ID:                      id,
		ReceiverCredentialID:    pgmapper.UUIDPtrToPgUUID(params.ReceiverCredentialID),
		ReceiverEmail:           receiverEmailEnc,
		Name:                    nameEnc,
		AcademicInstitution:     academicInstitutionEnc,
		CertificateTitle:        pgmapper.StringPtrToPgText(params.CertificateTitle),
		CertificateSubtitle:     pgmapper.StringPtrToPgText(params.CertificateSubtitle),
		EventContractAddress:    pgmapper.StringPtrToPgText(params.EventContractAddress),
		EventCertificateAddress: pgmapper.StringPtrToPgText(params.EventCertificateAddress),
		CertificateTokenID:      pgmapper.StringPtrToPgText(params.CertificateTokenID),
		UserClaimSignatureID:    pgmapper.UUIDPtrToPgUUID(params.UserClaimSignatureId),
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

	receiverEmailDec, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventCertificate{
		Id:                      result.ID,
		EventId:                 result.EventID,
		ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:           receiverEmailDec,
		Name:                    nameDec,
		AcademicInstitution:     academicInstitutionDec,
		CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
		CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
		EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
		EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
		CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
		CertificateDigest:       pgmapper.PgTextToStringPtr(result.CertificateDigest),
		InboxMessageId:          pgmapper.PgUUIDToUUIDPtr(result.InboxMessageID),
		UserClaimSignatureId:    pgmapper.PgUUIDToUUIDPtr(result.UserClaimSignatureID),
		CreatedAt:               result.CreatedAt.Time,
		RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
	}, nil
}

func (r *Repository) UpdateEventCertificateInboxMessageID(ctx context.Context, id uuid.UUID, inboxMessageID uuid.UUID) (*entity.EventCertificate, error) {
	result, err := r.queries.UpdateEventCertificateInboxMessageID(ctx, generated.UpdateEventCertificateInboxMessageIDParams{
		ID:             id,
		InboxMessageID: pgmapper.UUIDToPgUUID(inboxMessageID),
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

	receiverEmailDec, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventCertificate{
		Id:                      result.ID,
		EventId:                 result.EventID,
		ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:           receiverEmailDec,
		Name:                    nameDec,
		AcademicInstitution:     academicInstitutionDec,
		CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
		CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
		EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
		EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
		CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
		CertificateDigest:       pgmapper.PgTextToStringPtr(result.CertificateDigest),
		InboxMessageId:          pgmapper.PgUUIDToUUIDPtr(result.InboxMessageID),
		UserClaimSignatureId:    pgmapper.PgUUIDToUUIDPtr(result.UserClaimSignatureID),
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

func (r *Repository) GetClaimedCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	results, err := r.queries.GetClaimedCertificatesByEventID(ctx, eventID)
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

		receiverEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		certificates[i] = &entity.EventCertificate{
			Id:                      result.ID,
			EventId:                 result.EventID,
			ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:           receiverEmail,
			Name:                    name,
			AcademicInstitution:     academicInstitution,
			CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
			CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
			EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
			EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
			CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
			InboxMessageId:          pgmapper.PgUUIDToUUIDPtr(result.InboxMessageID),
			UserClaimSignatureId:    pgmapper.PgUUIDToUUIDPtr(result.UserClaimSignatureID),
			CreatedAt:               result.CreatedAt.Time,
			RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
		}
	}

	return certificates, nil
}

func (r *Repository) GetUnclaimedReadyCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	results, err := r.queries.GetUnclaimedReadyCertificatesByEventID(ctx, eventID)
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

		receiverEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		certificates[i] = &entity.EventCertificate{
			Id:                      result.ID,
			EventId:                 result.EventID,
			ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:           receiverEmail,
			Name:                    name,
			AcademicInstitution:     academicInstitution,
			CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
			CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
			EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
			EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
			CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
			InboxMessageId:          pgmapper.PgUUIDToUUIDPtr(result.InboxMessageID),
			UserClaimSignatureId:    pgmapper.PgUUIDToUUIDPtr(result.UserClaimSignatureID),
			CreatedAt:               result.CreatedAt.Time,
			RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
		}
	}

	return certificates, nil
}

func (r *Repository) GetClaimedCertificatesByCredentialID(ctx context.Context, credentialID uuid.UUID, email *string) ([]*entity.EventCertificate, error) {
	// Normalize email to lowercase for case-insensitive comparison
	normalizedEmail := email
	if email != nil && *email != "" {
		lowercaseEmail := strings.ToLower(*email)
		normalizedEmail = &lowercaseEmail
	}

	// Encrypt email for query
	encryptedEmail, err := pgmapper.EncryptStringPtrToPgText(normalizedEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.GetClaimedCertificatesByCredentialID(ctx, generated.GetClaimedCertificatesByCredentialIDParams{
		ReceiverCredentialID: pgmapper.UUIDToPgUUID(credentialID),
		ReceiverEmail:        encryptedEmail,
	})
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

		receiverEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		eventName := result.EventName
		certificates[i] = &entity.EventCertificate{
			Id:                      result.ID,
			EventId:                 result.EventID,
			EventName:               &eventName,
			ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:           receiverEmail,
			Name:                    name,
			AcademicInstitution:     academicInstitution,
			CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
			CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
			EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
			EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
			CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
			InboxMessageId:          pgmapper.PgUUIDToUUIDPtr(result.InboxMessageID),
			UserClaimSignatureId:    pgmapper.PgUUIDToUUIDPtr(result.UserClaimSignatureID),
			BroadcastedAt:           pgmapper.PgTimestamptzToTimePtr(result.BroadcastedAt),
			SignatureCreatedAt:      pgmapper.PgTimestamptzToTimePtr(result.SignatureCreatedAt),
			EstimatedDeadline:       pgmapper.PgTimestamptzToTimePtr(result.EstimatedDeadline),
			AbortedAt:               pgmapper.PgTimestamptzToTimePtr(result.AbortedAt),
			CreatedAt:               result.CreatedAt.Time,
			RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
		}
	}

	return certificates, nil
}

func (r *Repository) GetUnclaimedReadyCertificatesByCredentialID(ctx context.Context, credentialID uuid.UUID, email *string) ([]*entity.EventCertificate, error) {
	// Normalize email to lowercase for case-insensitive comparison
	normalizedEmail := email
	if email != nil && *email != "" {
		lowercaseEmail := strings.ToLower(*email)
		normalizedEmail = &lowercaseEmail
	}

	// Encrypt email for query
	encryptedEmail, err := pgmapper.EncryptStringPtrToPgText(normalizedEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.GetUnclaimedReadyCertificatesByCredentialID(ctx, generated.GetUnclaimedReadyCertificatesByCredentialIDParams{
		ReceiverCredentialID: pgmapper.UUIDToPgUUID(credentialID),
		ReceiverEmail:        encryptedEmail,
	})
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

		receiverEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		eventName := result.EventName
		certificates[i] = &entity.EventCertificate{
			Id:                      result.ID,
			EventId:                 result.EventID,
			EventName:               &eventName,
			ReceiverCredentialId:    pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:           receiverEmail,
			Name:                    name,
			AcademicInstitution:     academicInstitution,
			CertificateTitle:        pgmapper.PgTextToStringPtr(result.CertificateTitle),
			CertificateSubtitle:     pgmapper.PgTextToStringPtr(result.CertificateSubtitle),
			EventContractAddress:    *pgmapper.PgTextToStringPtr(result.EventContractAddress),
			EventCertificateAddress: pgmapper.PgTextToStringPtr(result.EventCertificateAddress),
			CertificateTokenId:      pgmapper.PgTextToStringPtr(result.CertificateTokenID),
			InboxMessageId:          pgmapper.PgUUIDToUUIDPtr(result.InboxMessageID),
			UserClaimSignatureId:    pgmapper.PgUUIDToUUIDPtr(result.UserClaimSignatureID),
			CreatedAt:               result.CreatedAt.Time,
			RevokedAt:               pgmapper.PgTimestampzToTimePtr(result.RevokedAt),
		}
	}

	return certificates, nil
}
