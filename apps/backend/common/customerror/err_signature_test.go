package customerror

import (
	"log/slog"
	"net/http"
	"testing"
)

func TestErrSignature_Presets(t *testing.T) {
	tests := []struct {
		name          string
		preset        ErrSignature
		expectCode    ErrCode
		expectStatus  int
		expectLevel   slog.Level
		expectMessage string
	}{
		{
			name:          "ErrNotFound preset",
			preset:        ErrNotFound,
			expectCode:    "NOT_FOUND",
			expectStatus:  http.StatusNotFound,
			expectLevel:   slog.LevelWarn,
			expectMessage: "Resource not found.",
		},
		{
			name:          "ErrInternalServer preset",
			preset:        ErrInternalServer,
			expectCode:    "INTERNAL_SERVER_ERROR",
			expectStatus:  http.StatusInternalServerError,
			expectLevel:   slog.LevelError,
			expectMessage: "Something went wrong. Please try again later.",
		},
		{
			name:          "ErrInvalidArgument preset",
			preset:        ErrInvalidArgument,
			expectCode:    "INVALID_ARGUMENT",
			expectStatus:  http.StatusBadRequest,
			expectLevel:   slog.LevelWarn,
			expectMessage: "Invalid request. Please check your request and try again.",
		},
		{
			name:          "ErrUnauthenticated preset",
			preset:        ErrUnauthenticated,
			expectCode:    "UNAUTHORIZED",
			expectStatus:  http.StatusUnauthorized,
			expectLevel:   slog.LevelWarn,
			expectMessage: "Unauthorized. Please login to continue.",
		},
		{
			name:          "ErrUnauthorized preset",
			preset:        ErrUnauthorized,
			expectCode:    "INSUFFICIENT_PERMISSION",
			expectStatus:  http.StatusForbidden,
			expectLevel:   slog.LevelWarn,
			expectMessage: "Insufficient permission. Please contact the administrator.",
		},
		{
			name:          "ErrDuplicateEntry preset",
			preset:        ErrDuplicateEntry,
			expectCode:    "DUPLICATE_ENTRY",
			expectStatus:  http.StatusConflict,
			expectLevel:   slog.LevelWarn,
			expectMessage: "Duplicate entry. Your input already exists.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preset.Code != tt.expectCode {
				t.Errorf("Code = %v, want %v", tt.preset.Code, tt.expectCode)
			}

			if tt.preset.HttpStatus != tt.expectStatus {
				t.Errorf("HttpStatus = %v, want %v", tt.preset.HttpStatus, tt.expectStatus)
			}

			if tt.preset.LoggerLevel != tt.expectLevel {
				t.Errorf("LoggerLevel = %v, want %v", tt.preset.LoggerLevel, tt.expectLevel)
			}

			if tt.preset.DefaultMessage != tt.expectMessage {
				t.Errorf("DefaultMessage = %v, want %v", tt.preset.DefaultMessage, tt.expectMessage)
			}
		})
	}
}

func TestErrCode_Type(t *testing.T) {
	// Test that ErrCode is a string type
	var code ErrCode = "TEST_CODE"

	if string(code) != "TEST_CODE" {
		t.Errorf("ErrCode string conversion failed, got %v", code)
	}
}

func TestErrSignature_Immutability(t *testing.T) {
	// Test that preset error signatures are properly defined
	// and can be used to create errors without modification

	originalNotFound := ErrNotFound

	// Use the preset to create an error
	_ = AsPresetError(ErrNotFound, nil)

	// Verify the preset wasn't modified
	if ErrNotFound.Code != originalNotFound.Code {
		t.Error("ErrNotFound.Code was modified")
	}

	if ErrNotFound.DefaultMessage != originalNotFound.DefaultMessage {
		t.Error("ErrNotFound.DefaultMessage was modified")
	}

	if ErrNotFound.HttpStatus != originalNotFound.HttpStatus {
		t.Error("ErrNotFound.HttpStatus was modified")
	}

	if ErrNotFound.LoggerLevel != originalNotFound.LoggerLevel {
		t.Error("ErrNotFound.LoggerLevel was modified")
	}
}

func TestErrSignature_CustomCreation(t *testing.T) {
	// Test creating a custom error signature
	customErr := ErrSignature{
		Code:           "CUSTOM_ERROR",
		DefaultMessage: "This is a custom error",
		HttpStatus:     http.StatusTeapot,
		LoggerLevel:    slog.LevelInfo,
	}

	if customErr.Code != "CUSTOM_ERROR" {
		t.Errorf("Custom Code = %v, want CUSTOM_ERROR", customErr.Code)
	}

	if customErr.DefaultMessage != "This is a custom error" {
		t.Errorf("Custom DefaultMessage = %v, want 'This is a custom error'", customErr.DefaultMessage)
	}

	if customErr.HttpStatus != http.StatusTeapot {
		t.Errorf("Custom HttpStatus = %v, want %v", customErr.HttpStatus, http.StatusTeapot)
	}

	if customErr.LoggerLevel != slog.LevelInfo {
		t.Errorf("Custom LoggerLevel = %v, want %v", customErr.LoggerLevel, slog.LevelInfo)
	}
}
