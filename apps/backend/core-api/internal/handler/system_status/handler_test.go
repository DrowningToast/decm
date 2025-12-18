package system_status

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/system_status"

	"github.com/gofiber/fiber/v2"
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

func setupApp(h *Handler) *fiber.App {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	app := fiber.New(fiber.Config{
		ErrorHandler: customerror.GetErrFiberHandler(logger),
	})
	return app
}

func TestGetLatestSchedulesHandler(t *testing.T) {
	t.Run("should return latest schedules", func(t *testing.T) {
		mockRepo := new(MockSystemStatusScheduleDataGateway)
		uc := system_status.NewSystemStatusUsecase(mockRepo)
		h := NewHandler(uc, nil)
		app := setupApp(h)
		app.Get("/latest", h.GetLatestSchedules)

		expectedSchedules := []*entity.SystemStatusSchedule{
			{ID: 1, OrderId: 1, Status: entity.SystemStatusOperating},
		}
		mockRepo.On("GetSystemStatusScheduleHistory", mock.Anything, int32(10), int32(0)).Return(expectedSchedules, nil)

		req := httptest.NewRequest(http.MethodGet, "/latest?page_size=10", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result GetLatestSchedulesResponse
		json.Unmarshal(body, &result)

		assert.Len(t, result.Schedules, 1)
		assert.Equal(t, int32(1), result.Schedules[0].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return 400 when page_size is missing", func(t *testing.T) {
		mockRepo := new(MockSystemStatusScheduleDataGateway)
		uc := system_status.NewSystemStatusUsecase(mockRepo)
		h := NewHandler(uc, nil)
		app := setupApp(h)
		app.Get("/latest", h.GetLatestSchedules)

		req := httptest.NewRequest(http.MethodGet, "/latest", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetSchedulesBetweenHandler(t *testing.T) {
	t.Run("should return schedules between period", func(t *testing.T) {
		mockRepo := new(MockSystemStatusScheduleDataGateway)
		uc := system_status.NewSystemStatusUsecase(mockRepo)
		h := NewHandler(uc, nil)
		app := setupApp(h)
		app.Get("/period", h.GetSchedulesBetween)

		now := time.Now().Unix()
		start := now - 3600
		end := now

		expectedSchedules := []*entity.SystemStatusSchedule{
			{ID: 1, OrderId: 1, Status: entity.SystemStatusOperating},
		}

		mockRepo.On("GetSystemStatusSchedulesUpdatedBetween", mock.Anything, mock.Anything, mock.Anything).Return(expectedSchedules, nil)

		req := httptest.NewRequest(http.MethodGet, "/period?start_time="+fmt.Sprintf("%d", start)+"&end_time="+fmt.Sprintf("%d", end), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result GetSchedulesBetweenResponse
		json.Unmarshal(body, &result)

		assert.Len(t, result.Schedules, 1)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetClosestIncomingScheduleHandler(t *testing.T) {
	t.Run("should return closest incoming schedule", func(t *testing.T) {
		mockRepo := new(MockSystemStatusScheduleDataGateway)
		uc := system_status.NewSystemStatusUsecase(mockRepo)
		h := NewHandler(uc, nil)
		app := setupApp(h)
		app.Get("/closest-incoming", h.GetClosestIncomingSchedule)

		expectedSchedule := &entity.SystemStatusSchedule{ID: 2, OrderId: 2, Status: entity.SystemStatusMaintenance}
		mockRepo.On("GetUpcomingSystemStatusSchedules", mock.Anything, int32(1)).Return([]*entity.SystemStatusSchedule{expectedSchedule}, nil)

		req := httptest.NewRequest(http.MethodGet, "/closest-incoming", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result GetClosestIncomingScheduleResponse
		json.Unmarshal(body, &result)

		assert.NotNil(t, result.Schedule)
		assert.Equal(t, int32(2), result.Schedule.ID)
		mockRepo.AssertExpectations(t)
	})
}
