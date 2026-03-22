package certificate_share_handler

import (
	customerror "apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	certificate_share_usecase "apps/backend/core-api/internal/usecase/certificate_share"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateCertificateShareBody struct {
	Password *string `json:"password,omitempty"`
}

func (b *CreateCertificateShareBody) Parse(ctx *fiber.Ctx) error {
	return ctx.BodyParser(b)
}

// IsValid validates the body. Password is optional, so any parsed body is valid.
func (b *CreateCertificateShareBody) IsValid() *customerror.Err {
	return nil
}

type CreateCertificateShareResponse struct {
	Share *CertificateShareViewModel `json:"share"`
}

// @Summary Create certificate share link
// @Description Create a shareable link for a certificate. Only the certificate owner may call this endpoint. An optional password can be set to restrict access.
// @ID create-certificate-share
// @Tags CertificateShares
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param certificate_id path string true "Certificate ID"
// @Param body body CreateCertificateShareBody false "Optional password configuration"
// @Success 201 {object} CreateCertificateShareResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificate-shares/config/{certificate_id} [post]
func (h *Handler) CreateCertificateShare(ctx *fiber.Ctx) error {
	// 1. Auth check first
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	// 2. Parse and validate path param
	certificateId, err := uuid.Parse(ctx.Params("certificate_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.Wrap(err, "invalid certificate_id format"))
	}

	// 3. Parse optional body (password is optional — no body is valid)
	var body CreateCertificateShareBody
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

	share, err := h.CertificateShareUc.CreateCertificateShare(ctx.UserContext(), currentUser, certificateId, hashedPassword)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(CreateCertificateShareResponse{Share: newCertificateShareViewModel(certificate_share_usecase.NewCertificateShareViewModel(share))})
}
