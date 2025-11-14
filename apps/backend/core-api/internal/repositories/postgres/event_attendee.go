package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

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
	query, err := r.queries.AddParticipant(ctx, generated.AddParticipantParams{
		EventID:              params.EventId,
		AttendeeCredentialID: params.CredentialId,
		ContractAddress:      params.ContractAddress,
		IsAttendeeAccepted:   pgmapper.BoolToInt32(params.IsParticipantAccepted),
		FirstName:            pgmapper.StringPtrToPgText(params.FirstName),
		LastName:             pgmapper.StringPtrToPgText(params.LastName),
		Email:                pgmapper.StringPtrToPgText(params.Email),
		PhoneNumber:          pgmapper.StringPtrToPgText(params.PhoneNumber),
		AcademicInstitution:  pgmapper.StringPtrToPgText(params.AcademicInstitution),
		AcademicEmail:        pgmapper.StringPtrToPgText(params.AcademicEmail),
		Address:              pgmapper.StringPtrToPgText(params.Address),
		Bio:                  pgmapper.StringPtrToPgText(params.Bio),
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
