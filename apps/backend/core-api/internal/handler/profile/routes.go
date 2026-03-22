package profile

import (
	"apps/backend/core-api/internal/handler/metrics"
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted profile routes")

	profileGroup := r.Group("/profile").Use(h.AuthenticationGuardMiddleware.Middleware)
	profileGroup.Get("/viewmodel", metrics.MeasureHandlerDurationWrapper("profile.GetProfileViewModel", h.GetProfileViewModel))
	profileGroup.Post("", metrics.MeasureHandlerDurationWrapper("profile.CreateProfile", h.CreateProfile))
	profileGroup.Patch("/credential/:credential_id", metrics.MeasureHandlerDurationWrapper("profile.UpdateProfileByCredentialId", h.UpdateProfileByCredentialId))
	profileGroup.Post("/password/verify", metrics.MeasureHandlerDurationWrapper("profile.VerifyPassword", h.VerifyPassword))
}
