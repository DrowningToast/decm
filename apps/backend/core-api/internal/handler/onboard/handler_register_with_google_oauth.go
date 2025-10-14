package onboard

import (
	customerror "apps/backend/common/customerror"
	"apps/backend/services/oauth"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type registerWithGoogleOAuthRequest struct {
	AccessToken string `json:"access_token"`
	Password    string `json:"password" validate:"required,min=6"`
}

// @Summary Register a new user with Google OAuth
// @Description Register a new user with Google OAuth
// @ID register-with-google-oauth
// @Tags Onboard
// @Param access_token body onboard.registerWithGoogleOAuthRequest.AccessToken true "Access token"
// @Param password body onboard.registerWithGoogleOAuthRequest.Password true "Password"
// @Accept json
// @Produce json
// @Success 200 {object} registerResponse
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/onboard/register-with-google-oauth [post]
func (h Handler) RegisterWithGoogleOAuth(ctx *fiber.Ctx) error {
	requestBody := registerWithGoogleOAuthRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	// Parse token
	token, err := oauth.ParseToken(requestBody.AccessToken)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	credentialId, jwt, err := h.OnboardUc.RegisterWithGoogle(ctx.UserContext(), token, requestBody.Password)
	if err != nil {
		return err
	}

	h.AuthService.SetJwtCookie(ctx, *jwt)

	response := registerResponse{
		CredentialId: credentialId.String(),
		Jwt:          *jwt,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (r *registerWithGoogleOAuthRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	return nil
}

func (r *registerWithGoogleOAuthRequest) IsValid() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	return nil
}
