package event

import (
	"apps/backend/common/customerror"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
func (h *Handler) GetEventIssuerByID(ctx *fiber.Ctx) error {
	issuerID, err := uuid.Parse(ctx.Params("issuer_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	issuer, err := h.EventUc.GetEventIssuerByID(ctx.UserContext(), issuerID)
	if err != nil {
		return customerror.Parse(&customerror.ErrNotFound, err)
	}

	return ctx.Status(http.StatusOK).JSON(EventIssuerResponse{
		ID:                 issuer.ID,
		EventID:            issuer.EventID,
		IssuerCredentialID: issuer.IssuerCredentialID,
		IsSigned:           issuer.IsSigned,
		CreatedAt:          issuer.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:          issuer.UpdatedAt.Time.Format(time.RFC3339),
	})
}
