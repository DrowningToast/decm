package event_registration_invitation

import (
	"errors"
	"fmt"

	customerror "apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetEventRegistrationInvitationByUserAndEventResponse struct {
	RegistrationInvitation *entity.EventRegistrationInvitation `json:"registration_invitation,omitempty"`
	Inbox                  *entity.InboxMessage                `json:"inbox,omitempty"`
}

// @Summary Get event registration invitation of user and by event id
// @Description Get event registration invitation of user and by event id
// @ID get-event-registration-invitation-by-user-and-event
// @Tags Event Registration Invitation
// @Accept json
// @Produce json
// @Success 200 {object} GetEventRegistrationInvitationByUserAndEventResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/event-registration-invitations/my/{event_id} [get]
func (h *Handler) GetEventRegistrationInvitationByUserAndEvent(ctx *fiber.Ctx) error {
	// 1. Convert event_id from string to UUID

	eventIdStr := ctx.Params("event_id")
	if eventIdStr == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("event_id is required"))
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid event_id: %w", err))
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}
	if currentUser == nil {
		return customerror.Parse(&customerror.ErrUnauthorized, errors.New("not logged in"))
	}

	// 2. Prepare parameters for usecase
	registrationInvitation, inbox, err := h.EventRegistrationInvitationUc.GetEventRegistrationByUserAndEvent(ctx.UserContext(), eventId, currentUser)
	if err != nil {
		return err
	}

	// 4. Return response
	return ctx.JSON(GetEventRegistrationInvitationByUserAndEventResponse{
		RegistrationInvitation: registrationInvitation,
		Inbox:                  inbox,
	})
}
