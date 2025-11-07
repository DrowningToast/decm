package event

import (
	"fmt"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/validatorutils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	eventUc "apps/backend/core-api/internal/usecase/event"
)

type GetEventRegistrationInvitationsByEventIDRequest struct {
	EventID string `json:"event_id" validate:"required,uuid"`
}

// @Summary Get event registration invitations by event ID
// @Description Get all event registration invitations for a specific event
// @ID get-event-registration-invitations-by-event-id
// @Tags Event Registration Invitation
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {array} entity.EventRegistrationInvitation
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{eventId}/registration/invitations [get]
func (h *Handler) GetEventRegistrationInvitationsByEventID(ctx *fiber.Ctx) error {
	// 1. Parse and validate request
	requestBody := GetEventRegistrationInvitationsByEventIDRequest{}
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

	// 3. Convert event_id from string to UUID
	eventID, err := uuid.Parse(requestBody.EventID)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid event_id: %w", err))
	}

	// 4. Prepare parameters for usecase
	params := eventUc.GetEventRegistrationInvitationsByEventIDParameters{
		EventID: eventID,
	}

	// 5. Call usecase
	invitations, err := h.EventUc.GetEventRegistrationInvitationsByEventID(ctx.UserContext(), params, currentUser)
	if err != nil {
		return err
	}

	// 6. Return response
	return ctx.JSON(invitations)
}

// Parse - Parse path parameter from request
func (r *GetEventRegistrationInvitationsByEventIDRequest) Parse(ctx *fiber.Ctx) error {
	// Get event_id from path parameter
	eventID := ctx.Params("eventId")
	if eventID == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("event_id is required"))
	}
	r.EventID = eventID
	return nil
}

// IsValid - Validate request fields
func (r *GetEventRegistrationInvitationsByEventIDRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}
