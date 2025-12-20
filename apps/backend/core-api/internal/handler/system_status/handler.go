package system_status

import (
	"log/slog"

	"apps/backend/core-api/internal/usecase/blockchain"
	"apps/backend/core-api/internal/usecase/system_status"
)

type Handler struct {
	systemStatusUc *system_status.SystemStatusUsecase
	blockchainUc   *blockchain.BlockchainUsecase
	logger         *slog.Logger
}

func NewHandler(systemStatusUc *system_status.SystemStatusUsecase, blockchainUc *blockchain.BlockchainUsecase, logger *slog.Logger) *Handler {
	return &Handler{
		systemStatusUc: systemStatusUc,
		blockchainUc:   blockchainUc,
		logger:         logger,
	}
}
