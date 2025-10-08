package onboard

import (
	"errors"

	customerror "apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
)

type RegistrationMethod string

const (
	RegistrationMethodGoogle RegistrationMethod = "google"
	RegistrationMethodWallet RegistrationMethod = "wallet"
)

type CheckOnboardStatusRequest struct {
	Method RegistrationMethod `json:"method" validate:"required"`

	AccessToken string `json:"access_token"`
	SignMessage string `json:"sign_message"`
}

type CheckOnboardStatusResponse struct {
	IsExists bool `json:"is_exists"`
}

// @Summary Check onboard status
// @Description Check onboard status
// @ID check-onboard-status
// @Tags Onboard
// @Param method body Check true "Method"
// @Param access_token body string false "Access token"
// @Param sign_message body string false "Sign message"
// @Accept json
// @Produce json
// @Success 200 {object} CheckOnboardStatusResponse
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/onboard/check-onboard-status [post]
func (h Handler) CheckOnboardStatus(ctx *fiber.Ctx) error {
	requestBody := CheckOnboardStatusRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	switch requestBody.Method {
	case RegistrationMethodGoogle:
		accessToken := requestBody.AccessToken
		if accessToken == "" {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("access token is required"))
		}
		isExists, err := h.OnboardUc.CheckOnboardStatusWithGoogleConnectorRef(ctx.UserContext(), accessToken)
		if err != nil {
			return err
		}
		response := CheckOnboardStatusResponse{
			IsExists: isExists,
		}
		return ctx.Status(fiber.StatusOK).JSON(response)
	case RegistrationMethodWallet:
		signMessage := requestBody.SignMessage
		if signMessage == "" {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("sign message is required"))
		}
		isExists, err := h.OnboardUc.CheckOnboardStatusWithWalletAddress(ctx.UserContext(), signMessage)
		if err != nil {
			return err
		}
		response := CheckOnboardStatusResponse{
			IsExists: isExists,
		}
		return ctx.Status(fiber.StatusOK).JSON(response)
	default:
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("method is invalid"))
	}
}

func (r *CheckOnboardStatusRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func (r *CheckOnboardStatusRequest) IsValid() error {
	if r.Method == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("method is required"))
	}
	return nil
}
