package auth

import (
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted auth routes")

	authGroup := r.Group("/auth")
	authGroup.Get("/request-google-oauth", h.RequestGoogleOAuth)
	authGroup.Get("/verify-google-oauth", h.VerifyGoogleOAuth)
	authGroup.Post("/logout", h.Logout)
	authGroup.Use(h.AuthenticationGuardMiddleware.Middleware).Get("/check-role", h.CheckRole)
}
