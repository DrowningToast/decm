package customerror

import (
	"log/slog"
	"net/http"
)

type ErrSignature struct {
	Code           ErrCode
	DefaultMessage string
	HttpStatus     int
	LoggerLevel    slog.Level
}

// @name ErrCode
// @description Error code type
type ErrCode string

var (
	ErrNotFound ErrSignature = ErrSignature{
		Code:           "NOT_FOUND",
		DefaultMessage: "Resource not found.",
		HttpStatus:     http.StatusNotFound,
		LoggerLevel:    slog.LevelWarn,
	}
	ErrInternalServer ErrSignature = ErrSignature{
		Code:           "INTERNAL_SERVER_ERROR",
		DefaultMessage: "Something went wrong. Please try again later.",
		HttpStatus:     http.StatusInternalServerError,
		LoggerLevel:    slog.LevelError,
	}
	ErrInvalidArgument ErrSignature = ErrSignature{
		Code:           "INVALID_ARGUMENT",
		DefaultMessage: "Invalid request. Please check your request and try again.",
		HttpStatus:     http.StatusBadRequest,
		LoggerLevel:    slog.LevelWarn,
	}
	ErrUnauthenticated ErrSignature = ErrSignature{
		Code:           "UNAUTHORIZED",
		DefaultMessage: "Unauthorized. Please login to continue.",
		HttpStatus:     http.StatusUnauthorized,
		LoggerLevel:    slog.LevelWarn,
	}
	ErrUnauthorized ErrSignature = ErrSignature{
		Code:           "INSUFFICIENT_PERMISSION",
		DefaultMessage: "Insufficient permission. Please contact the administrator.",
		HttpStatus:     http.StatusForbidden,
		LoggerLevel:    slog.LevelWarn,
	}
	ErrDuplicateEntry ErrSignature = ErrSignature{
		Code:           "DUPLICATE_ENTRY",
		DefaultMessage: "Duplicate entry. Your input already exists.",
		HttpStatus:     http.StatusConflict,
		LoggerLevel:    slog.LevelWarn,
	}
)
