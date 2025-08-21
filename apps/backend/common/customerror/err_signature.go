package customerror

import "net/http"

type ErrSignature struct {
	Code           ErrCode
	DefaultMessage string
	HttpStatus     int
}
type ErrCode string

var (
	ErrNotFound ErrSignature = ErrSignature{
		Code:           "NOT_FOUND",
		DefaultMessage: "Resource not found.",
		HttpStatus:     http.StatusNotFound,
	}
	ErrInternalServer ErrSignature = ErrSignature{
		Code:           "INTERNAL_SERVER_ERROR",
		DefaultMessage: "Something went wrong. Please try again later.",
		HttpStatus:     http.StatusInternalServerError,
	}
	ErrInvalidArgument ErrSignature = ErrSignature{
		Code:           "INVALID_ARGUMENT",
		DefaultMessage: "Invalid request. Please check your request and try again.",
		HttpStatus:     http.StatusBadRequest,
	}
	ErrUnauthorized ErrSignature = ErrSignature{
		Code:           "UNAUTHORIZED",
		DefaultMessage: "Unauthorized. Please login to continue.",
		HttpStatus:     http.StatusUnauthorized,
	}
	ErrInsufficientPermission ErrSignature = ErrSignature{
		Code:           "INSUFFICIENT_PERMISSION",
		DefaultMessage: "Insufficient permission. Please contact the administrator.",
		HttpStatus:     http.StatusForbidden,
	}
	ErrDuplicateEntry ErrSignature = ErrSignature{
		Code:           "DUPLICATE_ENTRY",
		DefaultMessage: "Duplicate entry. Your input already exists.",
		HttpStatus:     http.StatusConflict,
	}
)
