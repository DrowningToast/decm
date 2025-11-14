package profile

import (
	"context"
	"errors"

	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
)

type ProfileUsecase struct {
	ProfileDg                  datagateway.ProfileDataGateway
	AuthenticationCredentialDg datagateway.AuthenticationCredentialDataGateway

	AuthService *auth.AuthService
}

func NewProfileUsecase(profileDg datagateway.ProfileDataGateway, authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway, authService *auth.AuthService) *ProfileUsecase {
	return &ProfileUsecase{
		ProfileDg:                  profileDg,
		AuthenticationCredentialDg: authenticationCredentialDg,
		AuthService:                authService,
	}
}

func (u *ProfileUsecase) GetProfileById(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
	return u.ProfileDg.GetProfileById(ctx, id)
}

func (u *ProfileUsecase) GetProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID) (*entity.Profile, error) {
	return u.ProfileDg.GetProfileByAuthenticationCredentialId(ctx, authenticationCredentialId)
}

func (u *ProfileUsecase) GetProfileByEmail(ctx context.Context, email string) (*entity.Profile, error) {
	return u.ProfileDg.GetProfileByEmail(ctx, email)
}

type CreateProfileParameters struct {
	AuthenticationCredentialId  uuid.UUID `json:"authentication_credential_id" validate:"required"`
	IsProfilePicturePublic      bool      `json:"is_profile_picture_public,omitempty"`
	ProfilePictureUrl           *string   `json:"profile_picture_url,omitempty" validate:"omitempty,url"`
	IsFirstNamePublic           bool      `json:"is_first_name_public,omitempty"`
	FirstName                   *string   `json:"first_name,omitempty" validate:"omitempty,min=3,max=32"`
	IsLastNamePublic            bool      `json:"is_last_name_public,omitempty"`
	LastName                    *string   `json:"last_name,omitempty" validate:"omitempty,min=3,max=32"`
	IsEmailPublic               bool      `json:"is_email_public,omitempty"`
	Email                       *string   `json:"email,omitempty" validate:"omitempty,email,max=64"`
	IsBioPublic                 bool      `json:"is_bio_public,omitempty"`
	Bio                         *string   `json:"bio,omitempty" validate:"omitempty,min=10,max=255"`
	IsPhoneNumberPublic         bool      `json:"is_phone_number_public,omitempty"`
	PhoneNumber                 *string   `json:"phone_number,omitempty" validate:"omitempty,e164"`
	IsAddressPublic             bool      `json:"is_address_public,omitempty"`
	Address                     *string   `json:"address,omitempty" validate:"omitempty,min=10,max=255"`
	IsAcademicInstitutionPublic bool      `json:"is_academic_institution_public,omitempty"`
	AcademicInstitution         *string   `json:"academic_institution,omitempty" validate:"omitempty,min=3,max=255"`
	IsAcademicEmailPublic       bool      `json:"is_academic_email_public,omitempty"`
	AcademicEmail               *string   `json:"academic_email,omitempty" validate:"omitempty,email"`
}

func (p *CreateProfileParameters) IsValid() error {
	return validatorutils.ValidateStruct(p)
}

type UpdateProfileParameters struct {
	IsProfilePicturePublic      *bool   `json:"is_profile_picture_public,omitempty"`
	ProfilePictureUrl           *string `json:"profile_picture_url,omitempty" validate:"omitempty,url"`
	IsFirstNamePublic           *bool   `json:"is_first_name_public,omitempty"`
	FirstName                   *string `json:"first_name,omitempty" validate:"omitempty,min=3,max=32"`
	IsLastNamePublic            *bool   `json:"is_last_name_public,omitempty"`
	LastName                    *string `json:"last_name,omitempty" validate:"omitempty,min=3,max=32"`
	IsEmailPublic               *bool   `json:"is_email_public,omitempty"`
	Email                       *string `json:"email,omitempty" validate:"omitempty,email,max=64"`
	IsBioPublic                 *bool   `json:"is_bio_public,omitempty"`
	Bio                         *string `json:"bio,omitempty" validate:"omitempty,min=10,max=255"`
	IsPhoneNumberPublic         *bool   `json:"is_phone_number_public,omitempty"`
	PhoneNumber                 *string `json:"phone_number,omitempty" validate:"omitempty,e164"`
	IsAddressPublic             *bool   `json:"is_address_public,omitempty"`
	Address                     *string `json:"address,omitempty" validate:"omitempty,min=10,max=255"`
	IsAcademicInstitutionPublic *bool   `json:"is_academic_institution_public,omitempty"`
	AcademicInstitution         *string `json:"academic_institution,omitempty" validate:"omitempty,min=3,max=255"`
	IsAcademicEmailPublic       *bool   `json:"is_academic_email_public,omitempty"`
	AcademicEmail               *string `json:"academic_email,omitempty" validate:"omitempty,email"`
}

func (p *UpdateProfileParameters) IsValid() error {
	return validatorutils.ValidateStruct(p)
}

