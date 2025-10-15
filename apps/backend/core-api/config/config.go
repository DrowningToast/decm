package config

import (
	"sync"
	"time"

	"apps/backend/core-api/config/postgres"
	"apps/backend/core-api/config/s3"

	google "apps/backend/core-api/config/google"

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
	Name               string `env:"NAME" envDefault:"decm-core-api"`
	Env                string `env:"ENVIRONMENT,required"`
	Port               int    `env:"PORT,required" envDefault:"8080"`
	CorsAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" envDefault:"http://localhost:3000, http://127.0.0.1:3000"`

	Api      ApiConfig       `envPrefix:"API_"`
	Postgres postgres.Config `envPrefix:"DB_"`

	// PII Encryption Configuration
	PIIEncryptionKey string `env:"PII_ENCRYPTION_KEY,required"`
	// Jwt Configuration
	Jwt JwtConfig `envPrefix:"JWT_"`
	// Google OAuth Configuration
	GoogleOAuth google.GoogleOAuthConfig `envPrefix:"GOOGLE_OAUTH_"`
	// S3 Configuration
	S3 s3.S3Config `envPrefix:"S3_"`
}

type ApiConfig struct {
	Timeout           time.Duration `env:"TIMEOUT" envDefault:"60s"`
	MaxReadBufferSize int           `env:"MAX_READ_BUFFER_SIZE" envDefault:"4096"` // 4KB
}

type JwtConfig struct {
	Issuer     string        `env:"ISSUER" envDefault:"decm-service"`
	SecretKey  string        `env:"SECRET,required"`
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

			environment, err := IParseEnvironment(config.Env)
			if err != nil {
				panic(errors.Wrap(err, "failed to parse environment"))
			}

			config.Env = environment.String()
		},
	)

	return *config
}
