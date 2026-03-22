package event_registration

import (
	"apps/backend/core-api/internal/handler/metrics"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Group routes for event registration invitations
	eventRegistrationGroup := router.Group("/event-registration").
		Use(h.AuthenticationGuardMiddleware.Middleware)
	// Get event registration invitation of user and by event id
	eventRegistrationGroup.Get("/my/:event_id/", metrics.MeasureHandlerDurationWrapper("eventRegistration.GetEventRegistrationInvitationByUserAndEvent", h.GetEventRegistrationInvitationByUserAndEvent))

	// Get join event sign message
	eventRegistrationGroup.Get("/join/:event_id/sign-message", metrics.MeasureHandlerDurationWrapper("eventRegistration.GetJoinEventSignMessage", h.GetJoinEventSignMessage))
	// Join event
	eventRegistrationGroup.Post("/join/:event_id", metrics.MeasureHandlerDurationWrapper("eventRegistration.JoinEvent", h.JoinEvent))

	// Get event registration invitations by event id
	eventRegistrationGroup.
		Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).
		Get("/invitation/:event_id", metrics.MeasureHandlerDurationWrapper("eventRegistration.GetEventRegistrationInvitationsByEventId", h.GetEventRegistrationInvitationsByEventId))
	// Import event participants
	eventRegistrationGroup.
		Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).
		Post("/invitation/:event_id/import", metrics.MeasureHandlerDurationWrapper("eventRegistration.ImportEventParticipants", h.ImportEventParticipants))
	// Cancel event registration invitation
	eventRegistrationGroup.
		Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).
		Delete("/invitation", metrics.MeasureHandlerDurationWrapper("eventRegistration.CancelEventRegistrationInvitation", h.CancelEventRegistrationInvitation))
}
