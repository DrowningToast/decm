package metrics

import (
	"apps/backend/core-api/config"
	"strings"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	cfg        *config.Config
	prometheus *fiberprometheus.FiberPrometheus
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg:        cfg,
		prometheus: fiberprometheus.New(cfg.Name),
	}
}

// GetMiddleware returns the global prometheus middleware to be used by the app
func (h *Handler) GetMiddleware() fiber.Handler {
	return h.prometheus.Middleware
}

// Mount registers the metrics route and its authentication middleware
func (h *Handler) Mount(app *fiber.App) {
	// Protect /metrics endpoint with API Key
	app.Use("/metrics", h.authMiddleware)
	h.prometheus.RegisterAt(app, "/metrics")
}

func (h *Handler) authMiddleware(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-Key")
	if apiKey == "" {
		apiKey = c.Query("api_key")
	}

	// Also support standard Bearer token for easier Prometheus config
	if apiKey == "" {
		authHeader := c.Get("Authorization")
		const bearerPrefix = "Bearer "
		if token, ok := strings.CutPrefix(authHeader, bearerPrefix); ok {
			apiKey = token
		}
	}

	if apiKey != h.cfg.Metrics.ApiKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or missing API Key",
		})
	}
	return c.Next()
}
