package customerror

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestGetErrFiberHandler_CustomError(t *testing.T) {
	tests := []struct {
		name           string
		err            *Err
		expectStatus   int
		expectLogLevel slog.Level
	}{
		{
			name:           "custom error with error level",
			err:            Parse(&ErrInternalServer, errors.New("test error")),
			expectStatus:   http.StatusInternalServerError,
			expectLogLevel: slog.LevelError,
		},
		{
			name:           "custom error with warn level",
			err:            Parse(&ErrNotFound, errors.New("test error")),
			expectStatus:   http.StatusNotFound,
			expectLogLevel: slog.LevelWarn,
		},
		{
			name:           "custom error invalid argument",
			err:            Parse(&ErrInvalidArgument, errors.New("test error")),
			expectStatus:   http.StatusBadRequest,
			expectLogLevel: slog.LevelWarn,
		},
		{
			name:           "custom error unauthenticated",
			err:            Parse(&ErrUnauthenticated, errors.New("test error")),
			expectStatus:   http.StatusUnauthorized,
			expectLogLevel: slog.LevelWarn,
		},
		{
			name:           "custom error unauthorized",
			err:            Parse(&ErrUnauthorized, errors.New("test error")),
			expectStatus:   http.StatusForbidden,
			expectLogLevel: slog.LevelWarn,
		},
		{
			name:           "custom error duplicate entry",
			err:            Parse(&ErrDuplicateEntry, errors.New("test error")),
			expectStatus:   http.StatusConflict,
			expectLogLevel: slog.LevelWarn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture logs
			var logBuf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))

			// Create Fiber app with custom error handler
			app := fiber.New(fiber.Config{
				ErrorHandler: GetErrFiberHandler(logger),
			})

			// Create a test route that returns our error
			app.Get("/test", func(c *fiber.Ctx) error {
				return tt.err
			})

			// Make test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make test request: %v", err)
			}

			// Check status code
			if resp.StatusCode != tt.expectStatus {
				t.Errorf("Response status = %d, want %d", resp.StatusCode, tt.expectStatus)
			}

			// Check response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			_ = resp.Body.Close()

			var apiResponse map[string]interface{}
			if err := json.Unmarshal(body, &apiResponse); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if _, exists := apiResponse["message"]; !exists {
				t.Error("Response body missing 'message' field")
			}

			// Verify logs were written
			logOutput := logBuf.String()
			if logOutput == "" {
				t.Error("Expected log output but got none")
			}
		})
	}
}

func TestGetErrFiberHandler_FiberError(t *testing.T) {
	tests := []struct {
		name         string
		fiberErr     *fiber.Error
		expectStatus int
	}{
		{
			name:         "fiber bad request",
			fiberErr:     fiber.NewError(fiber.StatusBadRequest, "Bad request"),
			expectStatus: fiber.StatusBadRequest,
		},
		{
			name:         "fiber unauthorized",
			fiberErr:     fiber.NewError(fiber.StatusUnauthorized, "Unauthorized"),
			expectStatus: fiber.StatusUnauthorized,
		},
		{
			name:         "fiber not found",
			fiberErr:     fiber.NewError(fiber.StatusNotFound, "Not found"),
			expectStatus: fiber.StatusNotFound,
		},
		{
			name:         "fiber internal server error",
			fiberErr:     fiber.NewError(fiber.StatusInternalServerError, "Internal server error"),
			expectStatus: fiber.StatusInternalServerError,
		},
		{
			name:         "fiber service unavailable",
			fiberErr:     fiber.NewError(fiber.StatusServiceUnavailable, "Service unavailable"),
			expectStatus: fiber.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture logs
			var logBuf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))

			// Create Fiber app with custom error handler
			app := fiber.New(fiber.Config{
				ErrorHandler: GetErrFiberHandler(logger),
			})

			// Create a test route that returns our error
			app.Get("/test", func(c *fiber.Ctx) error {
				return tt.fiberErr
			})

			// Make test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make test request: %v", err)
			}

			// Check status code
			if resp.StatusCode != tt.expectStatus {
				t.Errorf("Response status = %d, want %d", resp.StatusCode, tt.expectStatus)
			}

			// Check response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			_ = resp.Body.Close()

			var apiResponse map[string]interface{}
			if err := json.Unmarshal(body, &apiResponse); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if message, exists := apiResponse["message"]; !exists {
				t.Error("Response body missing 'message' field")
			} else if message != tt.fiberErr.Message {
				t.Errorf("Response message = %v, want %v", message, tt.fiberErr.Message)
			}

			// Verify logs were written
			logOutput := logBuf.String()
			if logOutput == "" {
				t.Error("Expected log output but got none")
			}
		})
	}
}

