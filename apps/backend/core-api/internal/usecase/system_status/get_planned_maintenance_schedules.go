package system_status

import (
	"apps/backend/core-api/internal/entity"
	"context"
)

func (uc *SystemStatusUsecase) GetPlannedMaintenanceSchedules(ctx context.Context) ([]*entity.SystemStatusSchedule, error) {
	return uc.systemStatusRepo.GetPlannedMaintenanceSchedules(ctx)
}
