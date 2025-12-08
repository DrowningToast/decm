package inbox

import (
	"context"
	"errors"
	"testing"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInboxUsecase_GetMyInboxMessages(t *testing.T) {
	t.Run("should get user's inbox messages by credential ID", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		email := "user@example.com"

		user := auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		expectedMessages := []*entity.InboxMessage{
			{
				Id:                   uuid.New(),
				ReceiverCredentialId: &userID,
				Subject:              "Event Invitation",
				Body:                 "You're invited to our event",
				IsRead:               0,
			},
			{
				Id:                   uuid.New(),
				ReceiverCredentialId: &userID,
				Subject:              "Certificate Ready",
				Body:                 "Your certificate is ready",
				IsRead:               1,
			},
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessagesByCredentialID", ctx, userID).Return(expectedMessages, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		messages, err := uc.GetMyInboxMessages(ctx, user)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, messages)
		assert.Len(t, messages, 2)
		assert.Equal(t, "Event Invitation", messages[0].Subject)
		assert.Equal(t, "Certificate Ready", messages[1].Subject)
		assert.Equal(t, 0, messages[0].IsRead) // Unread
		assert.Equal(t, 1, messages[1].IsRead) // Read
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return empty list when user has no messages", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		email := "newuser@example.com"

		user := auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessagesByCredentialID", ctx, userID).Return([]*entity.InboxMessage{}, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		messages, err := uc.GetMyInboxMessages(ctx, user)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, messages)
		assert.Empty(t, messages)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return error when database fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		email := "user@example.com"

		user := auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessagesByCredentialID", ctx, userID).Return(nil, errors.New("database connection error"))

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		messages, err := uc.GetMyInboxMessages(ctx, user)

		// Assert
		require.Error(t, err)
		assert.Nil(t, messages)
		mockInboxDg.AssertExpectations(t)
	})
}

func TestInboxUsecase_GetInboxMessagesByCredentailID(t *testing.T) {
	t.Run("should get messages for specific credential", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()

		messages := []*entity.InboxMessage{
			{
				Id:                   uuid.New(),
				ReceiverCredentialId: &credentialID,
				Subject:              "Welcome",
				IsRead:               0,
			},
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessagesByCredentialID", ctx, credentialID).Return(messages, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.GetInboxMessagesByCredentailID(ctx, credentialID)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)
		mockInboxDg.AssertExpectations(t)
	})
}

func TestInboxUsecase_GetInboxMessagesByReceiverEmail(t *testing.T) {
	t.Run("should get messages sent to specific email", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "attendee@example.com"

		messages := []*entity.InboxMessage{
			{
				Id:            uuid.New(),
				ReceiverEmail: &email,
				Subject:       "Event Registration Confirmed",
				Body:          "Your registration is confirmed",
				IsRead:        0,
			},
			{
				Id:            uuid.New(),
				ReceiverEmail: &email,
				Subject:       "Event Reminder",
				Body:          "Event starts tomorrow",
				IsRead:        0,
			},
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessagesByReceiverEmail", ctx, email).Return(messages, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.GetInboxMessagesByReceiverEmail(ctx, email)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "Event Registration Confirmed", result[0].Subject)
		assert.Equal(t, "Event Reminder", result[1].Subject)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return error when email not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "nonexistent@example.com"

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessagesByReceiverEmail", ctx, email).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("no messages")))

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.GetInboxMessagesByReceiverEmail(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		mockInboxDg.AssertExpectations(t)
	})
}

func TestInboxUsecase_GetInboxMessagesByReceiverWalletAddress(t *testing.T) {
	t.Run("should get messages sent to specific wallet", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		walletAddress := "0x1234567890abcdef"

		messages := []*entity.InboxMessage{
			{
				Id:                    uuid.New(),
				ReceiverWalletAddress: &walletAddress,
				Subject:               "NFT Certificate Issued",
				Body:                  "Your certificate NFT is ready",
				IsRead:                0,
			},
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessagesByReceiverWalletAddress", ctx, walletAddress).Return(messages, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.GetInboxMessagesByReceiverWalletAddress(ctx, walletAddress)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, "NFT Certificate Issued", result[0].Subject)
		mockInboxDg.AssertExpectations(t)
	})
}

func TestInboxUsecase_GetInboxMessagesBySenderCredentialID(t *testing.T) {
	t.Run("should get messages sent by specific organizer", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		senderID := uuid.New()

		messages := []*entity.InboxMessage{
			{
				Id:                 uuid.New(),
				SenderCredentialId: &senderID,
				Subject:            "Event Update",
				Body:               "Event location changed",
				IsRead:             1,
			},
			{
				Id:                 uuid.New(),
				SenderCredentialId: &senderID,
				Subject:            "Event Cancelled",
				Body:               "Unfortunately we must cancel",
				IsRead:             1,
			},
		}

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessagesBySenderCredentialID", ctx, senderID).Return(messages, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.GetInboxMessagesBySenderCredentialID(ctx, senderID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "Event Update", result[0].Subject)
		assert.Equal(t, "Event Cancelled", result[1].Subject)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return empty when sender has sent no messages", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		senderID := uuid.New()

		mockInboxDg := new(MockInboxMessageDataGateway)
		mockInboxDg.On("GetInboxMessagesBySenderCredentialID", ctx, senderID).Return([]*entity.InboxMessage{}, nil)

		uc := &InboxUsecase{
			InboxMessageDg: mockInboxDg,
		}

		// Act
		result, err := uc.GetInboxMessagesBySenderCredentialID(ctx, senderID)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, result)
		mockInboxDg.AssertExpectations(t)
	})
}
