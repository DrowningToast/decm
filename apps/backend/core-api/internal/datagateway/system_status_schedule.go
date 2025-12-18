package datagateway

import (
	"context"
	"time"

	"apps/backend/core-api/internal/entity"
)

type SystemStatusScheduleDataGateway interface {
	GetSystemStatusScheduleById(ctx context.Context, id int32) (*entity.SystemStatusSchedule, error)
	GetSystemStatusScheduleByOrderId(ctx context.Context, orderId int32) (*entity.SystemStatusSchedule, error)
	GetCurrentSystemStatus(ctx context.Context) (*entity.SystemStatusSchedule, error)
	GetUpcomingSystemStatusSchedules(ctx context.Context, limit int32) ([]*entity.SystemStatusSchedule, error)
	GetSystemStatusScheduleHistory(ctx context.Context, limit int32, offset int32) ([]*entity.SystemStatusSchedule, error)
	GetPlannedMaintenanceSchedules(ctx context.Context) ([]*entity.SystemStatusSchedule, error)
	GetSystemStatusSchedulesUpdatedBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*entity.SystemStatusSchedule, error)
	CountSystemStatusSchedules(ctx context.Context) (int64, error)
}
