package system_status

import (
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/blockchain"

	"github.com/gofiber/fiber/v2"
)

type SystemStatusViewModel struct {
	LatestSchedule *entity.SystemStatusSchedule     `json:"latest_schedule,omitempty" extensions:"x-nullable"`
	GasPrice       *blockchain.GasPriceResponse     `json:"gas_price"`
}

// GetSystemStatusViewModel godoc
// @Summary Get system status viewmodel
// @Description Get the latest system status schedule combined with current gas price information
// @ID get-system-status-viewmodel
// @Tags SystemStatus
// @Produce json
// @Success 200 {object} system_status.SystemStatusViewModel
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/system-status/viewmodel [get]
func (h *Handler) GetSystemStatusViewModel(ctx *fiber.Ctx) error {
	// Set cache control header for 5 minutes
	ctx.Set("Cache-Control", "public, max-age=300")

	// Fetch latest system status schedule
	schedules, err := h.systemStatusUc.GetLatestSchedules(ctx.UserContext(), 1)
	if err != nil {
		return err
	}

	var latestSchedule *entity.SystemStatusSchedule
	if len(schedules) > 0 {
		latestSchedule = schedules[0]
	}

	// Fetch gas price
	gasPrice, err := h.blockchainUc.GetGasPrice(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(&SystemStatusViewModel{
		LatestSchedule: latestSchedule,
		GasPrice:       gasPrice,
	})
}
