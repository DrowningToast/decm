package certificate_share

import (
	customerror "apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/entity"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UpdateCertificateShareBody struct {
	Password *string `json:"password"`
}

func (b *UpdateCertificateShareBody) Parse(ctx *fiber.Ctx) error {
	return ctx.BodyParser(b)
}

type UpdateCertificateShareResponse struct {
	Share *entity.CertificateShare `json:"share"`
}

// @Summary Update certificate share link
// @Description Update the password of an existing share link. Only the certificate owner may call this endpoint. Set password to null or omit to remove password protection.
// @ID update-certificate-share
// @Tags CertificateShares
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param share_id path string true "Share ID"
// @Param body body UpdateCertificateShareBody false "Optional password configuration"
// @Success 200 {object} UpdateCertificateShareResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificate-shares/{share_id} [patch]
func (h *Handler) UpdateCertificateShare(ctx *fiber.Ctx) error {
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	shareIdStr := ctx.Params("share_id")
	shareId, err := uuid.Parse(shareIdStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "invalid share_id format"))
	}

	var body UpdateCertificateShareBody
	if len(ctx.Body()) > 0 {
		if err := body.Parse(ctx); err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "failed to parse request body"))
		}
	}

	var hashedPassword *string
	if body.Password != nil && *body.Password != "" {
		hp, err := hashutils.HashPassword(*body.Password)
		if err != nil {
			return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to hash password"))
		}
		hashedPassword = &hp
	}

	share, err := h.CertificateShareUc.UpdateCertificateShare(ctx.UserContext(), currentUser, shareId, hashedPassword)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(UpdateCertificateShareResponse{Share: share})
}
