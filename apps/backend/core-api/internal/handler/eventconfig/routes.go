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

	// Event Registration Config routes
	eventConfigGroup.Post("/registration", h.CreateEventRegistrationConfig)
	eventConfigGroup.Get("/registration", h.GetEventRegistrationConfig)
	eventConfigGroup.Put("/registration", h.UpdateEventRegistrationConfig)
	eventConfigGroup.Delete("/registration", h.DeleteEventRegistrationConfig)

	// Event Certificate Config routes
	eventConfigGroup.Post("/certificate", h.CreateEventCertificateConfig)
	eventConfigGroup.Get("/certificate", h.GetEventCertificateConfig)
	eventConfigGroup.Put("/certificate", h.UpdateEventCertificateConfig)
	eventConfigGroup.Delete("/certificate", h.DeleteEventCertificateConfig)
}
