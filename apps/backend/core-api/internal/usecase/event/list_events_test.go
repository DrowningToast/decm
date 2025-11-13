package event

import (
	"context"
	"errors"
	"testing"
	"time"

	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockListEventDataGateway struct {
	mock.Mock
}

func (m *MockListEventDataGateway) ListEvents(ctx context.Context, offset *int32, limit *int32) ([]*entity.Event, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func (m *MockListEventDataGateway) ListEventsByEventAttendeeCredentialID(ctx context.Context, credentialId uuid.UUID, offset *int32, limit *int32) ([]*entity.Event, error) {
	args := m.Called(ctx, credentialId, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func TestListEvents(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()

	mockEvents := []*entity.Event{
		{
			ID:          uuid.New(),
			Title:       "Active Event 1",
			EventStatus: entity.EventStatusActive,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "Inactive Event 1",
			EventStatus: entity.EventStatusInactive,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "Closed Event 1",
			EventStatus: entity.EventStatusClosed,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "Active Event 2",
			EventStatus: entity.EventStatusActive,
			CreatedAt:   time.Now(),
		},
	}

	t.Run("should return all active events when only include active is true", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockListEventDataGateway)
		mockEventDg.On("ListEvents", ctx, nil, nil).Return(mockEvents, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		params := ListEventsParameters{
			IncludeActiveEvents:   true,
			IncludeInactiveEvents: false,
			IncludeClosedEvents:   false,
		}

		// Act
		events, err := uc.ListEvents(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 2, len(events))
		for _, event := range events {
			assert.Equal(t, entity.EventStatusActive, event.EventStatus)
		}
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return all inactive events when only include inactive is true", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockListEventDataGateway)
		mockEventDg.On("ListEvents", ctx, nil, nil).Return(mockEvents, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		params := ListEventsParameters{
			IncludeActiveEvents:   false,
			IncludeInactiveEvents: true,
			IncludeClosedEvents:   false,
		}

		// Act
		events, err := uc.ListEvents(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(events))
		assert.Equal(t, entity.EventStatusInactive, events[0].EventStatus)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return all closed events when only include closed is true", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockListEventDataGateway)
		mockEventDg.On("ListEvents", ctx, nil, nil).Return(mockEvents, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		params := ListEventsParameters{
			IncludeActiveEvents:   false,
			IncludeInactiveEvents: false,
			IncludeClosedEvents:   true,
		}

		// Act
		events, err := uc.ListEvents(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(events))
		assert.Equal(t, entity.EventStatusClosed, events[0].EventStatus)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return multiple status events when multiple filters are true", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockListEventDataGateway)
		mockEventDg.On("ListEvents", ctx, nil, nil).Return(mockEvents, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		params := ListEventsParameters{
			IncludeActiveEvents:   true,
			IncludeInactiveEvents: true,
			IncludeClosedEvents:   false,
		}

		// Act
		events, err := uc.ListEvents(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 3, len(events))
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return empty list when no filters match", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockListEventDataGateway)
		mockEventDg.On("ListEvents", ctx, nil, nil).Return(mockEvents, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		params := ListEventsParameters{
			IncludeActiveEvents:   false,
			IncludeInactiveEvents: false,
			IncludeClosedEvents:   false,
		}

		// Act
		events, err := uc.ListEvents(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 0, len(events))
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return user joined events filtered by status", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockListEventDataGateway)
		mockEventDg.On("ListEventsByEventAttendeeCredentialID", ctx, userId, nil, nil).
			Return(mockEvents, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		params := ListEventsParameters{
			OnlyUserJoinedEvents:  currentUser,
			IncludeActiveEvents:   true,
			IncludeInactiveEvents: false,
			IncludeClosedEvents:   false,
		}

		// Act
		events, err := uc.ListEvents(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 2, len(events))
		for _, event := range events {
			assert.Equal(t, entity.EventStatusActive, event.EventStatus)
		}
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should handle error from data gateway", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockListEventDataGateway)
		mockEventDg.On("ListEvents", ctx, nil, nil).
			Return(nil, errors.New("database error"))

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		params := ListEventsParameters{
			IncludeActiveEvents:   true,
			IncludeInactiveEvents: true,
			IncludeClosedEvents:   true,
		}

		// Act
		events, err := uc.ListEvents(ctx, params)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, events)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should handle error from user joined events data gateway", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockListEventDataGateway)
		mockEventDg.On("ListEventsByEventAttendeeCredentialID", ctx, userId, nil, nil).
			Return(nil, errors.New("database error"))

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		params := ListEventsParameters{
			OnlyUserJoinedEvents:  currentUser,
			IncludeActiveEvents:   true,
			IncludeInactiveEvents: true,
			IncludeClosedEvents:   true,
		}

		// Act
		events, err := uc.ListEvents(ctx, params)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, events)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return empty list when no events exist", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockListEventDataGateway)
		mockEventDg.On("ListEvents", ctx, nil, nil).Return([]*entity.Event{}, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		params := ListEventsParameters{
			IncludeActiveEvents:   true,
			IncludeInactiveEvents: true,
			IncludeClosedEvents:   true,
		}

		// Act
		events, err := uc.ListEvents(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 0, len(events))
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should include all event statuses when all filters are true", func(t *testing.T) {
		// Arrange
		mockEventDg := new(MockListEventDataGateway)
		mockEventDg.On("ListEvents", ctx, nil, nil).Return(mockEvents, nil)

		uc := &EventUsecase{
			EventDataGateway: mockEventDg,
		}

		params := ListEventsParameters{
			IncludeActiveEvents:   true,
			IncludeInactiveEvents: true,
			IncludeClosedEvents:   true,
		}

		// Act
		events, err := uc.ListEvents(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, len(mockEvents), len(events))
		mockEventDg.AssertExpectations(t)
	})
}
