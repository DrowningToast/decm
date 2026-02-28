package certificate

import (
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GenerateCertificateImage generates and returns a certificate image as PNG
// @Summary Generate certificate image for participant
// @Description Generates a PNG certificate image for the authenticated user's certificate
// @ID generate-certificate-image
// @Tags certificates
// @Produce image/png
// @Param certificate_id path string true "Certificate ID" format(uuid)
// @Success 200 {file} binary "PNG certificate image"
// @Failure 400 {object} customerror.ErrResponse "Invalid certificate ID"
// @Failure 403 {object} customerror.ErrResponse "User does not own this certificate"
// @Failure 404 {object} customerror.ErrResponse "Certificate not found"
// @Failure 500 {object} customerror.ErrResponse "Internal server error"
// @Security CookieAuth
// @Router /api/v1/certificates/{certificate_id}/image [get]
func (h *Handler) GenerateCertificateImage(ctx *fiber.Ctx) error {
	// 1. Get current user from context (set by authentication middleware)
	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	// 2. Parse certificate ID from path
	certificateIDStr := ctx.Params("certificate_id")
	certificateID, err := uuid.Parse(certificateIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid certificate ID")
	}

	// 3. Generate the certificate image with authorization checks
	pngBytes, err := h.EventUc.GenerateCertificateImageForParticipant(
		ctx.Context(),
		certificateID,
		currentUser.UserId,
		currentUser.Email,
	)
	if err != nil {
		h.Logger.Error("Failed to generate certificate image",
			"error", err,
			"certificate_id", certificateID,
			"user_id", currentUser.UserId,
		)
		return err
	}

	// 4. Return the PNG image
	ctx.Set("Content-Type", "image/png")
	ctx.Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	ctx.Set("Content-Disposition", "inline; filename=certificate.png")

	return ctx.Send(pngBytes)
}
