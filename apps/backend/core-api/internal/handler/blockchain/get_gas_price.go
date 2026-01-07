package blockchain

import (
	"github.com/gofiber/fiber/v2"
)

// GetGasPrice godoc
// @Summary Get current blockchain gas price
// @Description Get current blockchain gas price and the configured maximum allowed fee
// @ID get-gas-price
// @Tags Blockchain
// @Produce json
// @Success 200 {object} blockchain.GasPriceResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/blockchain/gas-price [get]
func (h *Handler) GetGasPrice(ctx *fiber.Ctx) error {
	// Set cache control header for 5 minutes
	ctx.Set("Cache-Control", "public, max-age=300")

	gasPrice, err := h.blockchainUc.GetGasPrice(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(gasPrice)
}
