package system_status

import (
	"log/slog"

	"apps/backend/core-api/internal/usecase/system_status"
)

type Handler struct {
	systemStatusUc *system_status.SystemStatusUsecase
	logger         *slog.Logger
}

func NewHandler(systemStatusUc *system_status.SystemStatusUsecase, logger *slog.Logger) *Handler {
	return &Handler{
		systemStatusUc: systemStatusUc,
		logger:         logger,
	}
}
