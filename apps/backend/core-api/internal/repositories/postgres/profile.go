package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"

	"github.com/google/uuid"
)

var _ datagateway.ProfileDataGateway = (*Repository)(nil)

func (r *Repository) GetProfileById(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
	query, err := r.queries.GetProfileByID(ctx, generated.GetProfileByIDParams{
		EncryptionKey: r.piiEncryptionKey,
		ID:            id,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.Profile{
		Id:                          query.ID,
		AuthenticationCredentialId:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:           pgmapper.PgTextToStringPtr(query.ProfilePictureUrl),
		IsFirstNamePublic:           query.IsFirstNamePublic == 1,
		FirstName:                   pgmapper.PgTextToStringPtr(query.FirstName),
		IsLastNamePublic:            query.IsLastNamePublic == 1,
		LastName:                    pgmapper.PgTextToStringPtr(query.LastName),
		IsEmailPublic:               query.IsEmailPublic == 1,
		Email:                       pgmapper.PgTextToStringPtr(query.Email),
		IsBioPublic:                 query.IsBioPublic == 1,
		Bio:                         pgmapper.PgTextToStringPtr(query.Bio),
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic == 1,
		PhoneNumber:                 pgmapper.PgTextToStringPtr(query.PhoneNumber),
		IsAddressPublic:             query.IsAddressPublic == 1,
		Address:                     pgmapper.PgTextToStringPtr(query.Address),
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic == 1,
		AcademicInstitution:         pgmapper.PgTextToStringPtr(query.AcademicInstitution),
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic == 1,
		AcademicEmail:               pgmapper.PgTextToStringPtr(query.AcademicEmail),
		CreatedAt:                   query.CreatedAt.Time,
		UpdatedAt:                   query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID) (*entity.Profile, error) {
	query, err := r.queries.GetProfileByAuthCredentialID(ctx, generated.GetProfileByAuthCredentialIDParams{
		EncryptionKey:              r.piiEncryptionKey,
		AuthenticationCredentialID: authenticationCredentialId,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.Profile{
		Id:                          query.ID,
		AuthenticationCredentialId:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:           pgmapper.PgTextToStringPtr(query.ProfilePictureUrl),
		IsFirstNamePublic:           query.IsFirstNamePublic == 1,
		FirstName:                   pgmapper.PgTextToStringPtr(query.FirstName),
		IsLastNamePublic:            query.IsLastNamePublic == 1,
		LastName:                    pgmapper.PgTextToStringPtr(query.LastName),
		IsEmailPublic:               query.IsEmailPublic == 1,
		Email:                       pgmapper.PgTextToStringPtr(query.Email),
		IsBioPublic:                 query.IsBioPublic == 1,
		Bio:                         pgmapper.PgTextToStringPtr(query.Bio),
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic == 1,
		PhoneNumber:                 pgmapper.PgTextToStringPtr(query.PhoneNumber),
		IsAddressPublic:             query.IsAddressPublic == 1,
		Address:                     pgmapper.PgTextToStringPtr(query.Address),
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic == 1,
		AcademicInstitution:         pgmapper.PgTextToStringPtr(query.AcademicInstitution),
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic == 1,
		AcademicEmail:               pgmapper.PgTextToStringPtr(query.AcademicEmail),
		CreatedAt:                   query.CreatedAt.Time,
		UpdatedAt:                   query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetProfileByEmail(ctx context.Context, email string) (*entity.Profile, error) {
	query, err := r.queries.GetProfileByEmail(ctx, generated.GetProfileByEmailParams{
		EncryptionKey: r.piiEncryptionKey,
		EmailSearch:   []byte(email),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.Profile{
		Id:                         query.ID,
		AuthenticationCredentialId: query.AuthenticationCredentialID,
		IsProfilePicturePublic:     query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:          pgmapper.PgTextToStringPtr(query.ProfilePictureUrl),
		IsFirstNamePublic:          query.IsFirstNamePublic == 1,
		FirstName:                  pgmapper.PgTextToStringPtr(query.FirstName),
		IsLastNamePublic:           query.IsLastNamePublic == 1,
		LastName:                   pgmapper.PgTextToStringPtr(query.LastName),
		// Note: Email and other fields are not returned by this query
		CreatedAt: query.CreatedAt.Time,
		UpdatedAt: query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) CreateProfile(ctx context.Context, profile entity.Profile) (*entity.Profile, error) {
	query, err := r.queries.CreateProfile(ctx, generated.CreateProfileParams{
		EncryptionKey:               r.piiEncryptionKey,
		IsProfilePicturePublic:      pgmapper.BoolToInt32(profile.IsProfilePicturePublic),
		ProfilePictureUrl:           pgmapper.StringPtrToPgText(profile.ProfilePictureUrl),
		IsFirstNamePublic:           pgmapper.BoolToInt32(profile.IsFirstNamePublic),
		FirstName:                   pgmapper.StringPtrToPgText(profile.FirstName),
		IsLastNamePublic:            pgmapper.BoolToInt32(profile.IsLastNamePublic),
		LastName:                    pgmapper.StringPtrToPgText(profile.LastName),
		IsEmailPublic:               pgmapper.BoolToInt32(profile.IsEmailPublic),
		Email:                       pgmapper.StringPtrToPgText(profile.Email),
		IsBioPublic:                 pgmapper.BoolToInt32(profile.IsBioPublic),
		Bio:                         pgmapper.StringPtrToPgText(profile.Bio),
		IsPhoneNumberPublic:         pgmapper.BoolToInt32(profile.IsPhoneNumberPublic),
		PhoneNumber:                 pgmapper.StringPtrToPgText(profile.PhoneNumber),
		IsAddressPublic:             pgmapper.BoolToInt32(profile.IsAddressPublic),
		Address:                     pgmapper.StringPtrToPgText(profile.Address),
		IsAcademicInstitutionPublic: pgmapper.BoolToInt32(profile.IsAcademicInstitutionPublic),
		AcademicInstitution:         pgmapper.StringPtrToPgText(profile.AcademicInstitution),
		IsAcademicEmailPublic:       pgmapper.BoolToInt32(profile.IsAcademicEmailPublic),
		AcademicEmail:               pgmapper.StringPtrToPgText(profile.AcademicEmail),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.Profile{
		Id:                          query.ID,
		AuthenticationCredentialId:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:           pgmapper.PgTextToStringPtr(query.ProfilePictureUrl),
		IsFirstNamePublic:           query.IsFirstNamePublic == 1,
		FirstName:                   pgmapper.PgTextToStringPtr(query.FirstName),
		IsLastNamePublic:            query.IsLastNamePublic == 1,
		LastName:                    pgmapper.PgTextToStringPtr(query.LastName),
		IsEmailPublic:               query.IsEmailPublic == 1,
		Email:                       pgmapper.PgTextToStringPtr(query.Email),
		IsBioPublic:                 query.IsBioPublic == 1,
		Bio:                         pgmapper.PgTextToStringPtr(query.Bio),
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic == 1,
		PhoneNumber:                 pgmapper.PgTextToStringPtr(query.PhoneNumber),
		IsAddressPublic:             query.IsAddressPublic == 1,
		Address:                     pgmapper.PgTextToStringPtr(query.Address),
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic == 1,
		AcademicInstitution:         pgmapper.PgTextToStringPtr(query.AcademicInstitution),
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic == 1,
		AcademicEmail:               pgmapper.PgTextToStringPtr(query.AcademicEmail),
		CreatedAt:                   query.CreatedAt.Time,
		UpdatedAt:                   query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, id uuid.UUID, profile datagateway.UpdateProfileParameters) (*entity.Profile, error) {
	query, err := r.queries.UpdateProfile(ctx, generated.UpdateProfileParams{
		EncryptionKey:               r.piiEncryptionKey,
		ID:                          id,
		IsProfilePicturePublic:      pgmapper.BoolPtrToPgInt4(profile.IsProfilePicturePublic),
		ProfilePictureUrl:           pgmapper.StringPtrToPgText(profile.ProfilePictureUrl),
		IsFirstNamePublic:           pgmapper.BoolPtrToPgInt4(profile.IsFirstNamePublic),
		FirstName:                   pgmapper.StringPtrToPgText(profile.FirstName),
		IsLastNamePublic:            pgmapper.BoolPtrToPgInt4(profile.IsLastNamePublic),
		LastName:                    pgmapper.StringPtrToPgText(profile.LastName),
		IsEmailPublic:               pgmapper.BoolPtrToPgInt4(profile.IsEmailPublic),
		Email:                       pgmapper.StringPtrToPgText(profile.Email),
		IsBioPublic:                 pgmapper.BoolPtrToPgInt4(profile.IsBioPublic),
		Bio:                         pgmapper.StringPtrToPgText(profile.Bio),
		IsPhoneNumberPublic:         pgmapper.BoolPtrToPgInt4(profile.IsPhoneNumberPublic),
		PhoneNumber:                 pgmapper.StringPtrToPgText(profile.PhoneNumber),
		IsAddressPublic:             pgmapper.BoolPtrToPgInt4(profile.IsAddressPublic),
		Address:                     pgmapper.StringPtrToPgText(profile.Address),
		IsAcademicInstitutionPublic: pgmapper.BoolPtrToPgInt4(profile.IsAcademicInstitutionPublic),
		AcademicInstitution:         pgmapper.StringPtrToPgText(profile.AcademicInstitution),
		IsAcademicEmailPublic:       pgmapper.BoolPtrToPgInt4(profile.IsAcademicEmailPublic),
		AcademicEmail:               pgmapper.StringPtrToPgText(profile.AcademicEmail),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.Profile{
		Id:                         query.ID,
		AuthenticationCredentialId: query.AuthenticationCredentialID,
		IsProfilePicturePublic:     query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:          pgmapper.PgTextToStringPtr(query.ProfilePictureUrl),
		// Note: Only limited fields are returned by UpdateProfile query
		CreatedAt: query.CreatedAt.Time,
		UpdatedAt: query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) UpdateProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID, profile datagateway.UpdateProfileParameters) (*entity.Profile, error) {
	query, err := r.queries.UpdateProfileByAuthenticationCredentialId(ctx, generated.UpdateProfileByAuthenticationCredentialIdParams{
		EncryptionKey:               r.piiEncryptionKey,
		AuthenticationCredentialID:  authenticationCredentialId,
		IsProfilePicturePublic:      pgmapper.BoolPtrToPgInt4(profile.IsProfilePicturePublic),
		ProfilePictureUrl:           pgmapper.StringPtrToPgText(profile.ProfilePictureUrl),
		IsFirstNamePublic:           pgmapper.BoolPtrToPgInt4(profile.IsFirstNamePublic),
		FirstName:                   pgmapper.StringPtrToPgText(profile.FirstName),
		IsLastNamePublic:            pgmapper.BoolPtrToPgInt4(profile.IsLastNamePublic),
		LastName:                    pgmapper.StringPtrToPgText(profile.LastName),
		IsEmailPublic:               pgmapper.BoolPtrToPgInt4(profile.IsEmailPublic),
		Email:                       pgmapper.StringPtrToPgText(profile.Email),
		IsBioPublic:                 pgmapper.BoolPtrToPgInt4(profile.IsBioPublic),
		Bio:                         pgmapper.StringPtrToPgText(profile.Bio),
		IsPhoneNumberPublic:         pgmapper.BoolPtrToPgInt4(profile.IsPhoneNumberPublic),
		PhoneNumber:                 pgmapper.StringPtrToPgText(profile.PhoneNumber),
		IsAddressPublic:             pgmapper.BoolPtrToPgInt4(profile.IsAddressPublic),
		Address:                     pgmapper.StringPtrToPgText(profile.Address),
		IsAcademicInstitutionPublic: pgmapper.BoolPtrToPgInt4(profile.IsAcademicInstitutionPublic),
		AcademicInstitution:         pgmapper.StringPtrToPgText(profile.AcademicInstitution),
		IsAcademicEmailPublic:       pgmapper.BoolPtrToPgInt4(profile.IsAcademicEmailPublic),
		AcademicEmail:               pgmapper.StringPtrToPgText(profile.AcademicEmail),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.Profile{
		Id:                          query.ID,
		AuthenticationCredentialId:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:           pgmapper.PgTextToStringPtr(query.ProfilePictureUrl),
		IsFirstNamePublic:           query.IsFirstNamePublic == 1,
		FirstName:                   pgmapper.PgTextToStringPtr(query.FirstName),
		IsLastNamePublic:            query.IsLastNamePublic == 1,
		LastName:                    pgmapper.PgTextToStringPtr(query.LastName),
		IsEmailPublic:               query.IsEmailPublic == 1,
		Email:                       pgmapper.PgTextToStringPtr(query.Email),
		IsBioPublic:                 query.IsBioPublic == 1,
		Bio:                         pgmapper.PgTextToStringPtr(query.Bio),
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic == 1,
		PhoneNumber:                 pgmapper.PgTextToStringPtr(query.PhoneNumber),
		IsAddressPublic:             query.IsAddressPublic == 1,
		Address:                     pgmapper.PgTextToStringPtr(query.Address),
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic == 1,
		AcademicInstitution:         pgmapper.PgTextToStringPtr(query.AcademicInstitution),
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic == 1,
		AcademicEmail:               pgmapper.PgTextToStringPtr(query.AcademicEmail),
		CreatedAt:                   query.CreatedAt.Time,
		UpdatedAt:                   query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteProfile(ctx, id)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	return nil
}
