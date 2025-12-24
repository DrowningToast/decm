package blockchain

import (
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"context"

	"github.com/pkg/errors"
)

type GasPriceResponse struct {
	CurrentGasPriceGwei float64 `json:"current_gas_price_gwei"`
	BaseFeeGwei         float64 `json:"base_fee_gwei"`
	PriorityFeeGwei     float64 `json:"priority_fee_gwei"`
	SoftCapPriceGwei    float64 `json:"soft_cap_price_gwei"`
	HardCapPriceGwei    float64 `json:"hard_cap_price_gwei"`
	SoftSafetyMargin    float64 `json:"soft_safety_margin"` // Percentage remaining until soft cap
	HardSafetyMargin    float64 `json:"hard_safety_margin"` // Percentage remaining until hard cap
}

func (uc *BlockchainUsecase) GetGasPrice(ctx context.Context) (*GasPriceResponse, error) {
	client, err := cyptoutils.GetEthereumClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ethereum client")
	}
	defer client.Close()

	info, err := cyptoutils.GetCurrentGasPrice(ctx, client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current gas price")
	}

	softCap := uc.config.Blockchain.SoftCapGasPriceGwei
	hardCap := uc.config.Blockchain.MaxGasPriceGwei

	softMargin := 0.0
	if softCap > 0 {
		softMargin = ((softCap - info.MaxFeePerGasGwei) / softCap) * 100
	}

	hardMargin := 0.0
	if hardCap > 0 {
		hardMargin = ((hardCap - info.MaxFeePerGasGwei) / hardCap) * 100
	}

	return &GasPriceResponse{
		CurrentGasPriceGwei: info.MaxFeePerGasGwei,
		BaseFeeGwei:         info.BaseFeeGwei,
		PriorityFeeGwei:     info.MaxPriorityFeePerGasGwei,
		SoftCapPriceGwei:    softCap,
		HardCapPriceGwei:    hardCap,
		SoftSafetyMargin:    softMargin,
		HardSafetyMargin:    hardMargin,
	}, nil
}
