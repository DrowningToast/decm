package certificate_share_handler

import (
	customerror "apps/backend/common/customerror"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

// @Summary Get certificate share image
// @Description Returns a PNG certificate image for a share link. For password-protected shares, pass the password as a query parameter.
// @ID get-certificate-share-image
// @Tags CertificateShares
// @Produce image/png
// @Param handle path string true "Share handle"
// @Param password query string false "Password for password-protected shares"
// @Success 200 {file} binary "PNG certificate image"
// @Failure 400 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificate-shares/{handle}/image [get]
func (h *Handler) GetCertificateShareImage(ctx *fiber.Ctx) error {
	handle := ctx.Params("handle")
	if handle == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("handle is required"))
	}

	var password *string
	if pw := ctx.Query("password"); pw != "" {
		password = &pw
	}

	pngBytes, err := h.CertificateShareUc.GetCertificateShareImage(ctx.UserContext(), handle, password)
	if err != nil {
		return err
	}

	ctx.Set("Content-Type", "image/png")
	ctx.Set("Cache-Control", "public, max-age=86400")
	ctx.Set("Content-Disposition", "inline; filename=certificate.png")
	return ctx.Send(pngBytes)
}
