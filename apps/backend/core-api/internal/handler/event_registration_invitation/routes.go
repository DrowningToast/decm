package event_registration_invitation

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Group routes for event registration invitations
	eventRegistrationInvitationsGroup := router.Group("/event-registration-invitations").
		Use(h.AuthenticationGuardMiddleware.Middleware)
	// Get event registration invitation of user and by event id
	eventRegistrationInvitationsGroup.Get("/my/:event_id/", h.GetEventRegistrationInvitationByUserAndEvent)

	eventRegistrationInvitationsGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Get("/:event_id/", h.GetEventRegistrationInvitationsByEventId)
	// Import event participants
	eventRegistrationInvitationsGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Post("/import/:eventId", h.ImportEventParticipants)
	// Cancel event registration invitation
	eventRegistrationInvitationsGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Delete("/:eventRegistrationInvitationId", h.CancelEventRegistrationInvitation)
}
