package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	generated "decm-database/go/generated"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConfigFromEnv creates a new database config from environment variables
func NewConfigFromEnv() *Config {
	return &Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", ""),
		DBName:   getEnvOrDefault("DB_NAME", "decm"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
	}
}

// ConnectionString returns a PostgreSQL connection string
func (c *Config) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// DB wraps the database connection and queries
type DB struct {
	Pool    *pgxpool.Pool
	Queries *generated.Queries
}

// NewDB creates a new database connection
func NewDB(ctx context.Context, config *Config) (*DB, error) {
	pool, err := pgxpool.New(ctx, config.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := generated.New(pool)

	return &DB{
		Pool:    pool,
		Queries: queries,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Begin starts a new transaction
func (db *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.Pool.Begin(ctx)
}

// WithTx returns a new Queries instance with the given transaction
func (db *DB) WithTx(tx pgx.Tx) *generated.Queries {
	return db.Queries.WithTx(tx)
}

// RunMigrations runs database migrations
func RunMigrations(config *Config, migrationsPath string) error {
	// Create a sql.DB instance for the migrate package
	db, err := sql.Open("postgres", config.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to open database for migrations: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
