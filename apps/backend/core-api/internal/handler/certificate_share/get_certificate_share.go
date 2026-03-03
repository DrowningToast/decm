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

type UnlockCertificateShareBody struct {
	Password string `json:"password"`
}

func (b *UnlockCertificateShareBody) Parse(ctx *fiber.Ctx) error {
	return ctx.BodyParser(b)
}

// @Summary Get certificate share
// @Description Retrieve a shared certificate by its handle. Returns PASSWORD_LOCKED if password-protected, VALID_BUT_PENDING if the certificate has not been claimed yet, or READY with certificate data.
// @ID get-certificate-share
// @Tags CertificateShares
// @Produce json
// @Param handle path string true "Share handle"
// @Success 200 {object} CertificateShareStatusResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificate-shares/{handle} [get]
func (h *Handler) GetCertificateShare(ctx *fiber.Ctx) error {
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

// @Summary Get password-protected certificate share
// @Description Unlock a password-protected share link and retrieve certificate data. Returns PASSWORD_LOCKED if the password is incorrect.
// @ID get-certificate-share-with-password
// @Tags CertificateShares
// @Accept json
// @Produce json
// @Param handle path string true "Share handle"
// @Param body body UnlockCertificateShareBody true "Password"
// @Success 200 {object} CertificateShareStatusResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificate-shares/{handle}/unlock [post]
func (h *Handler) GetCertificateShareWithPassword(ctx *fiber.Ctx) error {
	handle := ctx.Params("handle")
	if handle == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("handle is required"))
	}

	var body UnlockCertificateShareBody
	if err := body.Parse(ctx); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "failed to parse request body"))
	}
	if body.Password == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("password is required"))
	}

	cert, status, err := h.CertificateShareUc.GetCertificateShareByHandleWithPassword(ctx.UserContext(), handle, body.Password)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(CertificateShareStatusResponse{
		Status:      *status,
		Certificate: cert,
	})
}
