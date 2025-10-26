package event

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

// DeleteEventIssuer godoc
// @Summary Delete event issuer
// @Description Delete an event issuer
// @Tags Events
// @ID delete-event-issuer
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param issuer_id path string true "Issuer ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/issuers/{issuer_id} [delete]
func (h *Handler) DeleteEventIssuer(ctx *fiber.Ctx) error {
	issuerID, err := uuid.Parse(ctx.Params("issuer_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	// Add a 30-second timeout for the usecase call
	ctxWithTimeout, cancel := context.WithTimeout(ctx.UserContext(), 30*time.Second)
	defer cancel()

	err = h.EventUc.DeleteEventIssuer(ctxWithTimeout, issuerID, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(map[string]string{"message": "Event issuer deleted successfully"})
}
