package logger

import (
	"apps/backend/common/log"
	"apps/backend/core-api/constants/ctxkey"
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
		// get the request id from the context
		requestId, _ := c.UserContext().Value("request_id").(string)
		logger := log.LoadLogger()

		// Create a logger with requestId as a default attribute
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
