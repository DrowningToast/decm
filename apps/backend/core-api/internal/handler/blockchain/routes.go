package blockchain

import (
	"apps/backend/common/log"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Mount(r fiber.Router) {
	logger := log.LoadLogger()
	defer logger.Info("Mounted blockchain routes")

	blockchainGroup := r.Group("/blockchain")

	blockchainGroup.Get("/gas-price", h.GetGasPrice)
}

