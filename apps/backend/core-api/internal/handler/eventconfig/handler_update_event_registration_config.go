package eventconfig

import (
	"apps/backend/common/customerror"
	"apps/backend/common/pgmapper"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/eventconfig"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UpdateEventRegistrationConfigRequest struct {
	FinalCallForRegistration             *time.Time       `json:"final_call_for_registration,omitempty" valid:"omitempty"`
	RegistrationPassword                 *string          `json:"registration_password,omitempty" valid:"omitempty"`
	FirstNameRequirementStatus           int32            `json:"first_name_requirement_status" valid:"required,range(0|2)"`
	LastNameRequirementStatus            int32            `json:"last_name_requirement_status" valid:"required,range(0|2)"`
	EmailRequirementStatus               int32            `json:"email_requirement_status" valid:"required,range(0|2)"`
	BioRequirementStatus                 int32            `json:"bio_requirement_status" valid:"required,range(0|2)"`
	PhoneNumberRequirementStatus         int32            `json:"phone_number_requirement_status" valid:"required,range(0|2)"`
	AddressRequirementStatus             int32            `json:"address_requirement_status" valid:"required,range(0|2)"`
	AcademicInstitutionRequirementStatus int32            `json:"academic_institution_requirement_status" valid:"required,range(0|2)"`
	AcademicEmailRequirementStatus       int32            `json:"academic_email_requirement_status" valid:"required,range(0|2)"`
	IsBookingRequestRequired             bool             `json:"is_booking_request_required"`
	IsTicketTransferable                 bool             `json:"is_ticket_transferable"`
	EventType                            entity.EventType `json:"event_type" valid:"required,oneof=public private invite"`
}

func (r *UpdateEventRegistrationConfigRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

// UpdateEventRegistrationConfig godoc
// @Summary Update event registration config
// @Description Update the event registration configuration for an event
// @ID update-event-registration-config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body UpdateEventRegistrationConfigRequest true "Event registration config data"
// @Success 200 {object} EventRegistrationConfigViewModel
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/config/registration [put]
func (h *Handler) UpdateEventRegistrationConfig(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req UpdateEventRegistrationConfigRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	if err := req.IsValid(); err != nil {
		return err
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	params := eventconfig.UpdateEventRegistrationConfigParams{
		FinalCallForRegistration:             pgmapper.TimePtrToPgTimestampz(req.FinalCallForRegistration),
		RegistrationPassword:                 pgmapper.StringPtrToPgText(req.RegistrationPassword),
		FirstNameRequirementStatus:           pgmapper.Int32ToPgInt4(req.FirstNameRequirementStatus),
		LastNameRequirementStatus:            pgmapper.Int32ToPgInt4(req.LastNameRequirementStatus),
		EmailRequirementStatus:               pgmapper.Int32ToPgInt4(req.EmailRequirementStatus),
		BioRequirementStatus:                 pgmapper.Int32ToPgInt4(req.BioRequirementStatus),
		PhoneNumberRequirementStatus:         pgmapper.Int32ToPgInt4(req.PhoneNumberRequirementStatus),
		AddressRequirementStatus:             pgmapper.Int32ToPgInt4(req.AddressRequirementStatus),
		AcademicInstitutionRequirementStatus: pgmapper.Int32ToPgInt4(req.AcademicInstitutionRequirementStatus),
		AcademicEmailRequirementStatus:       pgmapper.Int32ToPgInt4(req.AcademicEmailRequirementStatus),
		IsBookingRequestRequired:             &req.IsBookingRequestRequired,
		IsTicketTransferable:                 &req.IsTicketTransferable,
		EventType:                            &req.EventType,
	}

	config, err := h.EventConfigUc.UpdateEventRegistrationConfig(ctx.UserContext(), eventID, params, currentUser)
	if err != nil {
		// Check if err is already a customerror type
		if customErr := customerror.TryParseAsCustomErr(err); customErr != nil {
			return customErr
		}
		// For non-custom errors, wrap as internal error
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.Status(http.StatusOK).JSON(EventRegistrationConfigViewModel{
		ID:                                   config.ID,
		EventID:                              config.EventID,
		FinalCallForRegistration:             config.FinalCallForRegistration,
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
