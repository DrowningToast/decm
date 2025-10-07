package auth

import (
	customerror "apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
)

type requestGoogleOAuthResponse struct {
	URL string `json:"url" example:"https://accounts.google.com/o/oauth2/auth?client_id=YOUR_CLIENT_ID&redirect_uri=YOUR_REDIRECT_URI&response_type=code&scope=email profile"`
}

// @Summary Request Google OAuth
// @Tags Auth
// @Description Request Google OAuth
// @ID request-google-oauth
// @Accept json
// @Produce json
// @Success 200 {object} requestGoogleOAuthResponse
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

	return ctx.Status(fiber.StatusOK).JSON(requestGoogleOAuthResponse{URL: *url})
}
