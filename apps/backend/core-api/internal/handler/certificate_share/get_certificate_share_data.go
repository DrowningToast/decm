package certificate_share

import (
	customerror "apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

type CertificateShareDataResponse struct {
	Data *entity.CertificatePayload `json:"data"`
}

type UnlockCertificateShareBody struct {
	Password string `json:"password"`
}

func (b *UnlockCertificateShareBody) Parse(ctx *fiber.Ctx) error {
	return ctx.BodyParser(b)
}

// @Summary Get on-chain certificate share data
// @Description Retrieve the full on-chain Verifiable Credential data for a public share link. The certificate must be claimed.
// @ID get-certificate-share-data
// @Tags CertificateShares
// @Produce json
// @Param handle path string true "Share handle"
// @Success 200 {object} CertificateShareDataResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificate-shares/{handle}/data [get]
func (h *Handler) GetCertificateShareData(ctx *fiber.Ctx) error {
	handle := ctx.Params("handle")
	if handle == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("handle is required"))
	}

	data, err := h.CertificateShareUc.GetCertificateShareData(ctx.UserContext(), handle)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(CertificateShareDataResponse{Data: data})
}

// @Summary Get on-chain certificate share data (password-protected)
// @Description Retrieve the full on-chain Verifiable Credential data for a password-protected share link.
// @ID get-certificate-share-data-with-password
// @Tags CertificateShares
// @Accept json
// @Produce json
// @Param handle path string true "Share handle"
// @Param body body UnlockCertificateShareBody true "Password"
// @Success 200 {object} CertificateShareDataResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificate-shares/{handle}/data/unlock [post]
func (h *Handler) GetCertificateShareDataWithPassword(ctx *fiber.Ctx) error {
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

	data, err := h.CertificateShareUc.GetCertificateShareDataWithPassword(ctx.UserContext(), handle, body.Password)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(CertificateShareDataResponse{Data: data})
}
