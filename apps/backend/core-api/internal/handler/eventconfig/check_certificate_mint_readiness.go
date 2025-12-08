package eventconfig

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

// CheckCertificateMintReadiness godoc
// @Summary Check certificate mint readiness
// @Description Check if an event certificate is ready to be minted. Returns detailed status about configuration, signed issuers, and contract deployment. Accessible by verified organizers or issuers assigned to the event.
// @ID check-certificate-mint-readiness
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} CertificateMintReadinessResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/config/certificate/mint-readiness [get]
func (h *Handler) CheckCertificateMintReadiness(ctx *fiber.Ctx) error {
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
		isIssuer, err := h.EventUc.IsUserIssuerForEvent(
			ctx.UserContext(),
			eventID,
			currentUser.UserId,
		)
		if err != nil {
			return err
		}
		if !isIssuer {
			return customerror.Parse(
				&customerror.ErrForbidden,
				errors.New("user must be a verified organizer or an issuer assigned to this event"),
			)
		}
	}

	// Check certificate mint readiness
	readiness, err := h.EventConfigUc.CheckCertificateMintReadiness(ctx.UserContext(), eventID)
	if err != nil {
		return errors.Wrap(err, "failed to check certificate mint readiness")
	}

	return ctx.Status(http.StatusOK).JSON(readiness)
}
