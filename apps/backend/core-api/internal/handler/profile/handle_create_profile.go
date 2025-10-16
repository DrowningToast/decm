package profile

import (
	"errors"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/profile"

	"github.com/gofiber/fiber/v2"
)

type CreateProfileRequest struct {
	profile.CreateProfileParameters
}

type CreateProfileResponse struct {
	entity.Profile
}

func (r *CreateProfileRequest) Parse(c *fiber.Ctx) error {
	if err := c.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

// @Summary Create a profile
// @Description Create a profile
// @ID create-profile
// @Tags Profile
// @Accept json
// @Produce json
// @Param profile body CreateProfileRequest true "Profile"
// @Success 200 {object} CreateProfileResponse
// @Failure 400 {object} customerror.Err
// @Failure 403 {object} customerror.Err
// @Failure 404 {object} customerror.Err
// @Failure 500 {object} customerror.Err
// @Router /api/v1/profile [post]
func (h *Handler) CreateProfile(c *fiber.Ctx) error {
	var createProfileRequest CreateProfileRequest
	err := createProfileRequest.Parse(c)
	if err != nil {
		return err
	}
	if err := createProfileRequest.CreateProfileParameters.IsValid(); err != nil {
		return err
	}
	isCurrentUser, err := h.AuthenticationService.IsCurrentUser(c, createProfileRequest.CreateProfileParameters.AuthenticationCredentialId)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if !isCurrentUser {
		return customerror.Parse(&customerror.ErrUnauthorized, errors.New("mismatch authentication credential id"))
	}

	profile, err := h.ProfileUc.CreateProfile(c.Context(), createProfileRequest.CreateProfileParameters)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(CreateProfileResponse{
		Profile: *profile,
	})
}
