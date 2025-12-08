package eventconfig

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/usecase/eventconfig"
)

type UpdateEventCertificateConfigRequest struct {
	EventNamePosX           *float64 `form:"event_name_pos_x"`
	EventNamePosY           *float64 `form:"event_name_pos_y"`
	NamePosX                *float64 `form:"name_pos_x"`
	NamePosY                *float64 `form:"name_pos_y"`
	AcademicInstitutionPosX *float64 `form:"academic_institution_pos_x"`
	AcademicInstitutionPosY *float64 `form:"academic_institution_pos_y"`
	CertificateTitlePosX    *float64 `form:"certificate_title_pos_x"`
	CertificateTitlePosY    *float64 `form:"certificate_title_pos_y"`
	CertificateSubtitlePosX *float64 `form:"certificate_subtitle_pos_x"`
	CertificateSubtitlePosY *float64 `form:"certificate_subtitle_pos_y"`
}

func (r *UpdateEventCertificateConfigRequest) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

func (r *UpdateEventCertificateConfigRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

// UpdateEventCertificateConfig godoc
// @Summary Update event certificate config
// @Description Update the event certificate configuration for an event
// @ID update-event-certificate-config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param base_certificate_image formData file false "Base certificate image"
// @Param event_name_pos_x formData float64 false "Event name position x"
// @Param event_name_pos_y formData float64 false "Event name position y"
// @Param name_pos_x formData float64 true "Name position x"
// @Param name_pos_y formData float64 true "Name position y"
// @Param academic_institution_pos_x formData float64 false "Academic institution position x"
// @Param academic_institution_pos_y formData float64 false "Academic institution position y"
// @Param certificate_title_pos_x formData float64 false "Certificate title position x"
// @Param certificate_title_pos_y formData float64 false "Certificate title position y"
// @Param certificate_subtitle_pos_x formData float64 false "Certificate subtitle position x"
// @Param certificate_subtitle_pos_y formData float64 false "Certificate subtitle position y"
// @Success 200 {object} EventCertificateConfigResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/config/certificate [put]
func (h *Handler) UpdateEventCertificateConfig(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	requestBody := UpdateEventCertificateConfigRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	if err := requestBody.IsValid(); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	params := eventconfig.UpdateEventCertificateConfigParams{}

	if requestBody.EventNamePosX != nil {
		params.EventNamePosX = requestBody.EventNamePosX
	}
	if requestBody.EventNamePosY != nil {
		params.EventNamePosY = requestBody.EventNamePosY
	}
	if requestBody.NamePosX != nil {
		params.NamePosX = requestBody.NamePosX
	}
	if requestBody.NamePosY != nil {
		params.NamePosY = requestBody.NamePosY
	}
	if requestBody.AcademicInstitutionPosX != nil {
		params.AcademicInstitutionPosX = requestBody.AcademicInstitutionPosX
	}
	if requestBody.AcademicInstitutionPosY != nil {
		params.AcademicInstitutionPosY = requestBody.AcademicInstitutionPosY
	}
	if requestBody.CertificateTitlePosX != nil {
		params.CertificateTitlePosX = requestBody.CertificateTitlePosX
	}
	if requestBody.CertificateTitlePosY != nil {
		params.CertificateTitlePosY = requestBody.CertificateTitlePosY
	}
	if requestBody.CertificateSubtitlePosX != nil {
		params.CertificateSubtitlePosX = requestBody.CertificateSubtitlePosX
	}
	if requestBody.CertificateSubtitlePosY != nil {
		params.CertificateSubtitlePosY = requestBody.CertificateSubtitlePosY
	}

	baseCertificateImage, _ := ctx.FormFile("base_certificate_image")
	if baseCertificateImage != nil {
		if err := validatorutils.ValidateImageFile(baseCertificateImage); err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, err)
		}
		params.BaseCertificateImage = baseCertificateImage
	}

	dbEventCertConfig, _ := h.EventConfigUc.GetEventCertificateConfigByEventID(ctx.UserContext(), eventID)
	if dbEventCertConfig == nil {
		_, err = h.EventConfigUc.CreateEventCertificateConfig(ctx.UserContext(), eventID, eventconfig.CreateEventCertificateConfigParams{
			BaseCertificateImage:    *params.BaseCertificateImage,
			EventNamePosX:           *params.EventNamePosX,
			EventNamePosY:           *params.EventNamePosY,
			NamePosX:                *params.NamePosX,
			NamePosY:                *params.NamePosX,
			AcademicInstitutionPosX: params.AcademicInstitutionPosX,
			AcademicInstitutionPosY: params.AcademicInstitutionPosY,
			CertificateTitlePosX:    params.CertificateTitlePosX,
			CertificateTitlePosY:    params.CertificateTitlePosY,
			CertificateSubtitlePosX: params.CertificateSubtitlePosX,
			CertificateSubtitlePosY: params.CertificateSubtitlePosY,
		})
		if err != nil {
			return customerror.Parse(&customerror.ErrInternalServer, err)
		}
	} else {
		_, err := h.EventConfigUc.UpdateEventCertificateConfig(ctx.UserContext(), eventID, params)
		if err != nil {
			return customerror.Parse(&customerror.ErrInternalServer, err)
		}
	}

	dbEventCertConfig, err = h.EventConfigUc.GetEventCertificateConfigByEventID(ctx.UserContext(), eventID)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.Status(http.StatusOK).JSON(EventCertificateConfigResponse{
		ID:                              dbEventCertConfig.ID,
		EventID:                         dbEventCertConfig.EventID,
		BaseCertificateStorageKey:       dbEventCertConfig.BaseCertificateStorageKey,
		BaseCertificatePresignedURL:     dbEventCertConfig.BaseCertificatePresignedURL,
		EventNamePosX:                   dbEventCertConfig.EventNamePosX,
		EventNamePosY:                   dbEventCertConfig.EventNamePosY,
		NamePosX:                        dbEventCertConfig.NamePosX,
		NamePosY:                        dbEventCertConfig.NamePosY,
		AcademicInstitutionPosX:         dbEventCertConfig.AcademicInstitutionPosX,
		AcademicInstitutionPosY:         dbEventCertConfig.AcademicInstitutionPosY,
		CertificateTitlePosX:            dbEventCertConfig.CertificateTitlePosX,
		CertificateTitlePosY:            dbEventCertConfig.CertificateTitlePosY,
		CertificateSubtitlePosX:         dbEventCertConfig.CertificateSubtitlePosX,
		CertificateSubtitlePosY:         dbEventCertConfig.CertificateSubtitlePosY,
		EventNameFontFamilyID:           dbEventCertConfig.EventNameFontFamilyID,
		EventNameFontWeight:             dbEventCertConfig.EventNameFontWeight,
		NameFontFamilyID:                dbEventCertConfig.NameFontFamilyID,
		NameFontWeight:                  dbEventCertConfig.NameFontWeight,
		AcademicInstitutionFontFamilyID: dbEventCertConfig.AcademicInstitutionFontFamilyID,
		AcademicInstitutionFontWeight:   dbEventCertConfig.AcademicInstitutionFontWeight,
		CertificateTitleFontFamilyID:    dbEventCertConfig.CertificateTitleFontFamilyID,
		CertificateTitleFontWeight:      dbEventCertConfig.CertificateTitleFontWeight,
		CertificateSubtitleFontFamilyID: dbEventCertConfig.CertificateSubtitleFontFamilyID,
		CertificateSubtitleFontWeight:   dbEventCertConfig.CertificateSubtitleFontWeight,
		IsPublished:                     dbEventCertConfig.IsPublished,
		CreatedAt:                       dbEventCertConfig.CreatedAt,
		UpdatedAt:                       dbEventCertConfig.UpdatedAt,
	})
}
