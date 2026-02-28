package system_status

import (
	"apps/backend/core-api/internal/entity"
	"context"
)

func (uc *SystemStatusUsecase) GetLatestSchedules(ctx context.Context, pageSize int32) ([]*entity.SystemStatusSchedule, error) {
	// pageSize will be used as the limit for history querying which is ordered by start_time DESC
	return uc.systemStatusRepo.GetSystemStatusScheduleHistory(ctx, pageSize, 0)
}
