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

// Request/Response structures
type CreateEventCertificateConfigRequest struct {
	BaseCertificateStorageKey string `json:"base_certificate_storage_key"`
	EventNamePosX             int32  `json:"event_name_pos_x"`
	EventNamePosY             int32  `json:"event_name_pos_y"`
	NamePosX                  int32  `json:"name_pos_x"`
	NamePosY                  int32  `json:"name_pos_y"`
	AcademicInstitutionPosX   int32  `json:"academic_institution_pos_x"`
	AcademicInstitutionPosY   int32  `json:"academic_institution_pos_y"`
}

func (r *CreateEventCertificateConfigRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

// CreateEventCertificateConfig godoc
// @Summary Create event certificate config
// @Description Create a new event certificate configuration for an event
// @ID create-event-certificate-config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body CreateEventCertificateConfigRequest true "Event certificate config data"
// @Success 200 {object} EventCertificateConfigResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/certificate-config [post]
func (h *Handler) CreateEventCertificateConfig(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req CreateEventCertificateConfigRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	params := eventconfig.CreateEventCertificateConfigParams{
		BaseCertificateStorageKey: req.BaseCertificateStorageKey,
		EventNamePosX:             req.EventNamePosX,
		EventNamePosY:             req.EventNamePosY,
		NamePosX:                  req.NamePosX,
		NamePosY:                  req.NamePosY,
		AcademicInstitutionPosX:   pgtype.Int4{Int32: req.AcademicInstitutionPosX, Valid: req.AcademicInstitutionPosX != 0},
		AcademicInstitutionPosY:   pgtype.Int4{Int32: req.AcademicInstitutionPosY, Valid: req.AcademicInstitutionPosY != 0},
	}

	config, err := h.EventConfigUc.CreateEventCertificateConfig(ctx.UserContext(), eventID, params)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.Status(http.StatusOK).JSON(EventCertificateConfigResponse{
		ID:                        config.ID,
		EventID:                   config.EventID,
		BaseCertificateStorageKey: config.BaseCertificateStorageKey,
		EventNamePosX:             config.EventNamePosX,
		EventNamePosY:             config.EventNamePosY,
		NamePosX:                  config.NamePosX,
		NamePosY:                  config.NamePosY,
		AcademicInstitutionPosX:   config.AcademicInstitutionPosX.Int32,
		AcademicInstitutionPosY:   config.AcademicInstitutionPosY.Int32,
		CreatedAt:                 config.CreatedAt.Time.String(),
		UpdatedAt:                 config.UpdatedAt.Time.String(),
	})
}
