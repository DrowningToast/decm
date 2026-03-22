package certificate_share_handler

import (
	customerror "apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	certificate_share_usecase "apps/backend/core-api/internal/usecase/certificate_share"
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UpdateCertificateShareBody holds the optional fields for a share update.
// The password field is parsed manually to distinguish three cases:
//   - key absent: do not change the existing password
//   - key present, value "" or null: remove password protection
//   - key present, value non-empty: set a new password
type UpdateCertificateShareBody struct {
	passwordPresent bool
	Password        *string `json:"password,omitempty" validate:"omitempty"`
	Active          *bool   `json:"active,omitempty" validate:"omitempty"`
}

func (b *UpdateCertificateShareBody) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if rawActive, ok := raw["active"]; ok {
		var active bool
		if err := json.Unmarshal(rawActive, &active); err != nil {
			return err
		}
		b.Active = &active
	}
	if rawPw, ok := raw["password"]; ok {
		b.passwordPresent = true
		if string(rawPw) != "null" {
			var pw string
			if err := json.Unmarshal(rawPw, &pw); err != nil {
				return err
			}
			b.Password = &pw
		}
	}
	return nil
}

func (b *UpdateCertificateShareBody) Parse(ctx *fiber.Ctx) error {
	return ctx.BodyParser(b)
}

// IsValid validates the body. Both fields are optional, so any parsed body is valid.
func (b *UpdateCertificateShareBody) IsValid() *customerror.Err {
	return nil
}

type UpdateCertificateShareResponse struct {
	Share *CertificateShareViewModel `json:"share"`
}

// @Summary Update certificate share link
// @Description Update the password of an existing share link. Only the certificate owner may call this endpoint. Omit the password field to leave existing protection unchanged; set to empty string or null to remove protection; set to a non-empty string to set a new password.
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
// @Router /api/v1/certificate-shares/config/{share_id} [patch]
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
		if cerr := body.IsValid(); cerr != nil {
			return cerr
		}
	}

	passwordUpdate := event_datagateway.PasswordUpdate{}
	if body.passwordPresent {
		passwordUpdate.Changed = true
		if body.Password != nil && *body.Password != "" {
			hp, err := hashutils.HashPassword(*body.Password)
			if err != nil {
				return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to hash password"))
			}
			passwordUpdate.Value = &hp
		}
	}

	share, err := h.CertificateShareUc.UpdateCertificateShare(ctx.UserContext(), currentUser, shareId, passwordUpdate, body.Active)
	if err != nil {
		return err
	}

	// Invalidate cached share data so password changes take effect immediately.
	h.cache.Delete(share.Handle)

	return ctx.Status(fiber.StatusOK).JSON(UpdateCertificateShareResponse{Share: newCertificateShareViewModel(certificate_share_usecase.NewCertificateShareViewModel(share))})
}
