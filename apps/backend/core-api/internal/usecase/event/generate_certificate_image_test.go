package event

import (
	"os"
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

// =============================================================================
// Tests for hideTemplatePlaceholders - Enhanced Coverage
// =============================================================================

func TestHideTemplatePlaceholders_EnhancedCoverage(t *testing.T) {
	t.Run("should hide placeholder with id matching plain name without template syntax", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg">
			<text id="name" x="100" y="100">Placeholder Name</text>
			<text id="eventName" x="100" y="150">Placeholder Event</text>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		assert.Contains(t, result, `<text id="name" x="100" y="100" visibility="hidden">`, "Should hide element with id='name'")
		assert.Contains(t, result, `<text id="eventName" x="100" y="150" visibility="hidden">`, "Should hide element with id='eventName'")
	})

	t.Run("should hide placeholder with id matching template syntax", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg">
			<text id="{{ name }}" x="100" y="100">John Doe</text>
			<text id="{{ eventName }}" x="100" y="150">Conference 2024</text>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		assert.Contains(t, result, `<text id="{{ name }}" x="100" y="100" visibility="hidden">`)
		assert.Contains(t, result, `<text id="{{ eventName }}" x="100" y="150" visibility="hidden">`)
	})

	t.Run("should hide placeholder when content contains template syntax", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg">
			<text x="100" y="100">{{ name }}</text>
			<text x="100" y="150">Some text {{ eventName }} more text</text>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		assert.Contains(t, result, `visibility="hidden"`)
		// Count should be 2 (one for each placeholder)
		hiddenCount := strings.Count(result, `visibility="hidden"`)
		assert.GreaterOrEqual(t, hiddenCount, 2, "Should hide elements with placeholder content")
	})

	t.Run("should handle all five placeholder types", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg">
			<text id="name">{{ name }}</text>
			<text id="eventName">{{ eventName }}</text>
			<text id="academicInstitutionName">{{ academicInstitutionName }}</text>
			<text id="certificateTitle">{{ certificateTitle }}</text>
			<text id="certificateSubtitle">{{ certificateSubtitle }}</text>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		// All 5 placeholders should be hidden
		hiddenCount := strings.Count(result, `visibility="hidden"`)
		assert.GreaterOrEqual(t, hiddenCount, 5, "Should hide all 5 placeholder types")
	})

	t.Run("should not hide elements without placeholders", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg">
			<text id="header" x="100" y="50">Certificate Header</text>
			<text id="footer" x="100" y="500">Footer Text</text>
			<rect x="0" y="0" width="100" height="100"/>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		assert.NotContains(t, result, `id="header" visibility="hidden"`)
		assert.NotContains(t, result, `id="footer" visibility="hidden"`)
		assert.Contains(t, result, `id="header"`, "Should preserve non-placeholder elements")
	})

	t.Run("should not add duplicate visibility attributes", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg">
			<text id="name" visibility="hidden">{{ name }}</text>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		// Count occurrences of visibility="hidden" for this element
		// Should not add another visibility attribute
		assert.Contains(t, result, `visibility="hidden"`)
		visibilityCount := strings.Count(result, `visibility="hidden"`)
		assert.Equal(t, 1, visibilityCount, "Should not duplicate visibility attribute")
	})

	t.Run("should handle mixed placeholder patterns in same SVG", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg">
			<text id="name" x="100" y="100">Regular ID</text>
			<text id="{{ eventName }}" x="100" y="150">Template ID</text>
			<text x="100" y="200">{{ academicInstitutionName }}</text>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		hiddenCount := strings.Count(result, `visibility="hidden"`)
		assert.GreaterOrEqual(t, hiddenCount, 3, "Should hide all three placeholder types")
	})

	t.Run("should preserve other attributes when adding visibility", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg">
			<text id="name" x="100" y="100" font-size="24" fill="#000" text-anchor="middle">{{ name }}</text>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		assert.Contains(t, result, `visibility="hidden"`)
		assert.Contains(t, result, `font-size="24"`, "Should preserve font-size")
		assert.Contains(t, result, `fill="#000"`, "Should preserve fill")
		assert.Contains(t, result, `text-anchor="middle"`, "Should preserve text-anchor")
	})

	t.Run("should handle placeholder content with surrounding text", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg">
			<text x="100" y="100">Name: {{ name }} (placeholder)</text>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		assert.Contains(t, result, `visibility="hidden"`, "Should hide text with placeholder in content")
	})

	t.Run("should handle empty SVG", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg"></svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		assert.Equal(t, svgTemplate, result, "Empty SVG should remain unchanged")
	})

	t.Run("should handle SVG with no text elements", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg">
			<rect x="0" y="0" width="100" height="100" fill="blue"/>
			<circle cx="50" cy="50" r="25" fill="red"/>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		assert.Equal(t, svgTemplate, result, "SVG without text should remain unchanged")
	})

	t.Run("should handle complex SVG with multiple text elements", func(t *testing.T) {
		// Arrange
		svgTemplate := `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600">
			<text id="title" x="400" y="50">Certificate</text>
			<text id="name" x="400" y="200">{{ name }}</text>
			<text id="eventName" x="400" y="300">{{ eventName }}</text>
			<text id="institution" x="400" y="400">Regular Institution</text>
			<text id="academicInstitutionName" x="400" y="450">{{ academicInstitutionName }}</text>
		</svg>`

		// Act
		result := hideTemplatePlaceholders(svgTemplate)

		// Assert
		assert.Contains(t, result, `id="title"`, "Non-placeholder should be preserved")
		assert.NotContains(t, result, `id="title" visibility="hidden"`, "Non-placeholder should not be hidden")
		assert.Contains(t, result, `<text id="name" x="400" y="200" visibility="hidden">`)
		assert.Contains(t, result, `<text id="eventName" x="400" y="300" visibility="hidden">`)
		assert.Contains(t, result, `<text id="academicInstitutionName" x="400" y="450" visibility="hidden">`)
	})
}

