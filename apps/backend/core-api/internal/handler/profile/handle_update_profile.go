package profile

import (
	"errors"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/profile"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	profile.UpdateProfileParameters
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
// @Param IsProfilePicturePublic body profile.UpdateProfileRequest.IsProfilePicturePublic true "Is Profile Picture Public"
// @Param ProfilePictureUrl body profile.UpdateProfileRequest.ProfilePictureUrl true "Profile Picture URL"
// @Param IsFirstNamePublic body profile.UpdateProfileRequest.IsFirstNamePublic true "Is First Name Public"
// @Param FirstName body profile.UpdateProfileRequest.FirstName true "First Name"
// @Param IsLastNamePublic body profile.UpdateProfileRequest.IsLastNamePublic true "Is Last Name Public"
// @Param LastName body profile.UpdateProfileRequest.LastName true "Last Name"
// @Param IsEmailPublic body profile.UpdateProfileRequest.IsEmailPublic true "Is Email Public"
// @Param Email body profile.UpdateProfileRequest.Email true "Email"
// @Param IsPhoneNumberPublic body profile.UpdateProfileRequest.IsPhoneNumberPublic true "Is Phone Number Public"
// @Param PhoneNumber body profile.UpdateProfileRequest.PhoneNumber true "Phone Number"
// @Param IsAddressPublic body profile.UpdateProfileRequest.IsAddressPublic true "Is Address Public"
// @Param Address body profile.UpdateProfileRequest.Address true "Address"
// @Param IsAcademicInstitutionPublic body profile.UpdateProfileRequest.IsAcademicInstitutionPublic true "Is Academic Institution Public"
// @Param AcademicInstitution body profile.UpdateProfileRequest.AcademicInstitution true "Academic Institution"
// @Param IsAcademicEmailPublic body profile.UpdateProfileRequest.IsAcademicEmailPublic true "Is Academic Email Public"
// @Param AcademicEmail body profile.UpdateProfileRequest.AcademicEmail true "Academic Email"
// @Success 200 {object} UpdateProfileResponse
// @Failure 400 {object} customerror.Err
// @Failure 403 {object} customerror.Err
// @Failure 404 {object} customerror.Err
// @Failure 500 {object} customerror.Err
// @Router /api/v1/profile/credential/{credential_id} [patch]
func (h *Handler) UpdateProfileByCredentialId(c *fiber.Ctx) error {
	var upsertProfileRequest UpdateProfileRequest
	upsertProfileRequest.Parse(c)
	if err := upsertProfileRequest.IsValid(); err != nil {
		return err
	}
	credentialId := c.Params("credential_id")
	credentialIdUUID, err := uuid.Parse(credentialId)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	// Auth check
	if !auth.IsCurrentUser(c, credentialIdUUID) {
		return customerror.Parse(&customerror.ErrUnauthorized, errors.New("unauthorized"))
	}

	profile, customerr := h.ProfileUc.UpdateProfileByCredentialId(c.Context(), credentialIdUUID, upsertProfileRequest.UpdateProfileParameters)
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
