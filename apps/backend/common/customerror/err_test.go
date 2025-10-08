package customerror

import (
	"errors"
	"log/slog"
	"net/http"
	"testing"
)

func TestErr_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Err
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "nil custom error",
		},
		{
			name: "error without inner error",
			err: &Err{
				Message: "test error message",
			},
			expected: "test error message",
		},
		{
			name: "error with inner error",
			err: &Err{
				Message: "outer error",
				Inner:   errors.New("inner error"),
			},
			expected: "outer error: inner error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErr_ExtendWithError(t *testing.T) {
	httpStatus := http.StatusBadRequest
	code := ErrCode("TEST_CODE")

	tests := []struct {
		name         string
		err          *Err
		extendErr    error
		expectNil    bool
		expectCode   *ErrCode
		expectStatus *int
		expectInner  error
		expectLevel  slog.Level
	}{
		{
			name:      "extend nil error",
			err:       nil,
			extendErr: errors.New("test error"),
			expectNil: true,
		},
		{
			name: "extend with error",
			err: &Err{
				HttpStatus:  &httpStatus,
				Code:        &code,
				Message:     "original message",
				LoggerLevel: slog.LevelError,
			},
			extendErr:    errors.New("new error"),
			expectNil:    false,
			expectCode:   &code,
			expectStatus: &httpStatus,
			expectLevel:  slog.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.ExtendWithError(tt.extendErr)
			if tt.expectNil {
				if got != nil {
					t.Errorf("ExtendWithError() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("ExtendWithError() returned nil, expected non-nil")
			}

			if got.Code == nil || *got.Code != *tt.expectCode {
				t.Errorf("ExtendWithError() Code = %v, want %v", got.Code, tt.expectCode)
			}

			if got.HttpStatus == nil || *got.HttpStatus != *tt.expectStatus {
				t.Errorf("ExtendWithError() HttpStatus = %v, want %v", got.HttpStatus, tt.expectStatus)
			}

			if got.LoggerLevel != tt.expectLevel {
				t.Errorf("ExtendWithError() LoggerLevel = %v, want %v", got.LoggerLevel, tt.expectLevel)
			}

			if got.Inner != tt.extendErr {
				t.Errorf("ExtendWithError() Inner = %v, want %v", got.Inner, tt.extendErr)
			}

			if got.Message == "" {
				t.Error("ExtendWithError() Message is empty")
			}
		})
	}
}

func TestErr_Extend(t *testing.T) {
	httpStatus := http.StatusBadRequest
	code := ErrCode("TEST_CODE")

	tests := []struct {
		name         string
		err          *Err
		msg          string
		expectNil    bool
		expectCode   *ErrCode
		expectStatus *int
		expectLevel  slog.Level
	}{
		{
			name:      "extend nil error",
			err:       nil,
			msg:       "test message",
			expectNil: true,
		},
		{
			name: "extend with message",
			err: &Err{
				HttpStatus:  &httpStatus,
				Code:        &code,
				Message:     "original message",
				LoggerLevel: slog.LevelWarn,
			},
			msg:          "extended message",
			expectNil:    false,
			expectCode:   &code,
			expectStatus: &httpStatus,
			expectLevel:  slog.LevelWarn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Extend(tt.msg)
			if tt.expectNil {
				if got != nil {
					t.Errorf("Extend() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("Extend() returned nil, expected non-nil")
			}

			if got.Code == nil || *got.Code != *tt.expectCode {
				t.Errorf("Extend() Code = %v, want %v", got.Code, tt.expectCode)
			}

			if got.HttpStatus == nil || *got.HttpStatus != *tt.expectStatus {
				t.Errorf("Extend() HttpStatus = %v, want %v", got.HttpStatus, tt.expectStatus)
			}

			if got.LoggerLevel != tt.expectLevel {
				t.Errorf("Extend() LoggerLevel = %v, want %v", got.LoggerLevel, tt.expectLevel)
			}

			if got.Message == "" {
				t.Error("Extend() Message is empty")
			}
		})
	}
}

func TestParse(t *testing.T) {
	testErr := errors.New("test error")

	tests := []struct {
		name         string
		preset       *ErrSignature
		err          error
		expectCode   ErrCode
		expectStatus int
		expectLevel  slog.Level
	}{
		{
			name:         "parse with nil preset uses ErrInternalServer",
			preset:       nil,
			err:          testErr,
			expectCode:   ErrInternalServer.Code,
			expectStatus: ErrInternalServer.HttpStatus,
			expectLevel:  ErrInternalServer.LoggerLevel,
		},
		{
			name:         "parse with custom preset",
			preset:       &ErrNotFound,
			err:          testErr,
			expectCode:   ErrNotFound.Code,
			expectStatus: ErrNotFound.HttpStatus,
			expectLevel:  ErrNotFound.LoggerLevel,
		},
		{
			name:         "parse with invalid argument preset",
			preset:       &ErrInvalidArgument,
			err:          testErr,
			expectCode:   ErrInvalidArgument.Code,
			expectStatus: ErrInvalidArgument.HttpStatus,
			expectLevel:  ErrInvalidArgument.LoggerLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.preset, tt.err)

			if got == nil {
				t.Fatal("Parse() returned nil")
			}

			if got.Code == nil || *got.Code != tt.expectCode {
				t.Errorf("Parse() Code = %v, want %v", got.Code, tt.expectCode)
			}

			if got.HttpStatus == nil || *got.HttpStatus != tt.expectStatus {
				t.Errorf("Parse() HttpStatus = %v, want %v", got.HttpStatus, tt.expectStatus)
			}

			if got.LoggerLevel != tt.expectLevel {
				t.Errorf("Parse() LoggerLevel = %v, want %v", got.LoggerLevel, tt.expectLevel)
			}

			if got.Inner != tt.err {
				t.Errorf("Parse() Inner = %v, want %v", got.Inner, tt.err)
			}
		})
	}
}

