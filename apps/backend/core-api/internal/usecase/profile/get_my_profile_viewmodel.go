package profile

import (
	"context"
	"time"

	"apps/backend/common"
	"apps/backend/services/auth"
)

type GetMyProfileViewModel struct {
	AuthenticationCredentialId string `json:"authentication_credential_id"`
	ProfileId                  string `json:"profile_id"`

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

	ProfileCreatedAt time.Time `json:"profile_created_at"`
	ProfileUpdatedAt time.Time `json:"profile_updated_at"`

	AuthenticationCredentialCreatedAt time.Time `json:"authentication_credential_created_at"`
	AuthenticationCredentialUpdatedAt time.Time `json:"authentication_credential_updated_at"`

	WalletAddress      string                `json:"wallet_address"`
	SolutionStatus     common.SolutionStatus `json:"solution_status"`
	GoogleConnectorRef *string               `json:"google_connector_ref"`
	GithubConnectorRef *string               `json:"github_connector_ref"`
}

func (u *ProfileUsecase) GetMyProfileViewModel(ctx context.Context, user *auth.JwtClaims) (*GetMyProfileViewModel, error) {
	profile, authenticationCredential, err := u.ProfileDg.GetProfileAndCredentialWithCredentialId(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	return &GetMyProfileViewModel{
		AuthenticationCredentialId:        authenticationCredential.Id.String(),
		ProfileId:                         profile.Id.String(),
		IsProfilePicturePublic:            profile.IsProfilePicturePublic,
		ProfilePictureUrl:                 profile.ProfilePictureUrl,
		IsFirstNamePublic:                 profile.IsFirstNamePublic,
		FirstName:                         profile.FirstName,
		IsLastNamePublic:                  profile.IsLastNamePublic,
		LastName:                          profile.LastName,
		IsEmailPublic:                     profile.IsEmailPublic,
		Email:                             profile.Email,
		IsBioPublic:                       profile.IsBioPublic,
		Bio:                               profile.Bio,
		IsPhoneNumberPublic:               profile.IsPhoneNumberPublic,
		PhoneNumber:                       profile.PhoneNumber,
		IsAddressPublic:                   profile.IsAddressPublic,
		Address:                           profile.Address,
		IsAcademicInstitutionPublic:       profile.IsAcademicInstitutionPublic,
		AcademicInstitution:               profile.AcademicInstitution,
		IsAcademicEmailPublic:             profile.IsAcademicEmailPublic,
		AcademicEmail:                     profile.AcademicEmail,
		ProfileCreatedAt:                  profile.CreatedAt,
		ProfileUpdatedAt:                  profile.UpdatedAt,
		WalletAddress:                     authenticationCredential.WalletAddress,
		SolutionStatus:                    authenticationCredential.SolutionStatus,
		AuthenticationCredentialCreatedAt: authenticationCredential.CreatedAt,
		AuthenticationCredentialUpdatedAt: authenticationCredential.UpdatedAt,
		GoogleConnectorRef:                authenticationCredential.GoogleConnectorRef,
		GithubConnectorRef:                authenticationCredential.GithubConnectorRef,
	}, nil
}
