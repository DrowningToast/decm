package event

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

// DeleteEventContract godoc
// @Summary Delete event contract
// @Description Delete event contract for an event
// @Tags Event Contracts
// @ID delete-event-contract
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/contracts [delete]
func (h *Handler) DeleteEventContract(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Add a 30-second timeout for the usecase call
	ctxWithTimeout, cancel := context.WithTimeout(ctx.UserContext(), 30*time.Second)
	defer cancel()

	err = h.EventUc.DeleteEventContract(ctxWithTimeout, eventID)
	if err != nil {
		// Normalize the error before returning
		if _, isCustomErr := err.(*customerror.Err); isCustomErr {
			return err
		}
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.Status(http.StatusOK).JSON(map[string]string{"message": "Event contract deleted successfully"})
}
