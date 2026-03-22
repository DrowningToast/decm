package event

import (
	"apps/backend/core-api/internal/handler/metrics"
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted event routes")

	eventGroup := r.Group("/events").Use(
		h.AuthenticationGuardMiddleware.Middleware,
	)

	// PARTICIPANT
	eventGroup.Get("/", metrics.MeasureHandlerDurationWrapper("event.GetEventsList", h.GetEventsList))
	eventGroup.Get("/:event_id", metrics.MeasureHandlerDurationWrapper("event.GetEventById", h.GetEventById))
	eventGroup.Get("/:event_id/viewmodel", metrics.MeasureHandlerDurationWrapper("event.GetEventViewModel", h.GetEventViewModel))

	// HOST
	eventGroup.Post("/", metrics.MeasureHandlerDurationWrapper("event.CreateEvent", h.CreateEvent))

	eventGroup.Post("/:event_id/issuers", metrics.MeasureHandlerDurationWrapper("event.CreateEventIssuer", h.CreateEventIssuer))
	eventGroup.Get("/:event_id/participants", metrics.MeasureHandlerDurationWrapper("event.GetEventParticipants", h.GetEventParticipants))
	eventGroup.Post("/:event_id/certificates/import", metrics.MeasureHandlerDurationWrapper("event.ImportCertificateReceivers", h.ImportCertificateReceivers))
	eventGroup.Post("/:event_id/certificates/publish", metrics.MeasureHandlerDurationWrapper("event.PublishEventCertificates", h.PublishEventCertificates))
	eventGroup.Post("/:event_id/certificates/revoke", metrics.MeasureHandlerDurationWrapper("event.RevokeEventCertificates", h.RevokeEventCertificates))
	eventGroup.Post("/:event_id/certificates/revoke-all", metrics.MeasureHandlerDurationWrapper("event.RevokeAllEventCertificates", h.RevokeAllEventCertificates))
	eventGroup.Post("/:event_id/certificates/sign", metrics.MeasureHandlerDurationWrapper("event.SignEventCertificates", h.SignEventCertificates))
	eventGroup.Get("/:event_id/certificates", metrics.MeasureHandlerDurationWrapper("event.GetEventCertificates", h.GetEventCertificates))
	eventGroup.Get("/:event_id/certificates/list-viewmodel", metrics.MeasureHandlerDurationWrapper("event.GetCertificatesListViewModel", h.GetCertificatesListViewModel))
	eventGroup.Put("/:event_id/certificates/text-config", metrics.MeasureHandlerDurationWrapper("event.UpdateEventCertificateTextConfig", h.UpdateEventCertificateTextConfig))

	eventGroup.Get("/:event_id/contracts", metrics.MeasureHandlerDurationWrapper("event.GetEventContractByEventID", h.GetEventContractByEventID))
	eventGroup.Get("/:event_id/issuers", metrics.MeasureHandlerDurationWrapper("event.GetEventIssuersByEventID", h.GetEventIssuersByEventID))

	eventGroup.Get("/owner-credentials/:owner_credential_id", metrics.MeasureHandlerDurationWrapper("event.GetEventsByOwnerCredentialsId", h.GetEventsByOwnerCredentialsId))

	eventGroup.Put("/:event_id", metrics.MeasureHandlerDurationWrapper("event.UpdateEvent", h.UpdateEvent))
	eventGroup.Put("/:event_id/contracts", metrics.MeasureHandlerDurationWrapper("event.UpdateEventContract", h.UpdateEventContract))
	eventGroup.Put("/:event_id/issuers", metrics.MeasureHandlerDurationWrapper("event.UpdateEventIssuer", h.UpdateEventIssuer))

	eventGroup.Delete("/:event_id/contracts", metrics.MeasureHandlerDurationWrapper("event.DeleteEventContract", h.DeleteEventContract))
	eventGroup.Delete("/:event_id", metrics.MeasureHandlerDurationWrapper("event.DeleteEvent", h.DeleteEvent))
	eventGroup.Delete("/:event_id/issuers/:issuer_id", metrics.MeasureHandlerDurationWrapper("event.DeleteEventIssuer", h.DeleteEventIssuer))
}
