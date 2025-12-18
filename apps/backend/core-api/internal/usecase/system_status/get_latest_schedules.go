package system_status

import (
	"context"

	"apps/backend/core-api/internal/entity"
)

func (uc *SystemStatusUsecase) GetLatestSchedules(ctx context.Context, pageSize int32) ([]*entity.SystemStatusSchedule, error) {
	// pageSize will be used as the limit for history querying which is ordered by start_time DESC
	return uc.systemStatusRepo.GetSystemStatusScheduleHistory(ctx, pageSize, 0)
}
