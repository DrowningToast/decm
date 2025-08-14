package services

import (
	"context"
	"fmt"
	"log"

	"decm-backend/internal/config"
	database "decm-database"
)

// DatabaseService wraps the database connection and provides service methods
type DatabaseService struct {
	DB *database.DB
}

// NewDatabaseService creates a new database service
func NewDatabaseService(cfg *config.Config) (*DatabaseService, error) {
	ctx := context.Background()

	// Convert config to database config
	dbConfig := &database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	// Create database connection
	db, err := database.NewDB(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations if in development environment
	if cfg.Server.Environment == "development" {
		log.Println("Running database migrations...")
		if err := database.RunMigrations(dbConfig, "../../packages/database/migrations"); err != nil {
			log.Printf("Warning: Failed to run migrations: %v", err)
		} else {
			log.Println("Database migrations completed successfully")
		}
	}

	return &DatabaseService{
		DB: db,
	}, nil
}

// Close closes the database connection
func (s *DatabaseService) Close() {
	if s.DB != nil {
		s.DB.Close()
	}
}

// Health checks database connectivity
func (s *DatabaseService) Health(ctx context.Context) error {
	return s.DB.Pool.Ping(ctx)
}
