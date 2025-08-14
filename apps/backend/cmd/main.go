package main

import (
	"context"
	"log"
	"os"
	"time"

	"decm-backend/api/routes"
	_ "decm-backend/docs"
	"decm-backend/internal/config"
	"decm-backend/internal/handlers"
	"decm-backend/internal/middleware"
	"decm-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title DECM Backend API
// @version 1.0
// @description Backend for Frontend API for DECM (Decentralized Event Management) platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables from root directory
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found in root directory, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database service
	log.Println("Initializing database connection...")
	dbService, err := services.NewDatabaseService(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database service: ", err)
	}
	defer dbService.Close()
	log.Println("✅ Database connection established")

	// Initialize Fiber app with custom config
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "DECM-BFF-API",
		AppName:       "DECM Backend for Frontend API v1.0.0",
		ErrorHandler:  middleware.ErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	}))

	// CORS middleware - configure for frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Custom middleware
	app.Use(middleware.RequestID())

	// Initialize handlers
	h := handlers.NewHandlers(cfg, dbService)

	// Swagger documentation route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Setup routes
	routes.SetupRoutes(app, h)

	// Health check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		// Check database health
		dbHealthy := true
		if err := dbService.Health(ctx); err != nil {
			log.Printf("Database health check failed: %v", err)
			dbHealthy = false
		}

		status := "healthy"
		if !dbHealthy {
			status = "unhealthy"
		}

		return c.JSON(fiber.Map{
			"message":   "DECM Backend for Frontend API",
			"version":   "1.0.0",
			"status":    status,
			"docs":      "/swagger/index.html",
			"database":  dbHealthy,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("🚀 DECM BFF API server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
