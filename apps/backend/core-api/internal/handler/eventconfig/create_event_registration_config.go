package eventconfig

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/usecase/eventconfig"
)

// Request/Response structures
type CreateEventRegistrationConfigRequest struct {
	FinalCallForRegistration             *string `json:"final_call_for_registration,omitempty"`
	RegistrationPassword                 *string `json:"registration_password,omitempty"`
	FirstNameRequirementStatus           int32   `json:"first_name_requirement_status"`
	LastNameRequirementStatus            int32   `json:"last_name_requirement_status"`
	EmailRequirementStatus               int32   `json:"email_requirement_status"`
	BioRequirementStatus                 int32   `json:"bio_requirement_status"`
	PhoneNumberRequirementStatus         int32   `json:"phone_number_requirement_status"`
	AddressRequirementStatus             int32   `json:"address_requirement_status"`
	AcademicInstitutionRequirementStatus int32   `json:"academic_institution_requirement_status"`
	AcademicEmailRequirementStatus       int32   `json:"academic_email_requirement_status"`
}

func (r *CreateEventRegistrationConfigRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

// CreateEventRegistrationConfig godoc
// @Summary Create event registration config
// @Description Create a new event registration configuration for an event
// @ID create-event-registration-config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body CreateEventRegistrationConfigRequest true "Event registration config data"
// @Success 200 {object} EventRegistrationConfigResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/registration-config [post]
func (h *Handler) CreateEventRegistrationConfig(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req CreateEventRegistrationConfigRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	if err := req.IsValid(); err != nil {
		return err
	}

	// Parse final_call_for_registration if provided
	var finalCallForRegistration pgtype.Timestamptz
	if req.FinalCallForRegistration != nil && *req.FinalCallForRegistration != "" {
		parsedTime, err := time.Parse(time.RFC3339, *req.FinalCallForRegistration)
		if err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, err)
		}
		finalCallForRegistration = pgtype.Timestamptz{Time: parsedTime, Valid: true}
	}

	// Handle registration_password
	var registrationPassword pgtype.Text
	if req.RegistrationPassword != nil && *req.RegistrationPassword != "" {
		registrationPassword = pgtype.Text{String: *req.RegistrationPassword, Valid: true}
	}

	params := eventconfig.CreateEventRegistrationConfigParams{
		FinalCallForRegistration:             finalCallForRegistration,
		RegistrationPassword:                 registrationPassword,
		FirstNameRequirementStatus:           pgtype.Int4{Int32: req.FirstNameRequirementStatus, Valid: true},
		LastNameRequirementStatus:            pgtype.Int4{Int32: req.LastNameRequirementStatus, Valid: true},
		EmailRequirementStatus:               pgtype.Int4{Int32: req.EmailRequirementStatus, Valid: true},
		BioRequirementStatus:                 pgtype.Int4{Int32: req.BioRequirementStatus, Valid: true},
		PhoneNumberRequirementStatus:         pgtype.Int4{Int32: req.PhoneNumberRequirementStatus, Valid: true},
		AddressRequirementStatus:             pgtype.Int4{Int32: req.AddressRequirementStatus, Valid: true},
		AcademicInstitutionRequirementStatus: pgtype.Int4{Int32: req.AcademicInstitutionRequirementStatus, Valid: true},
		AcademicEmailRequirementStatus:       pgtype.Int4{Int32: req.AcademicEmailRequirementStatus, Valid: true},
	}

	config, err := h.EventConfigUc.CreateEventRegistrationConfig(ctx.UserContext(), eventID, params)
	if err != nil {
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
