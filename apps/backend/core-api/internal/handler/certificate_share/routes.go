package certificate_share

import (
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	defer log.Logger.Info("Mounted certificate share routes")

	shareGroup := r.Group("/certificate-shares")

	// Authenticated routes — auth guard applied per route
	shareGroup.Post("/:certificate_id", h.AuthenticationGuardMiddleware.Middleware, h.CreateCertificateShare)
	shareGroup.Patch("/:share_id", h.AuthenticationGuardMiddleware.Middleware, h.UpdateCertificateShare)

	// Public routes
	shareGroup.Get("/:handle", h.GetCertificateShareStatus)
	shareGroup.Get("/:handle/data", h.GetCertificateShareData)
	shareGroup.Post("/:handle/data/unlock", h.GetCertificateShareDataWithPassword)
}
