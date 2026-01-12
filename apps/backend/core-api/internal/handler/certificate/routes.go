package certificate

import (
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	defer log.Logger.Info("Mounted certificate routes")

	certificateGroup := r.Group("/certificates").Use(
		h.AuthenticationGuardMiddleware.Middleware,
	)
	certificateGroup.Get("/my-list-viewmodel", h.GetMyCertificatesListViewModel)
	certificateGroup.Get("/:certificate_id/image", h.GenerateCertificateImage)
	certificateGroup.Get("/claim/:certificate_id/sign-message", h.GetClaimCertificateSignMessage)
	certificateGroup.Post("/claim/:certificate_id", h.ClaimCertificate)
}
