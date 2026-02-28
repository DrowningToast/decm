package eventconfig

import (
	"apps/backend/common/customerror"
	"io"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetEventCertificateTemplate godoc
// @Summary Get event certificate template SVG
// @Description Proxy the base certificate SVG template from storage so the frontend can access it without CORS restrictions from presigned URLs. Accessible by the event owner or issuers assigned to the event.
// @ID get-event-certificate-template
// @Produce image/svg+xml
// @Param event_id path string true "Event ID"
// @Success 200 {string} string "SVG content"
// @Failure 400 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/config/certificate/template [get]
func (h *Handler) GetEventCertificateTemplate(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Get current user from context
	currentUser, err := h.AuthenticationService.RequireUserContext(ctx)
	if err != nil {
		return err
	}

	// Verify the caller is the event owner or an assigned issuer
	event, err := h.EventUc.GetEventById(ctx.UserContext(), eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return customerror.Parse(&customerror.ErrNotFound, errors.New("event not found"))
	}

	if currentUser.UserId != event.OwnerCredentialId {
		isIssuer, err := h.EventUc.IsUserIssuerForEvent(ctx.UserContext(), eventID, currentUser.UserId)
		if err != nil {
			return err
		}
		if !isIssuer {
			return customerror.Parse(
				&customerror.ErrForbidden,
				errors.New("user must be the event owner or an issuer assigned to this event"),
			)
		}
	}

	svgReader, err := h.EventConfigUc.GetEventCertificateTemplateSVG(ctx.UserContext(), eventID)
	if err != nil {
		return errors.Wrap(err, "failed to get event certificate template")
	}
	defer svgReader.Close()

	svgBytes, err := io.ReadAll(svgReader)
	if err != nil {
		return errors.Wrap(err, "failed to read event certificate template")
	}

	ctx.Set("Content-Type", "image/svg+xml")
	ctx.Set("Cache-Control", "private, max-age=300")
	return ctx.Status(http.StatusOK).Send(svgBytes)
}
