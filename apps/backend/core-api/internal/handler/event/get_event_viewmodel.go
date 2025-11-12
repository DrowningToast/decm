package event

import (
	"context"
	"time"

	"apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetEventById godoc
// @Summary Get event by ID
// @Description Get event by ID
// @Tags Events
// @ID get-event-viewmodel-by-id
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} event.EventViewModel
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/viewmodel [get]
func (h *Handler) GetEventViewModel(ctx *fiber.Ctx) error {
	eventId, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Create a context with 30-second timeout
	base := ctx.UserContext()
	ctxWithTimeout, cancel := context.WithTimeout(base, 30*time.Second)
	defer cancel()

	// User context
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if currentUser == nil {
		return customerror.Parse(&customerror.ErrUnauthorized, nil)
	}

	event, err := h.EventUc.GetEventViewModelByEventId(ctxWithTimeout, eventId, currentUser)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if event == nil {
		return customerror.Parse(&customerror.ErrNotFound, nil)
	}

	return ctx.JSON(event)
}
