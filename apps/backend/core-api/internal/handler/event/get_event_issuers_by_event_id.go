package event

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

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
func (h *Handler) GetEventIssuersByEventID(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	issuers, err := h.EventUc.GetEventIssuersByEventID(ctx.UserContext(), eventID)
	if err != nil {
		return customerror.Parse(&customerror.ErrNotFound, err)
	}

	var response []EventIssuerResponse
	for _, issuer := range issuers {
		response = append(response, EventIssuerResponse{
			ID:                 issuer.ID,
			EventID:            issuer.EventID,
			IssuerCredentialID: issuer.IssuerCredentialID,
			IsSigned:           issuer.IsSigned,
			Signature:          issuer.Signature.String,
			CreatedAt:          issuer.CreatedAt.Time.String(),
			UpdatedAt:          issuer.UpdatedAt.Time.String(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
