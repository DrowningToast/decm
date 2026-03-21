package certificate_share_handler

import (
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	defer log.Logger.Info("Mounted certificate share routes")

	shareGroup := r.Group("/certificate-shares")

	// Authenticated routes — auth guard applied per route
	configGroup := shareGroup.Group("/config")
	configGroup.Post("/:certificate_id", h.AuthenticationGuardMiddleware.Middleware, h.CreateCertificateShare)
	configGroup.Patch("/:share_id", h.AuthenticationGuardMiddleware.Middleware, h.UpdateCertificateShare)

	// Public routes
	shareGroup.Post("/:handle", h.GetCertificateShareData)
	shareGroup.Get("/:handle/image", h.GetCertificateShareImage)
}