func TestParseWithMessage(t *testing.T) {
	testErr := errors.New("test error")
	customMessage := "custom error message"

	tests := []struct {
		name         string
		preset       *ErrSignature
		err          error
		message      string
		expectCode   ErrCode
		expectStatus int
		expectLevel  slog.Level
		expectMsg    string
	}{
		{
			name:         "parse with message and nil preset",
			preset:       nil,
			err:          testErr,
			message:      customMessage,
			expectCode:   ErrInternalServer.Code,
			expectStatus: ErrInternalServer.HttpStatus,
			expectLevel:  ErrInternalServer.LoggerLevel,
			expectMsg:    customMessage,
		},
		{
			name:         "parse with message and custom preset",
			preset:       &ErrNotFound,
			err:          testErr,
			message:      customMessage,
			expectCode:   ErrNotFound.Code,
			expectStatus: ErrNotFound.HttpStatus,
			expectLevel:  ErrNotFound.LoggerLevel,
			expectMsg:    customMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseWithMessage(tt.preset, tt.err, tt.message)

			if got == nil {
				t.Fatal("ParseWithMessage() returned nil")
			}

			if got.Code == nil || *got.Code != tt.expectCode {
				t.Errorf("ParseWithMessage() Code = %v, want %v", got.Code, tt.expectCode)
			}

			if got.HttpStatus == nil || *got.HttpStatus != tt.expectStatus {
				t.Errorf("ParseWithMessage() HttpStatus = %v, want %v", got.HttpStatus, tt.expectStatus)
			}

			if got.LoggerLevel != tt.expectLevel {
				t.Errorf("ParseWithMessage() LoggerLevel = %v, want %v", got.LoggerLevel, tt.expectLevel)
			}

			if got.Message != tt.expectMsg {
				t.Errorf("ParseWithMessage() Message = %v, want %v", got.Message, tt.expectMsg)
			}

			if got.Inner != tt.err {
				t.Errorf("ParseWithMessage() Inner = %v, want %v", got.Inner, tt.err)
			}
		})
	}
}

func TestParseWithReasons(t *testing.T) {
	testErr := errors.New("test error")
	reasons := map[string]string{
		"field1": "error1",
		"field2": "error2",
	}

	tests := []struct {
		name         string
		preset       *ErrSignature
		err          error
		reasons      map[string]string
		expectCode   ErrCode
		expectStatus int
		expectLevel  slog.Level
	}{
		{
			name:         "parse with nil reasons falls back to Parse",
			preset:       &ErrNotFound,
			err:          testErr,
			reasons:      nil,
			expectCode:   ErrInternalServer.Code,
			expectStatus: ErrInternalServer.HttpStatus,
			expectLevel:  ErrInternalServer.LoggerLevel,
		},
		{
			name:         "parse with reasons and nil preset",
			preset:       nil,
			err:          testErr,
			reasons:      reasons,
			expectCode:   ErrInternalServer.Code,
			expectStatus: ErrInternalServer.HttpStatus,
			expectLevel:  ErrInternalServer.LoggerLevel,
		},
		{
			name:         "parse with reasons and custom preset",
			preset:       &ErrInvalidArgument,
			err:          testErr,
			reasons:      reasons,
			expectCode:   ErrInvalidArgument.Code,
			expectStatus: ErrInvalidArgument.HttpStatus,
			expectLevel:  ErrInvalidArgument.LoggerLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseWithReasons(tt.preset, tt.err, tt.reasons)

			if got == nil {
				t.Fatal("ParseWithReasons() returned nil")
			}

			if got.Code == nil || *got.Code != tt.expectCode {
				t.Errorf("ParseWithReasons() Code = %v, want %v", got.Code, tt.expectCode)
			}

			if got.HttpStatus == nil || *got.HttpStatus != tt.expectStatus {
				t.Errorf("ParseWithReasons() HttpStatus = %v, want %v", got.HttpStatus, tt.expectStatus)
			}

			if got.LoggerLevel != tt.expectLevel {
				t.Errorf("ParseWithReasons() LoggerLevel = %v, want %v", got.LoggerLevel, tt.expectLevel)
			}

			if tt.reasons != nil {
				if len(got.Reasons) != len(tt.reasons) {
					t.Errorf("ParseWithReasons() Reasons length = %v, want %v", len(got.Reasons), len(tt.reasons))
				}

				for key, value := range tt.reasons {
					if got.Reasons[key] != value {
						t.Errorf("ParseWithReasons() Reasons[%s] = %v, want %v", key, got.Reasons[key], value)
					}
				}
			}

			if got.Inner != tt.err {
				t.Errorf("ParseWithReasons() Inner = %v, want %v", got.Inner, tt.err)
			}
		})
	}
}

