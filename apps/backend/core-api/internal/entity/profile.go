package entity

import (
	"decm-database/go/generated"
	"time"

	"apps/backend/common/pgmapper"

	"github.com/google/uuid"
)

type Profile struct {
	Id                         uuid.UUID `json:"id"`
	AuthenticationCredentialId uuid.UUID `json:"authentication_credential_id"`

	IsProfilePicturePublic      bool    `json:"is_profile_picture_public"`
	ProfilePictureUrl           *string `json:"profile_picture_url"`
	IsFirstNamePublic           bool    `json:"is_first_name_public"`
	FirstName                   *string `json:"first_name"`
	IsLastNamePublic            bool    `json:"is_last_name_public"`
	LastName                    *string `json:"last_name"`
	IsEmailPublic               bool    `json:"is_email_public"`
	Email                       *string `json:"email"`
	IsBioPublic                 bool    `json:"is_bio_public"`
	Bio                         *string `json:"bio"`
	IsPhoneNumberPublic         bool    `json:"is_phone_number_public"`
	PhoneNumber                 *string `json:"phone_number"`
	IsAddressPublic             bool    `json:"is_address_public"`
	Address                     *string `json:"address"`
	IsAcademicInstitutionPublic bool    `json:"is_academic_institution_public"`
	AcademicInstitution         *string `json:"academic_institution"`
	IsAcademicEmailPublic       bool    `json:"is_academic_email_public"`
	AcademicEmail               *string `json:"academic_email"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (entity *Profile) ToPgModel() *generated.Profile {
	return &generated.Profile{
		ID:                          entity.Id,
		AuthenticationCredentialID:  entity.AuthenticationCredentialId,
		IsProfilePicturePublic:      pgmapper.BoolToInt32(entity.IsProfilePicturePublic),
		ProfilePictureUrl:           pgmapper.StringPtrToPgText(entity.ProfilePictureUrl),
		IsFirstNamePublic:           pgmapper.BoolToInt32(entity.IsFirstNamePublic),
		FirstName:                   pgmapper.StringPtrToPgText(entity.FirstName),
		IsLastNamePublic:            pgmapper.BoolToInt32(entity.IsLastNamePublic),
		LastName:                    pgmapper.StringPtrToPgText(entity.LastName),
		IsEmailPublic:               pgmapper.BoolToInt32(entity.IsEmailPublic),
		Email:                       pgmapper.StringPtrToPgText(entity.Email),
		IsBioPublic:                 pgmapper.BoolToInt32(entity.IsBioPublic),
		Bio:                         pgmapper.StringPtrToPgText(entity.Bio),
		IsPhoneNumberPublic:         pgmapper.BoolToInt32(entity.IsPhoneNumberPublic),
		PhoneNumber:                 pgmapper.StringPtrToPgText(entity.PhoneNumber),
		IsAddressPublic:             pgmapper.BoolToInt32(entity.IsAddressPublic),
		Address:                     pgmapper.StringPtrToPgText(entity.Address),
		IsAcademicInstitutionPublic: pgmapper.BoolToInt32(entity.IsAcademicInstitutionPublic),
		AcademicInstitution:         pgmapper.StringPtrToPgText(entity.AcademicInstitution),
		IsAcademicEmailPublic:       pgmapper.BoolToInt32(entity.IsAcademicEmailPublic),
		AcademicEmail:               pgmapper.StringPtrToPgText(entity.AcademicEmail),
		CreatedAt:                   pgmapper.TimePtrToPgTimestampz(&entity.CreatedAt),
		UpdatedAt:                   pgmapper.TimePtrToPgTimestampz(&entity.UpdatedAt),
	}
}

func MapProfilePgModelsToEntities(models []generated.Profile) []Profile {
	entities := make([]Profile, len(models))
	for i, model := range models {
		entities[i] = Profile{
			Id:                          model.ID,
			AuthenticationCredentialId:  model.AuthenticationCredentialID,
			IsProfilePicturePublic:      model.IsFirstNamePublic == 1,
			ProfilePictureUrl:           pgmapper.PgTextToStringPtr(model.ProfilePictureUrl),
			IsFirstNamePublic:           model.IsFirstNamePublic == 1,
			FirstName:                   pgmapper.PgTextToStringPtr(model.FirstName),
			IsLastNamePublic:            model.IsLastNamePublic == 1,
			LastName:                    pgmapper.PgTextToStringPtr(model.LastName),
			IsEmailPublic:               model.IsEmailPublic == 1,
			Email:                       pgmapper.PgTextToStringPtr(model.Email),
			IsBioPublic:                 model.IsBioPublic == 1,
			Bio:                         pgmapper.PgTextToStringPtr(model.Bio),
			IsPhoneNumberPublic:         model.IsPhoneNumberPublic == 1,
			PhoneNumber:                 pgmapper.PgTextToStringPtr(model.PhoneNumber),
			IsAddressPublic:             model.IsAddressPublic == 1,
			Address:                     pgmapper.PgTextToStringPtr(model.Address),
			IsAcademicInstitutionPublic: model.IsAcademicInstitutionPublic == 1,
			AcademicInstitution:         pgmapper.PgTextToStringPtr(model.AcademicInstitution),
			IsAcademicEmailPublic:       model.IsAcademicEmailPublic == 1,
			AcademicEmail:               pgmapper.PgTextToStringPtr(model.AcademicEmail),
			CreatedAt:                   model.CreatedAt.Time,
			UpdatedAt:                   model.UpdatedAt.Time,
		}
	}
	return entities
}
