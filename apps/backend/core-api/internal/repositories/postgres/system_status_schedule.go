package postgres

import (
	"context"
	"decm-database/go/generated"
	"time"

	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	datagateway "apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
)

var _ datagateway.SystemStatusScheduleDataGateway = (*Repository)(nil)

func (r *Repository) mapSystemStatusScheduleToEntity(row generated.SystemStatusSchedule) *entity.SystemStatusSchedule {
	return &entity.SystemStatusSchedule{
		ID:             row.ID,
		OrderId:        row.OrderID,
		StartTime:      row.StartTime,
		PlannedEndTime: pgmapper.PgTimestampzToTimePtr(row.PlannedEndTime),
		Status:         entity.SystemStatus(row.Status),
		IsPlanned:      row.IsPlanned,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
		DeletedAt:      pgmapper.PgTimestampzToTimePtr(row.DeletedAt),
	}
}

func (r *Repository) GetSystemStatusScheduleById(ctx context.Context, id int32) (*entity.SystemStatusSchedule, error) {
	result, err := r.queries.GetSystemStatusScheduleById(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return r.mapSystemStatusScheduleToEntity(result), nil
}

func (r *Repository) GetSystemStatusScheduleByOrderId(ctx context.Context, orderId int32) (*entity.SystemStatusSchedule, error) {
	result, err := r.queries.GetSystemStatusScheduleByOrderId(ctx, orderId)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return r.mapSystemStatusScheduleToEntity(result), nil
}

func (r *Repository) GetCurrentSystemStatus(ctx context.Context) (*entity.SystemStatusSchedule, error) {
	result, err := r.queries.GetCurrentSystemStatus(ctx)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return r.mapSystemStatusScheduleToEntity(result), nil
}

func (r *Repository) GetUpcomingSystemStatusSchedules(ctx context.Context, limit int32) ([]*entity.SystemStatusSchedule, error) {
	results, err := r.queries.GetUpcomingSystemStatusSchedules(ctx, limit)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	schedules := make([]*entity.SystemStatusSchedule, len(results))
	for i, row := range results {
		schedules[i] = r.mapSystemStatusScheduleToEntity(row)
	}

	return schedules, nil
}

func (r *Repository) GetSystemStatusScheduleHistory(ctx context.Context, limit int32, offset int32) ([]*entity.SystemStatusSchedule, error) {
	results, err := r.queries.GetSystemStatusScheduleHistory(ctx, generated.GetSystemStatusScheduleHistoryParams{
		LimitCount:  limit,
		OffsetCount: offset,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	schedules := make([]*entity.SystemStatusSchedule, len(results))
	for i, row := range results {
		schedules[i] = r.mapSystemStatusScheduleToEntity(row)
	}

	return schedules, nil
}

func (r *Repository) GetPlannedMaintenanceSchedules(ctx context.Context) ([]*entity.SystemStatusSchedule, error) {
	results, err := r.queries.GetPlannedMaintenanceSchedules(ctx)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	schedules := make([]*entity.SystemStatusSchedule, len(results))
	for i, row := range results {
		schedules[i] = r.mapSystemStatusScheduleToEntity(row)
	}

	return schedules, nil
}

func (r *Repository) GetSystemStatusSchedulesUpdatedBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*entity.SystemStatusSchedule, error) {
	results, err := r.queries.GetSystemStatusSchedulesUpdatedBetween(ctx, generated.GetSystemStatusSchedulesUpdatedBetweenParams{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	schedules := make([]*entity.SystemStatusSchedule, len(results))
	for i, row := range results {
		schedules[i] = r.mapSystemStatusScheduleToEntity(row)
	}

	return schedules, nil
}

func (r *Repository) CountSystemStatusSchedules(ctx context.Context) (int64, error) {
	count, err := r.queries.CountSystemStatusSchedules(ctx)
	if err != nil {
		return 0, pgerrutils.ParsePgError(err)
	}

	return count, nil
}
