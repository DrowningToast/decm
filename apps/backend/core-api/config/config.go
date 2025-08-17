package config

import (
	"sync"

	"apps/backend/core-api/config/postgres"
)

const (
	DefaultEnvPath = "../../.env"
)

var (
	configOnce sync.Once
	config     = &Config{}
)

type Config struct {
	Name string `env:"NAME" envDefault:"gaze-network-system"`
	ENV  string `env:"ENV,required" envDefault:"production"`
	Port int    `env:"PORT,required" envDefault:"8080"`

	Postgres postgres.Config `envPrefix:"DB_"`

	// PII Encryption Configuration
	PIIEncryptionKey string `env:"PII_ENCRYPTION_KEY,required"`
}
