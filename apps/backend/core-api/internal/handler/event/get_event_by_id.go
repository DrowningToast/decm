package event

import (
	"apps/backend/common/customerror"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetEventById godoc
// @Summary Get event by ID
// @Description Get event by ID
// @Tags Events
// @ID get-event-by-id
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} EventResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id} [get]
func (h *Handler) GetEventById(ctx *fiber.Ctx) error {
	eventId, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Create a context with 30-second timeout
	base := ctx.UserContext()
	ctxWithTimeout, cancel := context.WithTimeout(base, 30*time.Second)
	defer cancel()

	event, err := h.EventUc.GetEventById(ctxWithTimeout, eventId)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	if event == nil {
		return customerror.Parse(&customerror.ErrNotFound, nil)
	}

	eventResponse, err := h.EventUc.ToEventResponse(ctxWithTimeout, event)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.JSON(eventResponse)
}
