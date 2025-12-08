package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetEventCertificatesResponse struct {
	Certificates []*entity.EventCertificate `json:"certificates"`
}

// @Summary Get event certificates
// @Description Get all certificates for an event
// @ID get-event-certificates
// @Tags Event Certificates
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} GetEventCertificatesResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/certificates [get]
func (h Handler) GetEventCertificates(ctx *fiber.Ctx) error {
	// Get event ID from path parameter
	eventIDStr := ctx.Params("event_id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Get current user from JWT
	currentUser := ctx.Locals("user").(*auth.JwtClaims)

	// Get certificates for the event
	certificates, err := h.EventUc.GetEventCertificatesByEventId(ctx.UserContext(), eventID, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(GetEventCertificatesResponse{
		Certificates: certificates,
	})
}
