package profile

import (
	"apps/backend/common"
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	"apps/backend/services/auth"
	"context"
	"time"
)

type GetMyProfileViewModel struct {
	AuthenticationCredentialId string `json:"authentication_credential_id"`
	ProfileId                  string `json:"profile_id"`

	IsProfilePicturePublic      bool    `json:"is_profile_picture_public"`
	ProfilePictureUrl           *string `json:"profile_picture_url,omitempty"`
	IsFirstNamePublic           bool    `json:"is_first_name_public"`
	FirstName                   *string `json:"first_name,omitempty"`
	IsLastNamePublic            bool    `json:"is_last_name_public"`
	LastName                    *string `json:"last_name,omitempty"`
	IsEmailPublic               bool    `json:"is_email_public"`
	Email                       *string `json:"email,omitempty"`
	IsBioPublic                 bool    `json:"is_bio_public"`
	Bio                         *string `json:"bio,omitempty"`
	IsPhoneNumberPublic         bool    `json:"is_phone_number_public"`
	PhoneNumber                 *string `json:"phone_number,omitempty"`
	IsAddressPublic             bool    `json:"is_address_public"`
	Address                     *string `json:"address,omitempty"`
	IsAcademicInstitutionPublic bool    `json:"is_academic_institution_public"`
	AcademicInstitution         *string `json:"academic_institution,omitempty"`
	IsAcademicEmailPublic       bool    `json:"is_academic_email_public"`
	AcademicEmail               *string `json:"academic_email,omitempty"`

	ProfileCreatedAt time.Time `json:"profile_created_at"`
	ProfileUpdatedAt time.Time `json:"profile_updated_at"`

	AuthenticationCredentialCreatedAt time.Time `json:"authentication_credential_created_at"`
	AuthenticationCredentialUpdatedAt time.Time `json:"authentication_credential_updated_at"`

	WalletAddress      string                `json:"wallet_address"`
	SolutionStatus     common.SolutionStatus `json:"solution_status"`
	GoogleConnectorRef *string               `json:"google_connector_ref,omitempty"`
	GithubConnectorRef *string               `json:"github_connector_ref,omitempty"`

	UnreadInboxMessageCount int `json:"unread_inbox_message_count"`
}

func (u *ProfileUsecase) GetMyProfileViewModel(ctx context.Context, user *auth.JwtClaims) (*GetMyProfileViewModel, error) {
	profile, authenticationCredential, err := u.ProfileDg.GetProfileAndCredentialWithCredentialId(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	unreadInboxMessageCount, err := u.InboxMessageDg.GetUnreadInboxMessageCountByCredentialID(ctx, offchain_datagateway.GetInboxMessagesByCredentialIDParameters{
		CredentialID:          user.UserId,
		ReceiverEmail:         authenticationCredential.GoogleConnectorRef,
		ReceiverWalletAddress: &authenticationCredential.WalletAddress,
	})
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
		UnreadInboxMessageCount:           unreadInboxMessageCount,
	}, nil
}
