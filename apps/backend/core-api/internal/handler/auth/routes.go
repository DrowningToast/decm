package auth

import (
	"apps/backend/core-api/internal/handler/metrics"
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted auth routes")

	authGroup := r.Group("/auth")
	authGroup.Get("/request-google-oauth", metrics.MeasureHandlerDurationWrapper("auth.RequestGoogleOAuth", h.RequestGoogleOAuth))
	authGroup.Get("/verify-google-oauth", metrics.MeasureHandlerDurationWrapper("auth.VerifyGoogleOAuth", h.VerifyGoogleOAuth))
	authGroup.Post("/logout", metrics.MeasureHandlerDurationWrapper("auth.Logout", h.Logout))
	authGroup.Use(h.AuthenticationGuardMiddleware.Middleware).Get("/check-role", metrics.MeasureHandlerDurationWrapper("auth.CheckRole", h.CheckRole))
}
