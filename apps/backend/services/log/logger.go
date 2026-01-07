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
	// Logger is the singleton logger instance.
	// Initialize once at application startup by calling Init() or NewLogger().
	Logger       *slog.Logger
	configOnce   sync.Once
	rotatingFile *RotatingFileWriter // Rotating file writer for head.log
)

// Init initializes the logger singleton.
// This should be called once at application startup.
// Subsequent calls are safe but have no effect (idempotent via sync.Once).
func Init() {
	cfg := config.LoadConfig()
	LoadLoggerInEnvironment(cfg.Env)
}

// NewLogger returns the singleton logger instance.
// If not already initialized, it will initialize on first call.
// For best performance, call Init() once at startup, then use Logger directly.
//
// Deprecated: Prefer calling Init() at startup and using log.Logger directly.
func NewLogger() *slog.Logger {
	if Logger == nil {
		Init()
	}
	return Logger
}

// FromContext returns a logger from the context, or the singleton if not found.
// This is useful in HTTP handlers where logger may be stored in request context.
func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(ctxkey.Logger{}).(*slog.Logger)
	if !ok {
		// Fall back to singleton
		if Logger == nil {
			Init()
		}
		return Logger
	}
	return logger
}

func LoadLoggerInEnvironment(environment string) *slog.Logger {
	configOnce.Do(
		func() {
			var handler slog.Handler
			var level slog.Level

			// Create rotating file writer for log rotation support
			var err error
			rotatingFile, err = NewRotatingFileWriter("logs/head.log")
			if err != nil {
				panic(err)
			}
			multiWriter := io.MultiWriter(os.Stdout, rotatingFile)
			// Wrap with HookWriter to enable write hooks
			SetWriteHook(BeforeLogFileWrite)
			logWriter := &HookWriter{writer: multiWriter}
			// Note: rotatingFile is kept open for application lifetime, OS closes on process exit

			if config.IsDevelopment() {
				level = slog.LevelDebug
				// Use JSON handler for consistent structured logging
				handler = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
					Level:     level,
					AddSource: true, // Include file path and line number for debugging
				})
			} else if config.IsTesting() {
				level = slog.LevelDebug
				// Use JSON handler for consistent structured logging
				handler = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
					Level:     level,
					AddSource: true, // Include file path and line number for debugging
				})
			} else if config.IsProduction() {
				level = slog.LevelInfo
				// Use JSON handler for production - structured logging
				handler = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
					Level:     level,
					AddSource: true, // Include source location for production debugging
				})
			} else {
				panic("invalid environment when loading environment level")
			}

			Logger = slog.New(handler)
			slog.SetDefault(Logger) // Set as default logger
		},
	)

	return Logger
}
