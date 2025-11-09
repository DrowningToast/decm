package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	datagateway "apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

var _ datagateway.EventRegistrationInvitationDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventRegistrationInvitation(ctx context.Context, params datagateway.CreateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error) {
	// Encrypt PII fields using pgmapper
	firstNameEnc, err := pgmapper.EncryptStringPtrToPgText(params.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameEnc, err := pgmapper.EncryptStringPtrToPgText(params.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailEnc, err := pgmapper.EncryptStringPtrToPgText(params.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberEnc, err := pgmapper.EncryptStringPtrToPgText(params.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionEnc, err := pgmapper.EncryptStringPtrToPgText(params.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.CreateEventRegistrationInvitation(ctx, generated.CreateEventRegistrationInvitationParams{
		EventID:             params.EventID,
		InboxMessageID:      params.InboxMessageID,
		ValidUntil:          pgmapper.TimePtrToPgTimestampz(params.ValidUntil),
		Code:                pgmapper.StringPtrToPgText(params.Code),
		FirstName:           firstNameEnc,
		LastName:            lastNameEnc,
		Email:               emailEnc,
		PhoneNumber:         phoneNumberEnc,
		AcademicInstitution: academicInstitutionEnc,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields for return using pgmapper
	firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailDec, err := pgmapper.DecryptPgTextToStringPtr(result.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(result.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(result.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventRegistrationInvitation{
		ID:                  result.ID,
		EventID:             result.EventID,
		InboxMessageID:      result.InboxMessageID,
		ValidUntil:          pgmapper.PgTimestampzToTimePtr(result.ValidUntil),
		Code:                pgmapper.PgTextToStringPtr(result.Code),
		FirstName:           firstNameDec,
		LastName:            lastNameDec,
		Email:               emailDec,
		PhoneNumber:         phoneNumberDec,
		AcademicInstitution: academicInstitutionDec,
		CreatedAt:           *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:           *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		CancelledAt:         pgmapper.PgTimestampzToTimePtr(result.CancelledAt),
	}, nil
}

func (r *Repository) GetEventRegistrationInvitationByID(ctx context.Context, id uuid.UUID) (*entity.EventRegistrationInvitation, error) {
	result, err := r.queries.GetEventRegistrationInvitationByID(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields using pgmapper
	firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailDec, err := pgmapper.DecryptPgTextToStringPtr(result.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(result.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(result.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventRegistrationInvitation{
		ID:                  result.ID,
		EventID:             result.EventID,
		InboxMessageID:      result.InboxMessageID,
		ValidUntil:          pgmapper.PgTimestampzToTimePtr(result.ValidUntil),
		Code:                pgmapper.PgTextToStringPtr(result.Code),
		FirstName:           firstNameDec,
		LastName:            lastNameDec,
		Email:               emailDec,
		PhoneNumber:         phoneNumberDec,
		AcademicInstitution: academicInstitutionDec,
		CreatedAt:           *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:           *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		CancelledAt:         pgmapper.PgTimestampzToTimePtr(result.CancelledAt),
	}, nil
}

func (r *Repository) GetEventRegistrationInvitationByInboxMessageID(ctx context.Context, inboxMessageID uuid.UUID) (*entity.EventRegistrationInvitation, error) {
	result, err := r.queries.GetEventRegistrationInvitationByInboxMessageID(ctx, inboxMessageID)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields using pgmapper
	firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailDec, err := pgmapper.DecryptPgTextToStringPtr(result.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(result.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(result.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	return &entity.EventRegistrationInvitation{
		ID:                  result.ID,
		EventID:             result.EventID,
		InboxMessageID:      result.InboxMessageID,
		ValidUntil:          pgmapper.PgTimestampzToTimePtr(result.ValidUntil),
		Code:                pgmapper.PgTextToStringPtr(result.Code),
		FirstName:           firstNameDec,
		LastName:            lastNameDec,
		Email:               emailDec,
		PhoneNumber:         phoneNumberDec,
		AcademicInstitution: academicInstitutionDec,
		CreatedAt:           *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:           *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		CancelledAt:         pgmapper.PgTimestampzToTimePtr(result.CancelledAt),
	}, nil
}

func (r *Repository) GetEventRegistrationInvitationsByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventRegistrationInvitation, error) {
	results, err := r.queries.GetEventRegistrationInvitationsByEventID(ctx, eventID)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	invitations := make([]*entity.EventRegistrationInvitation, len(results))
	for i, result := range results {
		// Decrypt PII fields using pgmapper
		firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.FirstName, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.LastName, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		emailDec, err := pgmapper.DecryptPgTextToStringPtr(result.Email, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(result.PhoneNumber, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(result.AcademicInstitution, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		invitations[i] = &entity.EventRegistrationInvitation{
			ID:                  result.ID,
			EventID:             result.EventID,
			InboxMessageID:      result.InboxMessageID,
			ValidUntil:          pgmapper.PgTimestampzToTimePtr(result.ValidUntil),
			Code:                pgmapper.PgTextToStringPtr(result.Code),
			FirstName:           firstNameDec,
			LastName:            lastNameDec,
			Email:               emailDec,
			PhoneNumber:         phoneNumberDec,
			AcademicInstitution: academicInstitutionDec,
			CreatedAt:           *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
			UpdatedAt:           *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
			CancelledAt:         pgmapper.PgTimestampzToTimePtr(result.CancelledAt),
		}
	}

	return invitations, nil
}

func (r *Repository) UpdateEventRegistrationInvitation(ctx context.Context, id uuid.UUID, params datagateway.UpdateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error) {
	// Encrypt PII fields using pgmapper
	firstNameEnc, err := pgmapper.EncryptStringPtrToPgText(params.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameEnc, err := pgmapper.EncryptStringPtrToPgText(params.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailEnc, err := pgmapper.EncryptStringPtrToPgText(params.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberEnc, err := pgmapper.EncryptStringPtrToPgText(params.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionEnc, err := pgmapper.EncryptStringPtrToPgText(params.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.UpdateEventRegistrationInvitation(ctx, generated.UpdateEventRegistrationInvitationParams{
		ID:                  id,
		ValidUntil:          pgmapper.TimePtrToPgTimestampz(params.ValidUntil),
		Code:                pgmapper.StringPtrToPgText(params.Code),
		CancelledAt:         pgmapper.TimePtrToPgTimestampz(params.CancelledAt),
		FirstName:           firstNameEnc,
		LastName:            lastNameEnc,
		Email:               emailEnc,
		PhoneNumber:         phoneNumberEnc,
		AcademicInstitution: academicInstitutionEnc,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields for return using pgmapper
	firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailDec, err := pgmapper.DecryptPgTextToStringPtr(result.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(result.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(result.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventRegistrationInvitation{
		ID:                  result.ID,
		EventID:             result.EventID,
		InboxMessageID:      result.InboxMessageID,
		ValidUntil:          pgmapper.PgTimestampzToTimePtr(result.ValidUntil),
		Code:                pgmapper.PgTextToStringPtr(result.Code),
		FirstName:           firstNameDec,
		LastName:            lastNameDec,
		Email:               emailDec,
		PhoneNumber:         phoneNumberDec,
		AcademicInstitution: academicInstitutionDec,
		CreatedAt:           *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:           *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		CancelledAt:         pgmapper.PgTimestampzToTimePtr(result.CancelledAt),
	}, nil
}

func (r *Repository) DeleteEventRegistrationInvitation(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteEventRegistrationInvitation(ctx, id)
	if err != nil {
		return pgerrutils.ParsePgError(err)
	}
	return nil
}
