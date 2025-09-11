package auth

import (
	"errors"
	"time"

	customerror "apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
)

type verifyGoogleOAuthRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type verifyGoogleOAuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary Verify Google OAuth code
// @Description Verify Google OAuth code
// @ID verify-google-oauth
// @Param code body string true "Code"
// @Param state body string true "State"
// @Accept json
// @Produce json
// @Success 200 {object} verifyGoogleOAuthResponse
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/auth/verify-google-oauth [post]
func (h Handler) VerifyGoogleOAuth(ctx *fiber.Ctx) error {
	requestBody := verifyGoogleOAuthRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return *err
	}
	if err := requestBody.IsValid(); err != nil {
		return *err
	}

	// Get fiber session
	session, err := h.GoogleOAuthService.SessionStore.Get(ctx)
	if err != nil {
		return *customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	token, err := h.AuthUc.VerifyGoogleOAuthCode(ctx.UserContext(), session, requestBody.Code, requestBody.State)
	if err != nil {
		return *customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	response := verifyGoogleOAuthResponse{
		AccessToken:  token.AccessToken,
		ExpiresIn:    int(token.Expiry.Sub(time.Now()).Seconds()),
		RefreshToken: token.RefreshToken,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (r *verifyGoogleOAuthRequest) Parse(ctx *fiber.Ctx) *customerror.Err {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func (r *verifyGoogleOAuthRequest) IsValid() *customerror.Err {
	if r.Code == "" {
		return customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, errors.New("code is required"))
	}
	if r.State == "" {
		return customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, errors.New("state is required"))
	}
	return nil
}
