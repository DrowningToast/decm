package event

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	eventUc "apps/backend/core-api/internal/usecase/event"
)

type CreateEventIssuerRequest struct {
	IssuerCredentialID uuid.UUID `json:"issuer_credential_id" validate:"required,uuid"`
	IsSigned           int32     `json:"is_signed" validate:"required,oneof=0 1"`
	Signature          string    `json:"signature" validate:"required_if=IsSigned 1"`
	SignMessage        string    `json:"sign_message"`
}

func (r *CreateEventIssuerRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
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
func (h *Handler) CreateEventIssuer(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req CreateEventIssuerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	if err := req.IsValid(); err != nil {
		// Ensure validation errors are properly wrapped
		var validationErr *validator.ValidationErrors
		if errors.As(err, &validationErr) {
			return customerror.ParseValidationErr(validationErr)
		}
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	params := eventUc.CreateEventIssuerParams{
		EventID:            eventID,
		IssuerCredentialID: req.IssuerCredentialID,
		IsSigned:           req.IsSigned,
		Signature:          pgtype.Text{String: req.Signature, Valid: req.Signature != ""},
		SignMessage:        pgtype.Text{String: req.SignMessage, Valid: req.SignMessage != ""},
	}

	issuer, err := h.EventUc.CreateEventIssuer(ctx.UserContext(), params)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(EventIssuerResponse{
		ID:                 issuer.ID,
		EventID:            issuer.EventID,
		IssuerCredentialID: issuer.IssuerCredentialID,
		IsSigned:           issuer.IsSigned,
		Signature:          issuer.Signature.String,
		SignMessage:        "",
		CreatedAt:          issuer.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:          issuer.UpdatedAt.Time.Format(time.RFC3339),
	})
}
