package event

import (
	"context"
	"net/http"
	"strconv"
	"time"

	customerror "apps/backend/common/customerror"
	event_uc "apps/backend/core-api/internal/usecase/event"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetEventsByOwnerCredentialsId godoc
// @Summary Get events by owner credentials ID
// @Description Get events by owner credentials ID
// @ID get-events-by-owner-credentials-id
// @Accept json
// @Produce json
// @Param owner_credential_id path string true "Owner Credentials ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} EventResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/owner-credentials/{owner_credential_id} [get]
func (h *Handler) GetEventsByOwnerCredentialsId(ctx *fiber.Ctx) error {
	ownerCredentialID, err := uuid.Parse(ctx.Params("owner_credential_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	limitCount, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	offsetCount, err := strconv.Atoi(ctx.Query("offset", "0"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Add a 30-second timeout for the usecase call
	ctxWithTimeout, cancel := context.WithTimeout(ctx.UserContext(), 30*time.Second)
	defer cancel()

	events, err := h.EventUc.ListEventsByOwnerCredentialID(ctxWithTimeout, ownerCredentialID, int32(limitCount), int32(offsetCount))
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	result := make([]*event_uc.EventResponse, len(events))
	for i, event := range events {
		eventResponse, err := h.EventUc.ToEventResponse(ctx.Context(), event)
		if err != nil {
			return customerror.Parse(&customerror.ErrInternalServer, err)
		}
		result[i] = eventResponse
	}

	return ctx.Status(http.StatusOK).JSON(result)
}
