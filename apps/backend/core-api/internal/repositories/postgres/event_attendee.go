package postgres

import (
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
)

var _ datagateway.EventAttendeeDataGateway = (*Repository)(nil)

func (r *Repository) GetEventAttendeeByEventIdAndCredentialId(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*entity.EventAttendee, error) {
	query, err := r.queries.GetEventAttendeeByEventIDAndCredentialID(ctx, generated.GetEventAttendeeByEventIDAndCredentialIDParams{
		EventID:              eventID,
		AttendeeCredentialID: credentialID,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	firstName, err := pgmapper.DecryptPgTextToStringPtr(query.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastName, err := pgmapper.DecryptPgTextToStringPtr(query.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	email, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bio, err := pgmapper.DecryptPgTextToStringPtr(query.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := pgmapper.DecryptPgTextToStringPtr(query.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	address, err := pgmapper.DecryptPgTextToStringPtr(query.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitution, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmail, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventAttendee{
		Id:                   query.ID,
		EventId:              query.EventID,
		AttendeeCredentialId: query.AttendeeCredentialID,
		ContractAddress:      query.ContractAddress,
		IsAttendeeAccepted:   query.IsAttendeeAccepted == 1,
		FirstName:            firstName,
		LastName:             lastName,
		Email:                email,
		Bio:                  bio,
		PhoneNumber:          phoneNumber,
		Address:              address,
		AcademicInstitution:  academicInstitution,
		AcademicEmail:        academicEmail,
		CreatedAt:            query.CreatedAt.Time,
		UpdatedAt:            query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) AddParticipant(ctx context.Context, params datagateway.AddParticipantParameters) (*entity.EventAttendee, error) {
	// Encrypt PII fields before storing
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
	academicEmailEnc, err := pgmapper.EncryptStringPtrToPgText(params.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	addressEnc, err := pgmapper.EncryptStringPtrToPgText(params.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bioEnc, err := pgmapper.EncryptStringPtrToPgText(params.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	query, err := r.queries.AddParticipant(ctx, generated.AddParticipantParams{
		EventID:              params.EventId,
		AttendeeCredentialID: params.CredentialId,
		ContractAddress:      params.ContractAddress,
		IsAttendeeAccepted:   pgmapper.BoolToInt32(params.IsParticipantAccepted),
		FirstName:            firstNameEnc,
		LastName:             lastNameEnc,
		Email:                emailEnc,
		PhoneNumber:          phoneNumberEnc,
		AcademicInstitution:  academicInstitutionEnc,
		AcademicEmail:        academicEmailEnc,
		Address:              addressEnc,
		Bio:                  bioEnc,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	firstName, err := pgmapper.DecryptPgTextToStringPtr(query.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastName, err := pgmapper.DecryptPgTextToStringPtr(query.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	email, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bio, err := pgmapper.DecryptPgTextToStringPtr(query.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := pgmapper.DecryptPgTextToStringPtr(query.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	address, err := pgmapper.DecryptPgTextToStringPtr(query.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitution, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmail, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.EventAttendee{
		Id:                   query.ID,
		EventId:              query.EventID,
		AttendeeCredentialId: query.AttendeeCredentialID,
		ContractAddress:      query.ContractAddress,
		IsAttendeeAccepted:   query.IsAttendeeAccepted == 1,
		FirstName:            firstName,
		LastName:             lastName,
		Email:                email,
		Bio:                  bio,
		PhoneNumber:          phoneNumber,
		Address:              address,
		AcademicInstitution:  academicInstitution,
		AcademicEmail:        academicEmail,
		CreatedAt:            query.CreatedAt.Time,
		UpdatedAt:            query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) ListEventAttendeesByEventID(ctx context.Context, eventId uuid.UUID) ([]entity.EventAttendee, error) {
	rows, err := r.queries.ListEventAttendeesByEventID(ctx, eventId)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	attendees := make([]entity.EventAttendee, 0, len(rows))
	for _, row := range rows {
		// Decrypt PII fields
		firstName, err := pgmapper.DecryptPgTextToStringPtr(row.FirstName, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		lastName, err := pgmapper.DecryptPgTextToStringPtr(row.LastName, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		email, err := pgmapper.DecryptPgTextToStringPtr(row.Email, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		phoneNumber, err := pgmapper.DecryptPgTextToStringPtr(row.PhoneNumber, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		academicInstitution, err := pgmapper.DecryptPgTextToStringPtr(row.AcademicInstitution, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		attendees = append(attendees, entity.EventAttendee{
			Id:                   row.ID,
			EventId:              row.EventID,
			AttendeeCredentialId: row.AttendeeCredentialID,
			ContractAddress:      row.ContractAddress,
			IsAttendeeAccepted:   row.IsAttendeeAccepted == 1,
			FirstName:            firstName,
			LastName:             lastName,
			Email:                email,
			PhoneNumber:          phoneNumber,
			AcademicInstitution:  academicInstitution,
			WalletAddress:        row.WalletAddress,
			CreatedAt:            row.CreatedAt.Time,
			UpdatedAt:            row.UpdatedAt.Time,
		})
	}

	return attendees, nil
}
