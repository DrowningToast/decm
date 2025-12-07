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

// @Summary Get issuer events view model
// @Description Get events assigned to the authenticated issuer for signing with view model data
// @ID get-issuer-events-viewmodel
// @Tags Issuer
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param issuer_credential_id query string false "Issuer credential ID"
// @Success 200 {array} issuer.IssuerEventViewModel
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/issuers/events/viewmodel [get]
func (h *Handler) GetIssuerEventsViewModel(ctx *fiber.Ctx) error {
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

	// Get events view model for issuer
	events, err := h.IssuerUc.GetIssuerEventsViewModel(ctxWithTimeout, issuerCredentialID, int32(limit), int32(offset))
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.Status(http.StatusOK).JSON(events)
}