func TestAsPresetError(t *testing.T) {
	testErr := errors.New("test error")

	tests := []struct {
		name         string
		preset       ErrSignature
		err          error
		expectCode   ErrCode
		expectStatus int
		expectLevel  slog.Level
		expectMsg    string
	}{
		{
			name:         "create error from ErrNotFound preset",
			preset:       ErrNotFound,
			err:          testErr,
			expectCode:   ErrNotFound.Code,
			expectStatus: ErrNotFound.HttpStatus,
			expectLevel:  ErrNotFound.LoggerLevel,
			expectMsg:    ErrNotFound.DefaultMessage,
		},
		{
			name:         "create error from ErrInternalServer preset",
			preset:       ErrInternalServer,
			err:          testErr,
			expectCode:   ErrInternalServer.Code,
			expectStatus: ErrInternalServer.HttpStatus,
			expectLevel:  ErrInternalServer.LoggerLevel,
			expectMsg:    ErrInternalServer.DefaultMessage,
		},
		{
			name:         "create error from ErrInvalidArgument preset",
			preset:       ErrInvalidArgument,
			err:          testErr,
			expectCode:   ErrInvalidArgument.Code,
			expectStatus: ErrInvalidArgument.HttpStatus,
			expectLevel:  ErrInvalidArgument.LoggerLevel,
			expectMsg:    ErrInvalidArgument.DefaultMessage,
		},
		{
			name:         "create error from ErrUnauthenticated preset",
			preset:       ErrUnauthenticated,
			err:          testErr,
			expectCode:   ErrUnauthenticated.Code,
			expectStatus: ErrUnauthenticated.HttpStatus,
			expectLevel:  ErrUnauthenticated.LoggerLevel,
			expectMsg:    ErrUnauthenticated.DefaultMessage,
		},
		{
			name:         "create error from ErrUnauthorized preset",
			preset:       ErrUnauthorized,
			err:          testErr,
			expectCode:   ErrUnauthorized.Code,
			expectStatus: ErrUnauthorized.HttpStatus,
			expectLevel:  ErrUnauthorized.LoggerLevel,
			expectMsg:    ErrUnauthorized.DefaultMessage,
		},
		{
			name:         "create error from ErrDuplicateEntry preset",
			preset:       ErrDuplicateEntry,
			err:          testErr,
			expectCode:   ErrDuplicateEntry.Code,
			expectStatus: ErrDuplicateEntry.HttpStatus,
			expectLevel:  ErrDuplicateEntry.LoggerLevel,
			expectMsg:    ErrDuplicateEntry.DefaultMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AsPresetError(tt.preset, tt.err)

			if got == nil {
				t.Fatal("AsPresetError() returned nil")
			}

			if got.Code == nil || *got.Code != tt.expectCode {
				t.Errorf("AsPresetError() Code = %v, want %v", got.Code, tt.expectCode)
			}

			if got.HttpStatus == nil || *got.HttpStatus != tt.expectStatus {
				t.Errorf("AsPresetError() HttpStatus = %v, want %v", got.HttpStatus, tt.expectStatus)
			}

			if got.LoggerLevel != tt.expectLevel {
				t.Errorf("AsPresetError() LoggerLevel = %v, want %v", got.LoggerLevel, tt.expectLevel)
			}

			if got.Message != tt.expectMsg {
				t.Errorf("AsPresetError() Message = %v, want %v", got.Message, tt.expectMsg)
			}

			if got.Inner != tt.err {
				t.Errorf("AsPresetError() Inner = %v, want %v", got.Inner, tt.err)
			}
		})
	}
}
