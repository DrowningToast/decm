package event

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/beevik/etree" // REQUIRED: Standard library for robust XML manipulation
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// Certificate Image Generation Usecase

// CertificateTemplateVariables defines the strictly-typed template variables
type CertificateTemplateVariables struct {
	Name                string  `json:"name" validate:"required"`
	EventName           string  `json:"event_name" validate:"required"`
	AcademicInstitution *string `json:"academic_institution,omitempty"`
	CertificateTitle    *string `json:"certificate_title,omitempty"`
	CertificateSubtitle *string `json:"certificate_subtitle,omitempty"`
}

// GenerateCertificateImageParams contains parameters for certificate image generation
type GenerateCertificateImageParams struct {
	TemplateVariables CertificateTemplateVariables

	// Optional position overrides (if nil, uses position defined in SVG)
	NamePosX                *float64
	NamePosY                *float64
	EventNamePosX           *float64
	EventNamePosY           *float64
	CertificateTitlePosX    *float64
	CertificateTitlePosY    *float64
	CertificateSubtitlePosX *float64
	CertificateSubtitlePosY *float64
	AcademicInstitutionPosX *float64
	AcademicInstitutionPosY *float64
}

// GenerateCertificateImage generates a PNG image from an SVG certificate template
func (uc *EventUsecase) GenerateCertificateImage(
	ctx context.Context,
	certificateConfigID uuid.UUID,
	params GenerateCertificateImageParams,
) ([]byte, error) {
	// 1. Get the certificate config
	certConfig, err := uc.EventCertificateConfigDg.GetEventCertificateConfigByID(ctx, certificateConfigID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get certificate config")
	}

	// 2. Download the SVG template from S3
	svgReader, err := uc.S3Service.GetFile(ctx, certConfig.BaseCertificateStorageKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download SVG template from S3")
	}
	defer svgReader.Close()

	svgBytes, err := io.ReadAll(svgReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read SVG content")
	}

	if len(svgBytes) == 0 {
		return nil, errors.New("SVG template is empty")
	}

	// 3. Process SVG (Replace text and positions)
	modifiedSVGBytes, err := uc.processSVGTemplate(svgBytes, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process SVG template")
	}

	// 4. Render SVG to PNG
	pngBytes, err := uc.renderSVGToPNG(modifiedSVGBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render SVG to PNG")
	}

	return pngBytes, nil
}

