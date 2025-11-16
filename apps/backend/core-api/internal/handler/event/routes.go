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

	// PARTICIPANT
	eventGroup.Get("/", h.GetEventsList)
	eventGroup.Get("/:event_id", h.GetEventById)
	eventGroup.Get("/:event_id/viewmodel", h.GetEventViewModel)

	// HOST
	eventGroup.Post("/", h.CreateEvent)
	eventGroup.Post("/:event_id/issuers", h.CreateEventIssuer)
	eventGroup.Post("/:event_id/certificates/import", h.ImportCertificateReceivers)
	eventGroup.Post("/:event_id/certificates/revoke", h.RevokeEventCertificates)
	eventGroup.Post("/:event_id/certificates/revoke-all", h.RevokeAllEventCertificates)
	eventGroup.Post("/:event_id/certificates/sign", h.SignEventCertificates)
	eventGroup.Get("/:event_id/certificates", h.GetEventCertificates)

	eventGroup.Get("/:event_id/contracts", h.GetEventContractByEventID)
	eventGroup.Get("/:event_id/issuers", h.GetEventIssuersByEventID)

	eventGroup.Get("/owner-credentials/:owner_credential_id", h.GetEventsByOwnerCredentialsId)

	eventGroup.Put("/:event_id", h.UpdateEvent)
	eventGroup.Put("/:event_id/contracts", h.UpdateEventContract)
	eventGroup.Put("/:event_id/issuers", h.UpdateEventIssuer)

	eventGroup.Delete("/:event_id/contracts", h.DeleteEventContract)
	eventGroup.Delete("/:event_id", h.DeleteEvent)
	eventGroup.Delete("/:event_id/issuers/:issuer_id", h.DeleteEventIssuer)
}
