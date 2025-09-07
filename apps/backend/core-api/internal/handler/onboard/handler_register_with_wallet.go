package onboard

import (
	"errors"
	"regexp"
	"time"

	customerror "apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
)

var ethAddressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

type registerWithWalletRequest struct {
	SignedMessage string `json:"signed_message"`
	WalletAddress string `json:"wallet_address"`
}

// @Summary Register a new user with wallet address
// @Description Register a new user with wallet address
// @ID register-with-wallet
// @Param signed_message body string true "Signed message"
// @Param wallet_address body string true "Wallet address"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/onboard/register-with-wallet [post]
func (h Handler) RegisterWithWallet(ctx *fiber.Ctx) error {
	requestBody := registerWithWalletRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return *err
	}
	if err := requestBody.IsValid(); err != nil {
		return *err
	}

	jwt, err := h.OnboardUc.RegisterWithWalletAddress(ctx.UserContext(), requestBody.SignedMessage, requestBody.WalletAddress)
	if err != nil {
		return *err
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = *jwt
	cookie.Expires = time.Now().Add(h.SessionExpiration)
	ctx.Cookie(cookie)

	return ctx.Status(fiber.StatusOK).Send([]byte(""))
}

func (r *registerWithWalletRequest) Parse(ctx *fiber.Ctx) *customerror.Err {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func (r *registerWithWalletRequest) IsValid() *customerror.Err {
	if len(r.SignedMessage) == 0 {
		return customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, errors.New("signed message is required"))
	}
	if len(r.WalletAddress) == 0 {
		return customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, errors.New("wallet address is required"))
	}
	if !ethAddressRegex.MatchString(r.WalletAddress) {
		return customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, errors.New("invalid wallet address"))
	}

	return nil
}
