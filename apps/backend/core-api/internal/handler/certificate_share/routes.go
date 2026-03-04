package certificate_share

import (
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	defer log.Logger.Info("Mounted certificate share routes")

	// Authenticated routes — creating share links requires certificate ownership
	certGroup := r.Group("/certificates").Use(h.AuthenticationGuardMiddleware.Middleware)
	certGroup.Post("/:certificate_id/shares", h.CreateCertificateShare)

	// Public routes — viewing shared certificates requires no authentication
	shareGroup := r.Group("/certificate-shares")
	shareGroup.Get("/:handle", h.GetCertificateShareStatus)
	shareGroup.Get("/:handle/data", h.GetCertificateShareData)
	shareGroup.Post("/:handle/data/unlock", h.GetCertificateShareDataWithPassword)
}
