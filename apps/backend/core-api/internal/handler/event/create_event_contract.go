package event

import (
	"apps/backend/common/customerror"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	eventUc "apps/backend/core-api/internal/usecase/event"
)

// CreateEventContract godoc
// @Summary Create event contract
// @Description Create a new event contract for an event
// @Tags Event
// @ID create-event-contract
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body CreateEventContractRequest true "Event contract data"
// @Success 201 {object} EventContractResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/contracts [post]
func (h *Handler) CreateEventContract(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req CreateEventContractRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	if err := req.IsValid(); err != nil {
		return err
	}

	// Add a 30-second timeout for the usecase call
	ctxWithTimeout, cancel := context.WithTimeout(ctx.UserContext(), 30*time.Second)
	defer cancel()

	params := eventUc.CreateEventContractParams{
		AccessManagerContractAddress: req.AccessManagerContractAddress,
		EventContractAddress:         req.EventContractAddress,
		TicketContractAddress:        pgtype.Text{String: req.TicketContractAddress, Valid: req.TicketContractAddress != ""},
		CertificateContractAddress:   pgtype.Text{String: req.CertificateContractAddress, Valid: req.CertificateContractAddress != ""},
	}

	contract, err := h.EventUc.CreateEventContract(ctxWithTimeout, eventID, params)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(EventContractResponse{
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
