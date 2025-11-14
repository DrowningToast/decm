package event_registration_invitation

import (
	"fmt"

	customerror "apps/backend/common/customerror"
	event_registration_invitationUc "apps/backend/core-api/internal/usecase/event_registration_invitation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary Get event registration invitations by event ID
// @Description Get all event registration invitations for a specific event
// @ID get-event-registration-invitations-by-event-id
// @Tags Event Registration Invitation
// @Accept json
// @Produce json
// @Success 200 {array} entity.EventRegistrationInvitation
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/event-registration-invitations/{eventId} [get]
func (h *Handler) GetEventRegistrationInvitationsByEventId(ctx *fiber.Ctx) error {
	// 1. Get current user
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	// 2. Convert event_id from string to UUID
	eventId, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid event_id: %w", err))
	}

	// 4. Prepare parameters for usecase
	params := event_registration_invitationUc.GetEventRegistrationInvitationsByEventIDParameters{
		EventId: eventId,
	}

	// 5. Call usecase
	invitations, err := h.EventRegistrationInvitationUc.GetEventRegistrationInvitationsByEventId(ctx.UserContext(), params, currentUser)
	if err != nil {
		return err
	}

	// 6. Return response
	return ctx.JSON(invitations)
}
