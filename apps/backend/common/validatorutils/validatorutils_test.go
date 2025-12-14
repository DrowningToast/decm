package validatorutils

import (
	"mime/multipart"
	"testing"

	"apps/backend/common/customerror"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required,min=3"`
	Age      int    `validate:"gte=0,lte=150"`
	Optional string `validate:"omitempty"`
}

func TestValidateStruct_Valid(t *testing.T) {
	s := TestStruct{
		Email: "test@example.com",
		Name:  "John Doe",
		Age:   25,
	}

	err := ValidateStruct(s)
	assert.NoError(t, err)
}

func TestValidateStruct_InvalidEmail(t *testing.T) {
	s := TestStruct{
		Email: "invalid-email",
		Name:  "John Doe",
		Age:   25,
	}

	err := ValidateStruct(s)
	assert.Error(t, err)
	assert.IsType(t, &customerror.Err{}, err)
}

func TestValidateStruct_MissingRequired(t *testing.T) {
	s := TestStruct{
		Email: "test@example.com",
		Name:  "", // Missing required
		Age:   25,
	}

	err := ValidateStruct(s)
	assert.Error(t, err)
}

func TestValidateStruct_InvalidAge(t *testing.T) {
	s := TestStruct{
		Email: "test@example.com",
		Name:  "John Doe",
		Age:   200, // Invalid: > 150
	}

	err := ValidateStruct(s)
	assert.Error(t, err)
}

func TestValidateStruct_OptionalField(t *testing.T) {
	s := TestStruct{
		Email:    "test@example.com",
		Name:     "John Doe",
		Age:      25,
		Optional: "", // Optional field can be empty
	}

	err := ValidateStruct(s)
	assert.NoError(t, err)
}

func TestValidateStruct_EmptyStruct(t *testing.T) {
	s := TestStruct{}

	err := ValidateStruct(s)
	assert.Error(t, err)
}

func TestValidateImageFile_ValidJPEG(t *testing.T) {
	file := &multipart.FileHeader{
		Size:   5 * 1024 * 1024, // 5MB
		Header: make(map[string][]string),
	}
	file.Header.Set("Content-Type", "image/jpeg")

	err := ValidateImageFile(file)
	assert.NoError(t, err)
}

func TestValidateImageFile_ValidPNG(t *testing.T) {
	file := &multipart.FileHeader{
		Size:   3 * 1024 * 1024, // 3MB
		Header: make(map[string][]string),
	}
	file.Header.Set("Content-Type", "image/png")

	err := ValidateImageFile(file)
	assert.NoError(t, err)
}

func TestValidateImageFile_ValidWebP(t *testing.T) {
	file := &multipart.FileHeader{
		Size:   2 * 1024 * 1024, // 2MB
		Header: make(map[string][]string),
	}
	file.Header.Set("Content-Type", "image/webp")

	err := ValidateImageFile(file)
	assert.NoError(t, err)
}

func TestValidateImageFile_ValidSVG(t *testing.T) {
	file := &multipart.FileHeader{
		Size:   1 * 1024 * 1024, // 1MB
		Header: make(map[string][]string),
	}
	file.Header.Set("Content-Type", "image/svg+xml")

	err := ValidateImageFile(file)
	assert.NoError(t, err)
}

func TestValidateImageFile_TooLarge(t *testing.T) {
	file := &multipart.FileHeader{
		Size:   15 * 1024 * 1024, // 15MB - exceeds 10MB limit
		Header: make(map[string][]string),
	}
	file.Header.Set("Content-Type", "image/jpeg")

	err := ValidateImageFile(file)
	assert.Error(t, err)
	assert.IsType(t, &customerror.Err{}, err)
	assert.Contains(t, err.Error(), "10MB")
}

