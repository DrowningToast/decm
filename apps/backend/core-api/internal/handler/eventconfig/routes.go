package eventconfig

import (
	"apps/backend/common/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted event config routes")

	// Create event config group with authentication middleware
	eventConfigGroup := r.Group("/events/:event_id/config").Use(
		h.AuthenticationGuardMiddleware.Middleware,
	)

	// Participant end
	eventConfigGroup.Post("/password-check", h.CheckEventPassword)

	// Event Certificate Config routes
	// GET: Handler checks if user is organizer or event-specific issuer (authorization in handler)
	// IMPORTANT: Must be registered BEFORE any .Use() calls to avoid inheriting role guard middleware
	eventConfigGroup.Get("/certificate", h.GetEventCertificateConfig)
	eventConfigGroup.Get("/certificate/mint-readiness", h.CheckCertificateMintReadiness)

	// Event Registration Config routes (with role guard)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Post("/registration", h.CreateEventRegistrationConfig)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Get("/registration", h.GetEventRegistrationConfig)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Put("/registration", h.UpdateEventRegistrationConfig)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Delete("/registration", h.DeleteEventRegistrationConfig)

	// Certificate Config modification routes (with role guard)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Put("/certificate", h.UpdateEventCertificateConfig)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Patch("/certificate/published", h.ToggleCertificatePublished)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Delete("/certificate", h.DeleteEventCertificateConfig)
}
