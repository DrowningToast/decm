package eventconfig

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

// GetEventRegistrationConfig godoc
// @Summary Get event registration config
// @Description Get event registration configuration for an event
// @Tags Events
// @ID get-event-registration-config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} EventRegistrationConfigResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/config/registration [get]
func (h *Handler) GetEventRegistrationConfig(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	config, err := h.EventConfigUc.GetEventRegistrationConfigByEventID(ctx.UserContext(), eventID)
	if err != nil {
		// Check if err is already a customerror type
		if customErr := customerror.TryParseAsCustomErr(err); customErr != nil {
			return customErr
		}
		// For non-custom errors, wrap as internal error
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.Status(http.StatusOK).JSON(EventRegistrationConfigResponse{
		ID:                                   config.ID,
		EventID:                              config.EventID,
		FinalCallForRegistration:             config.FinalCallForRegistration,
		RegistrationPassword:                 config.RegistrationPassword,
		FirstNameRequirementStatus:           EventRegistrationConfigRequirementStatus(config.FirstNameRequirementStatus),
		LastNameRequirementStatus:            EventRegistrationConfigRequirementStatus(config.LastNameRequirementStatus),
		EmailRequirementStatus:               EventRegistrationConfigRequirementStatus(config.EmailRequirementStatus),
		BioRequirementStatus:                 EventRegistrationConfigRequirementStatus(config.BioRequirementStatus),
		PhoneNumberRequirementStatus:         EventRegistrationConfigRequirementStatus(config.PhoneNumberRequirementStatus),
		AddressRequirementStatus:             EventRegistrationConfigRequirementStatus(config.AddressRequirementStatus),
		AcademicInstitutionRequirementStatus: EventRegistrationConfigRequirementStatus(config.AcademicInstitutionRequirementStatus),
		AcademicEmailRequirementStatus:       EventRegistrationConfigRequirementStatus(config.AcademicEmailRequirementStatus),
		CreatedAt:                            config.CreatedAt.String(),
		UpdatedAt:                            config.UpdatedAt.String(),
	})
}
