package blockchain

import (
	"apps/backend/core-api/internal/usecase/blockchain"
	"apps/backend/core-api/internal/usecase/system_status"
)

type Handler struct {
	blockchainUc   *blockchain.BlockchainUsecase
	systemStatusUc *system_status.SystemStatusUsecase
}

func NewHandler(blockchainUc *blockchain.BlockchainUsecase, systemStatusUc *system_status.SystemStatusUsecase) *Handler {
	return &Handler{
		blockchainUc:   blockchainUc,
		systemStatusUc: systemStatusUc,
	}
}

