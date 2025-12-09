package profile

import (
	"errors"

	"apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VerifyPasswordRequest struct {
	Password                   string    `json:"password"`
	AuthenticationCredentialId uuid.UUID `json:"authentication_credential_id"`
}

type VerifyPasswordResponse struct {
	Message   string `json:"message"`
	IsSuccess bool   `json:"is_success"`
}

// @Summary Verify password
// @Description Verify password
// @ID verify-password
// @Tags Profile
// @Accept json
// @Produce json
// @Param verifyPasswordRequest body VerifyPasswordRequest true "Verify Password Request"
// @Success 200 {object} VerifyPasswordResponse
// @Failure 400 {object} customerror.Err
// @Router /api/v1/profile/password/verify [post]
func (r *VerifyPasswordRequest) Parse(c *fiber.Ctx) error {
	if err := c.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func (r *VerifyPasswordRequest) IsValid() error {
	if len(r.Password) == 0 {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("password is required"))
	}
	return nil
}

func (h *Handler) VerifyPassword(c *fiber.Ctx) error {
	var verifyPasswordRequest VerifyPasswordRequest
	if err := verifyPasswordRequest.Parse(c); err != nil {
		return err
	}

	if err := verifyPasswordRequest.IsValid(); err != nil {
		return err
	}

	isValid, err := h.ProfileUc.VerifyUserPassword(c.Context(), verifyPasswordRequest.AuthenticationCredentialId, verifyPasswordRequest.Password)
	if err != nil {
		return err
	}

	message := "Password verified successfully"
	if !isValid {
		message = "Invalid password"
	}

	return c.JSON(VerifyPasswordResponse{
		Message:   message,
		IsSuccess: isValid,
	})
}
