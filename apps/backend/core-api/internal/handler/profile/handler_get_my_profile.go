package profile

import (
	"apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get my profile
// @Description Get my profile
// @ID get-my-profile
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} entity.Profile
// @Failure 400 {object} customerror.Err
// @Failure 404 {object} customerror.Err
// @Failure 500 {object} customerror.Err
// @Router /api/v1/profile/my [get]
func (h *Handler) GetMyProfile(c *fiber.Ctx) error {
	user, err := h.AuthenticationService.RequireUserContext(c)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	profile, err := h.ProfileUc.GetProfileByAuthenticationCredentialId(c.Context(), user.UserId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}
