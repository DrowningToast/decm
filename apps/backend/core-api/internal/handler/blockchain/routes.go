package blockchain

import (
	"apps/backend/services/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	// Logger singleton initialized in main.go
	defer log.Logger.Info("Mounted blockchain routes")

	blockchainGroup := r.Group("/blockchain")

	blockchainGroup.Get("/gas-price", h.GetGasPrice)
	blockchainGroup.Get("/status", h.GetSystemStatus)
}
