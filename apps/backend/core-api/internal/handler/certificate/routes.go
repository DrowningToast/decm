package certificate

import (
	"apps/backend/core-api/internal/handler/metrics"
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	defer log.Logger.Info("Mounted certificate routes")

	certificateGroup := r.Group("/certificates").Use(
		h.AuthenticationGuardMiddleware.Middleware,
	)
	certificateGroup.Get("/my-list-viewmodel", metrics.MeasureHandlerDurationWrapper("certificate.GetMyCertificatesListViewModel", h.GetMyCertificatesListViewModel))
	certificateGroup.Get("/:certificate_id/image", metrics.MeasureHandlerDurationWrapper("certificate.GenerateCertificateImage", h.GenerateCertificateImage))
	certificateGroup.Get("/claim/:certificate_id/sign-message", metrics.MeasureHandlerDurationWrapper("certificate.GetClaimCertificateSignMessage", h.GetClaimCertificateSignMessage))
	certificateGroup.Post("/claim/:certificate_id", metrics.MeasureHandlerDurationWrapper("certificate.ClaimCertificate", h.ClaimCertificate))
}
