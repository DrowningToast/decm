package offchain_datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"

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
	GetProfileById(ctx context.Context, id uuid.UUID) (*entity.Profile, error)
	GetProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID) (*entity.Profile, error)
	GetProfileByEmail(ctx context.Context, email string) (*entity.Profile, error)
	GetProfileAndCredentialWithCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID) (*entity.Profile, *entity.AuthenticationCredential, error)
	CreateProfile(ctx context.Context, profile entity.Profile) (*entity.Profile, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, profile UpdateProfileParameters) (*entity.Profile, error)
	UpdateProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID, profile UpdateProfileParameters) (*entity.Profile, error)
	DeleteProfile(ctx context.Context, id uuid.UUID) error
	ListVerifiedIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error)
}
