package event

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCertificateTemplateVariablesToMap(t *testing.T) {
	t.Run("Convert required fields only", func(t *testing.T) {
		vars := CertificateTemplateVariables{
			Name:      "John Doe",
			EventName: "Tech Conference 2024",
		}

		result := vars.ToMap()

		assert.Equal(t, "John Doe", result["name"])
		assert.Equal(t, "Tech Conference 2024", result["event_name"])
		assert.Len(t, result, 2)
	})

	t.Run("Convert with all optional fields", func(t *testing.T) {
		institution := "MIT"
		title := "Certificate of Achievement"
		subtitle := "For Outstanding Performance"

		vars := CertificateTemplateVariables{
			Name:                "Jane Smith",
			EventName:           "AI Summit",
			AcademicInstitution: &institution,
			CertificateTitle:    &title,
			CertificateSubtitle: &subtitle,
		}

		result := vars.ToMap()

		assert.Equal(t, "Jane Smith", result["name"])
		assert.Equal(t, "AI Summit", result["event_name"])
		assert.Equal(t, "MIT", result["academic_institution"])
		assert.Equal(t, "Certificate of Achievement", result["certificate_title"])
		assert.Equal(t, "For Outstanding Performance", result["certificate_subtitle"])
		assert.Len(t, result, 5)
	})

	t.Run("Skip empty optional fields", func(t *testing.T) {
		emptyStr := ""
		vars := CertificateTemplateVariables{
			Name:                "John Doe",
			EventName:           "Conference",
			AcademicInstitution: &emptyStr, // Should be skipped
		}

		result := vars.ToMap()

		assert.Equal(t, 2, len(result))
		assert.NotContains(t, result, "academic_institution")
	})
}

func TestReplaceSVGTemplateVariables(t *testing.T) {
	uc := &EventUsecase{}

	tests := []struct {
		name      string
		svgInput  string
		variables map[string]string
		expected  string
	}{
		{
			name: "Replace simple text element with id",
			svgInput: `<svg>
				<text id="name" x="100" y="100">Placeholder Name</text>
			</svg>`,
			variables: map[string]string{
				"name": "John Doe",
			},
			expected: "John Doe",
		},
		{
			name: "Replace text element with template syntax",
			svgInput: `<svg>
				<text id="{{ name }}" x="100" y="100">Placeholder Name</text>
			</svg>`,
			variables: map[string]string{
				"name": "Jane Smith",
			},
			expected: "Jane Smith",
		},
		{
			name: "Replace multiple certificate elements",
			svgInput: `<svg>
				<text id="name" x="100" y="100">Name</text>
				<text id="event_name" x="100" y="200">Event</text>
				<text id="academic_institution" x="100" y="250">Institution</text>
			</svg>`,
			variables: map[string]string{
				"name":                 "John Doe",
				"event_name":           "Tech Conference",
				"academic_institution": "MIT",
			},
			expected: "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.replaceSVGTemplateVariables(tt.svgInput, tt.variables)
			assert.NoError(t, err)
			assert.Contains(t, result, tt.expected)
		})
	}
}

func TestRenderSVGToPNG(t *testing.T) {
	uc := &EventUsecase{}

	t.Run("Render simple SVG to PNG", func(t *testing.T) {
		svgContent := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="400" height="300" viewBox="0 0 400 300">
	<rect width="400" height="300" fill="#f0f0f0"/>
	<text x="200" y="150" text-anchor="middle" font-size="24" fill="#000000">Certificate</text>
</svg>`

		pngBytes, err := uc.renderSVGToPNG(svgContent)
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
		assert.Greater(t, len(pngBytes), 0)

		// Check PNG magic bytes
		assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, pngBytes[:4])
	})

	t.Run("Handle invalid SVG", func(t *testing.T) {
		invalidSVG := `<not-valid-svg>`

		_, err := uc.renderSVGToPNG(invalidSVG)
		assert.Error(t, err)
	})
}

func TestReplaceTextContentByID(t *testing.T) {
	uc := &EventUsecase{}

	tests := []struct {
		name       string
		svgContent string
		idPattern  string
		newValue   string
		shouldFind bool
	}{
		{
			name:       "Replace text in simple text element",
			svgContent: `<text id="name">Original</text>`,
			idPattern:  `id="name"`,
			newValue:   "Replaced",
			shouldFind: true,
		},
		{
			name:       "Handle non-existent ID",
			svgContent: `<text id="other">Original</text>`,
			idPattern:  `id="name"`,
			newValue:   "Replaced",
			shouldFind: false,
		},
		{
			name:       "Replace text in element with attributes",
			svgContent: `<text id="name" x="100" y="200" font-size="20">Original</text>`,
			idPattern:  `id="name"`,
			newValue:   "New Value",
			shouldFind: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.replaceTextContentByID(tt.svgContent, tt.idPattern, tt.newValue)

			if tt.shouldFind {
				assert.Contains(t, result, tt.newValue)
				assert.NotContains(t, result, "Original")
			} else {
				assert.Equal(t, tt.svgContent, result)
			}
		})
	}
}

func TestIntegrationReplacementAndRender(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	uc := &EventUsecase{}

	t.Run("Complete flow: replace variables and render", func(t *testing.T) {
		svgTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600" viewBox="0 0 800 600">
	<rect width="800" height="600" fill="#ffffff"/>
	<text id="name" x="400" y="300" text-anchor="middle" font-size="32" fill="#000000">NAME_PLACEHOLDER</text>
	<text id="event_name" x="400" y="350" text-anchor="middle" font-size="24" fill="#666666">EVENT_PLACEHOLDER</text>
</svg>`

		// Use strictly-typed template variables
		templateVars := CertificateTemplateVariables{
			Name:      "John Doe",
			EventName: "Tech Conference 2024",
		}
		variables := templateVars.ToMap()

		// Replace variables
		modifiedSVG, err := uc.replaceSVGTemplateVariables(svgTemplate, variables)
		assert.NoError(t, err)
		assert.Contains(t, modifiedSVG, "John Doe")
		assert.Contains(t, modifiedSVG, "Tech Conference 2024")

		// Render to PNG
		pngBytes, err := uc.renderSVGToPNG(modifiedSVG)
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
		assert.Greater(t, len(pngBytes), 0)

		// Verify PNG signature
		assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, pngBytes[:4])
	})
}

// Benchmark tests
func BenchmarkReplaceSVGTemplateVariables(b *testing.B) {
	uc := &EventUsecase{}
	svgTemplate := strings.Repeat(`<text id="name">Placeholder</text>`, 10)

	templateVars := CertificateTemplateVariables{
		Name:      "John Doe",
		EventName: "Tech Conference",
	}
	variables := templateVars.ToMap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uc.replaceSVGTemplateVariables(svgTemplate, variables)
	}
}

func BenchmarkRenderSVGToPNG(b *testing.B) {
	uc := &EventUsecase{}
	svgContent := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600" viewBox="0 0 800 600">
	<rect width="800" height="600" fill="#f0f0f0"/>
	<text x="400" y="300" text-anchor="middle" font-size="24">Certificate</text>
</svg>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uc.renderSVGToPNG(svgContent)
	}
}

