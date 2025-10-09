package auth

import (
	customerror "apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
)

// @Summary Request Google OAuth
// @Tags Auth
// @Description Request Google OAuth
// @ID request-google-oauth
// @Accept json
// @Produce json
// @Success 302
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/auth/request-google-oauth [get]
func (h Handler) RequestGoogleOAuth(ctx *fiber.Ctx) error {
	session, err := h.GoogleOAuthService.SessionStore.Get(ctx)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	url, err := h.GoogleOAuthService.Login(session)
	if err != nil {
		return err
	}

	return ctx.Redirect(*url)
}
