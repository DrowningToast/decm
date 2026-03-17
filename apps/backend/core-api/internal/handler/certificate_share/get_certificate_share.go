package certificate_share

import (
	customerror "apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/certificate_share"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

type CertificateShareStatusResponse struct {
	Status      certificate_share.CertificateShareViewStatus `json:"status"`
	Certificate *entity.EventCertificate                     `json:"certificate,omitempty"`
}

// @Summary Get certificate share status
// @Description Retrieve a shared certificate by its handle. Returns PASSWORD_LOCKED if password-protected, VALID_BUT_PENDING if the certificate has not been claimed yet, or READY with certificate data.
// @ID get-certificate-share-status
// @Tags CertificateShares
// @Produce json
// @Param handle path string true "Share handle"
// @Success 200 {object} CertificateShareStatusResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificate-shares/{handle} [get]
func (h *Handler) GetCertificateShareStatus(ctx *fiber.Ctx) error {
	handle := ctx.Params("handle")
	if handle == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("handle is required"))
	}

	cert, status, err := h.CertificateShareUc.GetCertificateShareByHandle(ctx.UserContext(), handle)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(CertificateShareStatusResponse{
		Status:      *status,
		Certificate: cert,
	})
}
