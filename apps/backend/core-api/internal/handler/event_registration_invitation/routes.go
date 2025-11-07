package event_registration_invitation

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Group routes under /api/v1/events/{eventId}/participants
	eventParticipantsGroup := router.Group("/events/:eventId/participants")

	// Import event participants
	eventParticipantsGroup.Post("/import", h.ImportEventParticipants)

	// Group routes for event registration invitations
	eventRegistrationInvitationsGroup := router.Group("/event-registration-invitations")

	// Cancel event registration invitation
	eventRegistrationInvitationsGroup.Delete("/:eventRegistrationInvitationId", h.CancelEventRegistrationInvitation)
}
