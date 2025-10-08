package auth

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	customerror "apps/backend/common/customerror"
	"apps/backend/core-api/config"

	"github.com/gofiber/fiber/v2"
)

type verifyGoogleOAuthRequest struct {
	Code  string `query:"code"`
	State string `query:"state"`
}

// @Summary Verify Google OAuth code
// @Description Verify Google OAuth code
// @Tags Auth
// @ID verify-google-oauth
// @Param code query auth.verifyGoogleOAuthRequest.Code true "Code"
// @Param state query auth.verifyGoogleOAuthRequest.State true "State"
// @Accept json
// @Produce json
// @Success 302
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/auth/verify-google-oauth [get]
func (h Handler) VerifyGoogleOAuth(ctx *fiber.Ctx) error {
	requestQueries := verifyGoogleOAuthRequest{}
	err := ctx.QueryParser(&requestQueries)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	if err := requestQueries.IsValid(); err != nil {
		return err
	}

	// Get fiber session
	session, err := h.GoogleOAuthService.SessionStore.Get(ctx)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	token, err := h.AuthUc.VerifyGoogleOAuthCode(ctx.UserContext(), session, requestQueries.Code, requestQueries.State)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	cfg := config.LoadConfig()

	queries := make(url.Values)
	queries.Set("access_token", token.AccessToken)
	queries.Set("expires_in", strconv.Itoa(int(time.Until(token.Expiry).Seconds())))
	queries.Set("refresh_token", token.RefreshToken)

	return ctx.Redirect(cfg.GoogleOAuth.SuccessURL + "?" + queries.Encode())
}

func (r *verifyGoogleOAuthRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func (r *verifyGoogleOAuthRequest) IsValid() error {
	if r.Code == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("code is required"))
	}
	if r.State == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("state is required"))
	}
	return nil
}
