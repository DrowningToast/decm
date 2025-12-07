package eventconfig

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"apps/backend/common/customerror"
)

// GetEventCertificateConfig godoc
// @Summary Get event certificate config
// @Description Get the event certificate configuration for an event. Accessible by verified organizers or issuers assigned to the event.
// @ID get-event-certificate-config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} EventCertificateConfigResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/config/certificate [get]
func (h *Handler) GetEventCertificateConfig(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Get current user from context
	currentUser, err := h.AuthenticationService.RequireUserContext(ctx)
	if err != nil {
		return err
	}

	// Check if user is a verified organizer
	isVerifiedOrganizer := currentUser.IsVerifiedOrganizer != nil && *currentUser.IsVerifiedOrganizer

	// If not a verified organizer, check if user is an issuer for this specific event
	if !isVerifiedOrganizer {
		_, err := h.EventUc.EventIssuerDataGateway.GetEventIssuerByEventIDAndIssuerCredentialID(
			ctx.UserContext(),
			eventID,
			currentUser.UserId,
		)
		if err != nil {
			// Check if it's a "not found" error (user is not an issuer)
			if errors.Is(err, pgx.ErrNoRows) {
				return customerror.Parse(
					&customerror.ErrForbidden,
					errors.New("user must be a verified organizer or an issuer assigned to this event"),
				)
			}
			// For other database errors, return internal server error
			return errors.Wrap(err, "failed to check if user is an issuer for this event")
		}
		// Note: SQL query now filters out soft-deleted records, so no need to check DeletedAt
	}

	config, err := h.EventConfigUc.GetEventCertificateConfigByEventID(ctx.UserContext(), eventID)
	if err != nil {
		return errors.Wrap(err, "failed to get event certificate config")
	}

	return ctx.Status(http.StatusOK).JSON(config)
}
