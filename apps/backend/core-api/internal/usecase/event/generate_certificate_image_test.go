package event

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCertificateTemplateVariables(t *testing.T) {
	t.Run("Create with required fields only", func(t *testing.T) {
		vars := CertificateTemplateVariables{
			Name:      "John Doe",
			EventName: "Tech Conference 2024",
		}

		assert.Equal(t, "John Doe", vars.Name)
		assert.Equal(t, "Tech Conference 2024", vars.EventName)
		assert.Empty(t, vars.AcademicInstitution)
		assert.Empty(t, vars.CertificateTitle)
		assert.Empty(t, vars.CertificateSubtitle)
	})

	t.Run("Create with all fields", func(t *testing.T) {
		vars := CertificateTemplateVariables{
			Name:                "Jane Smith",
			EventName:           "AI Summit",
			AcademicInstitution: "MIT",
			CertificateTitle:    "Certificate of Achievement",
			CertificateSubtitle: "For Outstanding Performance",
		}

		assert.Equal(t, "Jane Smith", vars.Name)
		assert.Equal(t, "AI Summit", vars.EventName)
		assert.Equal(t, "MIT", vars.AcademicInstitution)
		assert.Equal(t, "Certificate of Achievement", vars.CertificateTitle)
		assert.Equal(t, "For Outstanding Performance", vars.CertificateSubtitle)
	})
}

func TestReplaceSVGTemplateVariables(t *testing.T) {
	tests := []struct {
		name      string
		svgInput  string
		variables CertificateTemplateVariables
		expected  string
	}{
		{
			name: "Replace simple text element with id",
			svgInput: `<svg>
				<text id="name" x="100" y="100">Placeholder Name</text>
			</svg>`,
			variables: CertificateTemplateVariables{
				Name: "John Doe",
			},
			expected: "John Doe",
		},
		{
			name: "Replace text element with template syntax",
			svgInput: `<svg>
				<text id="name" x="100" y="100">Placeholder Name</text>
			</svg>`,
			variables: CertificateTemplateVariables{
				Name: "Jane Smith",
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
			variables: CertificateTemplateVariables{
				Name:                "John Doe",
				EventName:           "Tech Conference",
				AcademicInstitution: "MIT",
			},
			expected: "John Doe",
		},
		{
			name: "Replace text element with nested tspan",
			svgInput: `<svg>
				<text id="name" x="100" y="100">
					<tspan>Old Name</tspan>
				</text>
			</svg>`,
			variables: CertificateTemplateVariables{
				Name: "New Name",
			},
			expected: "New Name",
		},
		{
			name: "Replace tspan element with id directly",
			svgInput: `<svg>
				<text x="100" y="100">
					<tspan id="name">Old Name</tspan>
				</text>
			</svg>`,
			variables: CertificateTemplateVariables{
				Name: "Direct Tspan Name",
			},
			expected: "Direct Tspan Name",
		},
		{
			name: "Replace placeholder syntax {{name}}",
			svgInput: `<svg>
				<text x="100" y="100">Certificate for {{name}}</text>
			</svg>`,
			variables: CertificateTemplateVariables{
				Name: "Alice",
			},
			expected: "Certificate for Alice",
		},
		{
			name: "Replace placeholder syntax {name}",
			svgInput: `<svg>
				<text x="100" y="100">Awarded to {name}</text>
			</svg>`,
			variables: CertificateTemplateVariables{
				Name: "Bob",
			},
			expected: "Awarded to Bob",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceTemplateVariables(tt.svgInput, tt.variables)
			assert.Contains(t, result, tt.expected)
		})
	}
}

