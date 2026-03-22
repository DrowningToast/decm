package certificate_share_handler

import (
	customerror "apps/backend/common/customerror"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

type GetCertificateShareImageBody struct {
	Password *string `json:"password,omitempty"`
}

// @Summary Get certificate share image
// @Description Returns a PNG certificate image for a share link. For password-protected shares, include the password in the request body.
// @ID get-certificate-share-image
// @Tags CertificateShares
// @Accept json
// @Produce image/png
// @Param handle path string true "Share handle"
// @Param body body GetCertificateShareImageBody false "Optional password for password-protected shares"
// @Success 200 {file} binary "PNG certificate image"
// @Failure 400 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificate-shares/{handle}/image [post]
func (h *Handler) GetCertificateShareImage(ctx *fiber.Ctx) error {
	handle := ctx.Params("handle")
	if handle == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("handle is required"))
	}

	var body GetCertificateShareImageBody
	if len(ctx.Body()) > 0 {
		if err := ctx.BodyParser(&body); err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "failed to parse request body"))
		}
	}

	pngBytes, err := h.CertificateShareUc.GetCertificateShareImage(ctx.UserContext(), handle, body.Password)
	if err != nil {
		return err
	}

	ctx.Set("Content-Type", "image/png")
	if body.Password == nil {
		ctx.Set("Cache-Control", "public, max-age=86400")
	} else {
		ctx.Set("Cache-Control", "private, no-store")
	}
	ctx.Set("Content-Disposition", "inline; filename=certificate.png")
	return ctx.Send(pngBytes)
}
