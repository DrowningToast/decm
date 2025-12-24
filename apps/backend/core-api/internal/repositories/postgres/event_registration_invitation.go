package postgres

import (
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"
	"time"

	datagateway "apps/backend/core-api/internal/datagateway"

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
		Id:                  result.ID,
		EventId:             result.EventID,
		InboxMessageId:      result.InboxMessageID,
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
		AcceptedAt:          pgmapper.PgTimestampzToTimePtr(result.AcceptedAt),
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
		Id:                  result.ID,
		EventId:             result.EventID,
		InboxMessageId:      result.InboxMessageID,
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
		AcceptedAt:          pgmapper.PgTimestampzToTimePtr(result.AcceptedAt),
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
		Id:                  result.ID,
		EventId:             result.EventID,
		InboxMessageId:      result.InboxMessageID,
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
		AcceptedAt:          pgmapper.PgTimestampzToTimePtr(result.AcceptedAt),
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
			Id:                  result.ID,
			EventId:             result.EventID,
			InboxMessageId:      result.InboxMessageID,
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
			AcceptedAt:          pgmapper.PgTimestampzToTimePtr(result.AcceptedAt),
		}
	}

	return invitations, nil
}

func (r *Repository) GetEventRegistrationInvitationByEventIDAndCredential(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID, email *string, walletAddress *string) (*entity.EventRegistrationInvitation, *entity.InboxMessage, error) {
	encryptedEmail, err := pgmapper.EncryptStringPtrToPgText(email, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	params := generated.GetEventRegistrationInvitationByEventIdAndCredentialIdParams{
		EventID:       eventId,
		CredentialID:  pgmapper.UUIDToPgUUID(&credentialId),
		Email:         encryptedEmail,
		WalletAddress: pgmapper.StringPtrToPgText(walletAddress),
	}
	result, err := r.queries.GetEventRegistrationInvitationByEventIdAndCredentialId(ctx, params)
	if err != nil {
		return nil, nil, pgerrutils.ParsePgError(err)
	}

	// decrypt
	firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(result.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	emailDec, err := pgmapper.DecryptPgTextToStringPtr(result.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(result.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(result.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}

	registrationInvitation := &entity.EventRegistrationInvitation{
		Id:                  result.EventRegistrationInvitationID,
		EventId:             result.EventID,
		InboxMessageId:      result.InboxMessageID,
		ValidUntil:          pgmapper.PgTimestampzToTimePtr(result.ValidUntil),
		Code:                pgmapper.PgTextToStringPtr(result.Code),
		FirstName:           firstNameDec,
		LastName:            lastNameDec,
		Email:               emailDec,
		PhoneNumber:         phoneNumberDec,
		AcademicInstitution: academicInstitutionDec,
		CreatedAt:           *pgmapper.PgTimestampzToTimePtr(result.EventRegistrationInvitationCreatedAt),
		UpdatedAt:           *pgmapper.PgTimestampzToTimePtr(result.EventRegistrationInvitationUpdatedAt),
		CancelledAt:         pgmapper.PgTimestampzToTimePtr(result.EventRegistrationInvitationCancelledAt),
		AcceptedAt:          pgmapper.PgTimestampzToTimePtr(result.EventRegistrationInvitationAcceptedAt),
	}

	inboxMessage := &entity.InboxMessage{
		Id:                     result.InboxMessageID,
		SenderCredentialId:     &result.SenderCredentialID,
		ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:          emailDec,
		ReceiverWalletAddress:  pgmapper.PgTextToStringPtr(result.ReceiverWalletAddress),
		MessageType:            entity.InboxMessageType(result.MessageType),
		MessageContent:         string(result.MessageContent),
		FallbackMessageContent: pgmapper.PgTextToStringPtr(result.FallbackMessageContent),
		IsRead:                 int(result.InboxMessageIsRead.Int32),
		CreatedAt:              *pgmapper.PgTimestampzToTimePtr(result.InboxMessageCreatedAt),
		UpdatedAt:              *pgmapper.PgTimestampzToTimePtr(result.InboxMessageUpdatedAt),
		HiddenAt:               pgmapper.PgTimestampzToTimePtr(result.InboxMessageHiddenAt),
		DeletedAt:              pgmapper.PgTimestampzToTimePtr(result.InboxMessageDeletedAt),
	}

	return registrationInvitation, inboxMessage, nil
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
		AcceptedAt:          pgmapper.TimePtrToPgTimestampz(params.AcceptedAt),
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
		Id:                  result.ID,
		EventId:             result.EventID,
		InboxMessageId:      result.InboxMessageID,
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
		AcceptedAt:          pgmapper.PgTimestampzToTimePtr(result.AcceptedAt),
	}, nil
}

func (r *Repository) UpdateEventRegistrationInvitationAcceptedStatus(ctx context.Context, id uuid.UUID, acceptedAt *time.Time) (*entity.EventRegistrationInvitation, error) {
	result, err := r.queries.UpdateEventRegistrationInvitationAcceptedStatus(ctx, generated.UpdateEventRegistrationInvitationAcceptedStatusParams{
		ID:         id,
		AcceptedAt: pgmapper.TimePtrToPgTimestampz(acceptedAt),
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
		Id:                  result.ID,
		EventId:             result.EventID,
		InboxMessageId:      result.InboxMessageID,
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
		AcceptedAt:          pgmapper.PgTimestampzToTimePtr(result.AcceptedAt),
	}, nil
}

func (r *Repository) DeleteEventRegistrationInvitation(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteEventRegistrationInvitation(ctx, id)
	if err != nil {
		return pgerrutils.ParsePgError(err)
	}
	return nil
}
