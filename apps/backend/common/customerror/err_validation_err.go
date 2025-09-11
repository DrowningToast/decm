package customerror

import (
	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
)

func ParseValidationErr(err *validator.ValidationErrors) *Err {
	kvErrs := make(map[string]string)
	for _, validationErr := range *err {
		kvErrs[validationErr.Field()] = validationErr.Error()
	}

	return ParseWithReasons(&ErrInvalidArgument, errors.Wrap(err, "validation error"), kvErrs)
}
