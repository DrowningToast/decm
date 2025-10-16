package eventconfig

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
)

// GetEventCertificateConfig godoc
// @Summary Get event certificate config
// @Description Get the event certificate configuration for an event
// @ID get-event-certificate-config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} EventCertificateConfigResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/certificate-config [get]
func (h *Handler) GetEventCertificateConfig(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	config, err := h.EventConfigUc.GetEventCertificateConfigByEventID(ctx.UserContext(), eventID)
	if err != nil {
		return customerror.Parse(&customerror.ErrNotFound, err)
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
