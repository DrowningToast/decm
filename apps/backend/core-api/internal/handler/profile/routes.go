package profile

import (
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted profile routes")

	profileGroup := r.Group("/profile").Use(h.AuthenticationGuardMiddleware.Middleware)
	profileGroup.Get("/viewmodel", h.GetProfileViewModel)
	profileGroup.Post("", h.CreateProfile)
	profileGroup.Patch("/credential/:credential_id", h.UpdateProfileByCredentialId)
	profileGroup.Post("/password/verify", h.VerifyPassword)
}
