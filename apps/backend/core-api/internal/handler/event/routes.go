package event

import (
	"apps/backend/common/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted event routes")

	eventGroup := r.Group("/events").Use(
		h.AuthenticationGuardMiddleware.Middleware,
	)

	eventGroup.Post("/", h.CreateEvent)
	eventGroup.Post("/:event_id/contracts", h.CreateEventContract)
	eventGroup.Post("/:event_id/issuers", h.CreateEventIssuer)

	eventGroup.Get("/:event_id/contracts", h.GetEventContractByEventID)
	eventGroup.Get("/:event_id/issuers", h.GetEventIssuersByEventID)
	eventGroup.Get("/:event_id/issuers/:issuer_id", h.GetEventIssuerByID)

	eventGroup.Put("/:event_id/contracts", h.UpdateEventContract)
	eventGroup.Put("/:event_id/issuers", h.UpdateEventIssuer)

	// eventGroup.Delete("/:event_id/contracts", h.DeleteEventContract)
	eventGroup.Delete("/:event_id/issuers", h.DeleteEventIssuer)
	// eventGroup.Delete("/:event_id/issuers/:issuer_id", h.DeleteEventContract)
}