func TestGetErrFiberHandler_UnknownError(t *testing.T) {
	// Create a buffer to capture logs
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create Fiber app with custom error handler
	app := fiber.New(fiber.Config{
		ErrorHandler: GetErrFiberHandler(logger),
	})

	// Create a test route that returns an unknown error
	app.Get("/test", func(c *fiber.Ctx) error {
		return errors.New("unknown error")
	})

	// Make test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make test request: %v", err)
	}

	// Check status code - should be internal server error
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("Response status = %d, want %d", resp.StatusCode, fiber.StatusInternalServerError)
	}

	// Check response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	_ = resp.Body.Close()

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if message, exists := apiResponse["message"]; !exists {
		t.Error("Response body missing 'message' field")
	} else {
		expectedMsg := "An unknown error has occurred. Please try again later."
		if message != expectedMsg {
			t.Errorf("Response message = %v, want %v", message, expectedMsg)
		}
	}

	// Verify logs were written with error level
	logOutput := logBuf.String()
	if logOutput == "" {
		t.Error("Expected log output but got none")
	}
}

func TestGetErrFiberHandler_LogLevels(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  slog.Level
		errPreset ErrSignature
	}{
		{
			name:      "error level logging",
			logLevel:  slog.LevelError,
			errPreset: ErrInternalServer,
		},
		{
			name:      "warn level logging",
			logLevel:  slog.LevelWarn,
			errPreset: ErrNotFound,
		},
		{
			name:     "info level logging",
			logLevel: slog.LevelInfo,
			errPreset: ErrSignature{
				Code:           "INFO_TEST",
				DefaultMessage: "Info test message",
				HttpStatus:     http.StatusOK,
				LoggerLevel:    slog.LevelInfo,
			},
		},
		{
			name:     "debug level logging",
			logLevel: slog.LevelDebug,
			errPreset: ErrSignature{
				Code:           "DEBUG_TEST",
				DefaultMessage: "Debug test message",
				HttpStatus:     http.StatusOK,
				LoggerLevel:    slog.LevelDebug,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture logs
			var logBuf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{
				Level: slog.LevelDebug, // Capture all log levels
			}))

			// Create Fiber app with custom error handler
			app := fiber.New(fiber.Config{
				ErrorHandler: GetErrFiberHandler(logger),
			})

			// Create test error with specific log level
			testErr := Parse(&tt.errPreset, errors.New("test error"))

			// Create a test route that returns our error
			app.Get("/test", func(c *fiber.Ctx) error {
				return testErr
			})

			// Make test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			_, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make test request: %v", err)
			}

			// Verify logs were written
			logOutput := logBuf.String()
			if logOutput == "" {
				t.Error("Expected log output but got none")
			}

			// Parse log to verify level
			var logEntry map[string]interface{}
			if err := json.Unmarshal([]byte(logOutput), &logEntry); err == nil {
				if level, exists := logEntry["level"]; exists {
					levelStr := level.(string)
					expectedLevel := tt.logLevel.String()
					if levelStr != expectedLevel {
						t.Errorf("Log level = %v, want %v", levelStr, expectedLevel)
					}
				}
			}
		})
	}
}
