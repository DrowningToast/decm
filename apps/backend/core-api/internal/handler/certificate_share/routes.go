package certificate_share_handler

import (
	"apps/backend/services/log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func (h *Handler) Mount(r fiber.Router) {
	defer log.Logger.Info("Mounted certificate share routes")

	shareGroup := r.Group("/certificate-shares")

	// Authenticated routes — auth guard applied per route
	configGroup := shareGroup.Group("/config")
	configGroup.Post("/:certificate_id", h.AuthenticationGuardMiddleware.Middleware, h.CreateCertificateShare)
	configGroup.Patch("/:share_id", h.AuthenticationGuardMiddleware.Middleware, h.UpdateCertificateShare)

	// Rate limiter for public verification routes: 20 requests per IP per minute
	verifyLimiter := limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Please wait before trying again.",
			})
		},
	})

	// Public routes
	shareGroup.Post("/:handle", verifyLimiter, h.GetCertificateShareData)
	shareGroup.Get("/:handle/image", verifyLimiter, h.GetCertificateShareImage)
}
