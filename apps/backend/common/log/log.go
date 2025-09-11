package log

import (
	"log/slog"
	"os"
	"sync"

	"apps/backend/core-api/config"
)

var (
	Logger     *slog.Logger
	configOnce sync.Once
)

func LoadLogger() *slog.Logger {
	cfg := config.LoadConfig()
	return LoadLoggerInEnvironment(cfg.ENV)
}

func LoadLoggerInEnvironment(environment string) *slog.Logger {
	configOnce.Do(
		func() {
			var handler slog.Handler
			var level slog.Level

			if config.IsDevelopment() {
				level = slog.LevelDebug
				// Use text handler for development - easier to read
				handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
					Level:     level,
					AddSource: true,
				})
			} else if config.IsTesting() {
				level = slog.LevelDebug
				// Use text handler for testing - easier to debug
				handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
					Level:     level,
					AddSource: true,
				})
			} else if config.IsProduction() {
				level = slog.LevelInfo
				// Use JSON handler for production - structured logging
				handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
					Level: level,
				})
				// TODO: Add file output and log rotation
			} else {
				panic("invalid environment when loading environment level")
			}

			// Wrap the base handler with the RequestID handler
			requestIDHandler := NewRequestIDHandler(handler)

			Logger = slog.New(requestIDHandler)
			slog.SetDefault(Logger) // Set as default logger
		},
	)

	return Logger
}
