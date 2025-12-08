package event_registration

import (
	"errors"
	"fmt"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	eventRegistrationUc "apps/backend/core-api/internal/usecase/event_registration"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CancelEventRegistrationInvitationQuery struct {
	EventRegistrationInvitationID string `json:"event_registration_invitation_id" validate:"required,uuid"`
}

// @Summary Cancel event registration invitation
// @Description Cancel an event registration invitation by ID
// @ID cancel-event-registration-invitation
// @Tags Event Registration Invitation
// @Accept json
// @Produce json
// @Param event_registration_invitation_id query string true "Event Registration Invitation ID"
// @Success 200 {object} entity.EventRegistrationInvitation
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/event-registration/invitation [delete]
func (h *Handler) CancelEventRegistrationInvitation(ctx *fiber.Ctx) error {
	// 1. Parse and validate request
	query := CancelEventRegistrationInvitationQuery{}
	if err := query.Parse(ctx); err != nil {
		return err
	}
	if err := query.IsValid(); err != nil {
		return err
	}

	// 2. Get current user
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	// 3. Convert event_registration_invitation_id from string to UUID
	eventRegistrationInvitationID, err := uuid.Parse(query.EventRegistrationInvitationID)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid event_registration_invitation_id: %w", err))
	}

	// 4. Prepare parameters for usecase
	params := eventRegistrationUc.CancelEventRegistrationInvitationParameters{
		EventRegistrationInvitationID: eventRegistrationInvitationID,
	}

	// 5. Call usecase
	invitation, err := h.EventRegistrationUc.CancelEventRegistrationInvitation(ctx.UserContext(), params, currentUser)
	if err != nil {
		return err
	}

	// 6. Return response
	return ctx.JSON(invitation)
}

// Parse - Parse path parameter from request
func (r *CancelEventRegistrationInvitationQuery) Parse(ctx *fiber.Ctx) error {
	// Get event_registration_invitation_id from path parameter
	eventRegistrationInvitationID := ctx.Query("event_registration_invitation_id")
	if eventRegistrationInvitationID == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("event_registration_invitation_id is required"))
	}
	r.EventRegistrationInvitationID = eventRegistrationInvitationID
	return nil
}

// IsValid - Validate request fields
func (r *CancelEventRegistrationInvitationQuery) IsValid() error {
	return validatorutils.ValidateStruct(r)
}
