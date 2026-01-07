package system_status

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"time"
)

func (uc *SystemStatusUsecase) GetSchedulesBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*entity.SystemStatusSchedule, error) {
	return uc.systemStatusRepo.GetSystemStatusSchedulesUpdatedBetween(ctx, startDate, endDate)
}