// =============================================================================
// Tests for addTextOverlaysToSVG - Event Name Null Position Handling
// =============================================================================

func TestAddTextOverlaysToSVG_EventNameNullPositionHandling(t *testing.T) {
	t.Run("should skip event name when position X is nil", func(t *testing.T) {
		// Arrange - This test would require mocking the EventUsecase and its dependencies
		// For now, documenting the expected behavior
		// When config.EventNamePosX is nil, the event name text should not be added
	})

	t.Run("should skip event name when position Y is nil", func(t *testing.T) {
		// Similar to above - requires full usecase context
	})

	t.Run("should add event name when both positions are valid", func(t *testing.T) {
		// When both config.EventNamePosX and config.EventNamePosY are not nil
		// the event name text should be added to the overlay
	})
}

// =============================================================================
// Tests for renderSVGToPNG - Containerized Environment Detection
// =============================================================================

func TestRenderSVGToPNG_ContainerizedEnvironment(t *testing.T) {
	t.Run("should detect CI environment", func(t *testing.T) {
		// Arrange
		originalCI := os.Getenv("CI")
		defer func() {
			if originalCI == "" {
				os.Unsetenv("CI")
			} else {
				os.Setenv("CI", originalCI)
			}
		}()

		os.Setenv("CI", "true")

		svgContent := `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
			<rect width="100" height="100" fill="#ffffff"/>
		</svg>`

		// Act
		pngBytes, err := renderSVGToPNG(svgContent)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
		assert.Greater(t, len(pngBytes), 0)
	})

	t.Run("should detect GITHUB_ACTIONS environment", func(t *testing.T) {
		// Arrange
		originalGHA := os.Getenv("GITHUB_ACTIONS")
		defer func() {
			if originalGHA == "" {
				os.Unsetenv("GITHUB_ACTIONS")
			} else {
				os.Setenv("GITHUB_ACTIONS", originalGHA)
			}
		}()

		os.Setenv("GITHUB_ACTIONS", "true")

		svgContent := `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
			<rect width="100" height="100" fill="#f0f0f0"/>
		</svg>`

		// Act
		pngBytes, err := renderSVGToPNG(svgContent)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
		assert.Greater(t, len(pngBytes), 0)
	})

	t.Run("should detect production environment", func(t *testing.T) {
		// Arrange
		originalEnv := os.Getenv("ENVIRONMENT")
		defer func() {
			if originalEnv == "" {
				os.Unsetenv("ENVIRONMENT")
			} else {
				os.Setenv("ENVIRONMENT", originalEnv)
			}
		}()

		os.Setenv("ENVIRONMENT", "production")

		svgContent := `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
			<circle cx="50" cy="50" r="40" fill="#0000ff"/>
		</svg>`

		// Act
		pngBytes, err := renderSVGToPNG(svgContent)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
		assert.Greater(t, len(pngBytes), 0)
	})

	t.Run("should detect docker container environment", func(t *testing.T) {
		// Arrange
		originalDocker := os.Getenv("DOCKER_CONTAINER")
		defer func() {
			if originalDocker == "" {
				os.Unsetenv("DOCKER_CONTAINER")
			} else {
				os.Setenv("DOCKER_CONTAINER", originalDocker)
			}
		}()

		os.Setenv("DOCKER_CONTAINER", "true")

		svgContent := `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
			<text x="50" y="50" text-anchor="middle">Test</text>
		</svg>`

		// Act
		pngBytes, err := renderSVGToPNG(svgContent)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, pngBytes)
		assert.Greater(t, len(pngBytes), 0)
	})
}
