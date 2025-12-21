package blockchain

import (
	"log/slog"

	"apps/backend/core-api/config"
)

type BlockchainUsecase struct {
	logger *slog.Logger
	config *config.Config
}

func NewBlockchainUsecase(logger *slog.Logger, config *config.Config) *BlockchainUsecase {
	return &BlockchainUsecase{
		logger: logger,
		config: config,
	}
}

