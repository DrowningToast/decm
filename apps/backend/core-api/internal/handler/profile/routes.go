package profile

import (
	"apps/backend/common/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted profile routes")

	profileGroup := r.Group("/profile").Use(h.AuthenticationGuardMiddleware.Middleware)
	profileGroup.Get("/my", h.GetMyProfile)
	profileGroup.Post("", h.CreateProfile)
	profileGroup.Patch("/credential/:credential_id", h.UpdateProfileByCredentialId)
	profileGroup.Post("/password/verify", h.VerifyPassword)
}
