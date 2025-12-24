package eventconfig

import (
	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock for EventRegistrationConfigDataGateway
type MockEventRegistrationConfigDataGateway struct {
	mock.Mock
}

func (m *MockEventRegistrationConfigDataGateway) CreateEventRegistrationConfig(ctx context.Context, params generated.CreateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationConfig), args.Error(1)
}

func (m *MockEventRegistrationConfigDataGateway) GetEventRegistrationConfigByEventId(ctx context.Context, eventID uuid.UUID) (*entity.EventRegistrationConfig, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationConfig), args.Error(1)
}

func (m *MockEventRegistrationConfigDataGateway) GetEventRegistrationConfigPasswordByEventId(ctx context.Context, eventID uuid.UUID) (*entity.EventRegistrationConfig, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationConfig), args.Error(1)
}

func (m *MockEventRegistrationConfigDataGateway) UpdateEventRegistrationConfig(ctx context.Context, params generated.UpdateEventRegistrationConfigParams) (*entity.EventRegistrationConfig, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationConfig), args.Error(1)
}

func (m *MockEventRegistrationConfigDataGateway) DeleteEventRegistrationConfig(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func TestEventConfigUsecase_CheckEventPassword(t *testing.T) {
	t.Run("should return true when password matches", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		plainPassword := "correct-password"

		// Hash the password
		hashedPassword, err := hashutils.HashPassword(plainPassword)
		require.NoError(t, err)

		config := &entity.EventRegistrationConfig{
			EventID:              eventID,
			RegistrationPassword: &hashedPassword,
		}

		mockConfigDg := new(MockEventRegistrationConfigDataGateway)
		mockConfigDg.On("GetEventRegistrationConfigPasswordByEventId", ctx, eventID).Return(config, nil)

		uc := &EventConfigUsecase{
			EventRegistrationDg: mockConfigDg,
		}

		// Act
		match, err := uc.CheckEventPassword(ctx, eventID, plainPassword)

		// Assert
		require.NoError(t, err)
		assert.True(t, match)
		mockConfigDg.AssertExpectations(t)
	})

	t.Run("should return false when password does not match", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		correctPassword := "correct-password"
		wrongPassword := "wrong-password"

		// Hash the correct password
		hashedPassword, err := hashutils.HashPassword(correctPassword)
		require.NoError(t, err)

		config := &entity.EventRegistrationConfig{
			EventID:              eventID,
			RegistrationPassword: &hashedPassword,
		}

		mockConfigDg := new(MockEventRegistrationConfigDataGateway)
		mockConfigDg.On("GetEventRegistrationConfigPasswordByEventId", ctx, eventID).Return(config, nil)

		uc := &EventConfigUsecase{
			EventRegistrationDg: mockConfigDg,
		}

		// Act
		match, err := uc.CheckEventPassword(ctx, eventID, wrongPassword)

		// Assert
		require.NoError(t, err)
		assert.False(t, match)
		mockConfigDg.AssertExpectations(t)
	})

	t.Run("should return error when config not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		password := "any-password"

		mockConfigDg := new(MockEventRegistrationConfigDataGateway)
		mockConfigDg.On("GetEventRegistrationConfigPasswordByEventId", ctx, eventID).Return(nil, errors.New("config not found"))

		uc := &EventConfigUsecase{
			EventRegistrationDg: mockConfigDg,
		}

		// Act
		match, err := uc.CheckEventPassword(ctx, eventID, password)

		// Assert
		require.Error(t, err)
		assert.False(t, match)
		mockConfigDg.AssertExpectations(t)
	})

	t.Run("should return error when registration password is not set", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		password := "any-password"

		config := &entity.EventRegistrationConfig{
			EventID:              eventID,
			RegistrationPassword: nil, // Password not set
		}

		mockConfigDg := new(MockEventRegistrationConfigDataGateway)
		mockConfigDg.On("GetEventRegistrationConfigPasswordByEventId", ctx, eventID).Return(config, nil)

		uc := &EventConfigUsecase{
			EventRegistrationDg: mockConfigDg,
		}

		// Act
		match, err := uc.CheckEventPassword(ctx, eventID, password)

		// Assert
		require.Error(t, err)
		assert.False(t, match)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customError.Code)
		mockConfigDg.AssertExpectations(t)
	})

	t.Run("should return error when config is nil", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		eventID := uuid.New()
		password := "any-password"

		mockConfigDg := new(MockEventRegistrationConfigDataGateway)
		mockConfigDg.On("GetEventRegistrationConfigPasswordByEventId", ctx, eventID).Return(nil, nil)

		uc := &EventConfigUsecase{
			EventRegistrationDg: mockConfigDg,
		}

		// Act
		match, err := uc.CheckEventPassword(ctx, eventID, password)

		// Assert
		require.Error(t, err)
		assert.False(t, match)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customError.Code)
		mockConfigDg.AssertExpectations(t)
	})
}
