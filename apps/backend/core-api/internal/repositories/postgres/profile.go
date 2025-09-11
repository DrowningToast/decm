package postgres

import (
	"context"
	"decm-database/go/generated"
	"errors"

	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/pgmapper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ datagateway.ProfileDataGateway = (*Repository)(nil)

func (r *Repository) GetProfileById(ctx context.Context, id uuid.UUID) (*entity.Profile, *customerror.Err) {
	query, err := r.queries.GetProfileByID(ctx, generated.GetProfileByIDParams{
		EncryptionKey: r.piiEncryptionKey,
		ID:            id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customerror.TryParseAsCustomErr(&customerror.ErrNotFound, err)
		}
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	model := generated.Profile{
		ID:                          query.ID,
		AuthenticationCredentialID:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic,
		ProfilePictureUrl:           query.ProfilePictureUrl,
		IsFirstNamePublic:           query.IsFirstNamePublic,
		FirstName:                   query.FirstName,
		IsLastNamePublic:            query.IsLastNamePublic,
		LastName:                    query.LastName,
		IsEmailPublic:               query.IsEmailPublic,
		Email:                       query.Email,
		IsBioPublic:                 query.IsBioPublic,
		Bio:                         query.Bio,
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic,
		PhoneNumber:                 query.PhoneNumber,
		IsAddressPublic:             query.IsAddressPublic,
		Address:                     query.Address,
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic,
		AcademicInstitution:         query.AcademicInstitution,
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic,
		AcademicEmail:               query.AcademicEmail,
		CreatedAt:                   query.CreatedAt,
		UpdatedAt:                   query.UpdatedAt,
	}

	entity := entity.MapProfilePgModelsToEntities([]generated.Profile{model})[0]
	return &entity, nil
}

func (r *Repository) GetProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID) (*entity.Profile, *customerror.Err) {
	query, err := r.queries.GetProfileByAuthCredentialID(ctx, generated.GetProfileByAuthCredentialIDParams{
		EncryptionKey:              r.piiEncryptionKey,
		AuthenticationCredentialID: authenticationCredentialId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customerror.TryParseAsCustomErr(&customerror.ErrNotFound, err)
		}
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	model := generated.Profile{
		ID:                          query.ID,
		AuthenticationCredentialID:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic,
		ProfilePictureUrl:           query.ProfilePictureUrl,
		IsFirstNamePublic:           query.IsFirstNamePublic,
		FirstName:                   query.FirstName,
		IsLastNamePublic:            query.IsLastNamePublic,
		LastName:                    query.LastName,
		IsEmailPublic:               query.IsEmailPublic,
		Email:                       query.Email,
		IsBioPublic:                 query.IsBioPublic,
		Bio:                         query.Bio,
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic,
		PhoneNumber:                 query.PhoneNumber,
		IsAddressPublic:             query.IsAddressPublic,
		Address:                     query.Address,
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic,
		AcademicInstitution:         query.AcademicInstitution,
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic,
		AcademicEmail:               query.AcademicEmail,
		CreatedAt:                   query.CreatedAt,
		UpdatedAt:                   query.UpdatedAt,
	}

	entity := entity.MapProfilePgModelsToEntities([]generated.Profile{model})[0]
	return &entity, nil
}

func (r *Repository) GetProfileByEmail(ctx context.Context, email string) (*entity.Profile, *customerror.Err) {
	query, err := r.queries.GetProfileByEmail(ctx, generated.GetProfileByEmailParams{
		EncryptionKey: r.piiEncryptionKey,
		EmailSearch:   pgmapper.StringPtrToPgText(&email),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customerror.TryParseAsCustomErr(&customerror.ErrNotFound, err)
		}
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	model := generated.Profile{
		ID:                         query.ID,
		AuthenticationCredentialID: query.AuthenticationCredentialID,
		IsProfilePicturePublic:     query.IsProfilePicturePublic,
		ProfilePictureUrl:          query.ProfilePictureUrl,
		IsFirstNamePublic:          query.IsFirstNamePublic,
		FirstName:                  query.FirstName,
		IsLastNamePublic:           query.IsLastNamePublic,
		LastName:                   query.LastName,
	}

	entity := entity.MapProfilePgModelsToEntities([]generated.Profile{model})[0]
	return &entity, nil
}

func (r *Repository) CreateProfile(ctx context.Context, profile entity.Profile) (*entity.Profile, *customerror.Err) {
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
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}
	model := generated.Profile{
		ID:                          query.ID,
		AuthenticationCredentialID:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic,
		ProfilePictureUrl:           query.ProfilePictureUrl,
		IsFirstNamePublic:           query.IsFirstNamePublic,
		FirstName:                   query.FirstName,
		IsLastNamePublic:            query.IsLastNamePublic,
		LastName:                    query.LastName,
		IsEmailPublic:               query.IsEmailPublic,
		Email:                       query.Email,
		IsBioPublic:                 query.IsBioPublic,
		Bio:                         query.Bio,
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic,
		PhoneNumber:                 query.PhoneNumber,
		IsAddressPublic:             query.IsAddressPublic,
		Address:                     query.Address,
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic,
		AcademicInstitution:         query.AcademicInstitution,
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic,
		AcademicEmail:               query.AcademicEmail,
		CreatedAt:                   query.CreatedAt,
		UpdatedAt:                   query.UpdatedAt,
	}
	entity := entity.MapProfilePgModelsToEntities([]generated.Profile{model})[0]
	return &entity, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, id uuid.UUID, profile datagateway.UpdateProfileParameters) (*entity.Profile, *customerror.Err) {
	query, err := r.queries.UpdateProfile(ctx, generated.UpdateProfileParams{
		EncryptionKey:               r.piiEncryptionKey,
		ID:                          id,
		IsProfilePicturePublic:      pgmapper.BoolToPgInt4(profile.IsProfilePicturePublic),
		ProfilePictureUrl:           pgmapper.StringPtrToPgText(profile.ProfilePictureUrl),
		IsFirstNamePublic:           pgmapper.BoolToPgInt4(profile.IsFirstNamePublic),
		FirstName:                   pgmapper.StringPtrToPgText(profile.FirstName),
		IsLastNamePublic:            pgmapper.BoolToPgInt4(profile.IsLastNamePublic),
		LastName:                    pgmapper.StringPtrToPgText(profile.LastName),
		IsEmailPublic:               pgmapper.BoolToPgInt4(profile.IsEmailPublic),
		Email:                       pgmapper.StringPtrToPgText(profile.Email),
		IsBioPublic:                 pgmapper.BoolToPgInt4(profile.IsBioPublic),
		Bio:                         pgmapper.StringPtrToPgText(profile.Bio),
		IsPhoneNumberPublic:         pgmapper.BoolToPgInt4(profile.IsPhoneNumberPublic),
		PhoneNumber:                 pgmapper.StringPtrToPgText(profile.PhoneNumber),
		IsAddressPublic:             pgmapper.BoolToPgInt4(profile.IsAddressPublic),
		Address:                     pgmapper.StringPtrToPgText(profile.Address),
		IsAcademicInstitutionPublic: pgmapper.BoolToPgInt4(profile.IsAcademicInstitutionPublic),
		AcademicInstitution:         pgmapper.StringPtrToPgText(profile.AcademicInstitution),
		IsAcademicEmailPublic:       pgmapper.BoolToPgInt4(profile.IsAcademicEmailPublic),
		AcademicEmail:               pgmapper.StringPtrToPgText(profile.AcademicEmail),
	})
	if err != nil {
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}
	model := generated.Profile{
		ID:                         query.ID,
		AuthenticationCredentialID: query.AuthenticationCredentialID,
		IsProfilePicturePublic:     query.IsProfilePicturePublic,
		ProfilePictureUrl:          query.ProfilePictureUrl,
	}
	entity := entity.MapProfilePgModelsToEntities([]generated.Profile{model})[0]
	return &entity, nil
}

func (r *Repository) UpdateProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID, profile datagateway.UpdateProfileParameters) (*entity.Profile, *customerror.Err) {
	query, err := r.queries.UpdateProfileByAuthenticationCredentialId(ctx, generated.UpdateProfileByAuthenticationCredentialIdParams{
		EncryptionKey:               r.piiEncryptionKey,
		AuthenticationCredentialID:  authenticationCredentialId,
		IsProfilePicturePublic:      pgmapper.BoolToPgInt4(profile.IsProfilePicturePublic),
		ProfilePictureUrl:           pgmapper.StringPtrToPgText(profile.ProfilePictureUrl),
		IsFirstNamePublic:           pgmapper.BoolToPgInt4(profile.IsFirstNamePublic),
		FirstName:                   pgmapper.StringPtrToPgText(profile.FirstName),
		IsLastNamePublic:            pgmapper.BoolToPgInt4(profile.IsLastNamePublic),
		LastName:                    pgmapper.StringPtrToPgText(profile.LastName),
		IsEmailPublic:               pgmapper.BoolToPgInt4(profile.IsEmailPublic),
		Email:                       pgmapper.StringPtrToPgText(profile.Email),
		IsBioPublic:                 pgmapper.BoolToPgInt4(profile.IsBioPublic),
		Bio:                         pgmapper.StringPtrToPgText(profile.Bio),
		IsPhoneNumberPublic:         pgmapper.BoolToPgInt4(profile.IsPhoneNumberPublic),
		PhoneNumber:                 pgmapper.StringPtrToPgText(profile.PhoneNumber),
		IsAddressPublic:             pgmapper.BoolToPgInt4(profile.IsAddressPublic),
		Address:                     pgmapper.StringPtrToPgText(profile.Address),
		IsAcademicInstitutionPublic: pgmapper.BoolToPgInt4(profile.IsAcademicInstitutionPublic),
		AcademicInstitution:         pgmapper.StringPtrToPgText(profile.AcademicInstitution),
		IsAcademicEmailPublic:       pgmapper.BoolToPgInt4(profile.IsAcademicEmailPublic),
		AcademicEmail:               pgmapper.StringPtrToPgText(profile.AcademicEmail),
	})
	if err != nil {
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}
	model := generated.Profile{
		ID:                          query.ID,
		AuthenticationCredentialID:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic,
		ProfilePictureUrl:           query.ProfilePictureUrl,
		IsFirstNamePublic:           query.IsFirstNamePublic,
		FirstName:                   query.FirstName,
		IsLastNamePublic:            query.IsLastNamePublic,
		LastName:                    query.LastName,
		IsEmailPublic:               query.IsEmailPublic,
		Email:                       query.Email,
		IsBioPublic:                 query.IsBioPublic,
		Bio:                         query.Bio,
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic,
		PhoneNumber:                 query.PhoneNumber,
		IsAddressPublic:             query.IsAddressPublic,
		Address:                     query.Address,
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic,
		AcademicInstitution:         query.AcademicInstitution,
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic,
		AcademicEmail:               query.AcademicEmail,
	}
	entity := entity.MapProfilePgModelsToEntities([]generated.Profile{model})[0]
	return &entity, nil
}

func (r *Repository) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteProfile(ctx, id)
	if err != nil {
		return customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}
	return nil
}