// processSVGTemplate parses the XML, updates text content, and applies position overrides
func (uc *EventUsecase) processSVGTemplate(originalSVG []byte, params GenerateCertificateImageParams) ([]byte, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(originalSVG); err != nil {
		return nil, errors.Wrap(err, "failed to parse SVG XML")
	}

	// Define the mapping of IDs to Values and Optional Position Overrides
	type fieldConfig struct {
		Value string
		PosX  *float64
		PosY  *float64
	}

	// Helper to safely dereference string pointers
	safeStr := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	// Map the SVG element IDs to their specific values and coordinate params
	fields := map[string]fieldConfig{
		"name": {
			Value: params.TemplateVariables.Name,
			PosX:  params.NamePosX,
			PosY:  params.NamePosY,
		},
		"event_name": {
			Value: params.TemplateVariables.EventName,
			PosX:  params.EventNamePosX,
			PosY:  params.EventNamePosY,
		},
		"academic_institution": {
			Value: safeStr(params.TemplateVariables.AcademicInstitution),
			PosX:  params.AcademicInstitutionPosX,
			PosY:  params.AcademicInstitutionPosY,
		},
		"certificate_title": {
			Value: safeStr(params.TemplateVariables.CertificateTitle),
			PosX:  params.CertificateTitlePosX,
			PosY:  params.CertificateTitlePosY,
		},
		"certificate_subtitle": {
			Value: safeStr(params.TemplateVariables.CertificateSubtitle),
			PosX:  params.CertificateSubtitlePosX,
			PosY:  params.CertificateSubtitlePosY,
		},
	}

	// Iterate through fields and update the SVG document
	for id, config := range fields {
		if config.Value == "" {
			continue // Skip empty values
		}

		// Find element by ID - try multiple ID formats
		var element *etree.Element

		// Try standard id="name"
		element = doc.FindElement(fmt.Sprintf("//*[@id='%s']", id))

		// Try template formats if not found
		if element == nil {
			element = doc.FindElement(fmt.Sprintf("//*[@id='{{ %s }}']", id))
		}
		if element == nil {
			element = doc.FindElement(fmt.Sprintf("//*[@id='{{%s}}']", id))
		}

		if element != nil {
			parent := element.Parent()
			if parent == nil {
				uc.logger.Warn("Element has no parent, skipping", "id", id)
				continue
			}

			// Get original position if exists
			xPos := element.SelectAttrValue("x", "400")
			yPos := element.SelectAttrValue("y", "300")

			// Apply position overrides if provided
			if config.PosX != nil {
				xPos = fmt.Sprintf("%f", *config.PosX)
			}
			if config.PosY != nil {
				yPos = fmt.Sprintf("%f", *config.PosY)
			}

			// Preserve original styling attributes from the template
			fontSize := element.SelectAttrValue("font-size", "24")
			fontFamily := element.SelectAttrValue("font-family", "Arial, sans-serif")
			fontWeight := element.SelectAttrValue("font-weight", "normal")
			fill := element.SelectAttrValue("fill", "#000000")
			textAnchor := element.SelectAttrValue("text-anchor", "middle")

			// Create new <text> element with the actual content
			newTextElement := etree.NewElement("text")
			newTextElement.CreateAttr("id", id)
			newTextElement.CreateAttr("x", xPos)
			newTextElement.CreateAttr("y", yPos)
			newTextElement.CreateAttr("text-anchor", textAnchor)
			newTextElement.CreateAttr("font-size", fontSize)
			newTextElement.CreateAttr("font-family", fontFamily)
			newTextElement.CreateAttr("font-weight", fontWeight)
			newTextElement.CreateAttr("fill", fill)

			// Copy other attributes that might be present (opacity, transform, etc.)
			// but skip path-specific attributes like 'd', 'stroke', 'stroke-width'
			for _, attr := range element.Attr {
				attrName := attr.Key
				// Skip attributes we've already handled
				if attrName == "id" || attrName == "x" || attrName == "y" ||
					attrName == "text-anchor" || attrName == "font-size" ||
					attrName == "font-family" || attrName == "font-weight" || attrName == "fill" {
					continue
				}
				// Skip path-specific attributes that don't belong on text elements
				if attrName == "d" || attrName == "stroke" || attrName == "stroke-width" ||
					attrName == "stroke-linecap" || attrName == "stroke-linejoin" {
					continue
				}
				newTextElement.CreateAttr(attrName, attr.Value)
			}

			newTextElement.SetText(config.Value)

			// Replace the old element with the new text element
			// Remove old element and insert new one in its place
			children := parent.ChildElements()
			for i, child := range children {
				if child == element {
					parent.RemoveChild(element)
					parent.InsertChildAt(i, newTextElement)

					uc.logger.Info("Replaced element with text",
						"id", id,
						"value", config.Value,
						"x", xPos,
						"y", yPos,
						"fontSize", fontSize,
						"fontFamily", fontFamily,
					)
					break
				}
			}
		} else {
			uc.logger.Warn("SVG element not found for field", "id", id)
		}
	}

	// Serialize back to bytes
	modifiedSVG, err := doc.WriteToBytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize modified SVG")
	}

	// Log the full modified SVG for debugging
	uc.logger.Debug("Modified SVG", "svg", string(modifiedSVG))

	return modifiedSVG, nil
}

