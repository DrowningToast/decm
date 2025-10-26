package event

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	eventUc "apps/backend/core-api/internal/usecase/event"
)

type UpdateEventContractRequest struct {
	AccessManagerContractAddress string `json:"access_manager_contract_address"`
	EventContractAddress         string `json:"event_contract_address"`
	TicketContractAddress        string `json:"ticket_contract_address"`
	CertificateContractAddress   string `json:"certificate_contract_address"`
}

func (r *UpdateEventContractRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

// UpdateEventContract godoc
// @Summary Update event contract
// @Description Update the event contract for an event
// @ID update-event-contract
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body UpdateEventContractRequest true "Event contract data"
// @Success 200 {object} EventContractResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/contracts [put]
func (h *Handler) UpdateEventContract(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req UpdateEventContractRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	if err := req.IsValid(); err != nil {
		return err
	}

	params := eventUc.UpdateEventContractParams{
		AccessManagerContractAddress: req.AccessManagerContractAddress,
		EventContractAddress:         req.EventContractAddress,
		TicketContractAddress:        pgtype.Text{String: req.TicketContractAddress, Valid: req.TicketContractAddress != ""},
		CertificateContractAddress:   pgtype.Text{String: req.CertificateContractAddress, Valid: req.CertificateContractAddress != ""},
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.UserContext(), 30*time.Second)
	defer cancel()

	contract, err := h.EventUc.UpdateEventContract(ctxWithTimeout, eventID, params)
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
