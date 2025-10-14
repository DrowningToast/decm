package onboard

import (
	"errors"

	customerror "apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
)

type registerWithWalletRequest struct {
	SignedMessage string `json:"signed_message"`
}

// @Summary Register a new user with wallet address
// @Description Register a new user with wallet address
// @ID register-with-wallet
// @Tags Onboard
// @Param signed_message body onboard.registerWithWalletRequest.SignedMessage true "Signed message"
// @Accept json
// @Produce json
// @Success 200 {object} registerResponse
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/onboard/register-with-wallet [post]
func (h Handler) RegisterWithWallet(ctx *fiber.Ctx) error {
	requestBody := registerWithWalletRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	credentialId, jwt, err := h.OnboardUc.RegisterWithWalletAddress(ctx.UserContext(), requestBody.SignedMessage)
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

func (r *registerWithWalletRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func (r *registerWithWalletRequest) IsValid() error {
	if len(r.SignedMessage) == 0 {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("signed message is required"))
	}

	return nil
}
