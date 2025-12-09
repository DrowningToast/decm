package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/usecase/eventconfig"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UpdateEventCertificateTextConfigRequest represents the request body for updating certificate text configuration
type UpdateEventCertificateTextConfigRequest struct {
	EventNameFontFamilyID           *int32 `json:"event_name_font_family_id"`
	EventNameFontWeight             *int32 `json:"event_name_font_weight"`
	NameFontFamilyID                *int32 `json:"name_font_family_id"`
	NameFontWeight                  *int32 `json:"name_font_weight"`
	AcademicInstitutionFontFamilyID *int32 `json:"academic_institution_font_family_id"`
	AcademicInstitutionFontWeight   *int32 `json:"academic_institution_font_weight"`
	CertificateTitleFontFamilyID    *int32 `json:"certificate_title_font_family_id"`
	CertificateTitleFontWeight      *int32 `json:"certificate_title_font_weight"`
	CertificateSubtitleFontFamilyID *int32 `json:"certificate_subtitle_font_family_id"`
	CertificateSubtitleFontWeight   *int32 `json:"certificate_subtitle_font_weight"`
}

// Parse parses and validates the request body
func (r *UpdateEventCertificateTextConfigRequest) Parse(ctx *fiber.Ctx) *customerror.Err {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.ParseWithMessage(&customerror.ErrInvalidArgument, err, "Invalid request body")
	}
	return nil
}

// IsValid validates the request body
func (r *UpdateEventCertificateTextConfigRequest) IsValid() *customerror.Err {
	// TODO: Add validation to check if font family IDs exist in event_certificate_font_families table

	// Validate font weight values (if provided) - common values: 100, 200, 300, 400, 500, 600, 700, 800, 900
	validFontWeights := map[int32]bool{
		100: true, 200: true, 300: true, 400: true, 500: true,
		600: true, 700: true, 800: true, 900: true,
	}

	if r.EventNameFontWeight != nil {
		if !validFontWeights[*r.EventNameFontWeight] {
			return customerror.ParseWithMessage(&customerror.ErrInvalidArgument, nil, "Invalid event_name_font_weight (must be 100-900)")
		}
	}

	if r.NameFontWeight != nil {
		if !validFontWeights[*r.NameFontWeight] {
			return customerror.ParseWithMessage(&customerror.ErrInvalidArgument, nil, "Invalid name_font_weight (must be 100-900)")
		}
	}

	if r.AcademicInstitutionFontWeight != nil {
		if !validFontWeights[*r.AcademicInstitutionFontWeight] {
			return customerror.ParseWithMessage(&customerror.ErrInvalidArgument, nil, "Invalid academic_institution_font_weight (must be 100-900)")
		}
	}

	if r.CertificateTitleFontWeight != nil {
		if !validFontWeights[*r.CertificateTitleFontWeight] {
			return customerror.ParseWithMessage(&customerror.ErrInvalidArgument, nil, "Invalid certificate_title_font_weight (must be 100-900)")
		}
	}

	if r.CertificateSubtitleFontWeight != nil {
		if !validFontWeights[*r.CertificateSubtitleFontWeight] {
			return customerror.ParseWithMessage(&customerror.ErrInvalidArgument, nil, "Invalid certificate_subtitle_font_weight (must be 100-900)")
		}
	}

	return nil
}

