package eventconfig

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterEventConfigRoutes(
	app *fiber.App,
) {
	// Create event config group with authentication middleware
	eventConfigGroup := app.Group("/api/v1/events/:event_id/config").Use(
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
