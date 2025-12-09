package event_registration

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Group routes for event registration invitations
	eventRegistrationGroup := router.Group("/event-registration").
		Use(h.AuthenticationGuardMiddleware.Middleware)
	// Get event registration invitation of user and by event id
	eventRegistrationGroup.Get("/my/:event_id/", h.GetEventRegistrationInvitationByUserAndEvent)

	// Get join event sign message
	eventRegistrationGroup.Get("/join/:event_id/sign-message", h.GetJoinEventSignMessage)
	// Join event
	eventRegistrationGroup.Post("/join/:event_id", h.JoinEvent)

	// Fuck join event
	eventRegistrationGroup.Post("/fuck-join/:event_id", h.FuckJoinEvent)

	// Get event registration invitations by event id
	eventRegistrationGroup.
		// Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).
		Get("/invitation/:event_id", h.GetEventRegistrationInvitationsByEventId)
	// Import event participants
	eventRegistrationGroup.
		// Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).
		Post("/invitation/:event_id/import", h.ImportEventParticipants)
	// Cancel event registration invitation
	eventRegistrationGroup.
		// Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).
		Delete("/invitation", h.CancelEventRegistrationInvitation)
}
