package event

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	eventUc "apps/backend/core-api/internal/usecase/event"
)

type UpdateEventIssuerRequest struct {
	IsSigned    int32  `json:"is_signed"`
	Signature   string `json:"signature"`
	SignMessage string `json:"sign_message"`
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
// @Param issuer_id path string true "Issuer ID"
// @Param request body UpdateEventIssuerRequest true "Event issuer data"
// @Success 200 {object} EventIssuerResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/issuers/{issuer_id} [put]
func (h *Handler) UpdateEventIssuer(ctx *fiber.Ctx) error {
	issuerID, err := uuid.Parse(ctx.Params("issuer_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req UpdateEventIssuerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	if err := req.IsValid(); err != nil {
		return err
	}

	params := eventUc.UpdateEventIssuerParams{
		IsSigned:    req.IsSigned,
		Signature:   pgtype.Text{String: req.Signature, Valid: req.Signature != ""},
		SignMessage: pgtype.Text{String: req.SignMessage, Valid: req.SignMessage != ""},
	}

	issuer, err := h.EventUc.UpdateEventIssuer(ctx.UserContext(), issuerID, params)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(EventIssuerResponse{
		ID:                 issuer.ID,
		EventID:            issuer.EventID,
		IssuerCredentialID: issuer.IssuerCredentialID,
		IsSigned:           issuer.IsSigned,
		Signature:          issuer.Signature.String,
		SignMessage:        issuer.SignMessage.String,
		CreatedAt:          issuer.CreatedAt.Time.String(),
		UpdatedAt:          issuer.UpdatedAt.Time.String(),
	})
}
