package onboard

import (
	"time"

	customerror "apps/backend/common/customerror"
	"apps/backend/services/oauth"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type registerWithGoogleOAuthRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Password     string `json:"password" validate:"required,min=6"`
}

type registerWithGoogleOAuthResponse struct {
	Mnemonic []string `json:"mnemonic"`
}

// @Summary Register a new user with Google OAuth
// @Description Register a new user with Google OAuth
// @ID register-with-google-oauth
// @Tags Onboard
// @Param access_token body string true "Access token"
// @Param refresh_token body string true "Refresh token"
// @Param password body string true "Password"
// @Accept json
// @Produce json
// @Success 200 {object} registerWithGoogleOAuthResponse
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/onboard/register-with-google-oauth [post]
func (h Handler) RegisterWithGoogleOAuth(ctx *fiber.Ctx) error {
	requestBody := registerWithGoogleOAuthRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return *err
	}
	if err := requestBody.IsValid(); err != nil {
		return *err
	}

	// Parse token
	token, err := oauth.ParseToken(requestBody.AccessToken, requestBody.RefreshToken)
	if err != nil {
		return *customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	jwt, mnemonic, err := h.OnboardUc.RegisterWithGoogle(ctx.UserContext(), token, requestBody.Password)
	if err != nil {
		return *err
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = *jwt
	cookie.Expires = time.Now().Add(h.SessionExpiration)
	ctx.Cookie(cookie)

	response := registerWithGoogleOAuthResponse{
		Mnemonic: mnemonic,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (r *registerWithGoogleOAuthRequest) Parse(ctx *fiber.Ctx) *customerror.Err {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	return nil
}

func (r *registerWithGoogleOAuthRequest) IsValid() *customerror.Err {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	return nil
}
