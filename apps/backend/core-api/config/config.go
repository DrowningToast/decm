package config

import (
	"sync"
	"time"

	"apps/backend/core-api/config/blockchain"
	"apps/backend/core-api/config/postgres"
	"apps/backend/core-api/config/s3"

	google "apps/backend/core-api/config/google"

	"github.com/caarlos0/env/v10"
	"github.com/cockroachdb/errors"
	"github.com/joho/godotenv"
)

const (
	DefaultEnvPath = "../../.env" // For local development
	DockerEnvPath  = ".env"       // For Docker container
	AppEnvPath     = "/app/.env"  // For Docker absolute path
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
	CookieDomain       string `env:"COOKIE_DOMAIN"` // Empty for dev (host-only), domain for production

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
	// Blockchain Configuration
	Blockchain blockchain.BlockchainConfig `envPrefix:"BLOCKCHAIN_"`
	// Metrics Configuration
	Metrics MetricsConfig `envPrefix:"METRICS_"`
}

// Validate validates the configuration and returns an error if any required configuration is invalid
func (c *Config) Validate() error {
	if !c.S3.IsValid() {
		return errors.New("S3 configuration is invalid: AccessKeyID, SecretAccessKey, BucketName, and Endpoint are all required")
	}

	if err := c.Blockchain.Validate(); err != nil {
		return errors.Wrap(err, "blockchain configuration is invalid")
	}

	// Add additional component validators here as needed

	return nil
}

type ApiConfig struct {
	Timeout           time.Duration `env:"TIMEOUT" envDefault:"60s"`
	MaxReadBufferSize int           `env:"MAX_READ_BUFFER_SIZE" envDefault:"4096"` // 4KB
}

type JwtConfig struct {
	Issuer     string `env:"ISSUER" envDefault:"decm-service"`
	SecretKey  string `env:"SECRET,required"`
	Expiration string `env:"EXPIRATION" envDefault:"24h"`
}

type MetricsConfig struct {
	ApiKey string `env:"API_KEY" envDefault:"metrics-api-key"`
}

func LoadConfig() Config {
	configOnce.Do(
		func() {
			// Try to load .env file for local development
			// In production (Docker), environment variables are injected by docker-compose,
			// so it's okay if no .env file exists
			envPaths := []string{DefaultEnvPath, DockerEnvPath, AppEnvPath}

			fileLoaded := false
			for _, path := range envPaths {
				// Attempt to load, but don't panic if file doesn't exist
				// Environment variables might already be set by the container runtime
				if err := godotenv.Load(path); err == nil {
					fileLoaded = true
					break
				}
			}

			// Parse environment variables (whether from .env file or container runtime)
			// This will fail with a clear error if required environment variables are missing
			envOptions := env.Options{
				UseFieldNameByDefault: true,
				RequiredIfNoDef:       false,
			}

			if err := env.ParseWithOptions(config, envOptions); err != nil {
				// Provide context about where config should come from
				if !fileLoaded {
					panic(errors.Wrap(err, "failed to parse environment variables: no .env file found and required environment variables are not set. In Docker, ensure docker-compose is configured to inject environment variables"))
				}
				panic(errors.Wrap(err, "failed to parse environment variables from .env file"))
			}

			// No need for manual validation of required fields; env.ParseWithOptions enforces this.

			environment, err := IParseEnvironment(config.Env)
			if err != nil {
				panic(errors.Wrap(err, "failed to parse environment"))
			}

			config.Env = environment.String()
		},
	)

	return *config
}
