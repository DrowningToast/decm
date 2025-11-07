package issuer

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	customerror "apps/backend/common/customerror"
)

type IssuerEventResponse struct {
	ID                     string    `json:"id"`
	EventID                string    `json:"event_id"`
	EventTitle             string    `json:"event_title"`
	EventShortDescription  string    `json:"event_short_description"`
	EventStartDate         time.Time `json:"event_start_date"`
	EventEndDate           time.Time `json:"event_end_date"`
	EventLocation          string    `json:"event_location"`
	EventOwnerCredentialID string    `json:"event_owner_credential_id"`
	IssuerCredentialID     string    `json:"issuer_credential_id"`
	IsSigned               int32     `json:"is_signed"`
	Signature              string    `json:"signature"`
	SignMessage            string    `json:"sign_message"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// @Summary Get events for issuer signing
// @Description Get events assigned to the authenticated issuer for signing
// @ID get-issuer-events
// @Tags Issuer
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param issuer_credential_id query string false "Issuer credential ID"
// @Success 200 {array} IssuerEventResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/issuers/events [get]
func (h *Handler) GetIssuerEvents(ctx *fiber.Ctx) error {
	// Parse query parameters
	limitStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")
	issuerCredentialID := ctx.Query("issuer_credential_id")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	if issuerCredentialID == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("issuer credential ID is required"))
	}

	// Add timeout for usecase call
	ctxWithTimeout, cancel := context.WithTimeout(ctx.UserContext(), 30*time.Second)
	defer cancel()

	// Get events for issuer
	events, err := h.IssuerUc.GetEventsByIssuerCredentialID(ctxWithTimeout, issuerCredentialID, int32(limit), int32(offset))
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Convert to response format
	response := make([]IssuerEventResponse, len(events))
	for i, event := range events {
		response[i] = IssuerEventResponse{
			ID:                     event.ID.String(),
			EventID:                event.EventID.String(),
			EventTitle:             event.EventTitle,
			EventShortDescription:  event.EventShortDescription,
			EventStartDate:         event.CreatedAt.Time,
			EventEndDate:           event.UpdatedAt.Time,
			EventLocation:          event.EventLocation,
			EventOwnerCredentialID: event.EventOwnerCredentialID.String(),
			IssuerCredentialID:     event.IssuerCredentialID.String(),
			IsSigned:               event.IsSigned,
			Signature:              event.Signature.String,
			SignMessage:            event.SignMessage.String,
			CreatedAt:              event.CreatedAt.Time,
			UpdatedAt:              event.UpdatedAt.Time,
		}
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
