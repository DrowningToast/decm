package datagateway

import (
	"context"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type UpdateProfileParameters struct {
	IsProfilePicturePublic      *bool
	ProfilePictureUrl           *string
	IsFirstNamePublic           *bool
	FirstName                   *string
	IsLastNamePublic            *bool
	LastName                    *string
	IsEmailPublic               *bool
	Email                       *string
	IsBioPublic                 *bool
	Bio                         *string
	IsPhoneNumberPublic         *bool
	PhoneNumber                 *string
	IsAddressPublic             *bool
	Address                     *string
	IsAcademicInstitutionPublic *bool
	AcademicInstitution         *string
	IsAcademicEmailPublic       *bool
	AcademicEmail               *string
}

type ProfileDataGateway interface {
	GetProfileById(ctx context.Context, id uuid.UUID) (*entity.Profile, *customerror.Err)
	GetProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID) (*entity.Profile, *customerror.Err)
	GetProfileByEmail(ctx context.Context, email string) (*entity.Profile, *customerror.Err)
	CreateProfile(ctx context.Context, profile entity.Profile) (*entity.Profile, *customerror.Err)
	UpdateProfile(ctx context.Context, id uuid.UUID, profile UpdateProfileParameters) (*entity.Profile, *customerror.Err)
	UpdateProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID, profile UpdateProfileParameters) (*entity.Profile, *customerror.Err)
	DeleteProfile(ctx context.Context, id uuid.UUID) error
}
