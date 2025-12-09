package event

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

// GetEventIssuersByEventID godoc
// @Summary Get event issuers by event ID
// @Description Get all event issuers for an event
// @Tags Events
// @ID get-event-issuers-by-event-id
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {array} EventIssuerResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/issuers [get]
func (h *Handler) GetEventIssuersByEventID(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.UserContext(), 30*time.Second)
	defer cancel()

	issuers, err := h.EventUc.GetEventIssuersByEventID(ctxWithTimeout, eventID)
	if err != nil {
		// Check if err is already a customerror type
		var customErr *customerror.Err
		if errors.As(err, &customErr) {
			return customErr
		}
		// For non-custom errors, wrap as internal error
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	var response []EventIssuerResponse
	for _, issuer := range issuers {
		issuerCredentialId := issuer.IssuerCredentialID
		issuerProfile, authCredential, err := h.ProfileUc.GetProfileAndCredentialWithCredentialId(ctxWithTimeout, issuerCredentialId)
		if err != nil {
			// Check if err is already a customerror type
			var customErr *customerror.Err
			if errors.As(err, &customErr) {
				return customErr
			}
			// For non-custom errors, wrap as not found error
			return customerror.Parse(&customerror.ErrNotFound, err)
		}

		// Add connector references and wallet address to profile
		issuerProfile.GoogleConnectorRef = authCredential.GoogleConnectorRef
		issuerProfile.GithubConnectorRef = authCredential.GithubConnectorRef
		issuerProfile.WalletAddress = authCredential.WalletAddress

		response = append(response, EventIssuerResponse{
			ID:                 issuer.ID,
			EventID:            issuer.EventID,
			IssuerCredentialID: issuer.IssuerCredentialID,
			IsSigned:           issuer.IsSigned,
			CreatedAt:          issuer.CreatedAt.Time.String(),
			UpdatedAt:          issuer.UpdatedAt.Time.String(),
			IssuerProfile:      issuerProfile,
		})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
