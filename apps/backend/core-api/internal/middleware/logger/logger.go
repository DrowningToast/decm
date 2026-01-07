package logger

import (
	"apps/backend/core-api/constants/ctxkey"
	"apps/backend/services/log"
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
)

var ExcludedPaths = []string{"/swagger/*", "/metrics"}

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

		// Log with structured fields
		if err != nil {
			logger.ErrorContext(c.UserContext(), fmt.Sprintf("HTTP request error: %s %s", method, path),
				slog.String("method", method),
				slog.String("path", path),
				slog.Int("status", status),
				slog.Duration("duration", duration),
				slog.String("error", err.Error()),
			)
			return err
		}

		logger.InfoContext(c.UserContext(), fmt.Sprintf("HTTP request completed: %s %s", method, path),
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("duration", duration),
		)
		return nil
	}
}
