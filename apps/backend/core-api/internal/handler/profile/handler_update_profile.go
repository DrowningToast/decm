package profile

import (
	"apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/profile"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
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

func (r *UpdateProfileRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

type UpdateProfileResponse struct {
	entity.Profile
}

func (r *UpdateProfileRequest) Parse(c *fiber.Ctx) error {
	if err := c.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

// @Summary Update a profile by credential ID
// @Description Update a profile by credential ID
// @ID update-profile-by-credential-id
// @Tags Profile
// @Accept json
// @Produce json
// @Param credential_id path string true "Credential ID"
// @Param request body UpdateProfileRequest true "Profile update request"
// @Success 200 {object} UpdateProfileResponse
// @Failure 400 {object} customerror.Err
// @Failure 403 {object} customerror.Err
// @Failure 404 {object} customerror.Err
// @Failure 500 {object} customerror.Err
// @Router /api/v1/profile/credential/{credential_id} [patch]
func (h *Handler) UpdateProfileByCredentialId(c *fiber.Ctx) error {
	var upsertProfileRequest UpdateProfileRequest
	if err := upsertProfileRequest.Parse(c); err != nil {
		return err
	}
	if err := upsertProfileRequest.IsValid(); err != nil {
		return err
	}
	credentialId := c.Params("credential_id")
	credentialIdUUID, err := uuid.Parse(credentialId)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	// Auth check
	isCurrentUser, err := h.AuthenticationService.IsCurrentUser(c, credentialIdUUID)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if !isCurrentUser {
		return customerror.Parse(&customerror.ErrUnauthorized, errors.New("unauthorized"))
	}

	profile, customerr := h.ProfileUc.UpdateProfileByCredentialId(c.Context(), credentialIdUUID, profile.UpdateProfileParameters{
		IsProfilePicturePublic:      upsertProfileRequest.IsProfilePicturePublic,
		ProfilePictureUrl:           upsertProfileRequest.ProfilePictureUrl,
		IsFirstNamePublic:           upsertProfileRequest.IsFirstNamePublic,
		FirstName:                   upsertProfileRequest.FirstName,
		IsLastNamePublic:            upsertProfileRequest.IsLastNamePublic,
		LastName:                    upsertProfileRequest.LastName,
		IsEmailPublic:               upsertProfileRequest.IsEmailPublic,
		Email:                       upsertProfileRequest.Email,
		IsBioPublic:                 upsertProfileRequest.IsBioPublic,
		Bio:                         upsertProfileRequest.Bio,
		IsPhoneNumberPublic:         upsertProfileRequest.IsPhoneNumberPublic,
		PhoneNumber:                 upsertProfileRequest.PhoneNumber,
		IsAddressPublic:             upsertProfileRequest.IsAddressPublic,
		Address:                     upsertProfileRequest.Address,
		IsAcademicInstitutionPublic: upsertProfileRequest.IsAcademicInstitutionPublic,
		AcademicInstitution:         upsertProfileRequest.AcademicInstitution,
		IsAcademicEmailPublic:       upsertProfileRequest.IsAcademicEmailPublic,
		AcademicEmail:               upsertProfileRequest.AcademicEmail,
	})
	if customerr != nil {
		var customErr *customerror.Err
		if errors.As(customerr, &customErr) {
			return customErr.Extend("failed to update profile by credential ID")
		}
		return customerror.Parse(&customerror.ErrInternalServer, customerr)
	}

	return c.Status(fiber.StatusOK).JSON(UpdateProfileResponse{
		Profile: *profile,
	})
}
