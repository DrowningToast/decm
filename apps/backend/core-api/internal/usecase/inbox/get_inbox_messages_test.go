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

// MockInboxMessageDg is a mock for InboxMessageDataGateway (reusing from inbox_test.go)
type MockInboxMessageDg struct {
	mock.Mock
}

func (m *MockInboxMessageDg) CreateInboxMessage(ctx context.Context, params datagateway.CreateInboxMessageParameters) (*entity.InboxMessage, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetInboxMessagesByCredentialID(ctx context.Context, params datagateway.GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetInboxMessagesByReceiverEmail(ctx context.Context, email string) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetInboxMessagesByReceiverWalletAddress(ctx context.Context, address string) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetInboxMessagesBySenderCredentialID(ctx context.Context, senderID uuid.UUID) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, senderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetInboxMessageByID(ctx context.Context, messageID uuid.UUID) (*entity.InboxMessage, error) {
	args := m.Called(ctx, messageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) UpdateInboxMessageReadStatus(ctx context.Context, messageID uuid.UUID, isRead int) (*entity.InboxMessage, error) {
	args := m.Called(ctx, messageID, isRead)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) UpdateInboxMessageReadStatusAll(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	args := m.Called(ctx, credentialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.InboxMessage), args.Error(1)
}

func (m *MockInboxMessageDg) GetUnreadInboxMessageCountByCredentialID(ctx context.Context, params datagateway.GetInboxMessagesByCredentialIDParameters) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

// MockAuthenticationCredentialDg is a mock for AuthenticationCredentialDataGateway (reusing from inbox_test.go)
type MockAuthenticationCredentialDg struct {
	mock.Mock
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, walletAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByGoogleConnectorRef(ctx context.Context, googleConnectorRef string) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, googleConnectorRef)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddress(ctx context.Context, params datagateway.GetAuthenticationCredentialByGoogleConnectorRefOrWalletAddressParameters) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) UpdateAuthenticationCredential(ctx context.Context, id uuid.UUID, params datagateway.UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) DeleteAuthenticationCredential(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestInboxUsecase_GetMyInboxMessages(t *testing.T) {
	t.Run("should get user's inbox messages by credential ID", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		email := "user@example.com"
		walletAddress := "0x1234567890abcdef"

		user := auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		expectedMessages := []*entity.InboxMessage{
			{
				Id:                   uuid.New(),
				ReceiverCredentialId: &userID,
				MessageType:          entity.InboxMessageTypeEventRegistrationInvitation,
				MessageContent:       "You're invited to our event",
				IsRead:               0,
			},
			{
				Id:                   uuid.New(),
				ReceiverCredentialId: &userID,
				MessageType:          entity.InboxMessageTypeEventCertificateInvitation,
				MessageContent:       "Your certificate is ready",
				IsRead:               1,
			},
		}

		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, userID).Return(&entity.AuthenticationCredential{
			Id:                 userID,
			GoogleConnectorRef: &email,
			WalletAddress:      walletAddress,
		}, nil)

		mockInboxDg := new(MockInboxMessageDg)
		mockInboxDg.On("GetInboxMessagesByCredentialID", ctx, mock.Anything).Return(expectedMessages, nil)

		uc := &InboxUsecase{
			InboxMessageDg:             mockInboxDg,
			AuthenticationCredentialDg: mockAuthDg,
		}

		// Act
		messages, err := uc.GetMyInboxMessages(ctx, user)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, messages)
		assert.Len(t, messages, 2)
		assert.Equal(t, "You're invited to our event", messages[0].MessageContent)
		assert.Equal(t, "Your certificate is ready", messages[1].MessageContent)
		assert.Equal(t, 0, messages[0].IsRead) // Unread
		assert.Equal(t, 1, messages[1].IsRead) // Read
		mockInboxDg.AssertExpectations(t)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should return empty list when user has no messages", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		email := "newuser@example.com"
		walletAddress := "0xabcdef1234567890"

		user := auth.JwtClaims{
			UserId: userID,
			Email:  &email,
		}

		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, userID).Return(&entity.AuthenticationCredential{
			Id:                 userID,
			GoogleConnectorRef: &email,
			WalletAddress:      walletAddress,
		}, nil)

		mockInboxDg := new(MockInboxMessageDg)
		mockInboxDg.On("GetInboxMessagesByCredentialID", ctx, mock.Anything).Return([]*entity.InboxMessage{}, nil)

		uc := &InboxUsecase{
			InboxMessageDg:             mockInboxDg,
			AuthenticationCredentialDg: mockAuthDg,
		}

		// Act
		messages, err := uc.GetMyInboxMessages(ctx, user)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, messages)
		assert.Empty(t, messages)
		mockInboxDg.AssertExpectations(t)
		mockAuthDg.AssertExpectations(t)
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

		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, userID).Return(nil, errors.New("database connection error"))

		mockInboxDg := new(MockInboxMessageDg)

		uc := &InboxUsecase{
			InboxMessageDg:             mockInboxDg,
			AuthenticationCredentialDg: mockAuthDg,
		}

		// Act
		messages, err := uc.GetMyInboxMessages(ctx, user)

		// Assert
		require.Error(t, err)
		assert.Nil(t, messages)
		mockAuthDg.AssertExpectations(t)
	})
}

