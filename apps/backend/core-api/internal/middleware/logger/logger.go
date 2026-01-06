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

		duration := time.Since(start).Seconds()
		status := c.Response().StatusCode()

		prefix := fmt.Sprintf("[%d] %s, %s, %f", status, method, path, duration)

		// Skip logging for excluded paths
		for _, excludedPath := range ExcludedPaths {
			matched, matchErr := regexp.MatchString(excludedPath, path)
			if matchErr != nil {
				logger.Error("An error has occured", slog.String("error", matchErr.Error()))
				continue
			}
			if matched {
				return err
			}
		}

		if err != nil {
			logger.Error(prefix, "An error has occured", slog.String("error", err.Error()))
			return err
		}

		logger.Info(prefix, slog.Bool("success", true))
		return nil
	}
}
