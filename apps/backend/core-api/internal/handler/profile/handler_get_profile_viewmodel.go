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
// @Success 200 {object} profile.GetMyProfileViewModel
// @Failure 400 {object} customerror.Err
// @Failure 404 {object} customerror.Err
// @Failure 500 {object} customerror.Err
// @Router /api/v1/profile/viewmodel [get]
func (h *Handler) GetProfileViewModel(c *fiber.Ctx) error {
	user, err := h.AuthenticationService.RequireUserContext(c)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	viewmodel, err := h.ProfileUc.GetMyProfileViewModel(c.Context(), user)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return c.Status(fiber.StatusOK).JSON(viewmodel)
}
