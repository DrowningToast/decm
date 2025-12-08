package event_registration

import (
	"context"
	"errors"
	"testing"
	"time"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock for EventRegistrationInvitationDataGateway
type MockEventRegistrationInvitationDataGateway struct {
	mock.Mock
}

func (m *MockEventRegistrationInvitationDataGateway) CreateEventRegistrationInvitation(ctx context.Context, params datagateway.CreateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) GetEventRegistrationInvitationByID(ctx context.Context, invitationID uuid.UUID) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, invitationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) GetEventRegistrationInvitationsByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) GetEventRegistrationInvitationByInboxMessageID(ctx context.Context, inboxMessageID uuid.UUID) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, inboxMessageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) GetEventRegistrationInvitationByEventIDAndCredential(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID, email *string, walletAddress *string) (*entity.EventRegistrationInvitation, *entity.InboxMessage, error) {
	args := m.Called(ctx, eventID, credentialID, email, walletAddress)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), nil, args.Error(2)
}

func (m *MockEventRegistrationInvitationDataGateway) UpdateEventRegistrationInvitation(ctx context.Context, id uuid.UUID, params datagateway.UpdateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) UpdateEventRegistrationInvitationAcceptedStatus(ctx context.Context, invitationID uuid.UUID, acceptedAt *time.Time) (*entity.EventRegistrationInvitation, error) {
	args := m.Called(ctx, invitationID, acceptedAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventRegistrationInvitation), args.Error(1)
}

func (m *MockEventRegistrationInvitationDataGateway) DeleteEventRegistrationInvitation(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestEventRegistrationUsecase_AcceptEventRegistrationInvitation(t *testing.T) {
	t.Run("should accept invitation successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      invitationID,
			EventId: eventID,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, nil, nil)
		mockInvitationDg.On("UpdateEventRegistrationInvitationAcceptedStatus", ctx, invitationID, mock.AnythingOfType("*time.Time")).Return(invitation, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, currentUser)

		// Assert
		require.NoError(t, err)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when user is not authenticated", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()

		uc := &EventRegistrationUsecase{}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, nil)

		// Assert
		require.Error(t, err)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrUnauthenticated.Code, *customError.Code)
	})

	t.Run("should return error when invitation not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		email := "test@example.com"
		currentUser := &auth.JwtClaims{
			UserId:        uuid.New(),
			Email:         &email,
			WalletAddress: "0x1234567890",
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, currentUser)

		// Assert
		require.Error(t, err)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when invitation belongs to different user", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      invitationID,
			EventId: eventID,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		// Return nil to simulate invitation not belonging to current user
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(nil, nil, customerror.NewWithPreset(&customerror.ErrUnauthorized, errors.New("unauthorized")))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, currentUser)

		// Assert
		require.Error(t, err)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customError.Code)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when invitation ID mismatch", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		differentInvitationID := uuid.New()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      invitationID,
			EventId: eventID,
		}

		differentInvitation := &entity.EventRegistrationInvitation{
			Id:      differentInvitationID,
			EventId: eventID,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		// Return different invitation ID
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(differentInvitation, nil, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, currentUser)

		// Assert
		require.Error(t, err)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customError.Code)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invitationID := uuid.New()
		eventID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		currentUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		invitation := &entity.EventRegistrationInvitation{
			Id:      invitationID,
			EventId: eventID,
		}

		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockInvitationDg.On("GetEventRegistrationInvitationByID", ctx, invitationID).Return(invitation, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, nil, nil)
		mockInvitationDg.On("UpdateEventRegistrationInvitationAcceptedStatus", ctx, invitationID, mock.AnythingOfType("*time.Time")).Return(nil, errors.New("database error"))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		// Act
		err := uc.AcceptEventRegistrationInvitation(ctx, invitationID, currentUser)

		// Assert
		require.Error(t, err)
		mockInvitationDg.AssertExpectations(t)
	})
}
