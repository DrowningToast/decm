package event

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

// DeleteEventContract godoc
// @Summary Delete event contract
// @Description Delete the event contract for an event
// @ID delete-event-contract
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/contracts [delete]
func (h *Handler) DeleteEventContract(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	err = h.EventUc.DeleteEventContract(ctx.UserContext(), eventID)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(map[string]string{"message": "Event contract deleted successfully"})
}
