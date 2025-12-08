package issuer

import (
	"context"
	"decm-database/go/generated"
	"errors"
	"testing"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockEventIssuerDataGateway struct {
	mock.Mock
}

func (m *MockEventIssuerDataGateway) ListVerifiedIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventIssuerDataGateway) ListIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventIssuerDataGateway) SearchIssuerCredentialsByWalletAddress(ctx context.Context, searchQuery string, limitCount int, offsetCount int) ([]entity.AuthenticationCredential, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventIssuerDataGateway) ListAllIssuerCredentials(ctx context.Context, limitCount int) ([]entity.AuthenticationCredential, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventIssuerDataGateway) GetEventsByIssuerCredentialID(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetEventIssuersByCredentialIDRow, error) {
	args := m.Called(ctx, issuerCredentialID, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]generated.GetEventIssuersByCredentialIDRow), args.Error(1)
}

func (m *MockEventIssuerDataGateway) GetIssuerEventsWithDetails(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetIssuerEventsWithDetailsRow, error) {
	args := m.Called(ctx, issuerCredentialID, limitCount, offsetCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]generated.GetIssuerEventsWithDetailsRow), args.Error(1)
}

func TestGetEventsByIssuerCredentialID(t *testing.T) {
	ctx := context.Background()

	issuerCredId := uuid.New()
	mockEvents := []generated.GetEventIssuersByCredentialIDRow{
		{
			EventID:    uuid.New(),
			EventTitle: "Event 1",
		},
		{
			EventID:    uuid.New(),
			EventTitle: "Event 2",
		},
	}

	t.Run("should return events for valid issuer credential ID", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockEventIssuerDataGateway)
		mockIssuerDg.On("GetEventsByIssuerCredentialID", ctx, issuerCredId.String(), int32(10), int32(0)).
			Return(mockEvents, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		// Act
		result, err := uc.GetEventsByIssuerCredentialID(ctx, issuerCredId.String(), 10, 0)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "Event 1", result[0].EventTitle)
		assert.Equal(t, "Event 2", result[1].EventTitle)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should return empty list when no events found", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockEventIssuerDataGateway)
		mockIssuerDg.On("GetEventsByIssuerCredentialID", ctx, issuerCredId.String(), int32(10), int32(0)).
			Return([]generated.GetEventIssuersByCredentialIDRow{}, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		// Act
		result, err := uc.GetEventsByIssuerCredentialID(ctx, issuerCredId.String(), 10, 0)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 0, len(result))
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should handle error from data gateway", func(t *testing.T) {
		// Arrange
		mockIssuerDg := new(MockEventIssuerDataGateway)
		mockIssuerDg.On("GetEventsByIssuerCredentialID", ctx, issuerCredId.String(), int32(10), int32(0)).
			Return(nil, errors.New("database error"))

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		// Act
		result, err := uc.GetEventsByIssuerCredentialID(ctx, issuerCredId.String(), 10, 0)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should apply pagination correctly", func(t *testing.T) {
		// Arrange
		paginatedEvents := []generated.GetEventIssuersByCredentialIDRow{
			{
				EventID:    uuid.New(),
				EventTitle: "Event 3",
			},
		}

		mockIssuerDg := new(MockEventIssuerDataGateway)
		mockIssuerDg.On("GetEventsByIssuerCredentialID", ctx, issuerCredId.String(), int32(1), int32(2)).
			Return(paginatedEvents, nil)

		uc := &IssuerUsecase{
			IssuerDg: mockIssuerDg,
		}

		// Act
		result, err := uc.GetEventsByIssuerCredentialID(ctx, issuerCredId.String(), 1, 2)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "Event 3", result[0].EventTitle)
		mockIssuerDg.AssertExpectations(t)
	})
}
