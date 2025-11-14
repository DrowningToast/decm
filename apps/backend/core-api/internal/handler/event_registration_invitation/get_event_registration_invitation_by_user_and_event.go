package event_registration_invitation

import (
	"errors"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/entity"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetEventRegistrationInvitationByUserAndEventRequest struct {
	EventId uuid.UUID `json:"event_id" validate:"required,uuid"`
}

func (r *GetEventRegistrationInvitationByUserAndEventRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.ParamsParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func (r *GetEventRegistrationInvitationByUserAndEventRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

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
	// 1. Get current user
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}
	if currentUser == nil {
		return customerror.Parse(&customerror.ErrUnauthorized, errors.New("not logged in"))
	}

	// 2. Parse and validate request
	requestBody := GetEventRegistrationInvitationByUserAndEventRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	// 3. Prepare parameters for usecase
	registrationInvitation, inbox, err := h.EventRegistrationInvitationUc.GetEventRegistrationByUserAndEvent(ctx.UserContext(), requestBody.EventId, currentUser)
	if err != nil {
		return err
	}

	// 4. Return response
	return ctx.JSON(GetEventRegistrationInvitationByUserAndEventResponse{
		RegistrationInvitation: registrationInvitation,
		Inbox:                  inbox,
	})
}
