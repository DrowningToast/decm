package eventconfig

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

// DeleteEventRegistrationConfig godoc
// @Summary Delete event registration config
// @Description Delete the event registration configuration for an event
// @ID delete-event-registration-config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/registration-config [delete]
func (h *Handler) DeleteEventRegistrationConfig(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	err = h.EventConfigUc.DeleteEventRegistrationConfig(ctx.UserContext(), eventID)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.Status(http.StatusOK).JSON(map[string]string{"message": "Event registration config deleted successfully"})
}
