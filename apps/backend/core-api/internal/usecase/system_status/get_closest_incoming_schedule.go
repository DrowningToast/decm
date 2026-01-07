package system_status

import (
	"apps/backend/core-api/internal/entity"
	"context"
)

func (uc *SystemStatusUsecase) GetClosestIncomingSchedule(ctx context.Context) (*entity.SystemStatusSchedule, error) {
	// GetUpcomingSystemStatusSchedules returns schedules with start_time > NOW() ordered by start_time ASC
	schedules, err := uc.systemStatusRepo.GetUpcomingSystemStatusSchedules(ctx, 1)
	if err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		return nil, nil
	}

	return schedules[0], nil
}
