package blockchain

import (
	"apps/backend/core-api/internal/usecase/blockchain"
)

type Handler struct {
	blockchainUc *blockchain.BlockchainUsecase
}

func NewHandler(blockchainUc *blockchain.BlockchainUsecase) *Handler {
	return &Handler{
		blockchainUc: blockchainUc,
	}
}

