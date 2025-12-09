package event

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

// CertificateData contains all data to be overlaid on the certificate
type CertificateData struct {
	Name                string
	EventName           string
	AcademicInstitution string
	CertificateTitle    string
	CertificateSubtitle string
}

// GenerateCertificateImage generates a PNG image from certificate SVG template
// This function:
// 1. Retrieves the certificate by ID
// 2. Fetches the SVG template from S3
// 3. Overlays text elements with certificate data at absolute positions from config
// 4. Renders the SVG to PNG
// 5. Returns PNG as byte array
func (u *EventUsecase) GenerateCertificateImage(ctx context.Context, certificateID uuid.UUID) ([]byte, error) {
	// 1. Get certificate by ID
	certificate, err := u.EventCertificateDataGateway.GetEventCertificateByID(ctx, certificateID)
	if err != nil {
		u.logger.Error("failed to get certificate", "error", err, "certificate_id", certificateID)
		return nil, customerror.ParseWithMessage(&customerror.ErrNotFound, err, "Certificate not found")
	}

	// 2. Get certificate config to retrieve SVG template storage key
	certificateConfig, err := u.EventCertificateConfigDg.GetEventCertificateConfigByEventID(ctx, certificate.EventId)
	if err != nil {
		u.logger.Error("failed to get certificate config", "error", err, "event_id", certificate.EventId)
		return nil, customerror.ParseWithMessage(&customerror.ErrNotFound, err, "Certificate configuration not found")
	}

	// 3. Download SVG template from S3
	svgReader, err := u.S3Service.GetFile(ctx, certificateConfig.BaseCertificateStorageKey)
	if err != nil {
		u.logger.Error("failed to download SVG template from S3", "error", err, "storage_key", certificateConfig.BaseCertificateStorageKey)
		return nil, customerror.ParseWithMessage(&customerror.ErrInternalServer, err, "Failed to retrieve certificate template")
	}
	defer svgReader.Close()

	// Read SVG content
	svgBytes, err := io.ReadAll(svgReader)
	if err != nil {
		u.logger.Error("failed to read SVG content", "error", err)
		return nil, customerror.ParseWithMessage(&customerror.ErrInternalServer, err, "Failed to read certificate template")
	}

	// 4. Get event details for event name
	event, err := u.EventDataGateway.GetEventById(ctx, certificate.EventId)
	if err != nil {
		u.logger.Error("failed to get event", "error", err, "event_id", certificate.EventId)
		return nil, customerror.ParseWithMessage(&customerror.ErrNotFound, err, "Event not found")
	}

	// 5. Prepare certificate data
	certificateData := CertificateData{
		Name:                getValue(certificate.Name),
		EventName:           event.Title,
		AcademicInstitution: getValue(certificate.AcademicInstitution),
		CertificateTitle:    getValue(certificate.CertificateTitle),
		CertificateSubtitle: getValue(certificate.CertificateSubtitle),
	}

	// 6. Add text overlays to SVG using absolute positions from config
	originalSVG := string(svgBytes)
	processedSVG := u.addTextOverlaysToSVG(ctx, originalSVG, certificateData, certificateConfig)

	u.logger.Debug("text overlays added to certificate",
		"original_length", len(originalSVG),
		"processed_length", len(processedSVG),
		"data", fmt.Sprintf("%+v", certificateData))

	// 7. Render SVG to PNG
	pngBytes, err := renderSVGToPNG(processedSVG)
	if err != nil {
		u.logger.Error("failed to render SVG to PNG", "error", err)
		return nil, customerror.ParseWithMessage(&customerror.ErrInternalServer, err, "Failed to generate certificate image")
	}

	return pngBytes, nil
}

