package system_status

import (
	"context"
	"errors"
	"testing"
	"time"

	"apps/backend/core-api/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSystemStatusScheduleDataGateway struct {
	mock.Mock
}

func (m *MockSystemStatusScheduleDataGateway) GetSystemStatusScheduleById(ctx context.Context, id int32) (*entity.SystemStatusSchedule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.SystemStatusSchedule), args.Error(1)
}

func (m *MockSystemStatusScheduleDataGateway) GetSystemStatusScheduleByOrderId(ctx context.Context, orderId int32) (*entity.SystemStatusSchedule, error) {
	args := m.Called(ctx, orderId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.SystemStatusSchedule), args.Error(1)
}

func (m *MockSystemStatusScheduleDataGateway) GetCurrentSystemStatus(ctx context.Context) (*entity.SystemStatusSchedule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.SystemStatusSchedule), args.Error(1)
}

func (m *MockSystemStatusScheduleDataGateway) GetUpcomingSystemStatusSchedules(ctx context.Context, limit int32) ([]*entity.SystemStatusSchedule, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.SystemStatusSchedule), args.Error(1)
}

func (m *MockSystemStatusScheduleDataGateway) GetSystemStatusScheduleHistory(ctx context.Context, limit int32, offset int32) ([]*entity.SystemStatusSchedule, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.SystemStatusSchedule), args.Error(1)
}

func (m *MockSystemStatusScheduleDataGateway) GetPlannedMaintenanceSchedules(ctx context.Context) ([]*entity.SystemStatusSchedule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.SystemStatusSchedule), args.Error(1)
}

func (m *MockSystemStatusScheduleDataGateway) GetSystemStatusSchedulesUpdatedBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*entity.SystemStatusSchedule, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.SystemStatusSchedule), args.Error(1)
}

func (m *MockSystemStatusScheduleDataGateway) CountSystemStatusSchedules(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestGetLatestSchedules(t *testing.T) {
	ctx := context.Background()

	t.Run("should return latest schedules", func(t *testing.T) {
		mockRepo := new(MockSystemStatusScheduleDataGateway)
		uc := NewSystemStatusUsecase(mockRepo)
		pageSize := int32(10)
		expectedSchedules := []*entity.SystemStatusSchedule{
			{ID: 1, OrderId: 1, Status: entity.SystemStatusOperating},
		}

		mockRepo.On("GetSystemStatusScheduleHistory", ctx, pageSize, int32(0)).Return(expectedSchedules, nil)

		result, err := uc.GetLatestSchedules(ctx, pageSize)

		assert.NoError(t, err)
		assert.Equal(t, expectedSchedules, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repo fails", func(t *testing.T) {
		mockRepo := new(MockSystemStatusScheduleDataGateway)
		uc := NewSystemStatusUsecase(mockRepo)
		pageSize := int32(10)
		mockRepo.On("GetSystemStatusScheduleHistory", ctx, pageSize, int32(0)).Return(nil, errors.New("db error"))

		result, err := uc.GetLatestSchedules(ctx, pageSize)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetSchedulesBetween(t *testing.T) {
	ctx := context.Background()

	t.Run("should return schedules between period", func(t *testing.T) {
		mockRepo := new(MockSystemStatusScheduleDataGateway)
		uc := NewSystemStatusUsecase(mockRepo)
		start := time.Now().Add(-1 * time.Hour)
		end := time.Now()
		expectedSchedules := []*entity.SystemStatusSchedule{
			{ID: 1, OrderId: 1, Status: entity.SystemStatusOperating},
		}

		mockRepo.On("GetSystemStatusSchedulesUpdatedBetween", ctx, start, end).Return(expectedSchedules, nil)

		result, err := uc.GetSchedulesBetween(ctx, start, end)

		assert.NoError(t, err)
		assert.Equal(t, expectedSchedules, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetClosestIncomingSchedule(t *testing.T) {
	ctx := context.Background()

	t.Run("should return closest incoming schedule", func(t *testing.T) {
		mockRepo := new(MockSystemStatusScheduleDataGateway)
		uc := NewSystemStatusUsecase(mockRepo)
		expectedSchedule := &entity.SystemStatusSchedule{ID: 2, OrderId: 2, Status: entity.SystemStatusMaintenance}
		mockRepo.On("GetUpcomingSystemStatusSchedules", ctx, int32(1)).Return([]*entity.SystemStatusSchedule{expectedSchedule}, nil)

		result, err := uc.GetClosestIncomingSchedule(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedSchedule, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return nil when no upcoming schedules", func(t *testing.T) {
		mockRepo := new(MockSystemStatusScheduleDataGateway)
		uc := NewSystemStatusUsecase(mockRepo)
		mockRepo.On("GetUpcomingSystemStatusSchedules", ctx, int32(1)).Return([]*entity.SystemStatusSchedule{}, nil)

		result, err := uc.GetClosestIncomingSchedule(ctx)

		assert.NoError(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
