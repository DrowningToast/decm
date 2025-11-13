package event_registration_invitation

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Group routes for event registration invitations
	eventRegistrationInvitationsGroup := router.Group("/event-registration-invitations").
		Use(h.AuthenticationGuardMiddleware.Middleware)
	// Import event participants
	eventRegistrationInvitationsGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Post("/import/:eventId", h.ImportEventParticipants)
	// Cancel event registration invitation
	eventRegistrationInvitationsGroup.Use(h.RoleGuardMiddleware.RequireVerifiedOrganizer()).Delete("/:eventRegistrationInvitationId", h.CancelEventRegistrationInvitation)
}
