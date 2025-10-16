package eventconfig

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/usecase/eventconfig"
)

type UpdateEventRegistrationConfigRequest struct {
	FirstNameRequirementStatus           int32 `json:"first_name_requirement_status"`
	LastNameRequirementStatus            int32 `json:"last_name_requirement_status"`
	EmailRequirementStatus               int32 `json:"email_requirement_status"`
	BioRequirementStatus                 int32 `json:"bio_requirement_status"`
	PhoneNumberRequirementStatus         int32 `json:"phone_number_requirement_status"`
	AddressRequirementStatus             int32 `json:"address_requirement_status"`
	AcademicInstitutionRequirementStatus int32 `json:"academic_institution_requirement_status"`
	AcademicEmailRequirementStatus       int32 `json:"academic_email_requirement_status"`
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
// @Router /api/v1/events/{event_id}/registration-config [put]
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

	params := eventconfig.UpdateEventRegistrationConfigParams{
		FirstNameRequirementStatus:           pgtype.Int4{Int32: req.FirstNameRequirementStatus, Valid: true},
		LastNameRequirementStatus:            pgtype.Int4{Int32: req.LastNameRequirementStatus, Valid: true},
		EmailRequirementStatus:               pgtype.Int4{Int32: req.EmailRequirementStatus, Valid: true},
		BioRequirementStatus:                 pgtype.Int4{Int32: req.BioRequirementStatus, Valid: true},
		PhoneNumberRequirementStatus:         pgtype.Int4{Int32: req.PhoneNumberRequirementStatus, Valid: true},
		AddressRequirementStatus:             pgtype.Int4{Int32: req.AddressRequirementStatus, Valid: true},
		AcademicInstitutionRequirementStatus: pgtype.Int4{Int32: req.AcademicInstitutionRequirementStatus, Valid: true},
		AcademicEmailRequirementStatus:       pgtype.Int4{Int32: req.AcademicEmailRequirementStatus, Valid: true},
	}

	config, err := h.EventConfigUc.UpdateEventRegistrationConfig(ctx.UserContext(), eventID, params)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.Status(http.StatusOK).JSON(EventRegistrationConfigResponse{
		ID:                                   config.ID,
		EventID:                              config.EventID,
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
