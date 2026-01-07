package blockchain

import (
	"apps/backend/core-api/config"
	"log/slog"
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