func TestInboxUsecase_GetInboxMessagesByCredentailID(t *testing.T) {
	t.Run("should get messages for specific credential", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		credentialID := uuid.New()
		email := "credential@example.com"
		walletAddress := "0xfedcba0987654321"

		messages := []*entity.InboxMessage{
			{
				Id:                   uuid.New(),
				ReceiverCredentialId: &credentialID,
				MessageType:          entity.InboxMessageTypeGeneral,
				MessageContent:       "Welcome",
				IsRead:               0,
			},
		}

		mockInboxDg := new(MockInboxMessageDg)
		mockAuthDg := new(MockAuthenticationCredentialDg)

		// Mock AuthenticationCredentialDg.GetAuthenticationCredentialById
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, credentialID).Return(&entity.AuthenticationCredential{
			Id:                 credentialID,
			GoogleConnectorRef: &email,
			WalletAddress:      walletAddress,
		}, nil)

		mockInboxDg.On("GetInboxMessagesByCredentialID", ctx, mock.MatchedBy(func(params datagateway.GetInboxMessagesByCredentialIDParameters) bool {
			return params.CredentialID == credentialID
		})).Return(messages, nil)

		uc := &InboxUsecase{
			InboxMessageDg:             mockInboxDg,
			AuthenticationCredentialDg: mockAuthDg,
		}

		// Act
		result, err := uc.GetInboxMessagesByCredentailID(ctx, credentialID)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 1)
		mockInboxDg.AssertExpectations(t)
		mockAuthDg.AssertExpectations(t)
	})
}

func TestInboxUsecase_GetInboxMessagesByReceiverEmail(t *testing.T) {
	t.Run("should get messages sent to specific email", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "attendee@example.com"

		messages := []*entity.InboxMessage{
			{
				Id:             uuid.New(),
				ReceiverEmail:  &email,
				MessageType:    entity.InboxMessageTypeEventRegistrationInvitation,
				MessageContent: "Your registration is confirmed",
				IsRead:         0,
			},
			{
				Id:             uuid.New(),
				ReceiverEmail:  &email,
				MessageType:    entity.InboxMessageTypeGeneral,
				MessageContent: "Event starts tomorrow",
				IsRead:         0,
			},
		}

		mockInboxDg := new(MockInboxMessageDg)
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
		assert.Equal(t, "Your registration is confirmed", result[0].MessageContent)
		assert.Equal(t, "Event starts tomorrow", result[1].MessageContent)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return error when email not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "nonexistent@example.com"

		mockInboxDg := new(MockInboxMessageDg)
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
				MessageType:           entity.InboxMessageTypeEventCertificateInvitation,
				MessageContent:        "Your certificate NFT is ready",
				IsRead:                0,
			},
		}

		mockInboxDg := new(MockInboxMessageDg)
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
		assert.Equal(t, "Your certificate NFT is ready", result[0].MessageContent)
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
				MessageType:        entity.InboxMessageTypeGeneral,
				MessageContent:     "Event location changed",
				IsRead:             1,
			},
			{
				Id:                 uuid.New(),
				SenderCredentialId: &senderID,
				MessageType:        entity.InboxMessageTypeGeneral,
				MessageContent:     "Unfortunately we must cancel",
				IsRead:             1,
			},
		}

		mockInboxDg := new(MockInboxMessageDg)
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
		assert.Equal(t, "Event location changed", result[0].MessageContent)
		assert.Equal(t, "Unfortunately we must cancel", result[1].MessageContent)
		mockInboxDg.AssertExpectations(t)
	})

	t.Run("should return empty when sender has sent no messages", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		senderID := uuid.New()

		mockInboxDg := new(MockInboxMessageDg)
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
