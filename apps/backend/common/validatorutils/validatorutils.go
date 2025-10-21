package validatorutils

import (
	"errors"
	"fmt"
	"mime/multipart"

	"apps/backend/common/customerror"
	"apps/backend/common/utils"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(s interface{}) error {
	if err := validator.New().Struct(s); err != nil {
		var validationErr *validator.ValidationErrors
		if errors.As(err, &validationErr) {
			return customerror.ParseValidationErr(validationErr)
		}
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func ValidateImageFile(file *multipart.FileHeader) error {
	// Max 5MB
	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return customerror.Parse(
			&customerror.ErrInvalidArgument,
			errors.New("banner file must be less than 10MB"),
		)
	}

	// Check file type
	allowedTypes := map[string]bool{
		"image/jpeg":    true,
		"image/jpg":     true,
		"image/png":     true,
		"image/webp":    true,
		"image/svg+xml": true,
	}

	contentType := utils.GetFileContentType(file)
	if !allowedTypes[contentType] {
		return customerror.Parse(
			&customerror.ErrInvalidArgument,
			fmt.Errorf("invalid file type '%s': banner must be JPEG, PNG, WebP, or SVG", contentType),
		)
	}

	return nil
}
