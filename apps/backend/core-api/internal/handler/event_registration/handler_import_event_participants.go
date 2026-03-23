package event_registration

import (
	"apps/backend/common/validatorutils"
	"fmt"

	customerror "apps/backend/common/customerror"

	eventRegistrationInvitationUc "apps/backend/core-api/internal/usecase/event_registration"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ImportEventParticipantsRequest struct {
	EventID      string                   `json:"event_id" validate:"required,uuid"`
	Participants []ParticipantRequestItem `json:"participants" validate:"required,min=1,dive"`
}

type ParticipantRequestItem struct {
	FirstName           string  `json:"first_name" validate:"required"`
	LastName            string  `json:"last_name" validate:"required"`
	Email               string  `json:"email" validate:"omitempty,email"`
	WalletAddress       *string `json:"wallet_address,omitempty"`
	PhoneNumber         *string `json:"phone_number,omitempty"`
	AcademicInstitution *string `json:"academic_institution,omitempty"`
}

// @Summary Import event participants
// @Description Import a list of participants for an event, creating inbox messages and event registration invitations
// @ID import-event-participants
// @Tags Event Registration Invitation
// @Accept json
// @Produce json
// @Param request body ImportEventParticipantsRequest true "Import participants request"
// @Success 201 {array} EventRegistrationInvitationResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/event-registration/invitation/:event_id/import [post]
func (h *Handler) ImportEventParticipants(ctx *fiber.Ctx) error {
	// 1. Parse and validate request
	requestBody := ImportEventParticipantsRequest{}
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

	// 4. Convert participants to usecase format
	participants := make([]eventRegistrationInvitationUc.Participant, len(requestBody.Participants))
	for i, participant := range requestBody.Participants {
		participants[i] = eventRegistrationInvitationUc.Participant{
			Name:                fmt.Sprintf("%s %s", participant.FirstName, participant.LastName),
			FirstName:           participant.FirstName,
			LastName:            participant.LastName,
			Email:               participant.Email,
			WalletAddress:       participant.WalletAddress,
			PhoneNumber:         participant.PhoneNumber,
			AcademicInstitution: participant.AcademicInstitution,
		}
	}

	// 5. Prepare parameters for usecase
	params := eventRegistrationInvitationUc.ImportEventParticipantsParameters{
		EventID:      eventID,
		Participants: participants,
	}

	// 6. Call usecase
	invitations, err := h.EventRegistrationUc.ImportEventParticipants(ctx.UserContext(), params, currentUser)
	if err != nil {
		return err
	}

	// 7. Return response
	return ctx.Status(fiber.StatusCreated).JSON(invitations)
}

// Parse - Parse JSON from request
func (r *ImportEventParticipantsRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

// IsValid - Validate request fields
func (r *ImportEventParticipantsRequest) IsValid() error {
	if err := validatorutils.ValidateStruct(r); err != nil {
		return err
	}
	for _, p := range r.Participants {
		hasEmail := p.Email != ""
		hasWallet := p.WalletAddress != nil && *p.WalletAddress != ""
		if hasEmail == hasWallet {
			return customerror.Parse(&customerror.ErrInvalidArgument,
				fmt.Errorf("each participant must have exactly one of email or wallet_address, not both and not neither"))
		}
	}
	return nil
}