func TestValidateImageFile_InvalidType(t *testing.T) {
	file := &multipart.FileHeader{
		Size:   5 * 1024 * 1024,
		Header: make(map[string][]string),
	}
	file.Header.Set("Content-Type", "application/pdf")

	err := ValidateImageFile(file)
	assert.Error(t, err)
	assert.IsType(t, &customerror.Err{}, err)
	assert.Contains(t, err.Error(), "invalid file type")
}

func TestValidateImageFile_InvalidTypeGIF(t *testing.T) {
	file := &multipart.FileHeader{
		Size:   5 * 1024 * 1024,
		Header: make(map[string][]string),
	}
	file.Header.Set("Content-Type", "image/gif")

	err := ValidateImageFile(file)
	assert.Error(t, err)
}

func TestValidateImageFile_ExactSizeLimit(t *testing.T) {
	file := &multipart.FileHeader{
		Size:   10 * 1024 * 1024, // Exactly 10MB
		Header: make(map[string][]string),
	}
	file.Header.Set("Content-Type", "image/jpeg")

	err := ValidateImageFile(file)
	assert.NoError(t, err)
}

func TestValidateImageFile_JustOverLimit(t *testing.T) {
	file := &multipart.FileHeader{
		Size:   (10*1024*1024 + 1), // Just over 10MB
		Header: make(map[string][]string),
	}
	file.Header.Set("Content-Type", "image/jpeg")

	err := ValidateImageFile(file)
	assert.Error(t, err)
}

func TestValidateImageFile_EmptyContentType(t *testing.T) {
	file := &multipart.FileHeader{
		Size:   5 * 1024 * 1024,
		Header: make(map[string][]string),
		// No Content-Type header
	}

	err := ValidateImageFile(file)
	assert.Error(t, err)
}

func TestValidateStruct_ComplexValidation(t *testing.T) {
	type ComplexStruct struct {
		Email    string `validate:"required,email"`
		URL      string `validate:"url"`
		MinLen   string `validate:"min=5"`
		MaxLen   string `validate:"max=10"`
		Numeric  string `validate:"numeric"`
		AlphaNum string `validate:"alphanum"`
	}

	tests := []struct {
		name      string
		input     ComplexStruct
		hasError  bool
		errorType string
	}{
		{
			name: "all valid",
			input: ComplexStruct{
				Email:    "test@example.com",
				URL:      "https://example.com",
				MinLen:   "12345",
				MaxLen:   "1234567890",
				Numeric:  "12345",
				AlphaNum: "abc123",
			},
			hasError: false,
		},
		{
			name: "invalid email",
			input: ComplexStruct{
				Email:    "invalid-email",
				URL:      "https://example.com",
				MinLen:   "12345",
				MaxLen:   "1234567890",
				Numeric:  "12345",
				AlphaNum: "abc123",
			},
			hasError: true,
		},
		{
			name: "too short",
			input: ComplexStruct{
				Email:    "test@example.com",
				URL:      "https://example.com",
				MinLen:   "1234", // Too short
				MaxLen:   "1234567890",
				Numeric:  "12345",
				AlphaNum: "abc123",
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateStruct_NonStruct(t *testing.T) {
	// This should not panic, but may not validate properly
	var s = "test"
	err := ValidateStruct(s)
	// validator should handle this gracefully
	_ = err
}

func TestValidateStruct_Pointer(t *testing.T) {
	s := &TestStruct{
		Email: "test@example.com",
		Name:  "John Doe",
		Age:   25,
	}

	err := ValidateStruct(s)
	assert.NoError(t, err)
}

func TestValidateImageFile_AllAllowedTypes(t *testing.T) {
	allowedTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
		"image/svg+xml",
	}

	for _, contentType := range allowedTypes {
		t.Run(contentType, func(t *testing.T) {
			file := &multipart.FileHeader{
				Size:   5 * 1024 * 1024,
				Header: make(map[string][]string),
			}
			file.Header.Set("Content-Type", contentType)

			err := ValidateImageFile(file)
			assert.NoError(t, err, "Content-Type %s should be allowed", contentType)
		})
	}
}







