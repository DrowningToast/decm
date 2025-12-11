package event

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCertificateData(t *testing.T) {
	t.Run("Create with required fields only", func(t *testing.T) {
		data := CertificateData{
			Name:      "John Doe",
			EventName: "Tech Conference 2024",
		}

		assert.Equal(t, "John Doe", data.Name)
		assert.Equal(t, "Tech Conference 2024", data.EventName)
		assert.Empty(t, data.AcademicInstitution)
		assert.Empty(t, data.CertificateTitle)
		assert.Empty(t, data.CertificateSubtitle)
	})

	t.Run("Create with all fields", func(t *testing.T) {
		data := CertificateData{
			Name:                "Jane Smith",
			EventName:           "AI Summit",
			AcademicInstitution: "MIT",
			CertificateTitle:    "Certificate of Achievement",
			CertificateSubtitle: "For Outstanding Performance",
		}

		assert.Equal(t, "Jane Smith", data.Name)
		assert.Equal(t, "AI Summit", data.EventName)
		assert.Equal(t, "MIT", data.AcademicInstitution)
		assert.Equal(t, "Certificate of Achievement", data.CertificateTitle)
		assert.Equal(t, "For Outstanding Performance", data.CertificateSubtitle)
	})
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

	t.Run("Render SVG with certificate content", func(t *testing.T) {
		// Create SVG with certificate content
		svgTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600" viewBox="0 0 800 600">
	<rect width="800" height="600" fill="#ffffff"/>
	<text x="400" y="100" text-anchor="middle" font-size="32" fill="#000000">Certificate of Excellence</text>
	<text x="400" y="250" text-anchor="middle" font-size="48" fill="#1a1a1a">John Doe</text>
	<text x="400" y="350" text-anchor="middle" font-size="24" fill="#333333">Tech Conference 2024</text>
	<text x="400" y="450" text-anchor="middle" font-size="20" fill="#666666">Massachusetts Institute of Technology</text>
</svg>`

		// Render to PNG
		pngBytes, err := renderSVGToPNG(svgTemplate)
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
		assert.Greater(t, len(pngBytes), 10000, "PNG should be reasonably sized for 800x600 certificate")

		// Check PNG magic bytes
		assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, pngBytes[:4])
	})
}

func TestCreateTextElement(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		x          float64
		y          float64
		fontFamily string
		fontWeight string
		fontSize   int
		expected   []string
	}{
		{
			name:       "Create simple text element",
			text:       "John Doe",
			x:          400.0,
			y:          300.0,
			fontFamily: "Inter",
			fontWeight: "bold",
			fontSize:   16,
			expected:   []string{"John Doe", "x=\"400.00\"", "y=\"300.00\"", "font-family=\"Inter\"", "font-weight=\"bold\"", "font-size=\"16\""},
		},
		{
			name:       "Handle special XML characters",
			text:       "John & Jane <Smith>",
			x:          100.0,
			y:          200.0,
			fontFamily: "Prompt",
			fontWeight: "400",
			fontSize:   20,
			expected:   []string{"John &amp; Jane &lt;Smith&gt;", "x=\"100.00\"", "y=\"200.00\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createTextElement(tt.text, tt.x, tt.y, tt.fontFamily, tt.fontWeight, tt.fontSize)
			for _, exp := range tt.expected {
				assert.Contains(t, result, exp)
			}
		})
	}
}

func TestHideTemplatePlaceholders(t *testing.T) {
	t.Run("Hide template placeholder by ID", func(t *testing.T) {
		svgTemplate := `<svg>
	<text id="{{ name }}">Some text</text>
	<text id="name">Other text</text>
	<text id="eventName">Event</text>
</svg>`

		result := hideTemplatePlaceholders(svgTemplate)
		assert.Contains(t, result, `id="{{ name }}"`)
		assert.Contains(t, result, `visibility="hidden"`)
		// Check that both ID variations are hidden
		assert.Contains(t, result, `<text id="{{ name }}" visibility="hidden"`)
		assert.Contains(t, result, `<text id="name" visibility="hidden"`)
	})

	t.Run("Hide template placeholder by content", func(t *testing.T) {
		svgTemplate := `<svg>
	<text x="100" y="200">{{ name }}</text>
	<text x="200" y="300">{{ eventName }}</text>
	<text x="300" y="400">Regular text</text>
</svg>`

		result := hideTemplatePlaceholders(svgTemplate)
		// Check that placeholders in content are hidden
		assert.Contains(t, result, `visibility="hidden"`)
		// Regular text should not be hidden
		assert.Contains(t, result, "Regular text")
		assert.NotContains(t, result, `<text x="300" y="400" visibility="hidden"`)
	})

	t.Run("Hide template placeholder with both ID and content", func(t *testing.T) {
		svgTemplate := `<svg>
	<text id="name">{{ name }}</text>
	<text id="eventName">{{ eventName }}</text>
</svg>`

		result := hideTemplatePlaceholders(svgTemplate)
		// Should be hidden (by ID matching)
		assert.Contains(t, result, `visibility="hidden"`)
		// Should not be processed twice
		hiddenCount := strings.Count(result, `visibility="hidden"`)
		assert.Equal(t, 2, hiddenCount, "Should hide exactly 2 elements")
	})

	t.Run("Handle all placeholder types", func(t *testing.T) {
		svgTemplate := `<svg>
	<text id="name">{{ name }}</text>
	<text id="eventName">{{ eventName }}</text>
	<text id="academicInstitutionName">{{ academicInstitutionName }}</text>
	<text id="certificateTitle">{{ certificateTitle }}</text>
	<text id="certificateSubtitle">{{ certificateSubtitle }}</text>
</svg>`

		result := hideTemplatePlaceholders(svgTemplate)
		// All placeholders should be hidden
		hiddenCount := strings.Count(result, `visibility="hidden"`)
		assert.Equal(t, 5, hiddenCount, "Should hide all 5 placeholder elements")
	})
}

// Benchmark tests
func BenchmarkCreateTextElement(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = createTextElement("John Doe", 400.0, 300.0, "Inter", "bold", 16)
	}
}

func BenchmarkHideTemplatePlaceholders(b *testing.B) {
	svgTemplate := strings.Repeat(`<text id="name">{{ name }}</text>`, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hideTemplatePlaceholders(svgTemplate)
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