// renderSVGToPNG converts an SVG byte slice to PNG binary data
func (uc *EventUsecase) renderSVGToPNG(svgContent []byte) ([]byte, error) {
	// Log the SVG content for debugging (first 500 chars)
	svgPreview := string(svgContent)
	if len(svgPreview) > 500 {
		svgPreview = svgPreview[:500] + "..."
	}
	uc.logger.Debug("SVG content preview", "svg", svgPreview)

	// Parse the SVG using oksvg
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svgContent))
	if err != nil {
		uc.logger.Error("Failed to parse SVG", "error", err, "svgLength", len(svgContent))
		return nil, errors.Wrap(err, "failed to parse SVG for rendering")
	}

	// Determine dimensions
	width := int(icon.ViewBox.W)
	height := int(icon.ViewBox.H)

	if width <= 0 || height <= 0 {
		// Fallback dimensions if ViewBox is missing
		width = 1200
		height = 800
		uc.logger.Warn("SVG ViewBox missing, using default dimensions", "width", width, "height", height)
	}

	uc.logger.Info("Rendering SVG to PNG", "width", width, "height", height)

	// Setup the rasterizer with white background
	icon.SetTarget(0, 0, float64(width), float64(height))
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with white background to prevent black/transparent images
	// RGBA format: 4 bytes per pixel (Red, Green, Blue, Alpha)
	white := []byte{255, 255, 255, 255} // White with full opacity
	for i := 0; i < len(rgba.Pix); i += 4 {
		copy(rgba.Pix[i:i+4], white)
	}

	uc.logger.Debug("Initialized white background",
		"totalPixels", width*height,
		"bufferSize", len(rgba.Pix))

	scanner := rasterx.NewScannerGV(width, height, rgba, rgba.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)

	// Draw the SVG onto the white background
	icon.Draw(raster, 1.0)

	// Sample some pixels to see what was actually rendered
	samplePixels := make([]string, 0, 5)
	samplePositions := []int{0, len(rgba.Pix) / 4, len(rgba.Pix) / 2, 3 * len(rgba.Pix) / 4, len(rgba.Pix) - 4}
	for _, pos := range samplePositions {
		if pos >= 0 && pos < len(rgba.Pix)-3 {
			r, g, b, a := rgba.Pix[pos], rgba.Pix[pos+1], rgba.Pix[pos+2], rgba.Pix[pos+3]
			samplePixels = append(samplePixels, fmt.Sprintf("RGBA(%d,%d,%d,%d)", r, g, b, a))
		}
	}
	uc.logger.Debug("Sample pixels after SVG render", "samples", samplePixels)

	// Check if the image has any non-white pixels (to detect if SVG rendered)
	hasContent := false
	whitePixelCount := 0
	blackPixelCount := 0
	otherPixelCount := 0

	for i := 0; i < len(rgba.Pix); i += 4 {
		r, g, b := rgba.Pix[i], rgba.Pix[i+1], rgba.Pix[i+2]

		if r == 255 && g == 255 && b == 255 {
			whitePixelCount++
		} else if r == 0 && g == 0 && b == 0 {
			blackPixelCount++
			hasContent = true
		} else {
			otherPixelCount++
			hasContent = true
		}
	}

	totalPixels := width * height
	uc.logger.Info("Pixel analysis",
		"whitePixels", whitePixelCount,
		"blackPixels", blackPixelCount,
		"otherPixels", otherPixelCount,
		"totalPixels", totalPixels,
		"percentWhite", float64(whitePixelCount)/float64(totalPixels)*100,
	)

	if !hasContent {
		uc.logger.Warn("SVG rendered but resulted in blank white image - text may not be rendering properly")
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, rgba); err != nil {
		uc.logger.Error("Failed to encode PNG", "error", err)
		return nil, errors.Wrap(err, "failed to encode PNG")
	}

	uc.logger.Info("Successfully rendered PNG", "pngSize", buf.Len())
	return buf.Bytes(), nil
}

// GenerateCertificateImageByEventID is a convenience method
func (uc *EventUsecase) GenerateCertificateImageByEventID(
	ctx context.Context,
	eventID uuid.UUID,
	params GenerateCertificateImageParams,
) ([]byte, error) {
	certConfig, err := uc.EventCertificateConfigDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get certificate config")
	}

	return uc.GenerateCertificateImage(ctx, certConfig.ID, params)
}