// UpdateEventCertificateTextConfigResponse represents the response for updating certificate text configuration
type UpdateEventCertificateTextConfigResponse struct {
	ID                              string   `json:"id"`
	EventID                         string   `json:"event_id"`
	BaseCertificateStorageKey       string   `json:"base_certificate_storage_key"`
	BaseCertificatePresignedURL     string   `json:"base_certificate_presigned_url"`
	EventNamePosX                   float64  `json:"event_name_pos_x"`
	EventNamePosY                   float64  `json:"event_name_pos_y"`
	NamePosX                        float64  `json:"name_pos_x"`
	NamePosY                        float64  `json:"name_pos_y"`
	AcademicInstitutionPosX         *float64 `json:"academic_institution_pos_x,omitempty"`
	AcademicInstitutionPosY         *float64 `json:"academic_institution_pos_y,omitempty"`
	CertificateTitlePosX            *float64 `json:"certificate_title_pos_x,omitempty"`
	CertificateTitlePosY            *float64 `json:"certificate_title_pos_y,omitempty"`
	CertificateSubtitlePosX         *float64 `json:"certificate_subtitle_pos_x,omitempty"`
	CertificateSubtitlePosY         *float64 `json:"certificate_subtitle_pos_y,omitempty"`
	EventNameFontFamilyID           *int32   `json:"event_name_font_family_id,omitempty"`
	EventNameFontWeight             *int32   `json:"event_name_font_weight,omitempty"`
	NameFontFamilyID                *int32   `json:"name_font_family_id,omitempty"`
	NameFontWeight                  *int32   `json:"name_font_weight,omitempty"`
	AcademicInstitutionFontFamilyID *int32   `json:"academic_institution_font_family_id,omitempty"`
	AcademicInstitutionFontWeight   *int32   `json:"academic_institution_font_weight,omitempty"`
	CertificateTitleFontFamilyID    *int32   `json:"certificate_title_font_family_id,omitempty"`
	CertificateTitleFontWeight      *int32   `json:"certificate_title_font_weight,omitempty"`
	CertificateSubtitleFontFamilyID *int32   `json:"certificate_subtitle_font_family_id,omitempty"`
	CertificateSubtitleFontWeight   *int32   `json:"certificate_subtitle_font_weight,omitempty"`
	IsPublished                     bool     `json:"is_published"`
	CreatedAt                       string   `json:"created_at"`
	UpdatedAt                       string   `json:"updated_at"`
}

// UpdateEventCertificateTextConfig godoc
// @Summary Update certificate text configuration
// @Description Update font family and font weight for all text templates in the certificate. This endpoint allows customization of fonts for event name, participant name, academic institution, certificate title, and certificate subtitle.
// @ID update-event-certificate-text-config
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param body body UpdateEventCertificateTextConfigRequest true "Text configuration parameters"
// @Success 200 {object} UpdateEventCertificateTextConfigResponse "Updated certificate configuration"
// @Failure 400 {object} customerror.ErrResponse "Invalid request"
// @Failure 404 {object} customerror.ErrResponse "Certificate configuration not found"
// @Failure 500 {object} customerror.ErrResponse "Internal server error"
// @Router /api/v1/events/{event_id}/certificates/text-config [put]
func (h *Handler) UpdateEventCertificateTextConfig(ctx *fiber.Ctx) error {
	// 1. Parse event ID from path
	eventIDStr := ctx.Params("event_id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return customerror.ParseWithMessage(&customerror.ErrInvalidArgument, err, "Invalid event_id format")
	}

	// 2. Parse and validate request body
	requestBody := UpdateEventCertificateTextConfigRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	// 3. Call usecase to update text configuration
	result, usecaseErr := h.EventConfigUc.UpdateEventCertificateTextConfig(
		ctx.Context(),
		eventID,
		eventconfig.UpdateEventCertificateTextConfigParams{
			EventNameFontFamilyID:           requestBody.EventNameFontFamilyID,
			EventNameFontWeight:             requestBody.EventNameFontWeight,
			NameFontFamilyID:                requestBody.NameFontFamilyID,
			NameFontWeight:                  requestBody.NameFontWeight,
			AcademicInstitutionFontFamilyID: requestBody.AcademicInstitutionFontFamilyID,
			AcademicInstitutionFontWeight:   requestBody.AcademicInstitutionFontWeight,
			CertificateTitleFontFamilyID:    requestBody.CertificateTitleFontFamilyID,
			CertificateTitleFontWeight:      requestBody.CertificateTitleFontWeight,
			CertificateSubtitleFontFamilyID: requestBody.CertificateSubtitleFontFamilyID,
			CertificateSubtitleFontWeight:   requestBody.CertificateSubtitleFontWeight,
		},
	)
	if usecaseErr != nil {
		h.Logger.Error("failed to update certificate text config", "error", usecaseErr, "event_id", eventID)
		return customerror.ParseWithMessage(&customerror.ErrInternalServer, usecaseErr, "Failed to update certificate text configuration")
	}

	// 4. Return updated configuration
	return ctx.Status(fiber.StatusOK).JSON(result)
}
