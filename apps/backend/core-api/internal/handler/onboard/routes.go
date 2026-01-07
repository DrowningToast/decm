package onboard

import (
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted core api routes")

	onboardGroup := r.Group("/onboard")
	onboardGroup.Get("/sign-message", h.GetSignMessage)
	onboardGroup.Post("/register-with-wallet", h.RegisterWithWallet)
	onboardGroup.Post("/register-with-google-oauth", h.RegisterWithGoogleOAuth)
	onboardGroup.Post("/check-onboard-status", h.VerifyJwtMiddleware.Middleware, h.CheckOnboardStatus)
}
