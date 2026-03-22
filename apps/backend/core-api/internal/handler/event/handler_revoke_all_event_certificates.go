package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RevokeAllEventCertificatesResponse struct {
	RevokedCertificates []*entity.EventCertificate `json:"revoked_certificates"`
}

// @Summary Revoke all event certificates
// @Description Revoke all certificates for an event by event ID
// @ID revoke-all-event-certificates
// @Tags Event Certificates
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} RevokeAllEventCertificatesResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/certificates/revoke-all [post]
func (h Handler) RevokeAllEventCertificates(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Get current user from JWT
	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	response, err := h.EventUc.RevokeAllEventCertificates(ctx.UserContext(), eventID, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
