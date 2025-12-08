package issuer

import (
	"time"

	"github.com/google/uuid"
)

// VerifiedIssuerViewModel represents a verified issuer profile with all necessary information
type VerifiedIssuerViewModel struct {
	Id                          uuid.UUID `json:"id"`
	AuthenticationCredentialId  uuid.UUID `json:"authentication_credential_id"`
	IsProfilePicturePublic      bool      `json:"is_profile_picture_public"`
	ProfilePictureUrl           *string   `json:"profile_picture_url,omitempty"`
	IsFirstNamePublic           bool      `json:"is_first_name_public"`
	FirstName                   *string   `json:"first_name,omitempty"`
	IsLastNamePublic            bool      `json:"is_last_name_public"`
	LastName                    *string   `json:"last_name,omitempty"`
	IsEmailPublic               bool      `json:"is_email_public"`
	Email                       *string   `json:"email,omitempty"`
	IsBioPublic                 bool      `json:"is_bio_public"`
	Bio                         *string   `json:"bio,omitempty"`
	IsPhoneNumberPublic         bool      `json:"is_phone_number_public"`
	PhoneNumber                 *string   `json:"phone_number,omitempty"`
	IsAddressPublic             bool      `json:"is_address_public"`
	Address                     *string   `json:"address,omitempty"`
	IsAcademicInstitutionPublic bool      `json:"is_academic_institution_public"`
	AcademicInstitution         *string   `json:"academic_institution,omitempty"`
	IsAcademicEmailPublic       bool      `json:"is_academic_email_public"`
	AcademicEmail               *string   `json:"academic_email,omitempty"`
	GoogleConnectorRef          *string   `json:"google_connector_ref,omitempty"`
	WalletAddress               string    `json:"wallet_address"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}
