package eventconfig

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/usecase/event_config"
)

type EventContractHandler struct {
	usecase *event_config.EventConfigUsecase
}

func NewEventContractHandler(usecase *event_config.EventConfigUsecase) *EventContractHandler {
	return &EventContractHandler{
		usecase: usecase,
	}
}

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
func (h *EventContractHandler) CreateEventContract(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.NewBadRequestError("Invalid event ID")
	}

	var req CreateEventContractRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.NewBadRequestError("Invalid request body")
	}

	if err := req.IsValid(); err != nil {
		return err
	}

	params := event_config.CreateEventContractParams{
		AccessManagerContractAddress: req.AccessManagerContractAddress,
		EventContractAddress:         req.EventContractAddress,
		TicketContractAddress:        pgtype.Text{String: req.TicketContractAddress, Valid: req.TicketContractAddress != ""},
		CertificateContractAddress:   pgtype.Text{String: req.CertificateContractAddress, Valid: req.CertificateContractAddress != ""},
	}

	contract, err := h.usecase.EventContract.CreateEventContract(ctx.UserContext(), eventID, params)
	if err != nil {
		return customerror.NewInternalServerError(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(EventContractResponse{
		ID:                           contract.ID,
		EventID:                      contract.EventID,
		AccessManagerContractAddress: contract.AccessManagerContractAddress,
		EventContractAddress:         contract.EventContractAddress,
		TicketContractAddress:        contract.TicketContractAddress.String,
		CertificateContractAddress:   contract.CertificateContractAddress.String,
		CreatedAt:                    contract.CreatedAt,
		UpdatedAt:                    contract.UpdatedAt,
	})
}

// GetEventContractByEventID godoc
// @Summary Get event contract by event ID
// @Description Get the event contract for an event
// @ID get-event-contract-by-event-id
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} EventContractResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/contracts [get]
func (h *EventContractHandler) GetEventContractByEventID(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.NewBadRequestError("Invalid event ID")
	}

	contract, err := h.usecase.EventContract.GetEventContractByEventID(ctx.UserContext(), eventID)
	if err != nil {
		return customerror.NewNotFoundError("Event contract not found")
	}

	return ctx.Status(http.StatusOK).JSON(EventContractResponse{
		ID:                           contract.ID,
		EventID:                      contract.EventID,
		AccessManagerContractAddress: contract.AccessManagerContractAddress,
		EventContractAddress:         contract.EventContractAddress,
		TicketContractAddress:        contract.TicketContractAddress.String,
		CertificateContractAddress:   contract.CertificateContractAddress.String,
		CreatedAt:                    contract.CreatedAt,
		UpdatedAt:                    contract.UpdatedAt,
	})
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
func (h *EventContractHandler) UpdateEventContract(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.NewBadRequestError("Invalid event ID")
	}

	var req UpdateEventContractRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.NewBadRequestError("Invalid request body")
	}

	if err := req.IsValid(); err != nil {
		return err
	}

	params := event_config.UpdateEventContractParams{
		AccessManagerContractAddress: req.AccessManagerContractAddress,
		EventContractAddress:         req.EventContractAddress,
		TicketContractAddress:        pgtype.Text{String: req.TicketContractAddress, Valid: req.TicketContractAddress != ""},
		CertificateContractAddress:   pgtype.Text{String: req.CertificateContractAddress, Valid: req.CertificateContractAddress != ""},
	}

	contract, err := h.usecase.EventContract.UpdateEventContract(ctx.UserContext(), eventID, params)
	if err != nil {
		return customerror.NewInternalServerError(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(EventContractResponse{
		ID:                           contract.ID,
		EventID:                      contract.EventID,
		AccessManagerContractAddress: contract.AccessManagerContractAddress,
		EventContractAddress:         contract.EventContractAddress,
		TicketContractAddress:        contract.TicketContractAddress.String,
		CertificateContractAddress:   contract.CertificateContractAddress.String,
		CreatedAt:                    contract.CreatedAt,
		UpdatedAt:                    contract.UpdatedAt,
	})
}

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
func (h *EventContractHandler) DeleteEventContract(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.NewBadRequestError("Invalid event ID")
	}

	err = h.usecase.EventContract.DeleteEventContract(ctx.UserContext(), eventID)
	if err != nil {
		return customerror.NewInternalServerError(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(map[string]string{"message": "Event contract deleted successfully"})
}

// Request/Response structures
type CreateEventContractRequest struct {
	AccessManagerContractAddress string `json:"access_manager_contract_address"`
	EventContractAddress         string `json:"event_contract_address"`
	TicketContractAddress        string `json:"ticket_contract_address"`
	CertificateContractAddress   string `json:"certificate_contract_address"`
}

func (r *CreateEventContractRequest) IsValid() *customerror.ErrResponse {
	if r.AccessManagerContractAddress == "" {
		return customerror.NewBadRequestError("Access manager contract address is required")
	}
	if r.EventContractAddress == "" {
		return customerror.NewBadRequestError("Event contract address is required")
	}

	return nil
}

type UpdateEventContractRequest struct {
	AccessManagerContractAddress string `json:"access_manager_contract_address"`
	EventContractAddress         string `json:"event_contract_address"`
	TicketContractAddress        string `json:"ticket_contract_address"`
	CertificateContractAddress   string `json:"certificate_contract_address"`
}

func (r *UpdateEventContractRequest) IsValid() *customerror.ErrResponse {
	if r.AccessManagerContractAddress == "" {
		return customerror.NewBadRequestError("Access manager contract address is required")
	}
	if r.EventContractAddress == "" {
		return customerror.NewBadRequestError("Event contract address is required")
	}

	return nil
}

type EventContractResponse struct {
	ID                           uuid.UUID `json:"id"`
	EventID                      uuid.UUID `json:"event_id"`
	AccessManagerContractAddress string    `json:"access_manager_contract_address"`
	EventContractAddress         string    `json:"event_contract_address"`
	TicketContractAddress        string    `json:"ticket_contract_address"`
	CertificateContractAddress   string    `json:"certificate_contract_address"`
	CreatedAt                    string    `json:"created_at"`
	UpdatedAt                    string    `json:"updated_at"`
}
