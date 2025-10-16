package eventconfig

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/usecase/event_config"
)

type EventIssuerHandler struct {
	usecase *event_config.EventConfigUsecase
}

func NewEventIssuerHandler(usecase *event_config.EventConfigUsecase) *EventIssuerHandler {
	return &EventIssuerHandler{
		usecase: usecase,
	}
}

// CreateEventIssuer godoc
// @Summary Create event issuer
// @Description Create a new event issuer for an event
// @ID create-event-issuer
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body CreateEventIssuerRequest true "Event issuer data"
// @Success 200 {object} EventIssuerResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/issuers [post]
func (h *EventIssuerHandler) CreateEventIssuer(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.NewBadRequestError("Invalid event ID")
	}

	var req CreateEventIssuerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.NewBadRequestError("Invalid request body")
	}

	if err := req.IsValid(); err != nil {
		return err
	}

	params := event_config.CreateEventIssuerParams{
		EventID:            eventID,
		IssuerCredentialID: req.IssuerCredentialID,
		IsSigned:           req.IsSigned,
		Signature:          pgtype.Text{String: req.Signature, Valid: req.Signature != ""},
	}

	issuer, err := h.usecase.EventIssuer.CreateEventIssuer(ctx.UserContext(), params)
	if err != nil {
		return customerror.NewInternalServerError(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(EventIssuerResponse{
		ID:                 issuer.ID,
		EventID:            issuer.EventID,
		IssuerCredentialID: issuer.IssuerCredentialID,
		IsSigned:           issuer.IsSigned,
		Signature:          issuer.Signature.String,
		CreatedAt:          issuer.CreatedAt,
		UpdatedAt:          issuer.UpdatedAt,
	})
}

// GetEventIssuersByEventID godoc
// @Summary Get event issuers by event ID
// @Description Get all event issuers for an event
// @ID get-event-issuers-by-event-id
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {array} EventIssuerResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/issuers [get]
func (h *EventIssuerHandler) GetEventIssuersByEventID(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.NewBadRequestError("Invalid event ID")
	}

	issuers, err := h.usecase.EventIssuer.GetEventIssuersByEventID(ctx.UserContext(), eventID)
	if err != nil {
		return customerror.NewNotFoundError("Event issuers not found")
	}

	var response []EventIssuerResponse
	for _, issuer := range issuers {
		response = append(response, EventIssuerResponse{
			ID:                 issuer.ID,
			EventID:            issuer.EventID,
			IssuerCredentialID: issuer.IssuerCredentialID,
			IsSigned:           issuer.IsSigned,
			Signature:          issuer.Signature.String,
			CreatedAt:          issuer.CreatedAt,
			UpdatedAt:          issuer.UpdatedAt,
		})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// GetEventIssuerByID godoc
// @Summary Get event issuer by ID
// @Description Get an event issuer by its ID
// @ID get-event-issuer-by-id
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param issuer_id path string true "Issuer ID"
// @Success 200 {object} EventIssuerResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/issuers/{issuer_id} [get]
func (h *EventIssuerHandler) GetEventIssuerByID(ctx *fiber.Ctx) error {
	issuerID, err := uuid.Parse(ctx.Params("issuer_id"))
	if err != nil {
		return customerror.NewBadRequestError("Invalid issuer ID")
	}

	issuer, err := h.usecase.EventIssuer.GetEventIssuerByID(ctx.UserContext(), issuerID)
	if err != nil {
		return customerror.NewNotFoundError("Event issuer not found")
	}

	return ctx.Status(http.StatusOK).JSON(EventIssuerResponse{
		ID:                 issuer.ID,
		EventID:            issuer.EventID,
		IssuerCredentialID: issuer.IssuerCredentialID,
		IsSigned:           issuer.IsSigned,
		Signature:          issuer.Signature.String,
		CreatedAt:          issuer.CreatedAt,
		UpdatedAt:          issuer.UpdatedAt,
	})
}

// UpdateEventIssuer godoc
// @Summary Update event issuer
// @Description Update an event issuer
// @ID update-event-issuer
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param issuer_id path string true "Issuer ID"
// @Param request body UpdateEventIssuerRequest true "Event issuer data"
// @Success 200 {object} EventIssuerResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/issuers/{issuer_id} [put]
func (h *EventIssuerHandler) UpdateEventIssuer(ctx *fiber.Ctx) error {
	issuerID, err := uuid.Parse(ctx.Params("issuer_id"))
	if err != nil {
		return customerror.NewBadRequestError("Invalid issuer ID")
	}

	var req UpdateEventIssuerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.NewBadRequestError("Invalid request body")
	}

	if err := req.IsValid(); err != nil {
		return err
	}

	params := event_config.UpdateEventIssuerParams{
		IsSigned:  req.IsSigned,
		Signature: pgtype.Text{String: req.Signature, Valid: req.Signature != ""},
	}

	issuer, err := h.usecase.EventIssuer.UpdateEventIssuer(ctx.UserContext(), issuerID, params)
	if err != nil {
		return customerror.NewInternalServerError(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(EventIssuerResponse{
		ID:                 issuer.ID,
		EventID:            issuer.EventID,
		IssuerCredentialID: issuer.IssuerCredentialID,
		IsSigned:           issuer.IsSigned,
		Signature:          issuer.Signature.String,
		CreatedAt:          issuer.CreatedAt,
		UpdatedAt:          issuer.UpdatedAt,
	})
}

// DeleteEventIssuer godoc
// @Summary Delete event issuer
// @Description Delete an event issuer
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
func (h *EventIssuerHandler) DeleteEventIssuer(ctx *fiber.Ctx) error {
	issuerID, err := uuid.Parse(ctx.Params("issuer_id"))
	if err != nil {
		return customerror.NewBadRequestError("Invalid issuer ID")
	}

	err = h.usecase.EventIssuer.DeleteEventIssuer(ctx.UserContext(), issuerID)
	if err != nil {
		return customerror.NewInternalServerError(err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(map[string]string{"message": "Event issuer deleted successfully"})
}

// Request/Response structures
type CreateEventIssuerRequest struct {
	IssuerCredentialID uuid.UUID `json:"issuer_credential_id"`
	IsSigned           int32     `json:"is_signed"`
	Signature          string    `json:"signature"`
}

func (r *CreateEventIssuerRequest) IsValid() *customerror.ErrResponse {
	if r.IssuerCredentialID == uuid.Nil {
		return customerror.NewBadRequestError("Issuer credential ID is required")
	}
	if r.IsSigned < 0 || r.IsSigned > 1 {
		return customerror.NewBadRequestError("Is signed must be 0 (false) or 1 (true)")
	}

	return nil
}

type UpdateEventIssuerRequest struct {
	IsSigned  int32  `json:"is_signed"`
	Signature string `json:"signature"`
}

func (r *UpdateEventIssuerRequest) IsValid() *customerror.ErrResponse {
	if r.IsSigned < 0 || r.IsSigned > 1 {
		return customerror.NewBadRequestError("Is signed must be 0 (false) or 1 (true)")
	}

	return nil
}

type EventIssuerResponse struct {
	ID                 uuid.UUID `json:"id"`
	EventID            uuid.UUID `json:"event_id"`
	IssuerCredentialID uuid.UUID `json:"issuer_credential_id"`
	IsSigned           int32     `json:"is_signed"`
	Signature          string    `json:"signature"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
}
