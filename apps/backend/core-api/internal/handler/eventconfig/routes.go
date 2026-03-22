package eventconfig

import (
	"apps/backend/core-api/internal/handler/metrics"
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted event config routes")

	// Public routes for certificate font families (no authentication required)
	publicGroup := r.Group("/eventconfig")
	publicGroup.Get("/certificate-font-families", metrics.MeasureHandlerDurationWrapper("eventConfig.GetEventCertificateFontFamilies", h.GetEventCertificateFontFamilies))

	// Create event config group with authentication middleware
	eventConfigGroup := r.Group("/events/:event_id/config").Use(
		h.AuthenticationGuardMiddleware.Middleware,
	)

	// Participant end
	eventConfigGroup.Post("/password-check", metrics.MeasureHandlerDurationWrapper("eventConfig.CheckEventPassword", h.CheckEventPassword))

	// Event Certificate Config routes
	// GET: Handler checks if user is organizer or event-specific issuer (authorization in handler)
	// IMPORTANT: Must be registered BEFORE any .Use() calls to avoid inheriting role guard middleware
	eventConfigGroup.Get("/certificate", metrics.MeasureHandlerDurationWrapper("eventConfig.GetEventCertificateConfig", h.GetEventCertificateConfig))
	eventConfigGroup.Get("/certificate/template", metrics.MeasureHandlerDurationWrapper("eventConfig.GetEventCertificateTemplate", h.GetEventCertificateTemplate))
	eventConfigGroup.Get("/certificate/mint-readiness", metrics.MeasureHandlerDurationWrapper("eventConfig.CheckCertificateMintReadiness", h.CheckCertificateMintReadiness))

	// Event Registration Config routes (with role guard)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Post("/registration", metrics.MeasureHandlerDurationWrapper("eventConfig.CreateEventRegistrationConfig", h.CreateEventRegistrationConfig))
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Get("/registration", metrics.MeasureHandlerDurationWrapper("eventConfig.GetEventRegistrationConfig", h.GetEventRegistrationConfig))
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Put("/registration", metrics.MeasureHandlerDurationWrapper("eventConfig.UpdateEventRegistrationConfig", h.UpdateEventRegistrationConfig))
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Delete("/registration", metrics.MeasureHandlerDurationWrapper("eventConfig.DeleteEventRegistrationConfig", h.DeleteEventRegistrationConfig))

	// Certificate Config modification routes (with role guard)
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Put("/certificate", metrics.MeasureHandlerDurationWrapper("eventConfig.UpdateEventCertificateConfig", h.UpdateEventCertificateConfig))
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Patch("/certificate/published", metrics.MeasureHandlerDurationWrapper("eventConfig.ToggleCertificatePublished", h.ToggleCertificatePublished))
	eventConfigGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Delete("/certificate", metrics.MeasureHandlerDurationWrapper("eventConfig.DeleteEventCertificateConfig", h.DeleteEventCertificateConfig))
}