func TestRenderSVGToPNG(t *testing.T) {
	t.Run("Render simple SVG to PNG", func(t *testing.T) {
		svgContent := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="400" height="300" viewBox="0 0 400 300">
	<rect width="400" height="300" fill="#f0f0f0"/>
	<text x="200" y="150" text-anchor="middle" font-size="24" fill="#000000">Certificate</text>
</svg>`

		pngBytes, err := renderSVGToPNG(svgContent)
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
		assert.Greater(t, len(pngBytes), 0)

		// Check PNG magic bytes
		assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, pngBytes[:4])
	})

	t.Run("Handle invalid SVG", func(t *testing.T) {
		invalidSVG := `<not-valid-svg>`

		// Note: chromedp is forgiving and will still render something
		// even with invalid SVG, just testing it doesn't crash
		pngBytes, err := renderSVGToPNG(invalidSVG)
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
	})

	t.Run("Render SVG with replaced template variables", func(t *testing.T) {
		// Create SVG template with IDs
		svgTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600" viewBox="0 0 800 600">
	<rect width="800" height="600" fill="#ffffff"/>
	<text id="certificate_title" x="400" y="100" text-anchor="middle" font-size="32" fill="#000000">Title Placeholder</text>
	<text id="name" x="400" y="250" text-anchor="middle" font-size="48" fill="#1a1a1a">Name Placeholder</text>
	<text id="event_name" x="400" y="350" text-anchor="middle" font-size="24" fill="#333333">Event Placeholder</text>
	<text id="academic_institution" x="400" y="450" text-anchor="middle" font-size="20" fill="#666666">Institution Placeholder</text>
</svg>`

		// Replace template variables
		variables := CertificateTemplateVariables{
			Name:                "John Doe",
			EventName:           "Tech Conference 2024",
			AcademicInstitution: "Massachusetts Institute of Technology",
			CertificateTitle:    "Certificate of Excellence",
		}
		processedSVG := replaceTemplateVariables(svgTemplate, variables)

		// Verify replacements occurred
		assert.Contains(t, processedSVG, "John Doe")
		assert.Contains(t, processedSVG, "Tech Conference 2024")
		assert.Contains(t, processedSVG, "Massachusetts Institute of Technology")
		assert.Contains(t, processedSVG, "Certificate of Excellence")
		assert.NotContains(t, processedSVG, "Name Placeholder")
		assert.NotContains(t, processedSVG, "Event Placeholder")

		// Render to PNG
		pngBytes, err := renderSVGToPNG(processedSVG)
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
		assert.Greater(t, len(pngBytes), 10000, "PNG should be reasonably sized for 800x600 certificate")

		// Check PNG magic bytes
		assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, pngBytes[:4])
	})
}

func TestReplaceTextContentByID(t *testing.T) {
	tests := []struct {
		name       string
		svgContent string
		variables  CertificateTemplateVariables
		shouldFind string
		expected   string
	}{
		{
			name:       "Replace text in simple text element",
			svgContent: `<text id="name">Original</text>`,
			variables: CertificateTemplateVariables{
				Name: "Replaced",
			},
			shouldFind: "Replaced",
			expected:   "Replaced",
		},
		{
			name:       "Handle non-existent ID",
			svgContent: `<text id="other">Original</text>`,
			variables: CertificateTemplateVariables{
				Name: "Replaced",
			},
			shouldFind: "Original",
			expected:   "Original",
		},
		{
			name:       "Replace text in element with attributes",
			svgContent: `<text id="name" x="100" y="200" font-size="20">Original</text>`,
			variables: CertificateTemplateVariables{
				Name: "New Value",
			},
			shouldFind: "New Value",
			expected:   "New Value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceTemplateVariables(tt.svgContent, tt.variables)
			assert.Contains(t, result, tt.shouldFind)
		})
	}
}

func TestIntegrationReplacementAndRender(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

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

		// Replace variables
		modifiedSVG := replaceTemplateVariables(svgTemplate, templateVars)
		assert.Contains(t, modifiedSVG, "John Doe")
		assert.Contains(t, modifiedSVG, "Tech Conference 2024")

		// Render to PNG
		pngBytes, err := renderSVGToPNG(modifiedSVG)
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
		assert.Greater(t, len(pngBytes), 0)

		// Verify PNG signature
		assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, pngBytes[:4])
	})
}

// Benchmark tests
func BenchmarkReplaceSVGTemplateVariables(b *testing.B) {
	svgTemplate := strings.Repeat(`<text id="name">Placeholder</text>`, 10)

	templateVars := CertificateTemplateVariables{
		Name:      "John Doe",
		EventName: "Tech Conference",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = replaceTemplateVariables(svgTemplate, templateVars)
	}
}

func BenchmarkRenderSVGToPNG(b *testing.B) {
	svgContent := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600" viewBox="0 0 800 600">
	<rect width="800" height="600" fill="#f0f0f0"/>
	<text x="400" y="300" text-anchor="middle" font-size="24">Certificate</text>
</svg>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = renderSVGToPNG(svgContent)
	}
}