// addTextOverlaysToSVG adds new text elements to the SVG with absolute positioning
// from the database, instead of replacing template strings. This provides precise
// control over text positioning and styling.
func (uc *EventUsecase) addTextOverlaysToSVG(ctx context.Context, svgContent string, data CertificateData, config *entity.EventCertificateConfig) string {
	// Remove or hide existing template text elements to avoid duplicates
	result := hideTemplatePlaceholders(svgContent)

	// Build text overlay elements with absolute positioning from config
	var textOverlays strings.Builder

	// Add name text (required field, always has position)
	if data.Name != "" {
		fontFamily := uc.getFontFamilyName(ctx, config.NameFontFamilyID)
		fontWeight := int32PtrToString(config.NameFontWeight, "700")
		textOverlays.WriteString(createTextElement(data.Name, config.NamePosX, config.NamePosY, fontFamily, fontWeight, 16))
	}

	// Add event name text (required field, always has position)
	if data.EventName != "" {
		fontFamily := uc.getFontFamilyName(ctx, config.EventNameFontFamilyID)
		fontWeight := int32PtrToString(config.EventNameFontWeight, "700")
		textOverlays.WriteString(createTextElement(data.EventName, config.EventNamePosX, config.EventNamePosY, fontFamily, fontWeight, 16))
	}

	// Add academic institution text (optional field)
	if config.AcademicInstitutionPosX != nil && config.AcademicInstitutionPosY != nil && data.AcademicInstitution != "" {
		fontFamily := uc.getFontFamilyName(ctx, config.AcademicInstitutionFontFamilyID)
		fontWeight := int32PtrToString(config.AcademicInstitutionFontWeight, "700")
		textOverlays.WriteString(createTextElement(data.AcademicInstitution, *config.AcademicInstitutionPosX, *config.AcademicInstitutionPosY, fontFamily, fontWeight, 16))
	}

	// Add certificate title text (optional field)
	if config.CertificateTitlePosX != nil && config.CertificateTitlePosY != nil && data.CertificateTitle != "" {
		fontFamily := uc.getFontFamilyName(ctx, config.CertificateTitleFontFamilyID)
		fontWeight := int32PtrToString(config.CertificateTitleFontWeight, "700")
		textOverlays.WriteString(createTextElement(data.CertificateTitle, *config.CertificateTitlePosX, *config.CertificateTitlePosY, fontFamily, fontWeight, 16))
	}

	// Add certificate subtitle text (optional field)
	if config.CertificateSubtitlePosX != nil && config.CertificateSubtitlePosY != nil && data.CertificateSubtitle != "" {
		fontFamily := uc.getFontFamilyName(ctx, config.CertificateSubtitleFontFamilyID)
		fontWeight := int32PtrToString(config.CertificateSubtitleFontWeight, "700")
		textOverlays.WriteString(createTextElement(data.CertificateSubtitle, *config.CertificateSubtitlePosX, *config.CertificateSubtitlePosY, fontFamily, fontWeight, 16))
	}

	// Insert text overlays before closing </svg> tag
	closingSvgIndex := strings.LastIndex(result, "</svg>")
	if closingSvgIndex == -1 {
		// If no closing tag found, append at the end
		return result + textOverlays.String()
	}

	return result[:closingSvgIndex] + textOverlays.String() + result[closingSvgIndex:]
}

// getFontFamilyName fetches the CSS font name from the database by font family ID.
// Returns the default font family from the database if the ID is nil or if fetching fails.
func (uc *EventUsecase) getFontFamilyName(ctx context.Context, fontFamilyID *int32) string {
	// If no font family ID specified, use default from database
	if fontFamilyID == nil {
		defaultFont, err := uc.EventCertificateFontFamilyDg.GetDefaultEventCertificateFontFamily(ctx)
		if err != nil {
			uc.logger.Warn("failed to get default font family from database, using fallback",
				"error", err,
			)
			return "Prompt" // Fallback if database query fails
		}
		return defaultFont.CssFontName
	}

	// Fetch font family from database
	fontFamily, err := uc.EventCertificateFontFamilyDg.GetEventCertificateFontFamilyByID(ctx, *fontFamilyID)
	if err != nil {
		uc.logger.Warn("failed to get font family by ID, using default from database",
			"font_family_id", *fontFamilyID,
			"error", err,
		)
		// Try to get default font from database
		defaultFont, defaultErr := uc.EventCertificateFontFamilyDg.GetDefaultEventCertificateFontFamily(ctx)
		if defaultErr != nil {
			uc.logger.Warn("failed to get default font family from database, using fallback",
				"error", defaultErr,
			)
			return "Prompt" // Fallback if database query fails
		}
		return defaultFont.CssFontName
	}

	// Return the CSS font name (e.g., "Inter", "Prompt", "Sarabun")
	return fontFamily.CssFontName
}

// hideTemplatePlaceholders hides or removes template placeholder text elements
// by setting their opacity to 0 or removing them entirely
func hideTemplatePlaceholders(svgContent string) string {
	result := svgContent

	// Hide text elements that contain template placeholders by setting visibility="hidden"
	placeholders := []string{
		"{{ name }}", "{{ eventName }}", "{{ academicInstitutionName }}",
		"{{ certificateTitle }}", "{{ certificateSubtitle }}",
	}

	for _, placeholder := range placeholders {
		// Find text elements containing this placeholder
		escapedPlaceholder := regexp.QuoteMeta(placeholder)
		pattern := fmt.Sprintf(`(<text[^>]*id="%s"[^>]*>)`, escapedPlaceholder)
		re := regexp.MustCompile(pattern)

		result = re.ReplaceAllStringFunc(result, func(match string) string {
			// Add visibility="hidden" if not already present
			if !strings.Contains(match, "visibility") {
				return strings.TrimSuffix(match, ">") + ` visibility="hidden">`
			}
			return match
		})
	}

	return result
}

