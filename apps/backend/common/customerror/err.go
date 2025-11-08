package customerror

import (
	"log/slog"

	"github.com/cockroachdb/errors"
)

// @name Err
// @description Custom error type
type Err struct {
	HttpStatus  *int              `json:"http_status,omitempty"`
	Code        *ErrCode          `json:"code,omitempty"`
	LoggerLevel slog.Level        `json:"-"`
	Message     string            `json:"message,omitempty"`
	Reasons     map[string]string `json:"reasons,omitempty"`
	Inner       error             `json:"-"`
}

func (e *Err) Error() string {
	if e == nil {
		return "nil custom error"
	}
	if e.Inner != nil {
		return e.Message + ": " + e.Inner.Error()
	}
	return e.Message
}

func (err *Err) ExtendWithError(e error) *Err {
	if err == nil {
		return nil
	}
	return &Err{
		HttpStatus:  err.HttpStatus,
		Code:        err.Code,
		Message:     errors.Wrap(err, e.Error()).Error(),
		Inner:       e,
		LoggerLevel: err.LoggerLevel,
	}
}

func (err *Err) Extend(msg string) *Err {
	if err == nil {
		return nil
	}
	return &Err{
		HttpStatus:  err.HttpStatus,
		Code:        err.Code,
		Message:     errors.Wrap(err, msg).Error(),
		Inner:       err,
		LoggerLevel: err.LoggerLevel,
	}
}

func New(httpStatus int, code ErrCode, message string, inner error) *Err {
	return &Err{
		HttpStatus:  &httpStatus,
		Code:        &code,
		Message:     message,
		Inner:       inner,
		LoggerLevel: slog.LevelError,
	}
}

func NewWithPreset(preset *ErrSignature, err error) *Err {
	return &Err{
		HttpStatus:  &preset.HttpStatus,
		Code:        &preset.Code,
		Message:     preset.DefaultMessage,
		Inner:       err,
		LoggerLevel: preset.LoggerLevel,
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

func TryParseAsCustomErr(err error) *Err {
	// Check if err is already a customerror type
	var customErr *Err
	if errors.As(err, &customErr) {
		return customErr
	}
	// Return nil if it's not a custom error
	return nil
}
