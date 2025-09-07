package log

import (
	"os"
	"sync"

	"apps/backend/core-api/config"

	"github.com/sirupsen/logrus"
)

var (
	Logger     *logrus.Logger
	configOnce sync.Once
)

func LoadLogger(environment string) *logrus.Logger {
	configOnce.Do(
		func() {
			if config.IsDevelopment() {
				logger := logrus.New()
				logger.SetLevel(logrus.DebugLevel)
				logger.SetFormatter(&logrus.TextFormatter{
					FullTimestamp: true,
				})
				logger.SetOutput(os.Stdout)

				Logger = logger
			} else if config.IsTesting() {
				logger := logrus.New()
				logger.SetLevel(logrus.DebugLevel)
				logger.SetFormatter(&logrus.TextFormatter{
					FullTimestamp: true,
				})
				logger.SetOutput(os.Stdout)

				Logger = logger
			} else if config.IsProduction() {
				logger := logrus.New()
				logger.SetLevel(logrus.InfoLevel)
				logger.SetFormatter(&logrus.TextFormatter{
					FullTimestamp: true,
				})
				// TODO: Set output to file and stdout
				logger.SetOutput(os.Stdout)

				Logger = logger
			} else {
				panic("invalid environment when loading environment level")
			}
		},
	)

	return Logger
}
