package handlers

import (
	"decm-backend/internal/config"
	"decm-backend/internal/services"
)

// Handlers struct holds all handler dependencies
type Handlers struct {
	config *config.Config
	DB     *services.DatabaseService
}

// NewHandlers creates a new handlers instance
func NewHandlers(cfg *config.Config, dbService *services.DatabaseService) *Handlers {
	return &Handlers{
		config: cfg,
		DB:     dbService,
	}
}
