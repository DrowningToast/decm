package validatorutils

import (
	"errors"

	"apps/backend/common/customerror"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(s interface{}) *customerror.Err {
	if err := validator.New().Struct(s); err != nil {
		var validationErr *validator.ValidationErrors
		if errors.As(err, &validationErr) {
			return customerror.ParseValidationErr(validationErr)
		}
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}
