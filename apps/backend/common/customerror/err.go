package customerror

import (
	"log/slog"

	"github.com/cockroachdb/errors"
)

// @name Err
// @description Custom error type
type Err struct {
	HttpStatus  *int              `json:"http_status"`
	Code        *ErrCode          `json:"code"`
	LoggerLevel slog.Level        `json:"-"`
	Message     string            `json:"message,omitempty"`
	Reasons     map[string]string `json:"reasons,omitempty"`
	Inner       error             `json:"-"`
}

func (e Err) Error() string {
	return e.Message + e.Inner.Error()
}

func (err Err) ExtendWithError(e error) *Err {
	return &Err{
		HttpStatus:  err.HttpStatus,
		Code:        err.Code,
		Message:     errors.Wrap(err, e.Error()).Error(),
		Inner:       e,
		LoggerLevel: err.LoggerLevel,
	}
}

func (err Err) Extend(msg string) *Err {
	return &Err{
		HttpStatus:  err.HttpStatus,
		Code:        err.Code,
		Message:     errors.Wrap(err, msg).Error(),
		Inner:       err,
		LoggerLevel: err.LoggerLevel,
	}
}

func Parse(preset *ErrSignature, err error) *Err {
	var defaultErrSignature *ErrSignature = &ErrInternalServer
	if preset != nil {
		defaultErrSignature = &ErrSignature{
			Code:           preset.Code,
			DefaultMessage: preset.DefaultMessage,
			HttpStatus:     preset.HttpStatus,
			LoggerLevel:    preset.LoggerLevel,
		}
	}

	return AsPresetError(*defaultErrSignature, err)
}

func ParseWithMessage(preset *ErrSignature, err error, message string) *Err {
	var defaultErrSignature *ErrSignature = &ErrInternalServer
	if preset != nil {
		defaultErrSignature = &ErrSignature{
			Code:           preset.Code,
			DefaultMessage: preset.DefaultMessage,
			HttpStatus:     preset.HttpStatus,
			LoggerLevel:    preset.LoggerLevel,
		}
	}

	return &Err{
		HttpStatus:  &defaultErrSignature.HttpStatus,
		Code:        &defaultErrSignature.Code,
		Message:     message,
		Inner:       err,
		LoggerLevel: defaultErrSignature.LoggerLevel,
	}
}

func ParseWithReasons(preset *ErrSignature, err error, reasons map[string]string) *Err {
	if reasons == nil {
		return Parse(&ErrInternalServer, err)
	}
	var defaultErrSignature *ErrSignature = &ErrInternalServer
	if preset != nil {
		defaultErrSignature = &ErrSignature{
			Code:           preset.Code,
			DefaultMessage: preset.DefaultMessage,
			HttpStatus:     preset.HttpStatus,
			LoggerLevel:    preset.LoggerLevel,
		}
	}

	return &Err{
		HttpStatus:  &defaultErrSignature.HttpStatus,
		Code:        &defaultErrSignature.Code,
		Message:     defaultErrSignature.DefaultMessage,
		Reasons:     reasons,
		Inner:       err,
		LoggerLevel: defaultErrSignature.LoggerLevel,
	}
}

func AsPresetError(preset ErrSignature, err error) *Err {
	httpStatus := preset.HttpStatus
	code := preset.Code

	return &Err{
		HttpStatus:  &httpStatus,
		Code:        &code,
		Message:     preset.DefaultMessage,
		Inner:       err,
		LoggerLevel: preset.LoggerLevel,
	}
}
