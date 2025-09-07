package config

import (
	"sync"
	"time"

	"apps/backend/core-api/config/postgres"

	"github.com/caarlos0/env/v10"
	"github.com/cockroachdb/errors"
	"github.com/joho/godotenv"
)

const (
	DefaultEnvPath = "../../.env"
)

var (
	configOnce sync.Once
	config     = &Config{}
)

type Config struct {
	Name string `env:"NAME" envDefault:"decm-core-api"`
	ENV  string `env:"ENV,required" envDefault:"production"`
	Port int    `env:"PORT,required" envDefault:"8080"`

	Api      ApiConfig       `envPrefix:"API_"`
	Postgres postgres.Config `envPrefix:"DB_"`

	// PII Encryption Configuration
	PIIEncryptionKey string `env:"PII_ENCRYPTION_KEY,required"`
	// Jwt Configuration
	Jwt JwtConfig `envPrefix:"JWT_"`
}

type ApiConfig struct {
	Timeout           time.Duration `env:"TIMEOUT" envDefault:"60s"`
	MaxReadBufferSize int           `env:"MAX_READ_BUFFER_SIZE" envDefault:"4096"` // 4KB
}

type JwtConfig struct {
	SecretKey  string        `env:"SECRET_KEY,required"`
	Expiration time.Duration `env:"EXPIRATION" envDefault:"24h"`
}

func LoadConfig() Config {
	configOnce.Do(
		func() {
			if err := godotenv.Load(DefaultEnvPath); err != nil {
				panic(errors.Wrap(err, "failed to load environment variables"))
			}

			envOptions := env.Options{
				UseFieldNameByDefault: true,
				RequiredIfNoDef:       false,
			}

			if err := env.ParseWithOptions(config, envOptions); err != nil {
				panic(errors.Wrap(err, "failed to parse environment variables"))
			}

			environment, err := IParseEnvironment(config.ENV)
			if err != nil {
				panic(errors.Wrap(err, "failed to parse environment"))
			}

			config.ENV = environment.String()
		},
	)

	return *config
}
