package onboard

import (
	"apps/backend/common/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted core api routes")

	r.Group("/onboard")
	{
		r.Get("/sign-message", h.GetRegisterSignMessage)
		r.Post("/register-with-wallet", h.RegisterWithWallet)
	}
}
