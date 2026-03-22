package onboard

import (
	"apps/backend/core-api/internal/handler/metrics"
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted core api routes")

	onboardGroup := r.Group("/onboard")
	onboardGroup.Get("/sign-message", metrics.MeasureHandlerDurationWrapper("onboard.GetSignMessage", h.GetSignMessage))
	onboardGroup.Post("/register-with-wallet", metrics.MeasureHandlerDurationWrapper("onboard.RegisterWithWallet", h.RegisterWithWallet))
	onboardGroup.Post("/register-with-google-oauth", metrics.MeasureHandlerDurationWrapper("onboard.RegisterWithGoogleOAuth", h.RegisterWithGoogleOAuth))
	onboardGroup.Post("/check-onboard-status", h.VerifyJwtMiddleware.Middleware, metrics.MeasureHandlerDurationWrapper("onboard.CheckOnboardStatus", h.CheckOnboardStatus))
}
