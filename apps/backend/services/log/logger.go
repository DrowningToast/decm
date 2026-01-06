package log

import (
	"apps/backend/core-api/config"
	"apps/backend/core-api/constants/ctxkey"
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
)

var (
	Logger     *slog.Logger
	configOnce sync.Once
)

func NewLogger() *slog.Logger {
	cfg := config.LoadConfig()
	return LoadLoggerInEnvironment(cfg.Env)
}

func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(ctxkey.Logger{}).(*slog.Logger)
	if !ok {
		return NewLogger()
	}
	return logger
}

func LoadLoggerInEnvironment(environment string) *slog.Logger {
	configOnce.Do(
		func() {
			var handler slog.Handler
			var level slog.Level

			headLogFile, err := os.OpenFile("logs/head.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}
			multiWriter := io.MultiWriter(os.Stdout, headLogFile)
			// Wrap with HookWriter to enable write hooks
			SetWriteHook(BeforeLogFileWrite)
			logWriter := &HookWriter{writer: multiWriter}
			defer headLogFile.Close()

			if config.IsDevelopment() {
				level = slog.LevelDebug
				// Use text handler for development - easier to read
				handler = slog.NewTextHandler(logWriter, &slog.HandlerOptions{
					Level:     level,
					AddSource: true,
				})
			} else if config.IsTesting() {
				level = slog.LevelDebug
				// Use text handler for testing - easier to debug
				handler = slog.NewTextHandler(logWriter, &slog.HandlerOptions{
					Level:     level,
					AddSource: true,
				})
			} else if config.IsProduction() {
				level = slog.LevelInfo
				// Use JSON handler for production - structured logging
				handler = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
					Level: level,
				})
				// TODO: Add file output and log rotation
			} else {
				panic("invalid environment when loading environment level")
			}

			Logger = slog.New(handler)
			slog.SetDefault(Logger) // Set as default logger
		},
	)

	return Logger
}
