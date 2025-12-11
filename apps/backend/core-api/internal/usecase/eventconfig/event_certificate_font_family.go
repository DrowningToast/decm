package eventconfig

import (
	"context"
	"strings"
)

// EventCertificateFontFamilyResponse represents a font family available for certificates
type EventCertificateFontFamilyResponse struct {
	ID                   int32   `json:"id"`
	FontFamilyName       string  `json:"font_family_name"`
	CssFontName          string  `json:"css_font_name"`
	IsDefault            bool    `json:"is_default"`
	AvailableFontWeights []int32 `json:"available_font_weights"`
	IsSupportItalic      bool    `json:"is_support_italic"`
}

// GetAllEventCertificateFontFamilies retrieves all available font families for certificates
func (uc *EventConfigUsecase) GetAllEventCertificateFontFamilies(ctx context.Context) ([]EventCertificateFontFamilyResponse, error) {
	fontFamilies, err := uc.EventCertificateFontFamilyDg.GetAllEventCertificateFontFamilies(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]EventCertificateFontFamilyResponse, len(fontFamilies))
	for i, ff := range fontFamilies {
		// Parse available font weights from comma-separated string
		weightsStr := []string{}
		if ff.AvailableFontWeights.Valid {
			weightsStr = strings.Split(ff.AvailableFontWeights.String, ",")
		}
		weights := make([]int32, 0, len(weightsStr))
		for _, w := range weightsStr {
			w = strings.TrimSpace(w)
			var weight int32
			if _, err := parseIntSafe(w, &weight); err == nil {
				weights = append(weights, weight)
			}
		}

		response[i] = EventCertificateFontFamilyResponse{
			ID:                   ff.ID,
			FontFamilyName:       ff.FontFamilyName,
			CssFontName:          ff.CssFontName,
			IsDefault:            ff.IsDefault,
			AvailableFontWeights: weights,
			IsSupportItalic:      ff.IsSupportItalic,
		}
	}

	return response, nil
}

// parseIntSafe safely parses a string to int32
func parseIntSafe(s string, dest *int32) (int, error) {
	var val int
	n, err := scanInt(s, &val)
	if err != nil {
		return 0, err
	}
	*dest = int32(val)
	return n, nil
}

// scanInt is a simple integer parser
func scanInt(s string, dest *int) (int, error) {
	var result int
	for i, ch := range s {
		if ch < '0' || ch > '9' {
			if i == 0 {
				return 0, &parseError{s: s}
			}
			break
		}
		result = result*10 + int(ch-'0')
	}
	*dest = result
	return len(s), nil
}

type parseError struct {
	s string
}

func (e *parseError) Error() string {
	return "parse error: " + e.s
}







