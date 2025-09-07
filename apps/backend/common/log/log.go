package log

import (
	"log/slog"
	"os"
	"sync"

	"apps/backend/core-api/config"

	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
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
			logrusLogger := logrus.New()
			var level slog.Leveler

			if config.IsDevelopment() {
				level = slog.LevelDebug
				logrusLogger.SetLevel(logrus.DebugLevel)
				logrusLogger.SetFormatter(&logrus.TextFormatter{
					FullTimestamp: true,
				})
				logrusLogger.SetOutput(os.Stdout)

			} else if config.IsTesting() {
				level = slog.LevelDebug
				logrusLogger.SetLevel(logrus.DebugLevel)
				logrusLogger.SetFormatter(&logrus.TextFormatter{
					FullTimestamp: true,
				})
				logrusLogger.SetOutput(os.Stdout)

			} else if config.IsProduction() {
				level = slog.LevelInfo
				logrusLogger.SetLevel(logrus.InfoLevel)
				logrusLogger.SetFormatter(&logrus.TextFormatter{
					FullTimestamp: true,
				})
				// TODO: Set output to file and stdout
				logrusLogger.SetOutput(os.Stdout)

			} else {
				panic("invalid environment when loading environment level")
			}

			Logger = slog.New(sloglogrus.Option{
				Logger: logrusLogger,
				Level:  level,
			}.NewLogrusHandler())
		},
	)

	return Logger
}
