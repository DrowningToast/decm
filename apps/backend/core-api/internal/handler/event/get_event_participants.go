package event

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

type EventParticipantResponse struct {
	ID                  uuid.UUID  `json:"id"`
	EventID             uuid.UUID  `json:"event_id"`
	AttendeeCredentialID uuid.UUID `json:"attendee_credential_id"`
	WalletAddress       string     `json:"wallet_address"`
	ContractAddress     string     `json:"contract_address"`
	IsAttendeeAccepted  bool       `json:"is_attendee_accepted"`
	FirstName           *string    `json:"first_name,omitempty"`
	LastName            *string    `json:"last_name,omitempty"`
	Email               *string    `json:"email,omitempty"`
	PhoneNumber         *string    `json:"phone_number,omitempty"`
	AcademicInstitution *string    `json:"academic_institution,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// GetEventParticipants godoc
// @Summary Get event participants
// @Description Get all participants (attendees) for an event
// @Tags Events
// @ID get-event-participants
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {array} EventParticipantResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/participants [get]
func (h *Handler) GetEventParticipants(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.UserContext(), 30*time.Second)
	defer cancel()

	attendees, err := h.EventUc.ListEventAttendees(ctxWithTimeout, eventID)
	if err != nil {
		// Check if err is already a customerror type
		var customErr *customerror.Err
		if errors.As(err, &customErr) {
			return customErr
		}
		// For non-custom errors, wrap as internal error
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	var response []EventParticipantResponse
	for _, attendee := range attendees {
		response = append(response, EventParticipantResponse{
			ID:                   attendee.Id,
			EventID:              attendee.EventId,
			AttendeeCredentialID: attendee.AttendeeCredentialId,
			WalletAddress:        attendee.WalletAddress,
			ContractAddress:      attendee.ContractAddress,
			IsAttendeeAccepted:   attendee.IsAttendeeAccepted,
			FirstName:            attendee.FirstName,
			LastName:             attendee.LastName,
			Email:                attendee.Email,
			PhoneNumber:          attendee.PhoneNumber,
			AcademicInstitution:  attendee.AcademicInstitution,
			CreatedAt:            attendee.CreatedAt,
			UpdatedAt:            attendee.UpdatedAt,
		})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
