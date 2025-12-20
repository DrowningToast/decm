package blockchain

import (
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/blockchain"

	"github.com/gofiber/fiber/v2"
)

// GetSystemStatusResponse combines gas price and system status information
type GetSystemStatusResponse struct {
	GasPrice               *blockchain.GasPriceResponse  `json:"gas_price"`
	ClosestIncomingSchedule *entity.SystemStatusSchedule `json:"closest_incoming_schedule,omitempty" extensions:"x-nullable"`
}

// GetSystemStatus godoc
// @Summary Get system status including gas price
// @Description Get current blockchain gas price and upcoming system status schedule in a single request
// @ID get-system-status
// @Tags Blockchain
// @Produce json
// @Success 200 {object} blockchain.GetSystemStatusResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/blockchain/status [get]
func (h *Handler) GetSystemStatus(ctx *fiber.Ctx) error {
	// Set cache control header for 5 minutes
	ctx.Set("Cache-Control", "public, max-age=300")

	// Fetch gas price
	gasPrice, err := h.blockchainUc.GetGasPrice(ctx.UserContext())
	if err != nil {
		return err
	}

	// Fetch closest incoming system status schedule
	schedule, err := h.systemStatusUc.GetClosestIncomingSchedule(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(&GetSystemStatusResponse{
		GasPrice:               gasPrice,
		ClosestIncomingSchedule: schedule,
	})
}
