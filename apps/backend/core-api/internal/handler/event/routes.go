package event

import (
	"apps/backend/common/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted event routes")

	eventGroup := r.Group("/events")
	eventGroup.Post("/", h.CreateEvent)
}
