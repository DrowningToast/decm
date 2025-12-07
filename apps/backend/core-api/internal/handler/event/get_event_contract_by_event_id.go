package event

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

// GetEventContractByEventID godoc
// @Summary Get event contract by event ID
// @Description Get the event contract for an event
// @Tags Events
// @ID get-event-contract-by-event-id
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} EventContractResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/contracts [get]
func (h *Handler) GetEventContractByEventID(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	contract, err := h.EventUc.GetEventContractByEventID(ctx.UserContext(), eventID)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(EventContractResponse{
		ID:                           contract.ID,
		EventID:                      contract.EventID,
		AccessManagerContractAddress: contract.AccessManagerContractAddress,
		EventContractAddress:         contract.EventContractAddress,
		TicketContractAddress:        contract.TicketContractAddress,
		CertificateContractAddress:   contract.CertificateContractAddress,
		CreatedAt:                    contract.CreatedAt.Format(time.RFC3339),
		UpdatedAt:                    contract.UpdatedAt.Format(time.RFC3339),
	})
}
