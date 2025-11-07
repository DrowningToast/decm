package event_registration_invitation

import (
	"fmt"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	eventRegistrationInvitationUc "apps/backend/core-api/internal/usecase/event_registration_invitation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CancelEventRegistrationInvitationRequest struct {
	EventRegistrationInvitationID string `json:"event_registration_invitation_id" validate:"required,uuid"`
}

// @Summary Cancel event registration invitation
// @Description Cancel an event registration invitation by ID
// @ID cancel-event-registration-invitation
// @Tags Event Registration Invitation
// @Accept json
// @Produce json
// @Param eventRegistrationInvitationID path string true "Event Registration Invitation ID"
// @Success 200 {object} entity.EventRegistrationInvitation
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/event-registration-invitations/{eventRegistrationInvitationId} [delete]
func (h *Handler) CancelEventRegistrationInvitation(ctx *fiber.Ctx) error {
	// 1. Parse and validate request
	requestBody := CancelEventRegistrationInvitationRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	// 2. Get current user
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	// 3. Convert event_registration_invitation_id from string to UUID
	eventRegistrationInvitationID, err := uuid.Parse(requestBody.EventRegistrationInvitationID)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid event_registration_invitation_id: %w", err))
	}

	// 4. Prepare parameters for usecase
	params := eventRegistrationInvitationUc.CancelEventRegistrationInvitationParameters{
		EventRegistrationInvitationID: eventRegistrationInvitationID,
	}

	// 5. Call usecase
	invitation, err := h.EventRegistrationInvitationUc.CancelEventRegistrationInvitation(ctx.UserContext(), params, currentUser)
	if err != nil {
		return err
	}

	// 6. Return response
	return ctx.JSON(invitation)
}

// Parse - Parse path parameter from request
func (r *CancelEventRegistrationInvitationRequest) Parse(ctx *fiber.Ctx) error {
	// Get event_registration_invitation_id from path parameter
	eventRegistrationInvitationID := ctx.Params("eventRegistrationInvitationId")
	if eventRegistrationInvitationID == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("event_registration_invitation_id is required"))
	}
	r.EventRegistrationInvitationID = eventRegistrationInvitationID
	return nil
}

// IsValid - Validate request fields
func (r *CancelEventRegistrationInvitationRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}
