package eventconfig

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
	"apps/backend/common/pgmapper"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/eventconfig"
)

type UpdateEventRegistrationConfigRequest struct {
	FinalCallForRegistration             *time.Time       `json:"final_call_for_registration,omitempty"`
	RegistrationPassword                 *string          `json:"registration_password,omitempty"`
	FirstNameRequirementStatus           int32            `json:"first_name_requirement_status"`
	LastNameRequirementStatus            int32            `json:"last_name_requirement_status"`
	EmailRequirementStatus               int32            `json:"email_requirement_status"`
	BioRequirementStatus                 int32            `json:"bio_requirement_status"`
	PhoneNumberRequirementStatus         int32            `json:"phone_number_requirement_status"`
	AddressRequirementStatus             int32            `json:"address_requirement_status"`
	AcademicInstitutionRequirementStatus int32            `json:"academic_institution_requirement_status"`
	AcademicEmailRequirementStatus       int32            `json:"academic_email_requirement_status"`
	IsBookingRequestRequired             bool             `json:"is_booking_request_required"`
	IsTicketTransferable                 bool             `json:"is_ticket_transferable"`
	EventType                            entity.EventType `json:"event_type,omitempty"`
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
// @Success 200 {object} EventRegistrationConfigResponse
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
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Prepare response with nullable fields
	var finalCallForRegistrationResp *time.Time
	if config.FinalCallForRegistration.Valid {
		finalCallForRegistrationResp = &config.FinalCallForRegistration.Time
	}

	var registrationPasswordResp *string
	if config.RegistrationPassword.Valid {
		registrationPasswordResp = &config.RegistrationPassword.String
	}

	return ctx.Status(http.StatusOK).JSON(EventRegistrationConfigResponse{
		ID:                                   config.ID,
		EventID:                              config.EventID,
		FinalCallForRegistration:             finalCallForRegistrationResp,
		RegistrationPassword:                 registrationPasswordResp,
		FirstNameRequirementStatus:           config.FirstNameRequirementStatus.Int32,
		LastNameRequirementStatus:            config.LastNameRequirementStatus.Int32,
		EmailRequirementStatus:               config.EmailRequirementStatus.Int32,
		BioRequirementStatus:                 config.BioRequirementStatus.Int32,
		PhoneNumberRequirementStatus:         config.PhoneNumberRequirementStatus.Int32,
		AddressRequirementStatus:             config.AddressRequirementStatus.Int32,
		AcademicInstitutionRequirementStatus: config.AcademicInstitutionRequirementStatus.Int32,
		AcademicEmailRequirementStatus:       config.AcademicEmailRequirementStatus.Int32,
		CreatedAt:                            config.CreatedAt.Time.String(),
		UpdatedAt:                            config.UpdatedAt.Time.String(),
	})
}
