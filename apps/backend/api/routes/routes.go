package routes

import (
	v1 "decm-backend/api/v1"
	"decm-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App, h *handlers.Handlers) {
	// API version 1 group
	apiV1 := app.Group("/api/v1")

	// Health check routes
	apiV1.Get("/health", healthCheck)

	// User routes
	v1.UserRoutes(apiV1, h)

	// TODO: Add more DECM-specific routes here
	// - Authentication routes
	// - Event management routes
	// - NFT ticketing routes
	// - Credential routes
	// - Portfolio routes
	// - Academic identity routes
}

// healthCheck godoc
// @Summary Health check endpoint
// @Description Check API health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "DECM BFF API",
		"version": "1.0.0",
	})
}
