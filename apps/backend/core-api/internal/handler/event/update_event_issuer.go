package event

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	eventUc "apps/backend/core-api/internal/usecase/event"
)

type UpdateEventIssuerRequest struct {
	EventID            uuid.UUID `json:"event_id"`
	IssuerCredentialID uuid.UUID `json:"issuer_credential_id"`
}

func (r *UpdateEventIssuerRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

// UpdateEventIssuer godoc
// @Summary Update event issuer
// @Description Update an event issuer
// @ID update-event-issuer
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body []UpdateEventIssuerRequest true "Event issuer data"
// @Success 200 {object} EventIssuerResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/issuers [put]
func (h *Handler) UpdateEventIssuer(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req []UpdateEventIssuerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	for _, req := range req {
		if err := req.IsValid(); err != nil {
			return err
		}
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	params := make([]eventUc.UpdateEventIssuerParams, len(req))
	for i, req := range req {
		params[i] = eventUc.UpdateEventIssuerParams{
			EventID:            eventID,
			IssuerCredentialID: req.IssuerCredentialID,
		}
	}

	issuers, err := h.EventUc.UpdateEventIssuer(ctx.UserContext(), eventID, params, currentUser)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(issuers)
}
