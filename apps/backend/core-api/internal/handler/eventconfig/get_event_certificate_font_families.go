package eventconfig

import (
	"apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2"
)

// GetEventCertificateFontFamiliesResponse represents the response for getting font families
type GetEventCertificateFontFamiliesResponse struct {
	FontFamilies []FontFamilyItem `json:"font_families"`
}

// FontFamilyItem represents a single font family
type FontFamilyItem struct {
	ID                   int32   `json:"id"`
	FontFamilyName       string  `json:"font_family_name"`
	CssFontName          string  `json:"css_font_name"`
	IsDefault            bool    `json:"is_default"`
	AvailableFontWeights []int32 `json:"available_font_weights"`
	IsSupportItalic      bool    `json:"is_support_italic"`
}

// GetEventCertificateFontFamilies godoc
// @Summary Get all available font families for certificates
// @Description Retrieves all font families that can be used in certificate templates, including their available weights and italic support
// @ID get-event-certificate-font-families
// @Accept json
// @Produce json
// @Success 200 {object} GetEventCertificateFontFamiliesResponse "List of available font families"
// @Failure 500 {object} customerror.ErrResponse "Internal server error"
// @Router /eventconfig/certificate-font-families [get]
func (h *Handler) GetEventCertificateFontFamilies(ctx *fiber.Ctx) error {
	// Get all font families from usecase
	fontFamilies, err := h.EventConfigUc.GetAllEventCertificateFontFamilies(ctx.Context())
	if err != nil {
		h.Logger.Error("failed to get font families", "error", err)
		return customerror.ParseWithMessage(&customerror.ErrInternalServer, err, "Failed to retrieve font families")
	}

	// Convert to response format
	items := make([]FontFamilyItem, len(fontFamilies))
	for i, ff := range fontFamilies {
		items[i] = FontFamilyItem{
			ID:                   ff.ID,
			FontFamilyName:       ff.FontFamilyName,
			CssFontName:          ff.CssFontName,
			IsDefault:            ff.IsDefault,
			AvailableFontWeights: ff.AvailableFontWeights,
			IsSupportItalic:      ff.IsSupportItalic,
		}
	}

	response := GetEventCertificateFontFamiliesResponse{
		FontFamilies: items,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

