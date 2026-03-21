package certificate_share_handler

import (
	customerror "apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	stdjson "encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

type CertificateShareDataResponse struct {
	Data *entity.CertificatePayload `json:"data"`
}

type GetCertificateShareBody struct {
	Password *string `json:"password,omitempty"`
}

func (b *GetCertificateShareBody) Parse(ctx *fiber.Ctx) error {
	return ctx.BodyParser(b)
}

// @Summary Get on-chain certificate share data
// @Description Retrieve the full on-chain Verifiable Credential data for a share link. For password-protected shares, include the password in the request body.
// @ID get-certificate-share-data
// @Tags CertificateShares
// @Accept json
// @Produce json
// @Param handle path string true "Share handle"
// @Param body body GetCertificateShareBody false "Optional password for password-protected shares"
// @Success 200 {object} CertificateShareDataResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificate-shares/{handle} [post]
func (h *Handler) GetCertificateShareData(ctx *fiber.Ctx) error {
	handle := ctx.Params("handle")
	if handle == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("handle is required"))
	}

	var body GetCertificateShareBody
	if len(ctx.Body()) > 0 {
		if err := body.Parse(ctx); err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "failed to parse request body"))
		}
	}

	data, err := h.CertificateShareUc.GetCertificateShareData(ctx.UserContext(), handle, body.Password)
	if err != nil {
		return err
	}

	// goccy/go-json (Fiber's global encoder) panics on structs that contain
	// []string fields whose JSON tag starts with '@' (e.g. "@context" in the W3C
	// VC header). Use encoding/json directly to avoid the segfault.
	b, err := stdjson.Marshal(CertificateShareDataResponse{Data: data})
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to marshal response"))
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return ctx.Status(fiber.StatusOK).Send(b)
}
