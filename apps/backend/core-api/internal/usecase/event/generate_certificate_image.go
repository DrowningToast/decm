package event

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"apps/backend/common/customerror"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

// CertificateTemplateVariables contains all variables that can be replaced in the SVG template
type CertificateTemplateVariables struct {
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
// 3. Replaces template variables with certificate data
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

	// 5. Prepare template variables
	variables := CertificateTemplateVariables{
		Name:                getValue(certificate.Name),
		EventName:           event.Title,
		AcademicInstitution: getValue(certificate.AcademicInstitution),
		CertificateTitle:    getValue(certificate.CertificateTitle),
		CertificateSubtitle: getValue(certificate.CertificateSubtitle),
	}

	// 6. Replace variables in SVG
	originalSVG := string(svgBytes)
	processedSVG := replaceTemplateVariables(originalSVG, variables)

	// Log replacement details for debugging
	replacementOccurred := len(originalSVG) != len(processedSVG)
	u.logger.Debug("template replacement completed",
		"original_length", len(svgBytes),
		"processed_length", len(processedSVG),
		"replacement_occurred", replacementOccurred,
		"variables", fmt.Sprintf("%+v", variables))

	// Log which keywords were found in the original SVG (for debugging)
	keywordsFound := []string{}
	possibleKeywords := []string{
		"{{ name }}", "{{ eventName }}", "{{ academicInstitutionName }}",
		"{{ certificateTitle }}", "{{ certificateSubtitle }}",
		"{{name}}", "{{eventName}}", "{{academicInstitutionName}}",
		"{{certificateTitle}}", "{{certificateSubtitle}}",
	}
	for _, keyword := range possibleKeywords {
		if strings.Contains(originalSVG, keyword) {
			keywordsFound = append(keywordsFound, keyword)
		}
	}
	u.logger.Debug("keywords found in SVG template", "keywords", keywordsFound)

	// 7. Render SVG to PNG
	pngBytes, err := renderSVGToPNG(processedSVG)
	if err != nil {
		u.logger.Error("failed to render SVG to PNG", "error", err)
		return nil, customerror.ParseWithMessage(&customerror.ErrInternalServer, err, "Failed to generate certificate image")
	}

	return pngBytes, nil
}

// replaceTemplateVariables replaces text content in SVG based on element IDs or text patterns
// Supports multiple replacement strategies:
// 1. ID-based: <text id="name">...</text> - replaces entire text content (including nested tspan)
// 2. Placeholder: {{name}} or {name} - replaces inline placeholders
// 3. tspan-based: <tspan id="name">...</tspan> - replaces tspan content
func replaceTemplateVariables(svgContent string, variables CertificateTemplateVariables) string {
	result := svgContent

	// Strategy 1: Replace by element ID for <text> tags (e.g., <text id="name">placeholder</text>)
	// This handles nested tspan elements by using a more flexible regex
	// Support both camelCase (frontend) and snake_case (legacy) ID formats
	idReplacements := map[string]string{
		// CamelCase (matches frontend)
		"name":                    variables.Name,
		"eventName":               variables.EventName,
		"academicInstitutionName": variables.AcademicInstitution,
		"certificateTitle":        variables.CertificateTitle,
		"certificateSubtitle":     variables.CertificateSubtitle,
		// Snake_case (legacy support)
		"event_name":           variables.EventName,
		"academic_institution": variables.AcademicInstitution,
		"certificate_title":    variables.CertificateTitle,
		"certificate_subtitle": variables.CertificateSubtitle,
	}

	for id, value := range idReplacements {
		// Match <text id="name">...any content including nested tags...</text>
		// Using (?s) flag to make . match newlines
		pattern := fmt.Sprintf(`(?s)(<text[^>]*\bid="%s"[^>]*>)(.*?)(</text>)`, regexp.QuoteMeta(id))
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, fmt.Sprintf("${1}%s${3}", value))

		// Also try to match <tspan id="name">content</tspan> for nested cases
		tspanPattern := fmt.Sprintf(`(?s)(<tspan[^>]*\bid="%s"[^>]*>)(.*?)(</tspan>)`, regexp.QuoteMeta(id))
		tspanRe := regexp.MustCompile(tspanPattern)
		result = tspanRe.ReplaceAllString(result, fmt.Sprintf("${1}%s${3}", value))
	}

	// Strategy 2: Replace inline placeholders like {{ name }} or {{ eventName }}
	// Support both camelCase (frontend) and snake_case (legacy) formats
	// Note: Frontend uses spaces inside {{ }}, e.g., "{{ eventName }}"
	placeholderReplacements := map[string]string{
		// CamelCase with spaces (matches frontend format)
		"{{ name }}":                    variables.Name,
		"{{ eventName }}":               variables.EventName,
		"{{ academicInstitutionName }}": variables.AcademicInstitution,
		"{{ certificateTitle }}":        variables.CertificateTitle,
		"{{ certificateSubtitle }}":     variables.CertificateSubtitle,
		// CamelCase without spaces (alternative format)
		"{{name}}":                    variables.Name,
		"{{eventName}}":               variables.EventName,
		"{{academicInstitutionName}}": variables.AcademicInstitution,
		"{{certificateTitle}}":        variables.CertificateTitle,
		"{{certificateSubtitle}}":     variables.CertificateSubtitle,
		// Snake_case (legacy support)
		"{{event_name}}":           variables.EventName,
		"{event_name}":             variables.EventName,
		"{{academic_institution}}": variables.AcademicInstitution,
		"{academic_institution}":   variables.AcademicInstitution,
		"{{certificate_title}}":    variables.CertificateTitle,
		"{certificate_title}":      variables.CertificateTitle,
		"{{certificate_subtitle}}": variables.CertificateSubtitle,
		"{certificate_subtitle}":   variables.CertificateSubtitle,
	}

	for placeholder, value := range placeholderReplacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// renderSVGToPNG converts SVG string to PNG byte array
// Uses chromedp with headless Chrome for full SVG specification support
// This handles embedded images, patterns, custom fonts, and all SVG features
func renderSVGToPNG(svgContent string) ([]byte, error) {
	// Create context with timeout for rendering
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create chromedp context
	chromedpCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var buf []byte

	// Wrap SVG in HTML for proper rendering
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { 
            margin: 0; 
            padding: 0; 
            overflow: hidden;
        }
        svg { 
            display: block;
        }
    </style>
</head>
<body>
%s
</body>
</html>`, svgContent)

	// Render SVG using headless Chrome
	if err := chromedp.Run(chromedpCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.Sleep(1*time.Second),      // Wait for fonts and images to load
		chromedp.FullScreenshot(&buf, 100), // Quality 100
	); err != nil {
		return nil, fmt.Errorf("failed to render SVG: %w", err)
	}

	return buf, nil
}

// getValue safely extracts string value from pointer, returns empty string if nil
func getValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
