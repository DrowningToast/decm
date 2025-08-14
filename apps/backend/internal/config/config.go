package config

import (
	"os"
	"strings"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Blockchain BlockchainConfig
	CORS       CORSConfig
	JWT        JWTConfig
	LDAP       LDAPConfig
}

type ServerConfig struct {
	Port        string
	Environment string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type BlockchainConfig struct {
	Network    string
	RPCUrl     string
	PrivateKey string
	ChainID    string
}

type CORSConfig struct {
	AllowedOrigins string
}

type JWTConfig struct {
	Secret     string
	Expiration string
}

type LDAPConfig struct {
	Host     string
	Port     string
	BaseDN   string
	BindUser string
	BindPass string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "decm"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Blockchain: BlockchainConfig{
			Network:    getEnv("BLOCKCHAIN_NETWORK", "localhost"),
			RPCUrl:     getEnv("BLOCKCHAIN_RPC_URL", "http://localhost:8545"),
			PrivateKey: getEnv("BLOCKCHAIN_PRIVATE_KEY", ""),
			ChainID:    getEnv("BLOCKCHAIN_CHAIN_ID", "1337"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			Expiration: getEnv("JWT_EXPIRATION", "24h"),
		},
		LDAP: LDAPConfig{
			Host:     getEnv("LDAP_HOST", "ldap.example.com"),
			Port:     getEnv("LDAP_PORT", "389"),
			BaseDN:   getEnv("LDAP_BASE_DN", "dc=example,dc=com"),
			BindUser: getEnv("LDAP_BIND_USER", ""),
			BindPass: getEnv("LDAP_BIND_PASS", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *CORSConfig) GetAllowedOrigins() []string {
	return strings.Split(c.AllowedOrigins, ",")
}
