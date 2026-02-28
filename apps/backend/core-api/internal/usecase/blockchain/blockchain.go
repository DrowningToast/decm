package blockchain

import (
	"apps/backend/core-api/config"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	"log/slog"
)

type BlockchainUsecase struct {
	logger             *slog.Logger
	config             *config.Config
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway
}

func NewBlockchainUsecase(logger *slog.Logger, config *config.Config, blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway) *BlockchainUsecase {
	return &BlockchainUsecase{
		logger:             logger,
		config:             config,
		blockchainClientDg: blockchainClientDg,
	}
}