// createTextElement creates an SVG text element with the specified properties
// All text is center-aligned (text-anchor="middle") for consistent positioning
func createTextElement(text string, x, y float64, fontFamily, fontWeight string, fontSize int) string {
	// Escape special XML characters in text
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&apos;")

	return fmt.Sprintf(
		`<text x="%.2f" y="%.2f" font-family="%s" font-size="%d" font-weight="%s" fill="#090909" text-anchor="middle">%s</text>`,
		x, y, fontFamily, fontSize, fontWeight, text,
	)
}

// getValueOrDefault returns the value if not nil, otherwise returns default
func getValueOrDefault(value *string, defaultValue string) string {
	if value == nil || *value == "" {
		return defaultValue
	}
	return *value
}

// int32PtrToString converts *int32 font weight to string (e.g., "400", "700")
// If nil, returns the defaultValue
func int32PtrToString(weight *int32, defaultValue string) string {
	if weight == nil {
		return defaultValue
	}
	return fmt.Sprintf("%d", *weight)
}

// renderSVGToPNG converts SVG string to PNG byte array
// Uses chromedp with headless Chrome for full SVG specification support
// This handles embedded images, patterns, custom fonts, and all SVG features
// Renders at 2x resolution for high-DPI displays (Retina quality)
func renderSVGToPNG(svgContent string) ([]byte, error) {
	// Extract SVG dimensions from viewBox or width/height attributes
	width, height := extractSVGDimensions(svgContent)

	// Scale factor for high-resolution rendering (2x for Retina/high-DPI)
	const scaleFactor = 2
	renderWidth := width * scaleFactor
	renderHeight := height * scaleFactor

	// Create context with timeout for rendering
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create chromedp context with scaled viewport size for high-resolution rendering
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(renderWidth, renderHeight),
		chromedp.Flag("force-device-scale-factor", fmt.Sprintf("%d", scaleFactor)), // 2x device pixel ratio
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	chromedpCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var buf []byte

	// Wrap SVG in HTML for proper rendering with Google Fonts
	// Use original dimensions for CSS, but render at 2x scale
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Prompt:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&family=Inter:ital,wght@0,100..900;1,100..900&family=Noto+Sans+Thai:wght@100..900&family=Sarabun:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800&family=Kanit:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap" rel="stylesheet">
    <style>
        body {
            margin: 0;
            padding: 0;
            overflow: hidden;
            width: %dpx;
            height: %dpx;
        }
        svg {
            display: block;
            width: 100%%;
            height: 100%%;
        }
    </style>
</head>
<body>
%s
</body>
</html>`, renderWidth, renderHeight, svgContent)

	// Render SVG using headless Chrome at 2x resolution
	if err := chromedp.Run(chromedpCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.Sleep(2*time.Second),      // Wait for fonts and images to load
		chromedp.FullScreenshot(&buf, 100), // Quality 100
	); err != nil {
		return nil, fmt.Errorf("failed to render SVG: %w", err)
	}

	return buf, nil
}

// extractSVGDimensions extracts width and height from SVG viewBox or attributes
// Returns default 1920x1080 if dimensions cannot be determined
func extractSVGDimensions(svgContent string) (int, int) {
	// Default dimensions matching frontend template defaults
	defaultWidth := 1920
	defaultHeight := 1080

	// Try to extract from viewBox first (more reliable)
	viewBoxPattern := regexp.MustCompile(`viewBox=["']([^"']+)["']`)
	if matches := viewBoxPattern.FindStringSubmatch(svgContent); len(matches) > 1 {
		// viewBox format: "minX minY width height"
		var minX, minY, width, height float64
		if _, err := fmt.Sscanf(matches[1], "%f %f %f %f", &minX, &minY, &width, &height); err == nil {
			return int(width), int(height)
		}
	}

	// Fallback to width/height attributes
	widthPattern := regexp.MustCompile(`width=["']?(\d+)`)
	heightPattern := regexp.MustCompile(`height=["']?(\d+)`)

	var width, height int
	if matches := widthPattern.FindStringSubmatch(svgContent); len(matches) > 1 {
		fmt.Sscanf(matches[1], "%d", &width)
	}
	if matches := heightPattern.FindStringSubmatch(svgContent); len(matches) > 1 {
		fmt.Sscanf(matches[1], "%d", &height)
	}

	if width > 0 && height > 0 {
		return width, height
	}

	// Return defaults if no dimensions found
	return defaultWidth, defaultHeight
}

// getValue safely extracts string value from pointer, returns empty string if nil
func getValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
