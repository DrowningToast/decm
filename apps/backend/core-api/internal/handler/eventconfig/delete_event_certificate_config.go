package eventconfig

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

// DeleteEventCertificateConfig godoc
// @Summary Delete event certificate config
// @Description Delete the event certificate configuration for an event
// @ID delete-event-certificate-config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/certificate-config [delete]
func (h *Handler) DeleteEventCertificateConfig(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	err = h.EventConfigUc.DeleteEventCertificateConfig(ctx.UserContext(), eventID)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.Status(http.StatusOK).JSON(map[string]string{"message": "Event certificate config deleted successfully"})
}
