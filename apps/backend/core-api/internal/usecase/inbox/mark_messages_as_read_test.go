package inbox

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock for InboxMessageDataGateway
type MockInboxMessageDataGateway struct {
	mock.Mock
}

func (m *MockInboxMessageDataGateway) GetInboxMessageByID(ctx context.Context, id uuid.UUID) (*entity.InboxMessage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) CreateInboxMessage(ctx context.Context, params datagateway.CreateInboxMessageParameters) (*entity.InboxMessage, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetInboxMessagesByCredentialID(ctx context.Context, params datagateway.GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetInboxMessagesByReceiverEmail(ctx context.Context, receiverEmail string) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, receiverEmail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetInboxMessagesByReceiverWalletAddress(ctx context.Context, walletAddress string) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, walletAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetInboxMessagesBySenderCredentialID(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, credentialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) UpdateInboxMessageReadStatus(ctx context.Context, id uuid.UUID, isRead int) (*entity.InboxMessage, error) {
	args := m.Called(ctx, id, isRead)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) UpdateInboxMessageReadStatusAll(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, credentialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDataGateway) GetUnreadInboxMessageCountByCredentialID(ctx context.Context, params datagateway.GetInboxMessagesByCredentialIDParameters) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func TestInboxUsecase_MarkMessageAsRead(t *testing.T) {
	t.Run("should mark message as read when user is authorized", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		messageID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		user := auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		message := &entity.InboxMessage{
			Id:                   messageID,
			ReceiverCredentialId: &userID,
			IsRead:               0,
		}

		updatedMessage := &entity.InboxMessage{
			Id:                   messageID,
			ReceiverCredentialId: &userID,
			IsRead:               1,
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessageByID", ctx, messageID).Return(message, nil)
		mockInboxDg.On("UpdateInboxMessageReadStatus", ctx, messageID, 1).Return(updatedMessage, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.MarkMessageAsRead(ctx, user, messageID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.IsRead)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return error when message not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		messageID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		user := auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessageByID", ctx, messageID).Return(nil, errors.New("not found"))

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.MarkMessageAsRead(ctx, user, messageID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return error when message is nil", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		messageID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		user := auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessageByID", ctx, messageID).Return(nil, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.MarkMessageAsRead(ctx, user, messageID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrNotFound.Code, *customError.Code)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return error when user is not authorized", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		messageID := uuid.New()
		userID := uuid.New()
		differentUserID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		user := auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		message := &entity.InboxMessage{
			Id:                   messageID,
			ReceiverCredentialId: &differentUserID, // Different user ID
			IsRead:               0,
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessageByID", ctx, messageID).Return(message, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.MarkMessageAsRead(ctx, user, messageID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrForbidden.Code, *customError.Code)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should mark message as read when user email matches", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		messageID := uuid.New()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		user := auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		message := &entity.InboxMessage{
			Id:            messageID,
			ReceiverEmail: &email, // Authorized by email
			IsRead:        0,
		}

		updatedMessage := &entity.InboxMessage{
			Id:            messageID,
			ReceiverEmail: &email,
			IsRead:        1,
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessageByID", ctx, messageID).Return(message, nil)
		mockInboxDg.On("UpdateInboxMessageReadStatus", ctx, messageID, 1).Return(updatedMessage, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.MarkMessageAsRead(ctx, user, messageID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.IsRead)
		mockInboxDg.AssertExpectations(t)
	})
}

func TestInboxUsecase_MarkAllMessagesAsRead(t *testing.T) {
	t.Run("should mark all messages as read", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		user := auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		updatedMessages := []*entity.InboxMessage{
			{Id: uuid.New(), IsRead: 1},
			{Id: uuid.New(), IsRead: 1},
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("UpdateInboxMessageReadStatusAll", ctx, userID).Return(updatedMessages, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		results, err := uc.MarkAllMessagesAsRead(ctx, user)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 2)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		email := "test@example.com"
		walletAddress := "0x1234567890"

		user := auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: walletAddress,
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("UpdateInboxMessageReadStatusAll", ctx, userID).Return(nil, errors.New("database error"))

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		results, err := uc.MarkAllMessagesAsRead(ctx, user)

		// Assert
		require.Error(t, err)
		assert.Nil(t, results)
		mockInboxDg.AssertExpectations(t)
	})
}
