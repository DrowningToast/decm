package auth

import (
	"apps/backend/common/log"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted auth routes")

	authGroup := r.Group("/auth")
	authGroup.Get("/request-google-oauth", h.RequestGoogleOAuth)
	authGroup.Get("/verify-google-oauth", h.VerifyGoogleOAuth)
}
