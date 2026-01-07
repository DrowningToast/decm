package eventconfig

import (
	"apps/backend/common/customerror"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DeleteEventCertificateConfig godoc
// @Summary Delete event certificate config
// @Description Delete event certificate configuration for an event
// @Tags EventConfig
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
		return err
	}

	return ctx.Status(http.StatusOK).JSON(map[string]string{"message": "Event certificate config deleted successfully"})
}
