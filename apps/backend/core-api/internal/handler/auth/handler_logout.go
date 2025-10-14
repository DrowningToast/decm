package auth

import (
	"github.com/gofiber/fiber/v2"
)

// @Summary Logout
// @Tags Auth
// @Description Logout user by clearing session and OAuth cookies
// @ID logout
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Successfully logged out"
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/auth/logout [post]
func (h Handler) Logout(ctx *fiber.Ctx) error {
	// Call auth service to clear cookies
	h.AuthService.Logout(ctx)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
