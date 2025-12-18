package system_status

import (
	"context"

	"apps/backend/core-api/internal/entity"
)

func (uc *SystemStatusUsecase) GetPlannedMaintenanceSchedules(ctx context.Context) ([]*entity.SystemStatusSchedule, error) {
	return uc.systemStatusRepo.GetPlannedMaintenanceSchedules(ctx)
}
