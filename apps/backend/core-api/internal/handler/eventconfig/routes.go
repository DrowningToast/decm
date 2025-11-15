package eventconfig

import (
	"apps/backend/common/log"
	roleguard "apps/backend/core-api/internal/middleware/role_guard"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted event config routes")

	// Create event config group with authentication middleware
	eventConfigGroup := r.Group("/events/:event_id/config").Use(
		h.AuthenticationGuardMiddleware.Middleware,
	)

	// Event Certificate Config routes
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireRole(
		roleguard.RoleVerifiedIssuer,
		roleguard.RoleVerifiedOrganizer,
	)).Get("/certificate", h.GetEventCertificateConfig)

	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Put("/certificate", h.UpdateEventCertificateConfig)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Delete("/certificate", h.DeleteEventCertificateConfig)

	// Participant end
	eventConfigGroup.Post("/password-check", h.CheckEventPassword)

	// Event Registration Config routes
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Post("/registration", h.CreateEventRegistrationConfig)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Get("/registration", h.GetEventRegistrationConfig)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Put("/registration", h.UpdateEventRegistrationConfig)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Delete("/registration", h.DeleteEventRegistrationConfig)
}
