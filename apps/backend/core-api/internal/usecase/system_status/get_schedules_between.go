package system_status

import (
	"context"
	"time"

	"apps/backend/core-api/internal/entity"
)

func (uc *SystemStatusUsecase) GetSchedulesBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*entity.SystemStatusSchedule, error) {
	return uc.systemStatusRepo.GetSystemStatusSchedulesUpdatedBetween(ctx, startDate, endDate)
}
