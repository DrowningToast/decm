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

func (r *UpdateProfileRequest) Parse(c *fiber.Ctx) *customerror.Err {
	if err := c.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

// @Summary Update a profile by credential ID
// @Description Update a profile by credential ID
// @Tags Profile
// @Accept json
// @Produce json
// @Param credential_id path string true "Credential ID"
// @Param profile body UpdateProfileRequest true "Profile"
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
		return customerr.Extend("failed to update profile by credential ID")
	}

	return c.Status(fiber.StatusOK).JSON(UpdateProfileResponse{
		Profile: *profile,
	})
}
