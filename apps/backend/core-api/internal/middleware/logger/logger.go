package logger

import (
	"apps/backend/core-api/constants/ctxkey"
	"apps/backend/core-api/constants/security"
	"apps/backend/services/log"
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var ExcludedPaths = []string{"/swagger/*", "/metrics"}

// isSuspiciousPath checks if the request path matches known attack patterns.
// Patterns are defined in constants/security/patterns.go for easy maintenance.
// Supports prefix, suffix, exact, and contains matching for precise detection.
func isSuspiciousPath(path string) bool {
	pathLower := strings.ToLower(path)

	for _, sp := range security.SuspiciousPathPatterns {
		patternLower := strings.ToLower(sp.Pattern)

		switch sp.MatchType {
		case security.MatchPrefix:
			if strings.HasPrefix(pathLower, patternLower) {
				return true
			}
		case security.MatchSuffix:
			if strings.HasSuffix(pathLower, patternLower) {
				return true
			}
		case security.MatchExact:
			if pathLower == patternLower {
				return true
			}
		case security.MatchContains:
			if strings.Contains(pathLower, patternLower) {
				return true
			}
		}
	}
	return false
}

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get logger from context (creates new one if not found)
		logger := log.FromContext(c.UserContext())

		// Get the request id from the context and add it to logger
		requestId, _ := c.UserContext().Value("request_id").(string)
		if requestId != "" {
			logger = logger.With(slog.String("request_id", requestId))
		}

		ctx := context.WithValue(c.Context(), ctxkey.Logger{}, logger)
		c.SetUserContext(ctx)

		start := time.Now()
		method := c.Method()
		path := c.Path()

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		// Skip logging for excluded paths
		for _, excludedPath := range ExcludedPaths {
			matched, matchErr := regexp.MatchString(excludedPath, path)
			if matchErr != nil {
				logger.Error("Pattern match error", slog.String("error", matchErr.Error()))
				continue
			}
			if matched {
				return err
			}
		}

		// Check if this is a suspicious security probe
		suspicious := isSuspiciousPath(path)

		// Log with structured fields
		if err != nil {
			// If suspicious, log with security category and extra context
			if suspicious {
				logger.WarnContext(c.UserContext(), fmt.Sprintf("Security probe detected: %s %s", method, path),
					slog.String("category", "security"),
					slog.String("ip", c.IP()),
					slog.String("user_agent", c.Get("User-Agent")),
					slog.String("method", method),
					slog.String("path", path),
					slog.Int("status", status),
					slog.Duration("duration", duration),
					slog.String("error", err.Error()),
				)
			} else {
				logger.ErrorContext(c.UserContext(), fmt.Sprintf("HTTP request error: %s %s", method, path),
					slog.String("method", method),
					slog.String("path", path),
					slog.Int("status", status),
					slog.Duration("duration", duration),
					slog.String("error", err.Error()),
				)
			}
			return err
		}

		// For successful requests, only log if not suspicious (reduce noise)
		if !suspicious {
			logger.InfoContext(c.UserContext(), fmt.Sprintf("HTTP request completed: %s %s", method, path),
				slog.String("method", method),
				slog.String("path", path),
				slog.Int("status", status),
				slog.Duration("duration", duration),
			)
		}
		return nil
	}
}