func (u *ProfileUsecase) UpdateProfile(ctx context.Context, credentialId uuid.UUID, profile UpdateProfileParameters) (*entity.Profile, error) {
	if err := profile.IsValid(); err != nil {
		return nil, err
	}

	existingProfile, err := u.ProfileDg.GetProfileByAuthenticationCredentialId(ctx, credentialId)
	if err != nil {
		var customErr *customerror.Err
		if errors.As(err, &customErr) {
			if customErr.Code != &customerror.ErrNotFound.Code {
				return nil, customErr.Extend("failed to check for existing profile")
			}
		}
	}

	if existingProfile == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("profile not found"))
	}

	return u.ProfileDg.UpdateProfile(ctx, existingProfile.Id, datagateway.UpdateProfileParameters{
		IsProfilePicturePublic:      profile.IsProfilePicturePublic,
		ProfilePictureUrl:           profile.ProfilePictureUrl,
		IsFirstNamePublic:           profile.IsFirstNamePublic,
		FirstName:                   profile.FirstName,
		IsLastNamePublic:            profile.IsLastNamePublic,
		LastName:                    profile.LastName,
		IsEmailPublic:               profile.IsEmailPublic,
		Email:                       profile.Email,
		IsBioPublic:                 profile.IsBioPublic,
		Bio:                         profile.Bio,
		IsPhoneNumberPublic:         profile.IsPhoneNumberPublic,
		PhoneNumber:                 profile.PhoneNumber,
		IsAddressPublic:             profile.IsAddressPublic,
		Address:                     profile.Address,
		IsAcademicInstitutionPublic: profile.IsAcademicInstitutionPublic,
		AcademicInstitution:         profile.AcademicInstitution,
		IsAcademicEmailPublic:       profile.IsAcademicEmailPublic,
		AcademicEmail:               profile.AcademicEmail,
	})
}

func (u *ProfileUsecase) UpdateProfileByCredentialId(ctx context.Context, credentialId uuid.UUID, profile UpdateProfileParameters) (*entity.Profile, error) {
	if err := profile.IsValid(); err != nil {
		return nil, err
	}

	existingProfile, err := u.ProfileDg.GetProfileByAuthenticationCredentialId(ctx, credentialId)
	if err != nil {
		var customErr *customerror.Err
		if errors.As(err, &customErr) {
			if customErr.Code != &customerror.ErrNotFound.Code {
				// Profile not found, create new profile
				return u.ProfileDg.CreateProfile(ctx, entity.Profile{
					AuthenticationCredentialId: credentialId,
					ProfilePictureUrl:          profile.ProfilePictureUrl,
				})
			}
			return nil, customErr.Extend("failed to check for existing profile")
		}

	}

	if existingProfile == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("profile not found"))
	}

	return u.ProfileDg.UpdateProfileByAuthenticationCredentialId(ctx, credentialId, datagateway.UpdateProfileParameters{
		ProfilePictureUrl:           profile.ProfilePictureUrl,
		IsFirstNamePublic:           profile.IsFirstNamePublic,
		FirstName:                   profile.FirstName,
		IsLastNamePublic:            profile.IsLastNamePublic,
		LastName:                    profile.LastName,
		IsEmailPublic:               profile.IsEmailPublic,
		Email:                       profile.Email,
		IsBioPublic:                 profile.IsBioPublic,
		Bio:                         profile.Bio,
		IsPhoneNumberPublic:         profile.IsPhoneNumberPublic,
		PhoneNumber:                 profile.PhoneNumber,
		IsAddressPublic:             profile.IsAddressPublic,
		Address:                     profile.Address,
		IsAcademicInstitutionPublic: profile.IsAcademicInstitutionPublic,
		AcademicInstitution:         profile.AcademicInstitution,
		IsAcademicEmailPublic:       profile.IsAcademicEmailPublic,
		AcademicEmail:               profile.AcademicEmail,
		IsProfilePicturePublic:      profile.IsProfilePicturePublic,
	})
}

func (u *ProfileUsecase) CreateProfile(ctx context.Context, profile CreateProfileParameters) (*entity.Profile, error) {
	return u.ProfileDg.CreateProfile(ctx, entity.Profile{
		AuthenticationCredentialId:  profile.AuthenticationCredentialId,
		IsProfilePicturePublic:      profile.IsProfilePicturePublic,
		ProfilePictureUrl:           profile.ProfilePictureUrl,
		IsFirstNamePublic:           profile.IsFirstNamePublic,
		FirstName:                   profile.FirstName,
		IsLastNamePublic:            profile.IsLastNamePublic,
		LastName:                    profile.LastName,
		IsEmailPublic:               profile.IsEmailPublic,
		Email:                       profile.Email,
		IsBioPublic:                 profile.IsBioPublic,
		Bio:                         profile.Bio,
		IsPhoneNumberPublic:         profile.IsPhoneNumberPublic,
		PhoneNumber:                 profile.PhoneNumber,
		IsAddressPublic:             profile.IsAddressPublic,
		Address:                     profile.Address,
		IsAcademicInstitutionPublic: profile.IsAcademicInstitutionPublic,
		AcademicInstitution:         profile.AcademicInstitution,
		IsAcademicEmailPublic:       profile.IsAcademicEmailPublic,
		AcademicEmail:               profile.AcademicEmail,
	})
}

func (u *ProfileUsecase) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	return u.ProfileDg.DeleteProfile(ctx, id)
}

func (u *ProfileUsecase) VerifyUserPassword(ctx context.Context, credentialId uuid.UUID, password string) (bool, error) {
	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, credentialId)
	if err != nil {
		return false, err
	}

	dbHashedPassword := credential.HashedPassword
	if dbHashedPassword == nil {
		return false, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("password not set"))
	}

	isValid, err := hashutils.CompareHash(password, *dbHashedPassword)
	if err != nil {
		return false, err
	}

	return isValid, nil
}
