package event

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/common/customerror"
	eventUc "apps/backend/core-api/internal/usecase/event"
)

// CreateEventContract godoc
// @Summary Create event contract
// @Description Create a new event contract for an event
// @ID create-event-contract
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body CreateEventContractRequest true "Event contract data"
// @Success 200 {object} EventContractResponse
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

	params := eventUc.CreateEventContractParams{
		AccessManagerContractAddress: req.AccessManagerContractAddress,
		EventContractAddress:         req.EventContractAddress,
		TicketContractAddress:        pgtype.Text{String: req.TicketContractAddress, Valid: req.TicketContractAddress != ""},
		CertificateContractAddress:   pgtype.Text{String: req.CertificateContractAddress, Valid: req.CertificateContractAddress != ""},
	}

	contract, err := h.EventUc.CreateEventContract(ctx.UserContext(), eventID, params)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(EventContractResponse{
		ID:                           contract.ID,
		EventID:                      contract.EventID,
		AccessManagerContractAddress: contract.AccessManagerContractAddress,
		EventContractAddress:         contract.EventContractAddress,
		TicketContractAddress:        contract.TicketContractAddress.String,
		CertificateContractAddress:   contract.CertificateContractAddress.String,
		CreatedAt:                    contract.CreatedAt.Time.String(),
		UpdatedAt:                    contract.UpdatedAt.Time.String(),
	})
}
