package issuer

import (
	roleguard "apps/backend/core-api/internal/middleware/role_guard"
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted issuer routes")

	issuerGroup := r.Group("/issuers").Use(h.AuthenticationGuardMiddleware.Middleware)
	issuerGroup.Use(h.RoleGuardMiddleware.RequireRole(roleguard.RoleVerifiedIssuer, roleguard.RoleVerifiedOrganizer)).Get("/", h.GetVerifiedIssuers)
	issuerGroup.Get("/events", h.GetIssuerEvents)
	issuerGroup.Get("/events/viewmodel", h.GetIssuerEventsViewModel)
}
