package onboard

import (
	"apps/backend/common/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted core api routes")

	onboardGroup := r.Group("/onboard")
	onboardGroup.Get("/sign-message", h.GetRegisterSignMessage)
	onboardGroup.Post("/register-with-wallet", h.RegisterWithWallet)
	onboardGroup.Post("/check-onboard-status", h.CheckOnboardStatus)
}
