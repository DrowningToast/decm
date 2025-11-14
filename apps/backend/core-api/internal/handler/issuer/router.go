package issuer

import (
	"apps/backend/common/log"
	roleguard "apps/backend/core-api/internal/middleware/role_guard"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted issuer routes")

	issuerGroup := r.Group("/issuers").Use(h.AuthenticationGuardMiddleware.Middleware)
	issuerGroup.Use(h.RoleGuardMiddleware.RequireRole(roleguard.RoleVerifiedIssuer, roleguard.RoleVerifiedOrganizer)).Get("/", h.GetVerifiedIssuers)
	issuerGroup.Get("/events", h.GetIssuerEvents)
}
