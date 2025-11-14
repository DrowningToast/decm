package event

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"apps/backend/common/customerror"
	"apps/backend/core-api/config"
	"apps/backend/core-api/config/blockchain"
	authDg "apps/backend/core-api/internal/datagateway"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"apps/backend/services/s3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementations for testing
type MockEventDataGateway struct {
	mock.Mock
}

func (m *MockEventDataGateway) CreateEvent(ctx context.Context, params datagateway.CreateEventParameters) (*entity.Event, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDataGateway) GetEventById(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDataGateway) GetViewModelById(ctx context.Context, id uuid.UUID) (*entity.Event, *entity.EventRegistrationConfig, *entity.EventContract, error) {
	return nil, nil, nil, errors.New("not implemented")
}

func (m *MockEventDataGateway) ListEventsByOwnerCredentialID(ctx context.Context, ownerCredentialID uuid.UUID, limitCount int32, offsetCount int32) ([]*entity.Event, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventDataGateway) UpdateEvent(ctx context.Context, id uuid.UUID, params datagateway.UpdateEventParameters) (*entity.Event, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventDataGateway) DeleteEvent(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventDataGateway) ListEvents(ctx context.Context, limitCount *int32, offsetCount *int32) ([]*entity.Event, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventDataGateway) ListEventsByEventAttendeeCredentialID(ctx context.Context, eventAttendeeCredentialID uuid.UUID, limitCount *int32, offsetCount *int32) ([]*entity.Event, error) {
	return nil, errors.New("not implemented")
}

type MockAuthenticationCredentialDg struct {
	mock.Mock
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx context.Context, credentialId uuid.UUID) (*entity.AuthenticationCredential, error) {
	args := m.Called(ctx, credentialId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthenticationCredential), args.Error(1)
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialById(ctx context.Context, id uuid.UUID) (*entity.AuthenticationCredential, error) {
	return nil, errors.New("not implemented")
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByWalletAddress(ctx context.Context, walletAddress string) (*entity.AuthenticationCredential, error) {
	return nil, errors.New("not implemented")
}

func (m *MockAuthenticationCredentialDg) GetAuthenticationCredentialByGoogleConnectorRef(ctx context.Context, googleConnectorRef string) (*entity.AuthenticationCredential, error) {
	return nil, errors.New("not implemented")
}

func (m *MockAuthenticationCredentialDg) CreateAuthenticationCredential(ctx context.Context, credential entity.AuthenticationCredential) (*entity.AuthenticationCredential, error) {
	return nil, errors.New("not implemented")
}

func (m *MockAuthenticationCredentialDg) UpdateAuthenticationCredential(ctx context.Context, id uuid.UUID, params authDg.UpdateAuthenticationCredentialParameters) (*entity.AuthenticationCredential, error) {
	return nil, errors.New("not implemented")
}

func (m *MockAuthenticationCredentialDg) DeleteAuthenticationCredential(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented")
}

type MockS3Service struct {
	mock.Mock
}

func (m *MockS3Service) UploadFile(ctx context.Context, key string, file *multipart.FileHeader) (string, error) {
	args := m.Called(ctx, key, file)
	return args.String(0), args.Error(1)
}

func (m *MockS3Service) DeleteFile(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockS3Service) PutFile(ctx context.Context, requestObject interface{}) (string, error) {
	return "", errors.New("not implemented")
}

func (m *MockS3Service) GetS3UploadRequestObject(storageKeyType interface{}, entityID uuid.UUID, file *multipart.FileHeader) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// Helper function to create a mock config for testing
func createMockConfig() *config.Config {
	return &config.Config{
		Blockchain: blockchain.BlockchainConfig{
			ChainID:                    1,
			DecmAccessManagerAddress:   "0x1234567890123456789012345678901234567890",
		},
	}
}

func TestCreateEvent(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	_ = common.HexToAddress("0x1234567890123456789012345678901234567890") // unused in most tests but kept for future use

	t.Run("should fail when user is not authenticated", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(nil, errors.New("not found"))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfig(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		params := CreateEventParameters{
			Name: "Test Event",
		}

		// Act
		event, _, _, _, err := uc.CreateEvent(ctx, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when user is not a verified organizer", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: false,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: nil,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfig(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		params := CreateEventParameters{
			Name: "Test Event",
		}

		// Act
		event, _, _, _, err := uc.CreateEvent(ctx, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
		// Verify it's an unauthorized error
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when banner upload fails", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: nil,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfig(),
		}
		uc.UploadEventBanner = func(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error) {
			return "", errors.New("upload failed")
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		params := CreateEventParameters{
			Name: "Test Event",
		}

		// Act
		event, _, _, _, err := uc.CreateEvent(ctx, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when icon upload fails", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: nil,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfig(),
		}
		uc.UploadEventBanner = func(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error) {
			return "banner-key", nil
		}
		uc.UploadEventIcon = func(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error) {
			return "", errors.New("upload failed")
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		params := CreateEventParameters{
			Name: "Test Event",
		}

		// Act
		event, _, _, _, err := uc.CreateEvent(ctx, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should create event with valid parameters", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: nil,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		expectedEvent := &entity.Event{
			ID:           uuid.New(),
			Title:        "Test Event",
			CreatedAt:    time.Now(),
			MaxAttendees: 100,
		}
		mockEventDg.On("CreateEvent", ctx, mock.MatchedBy(func(params datagateway.CreateEventParameters) bool {
			return params.Name == "Test Event" &&
				params.SeatsCount == 100 &&
				params.OwnerCredentialID == userId
		})).Return(expectedEvent, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			cfg:                        createMockConfig(),
		}
		uc.UploadEventBanner = func(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error) {
			return "banner-key", nil
		}
		uc.UploadEventIcon = func(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error) {
			return "icon-key", nil
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		params := CreateEventParameters{
			Name:       "Test Event",
			SeatsCount: 100,
			StartDate:  time.Now(),
			EndDate:    time.Now().Add(24 * time.Hour),
		}

		// Act
		event, _, _, _, err := uc.CreateEvent(ctx, params, currentUser)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, expectedEvent.ID, event.ID)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should cleanup files when database creation fails", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: nil,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("CreateEvent", ctx, mock.Anything).
			Return(nil, errors.New("database error"))

		// Create a mock S3 service
		mockS3 := new(MockS3Service)
		mockS3.On("DeleteFile", ctx, "banner-key").Return(nil)
		mockS3.On("DeleteFile", ctx, "icon-key").Return(nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			cfg:                        createMockConfig(),
		}
		// Set S3Service field via interface{} to avoid type assertion issues
		uc.S3Service = (*s3.S3Service)(nil) // Will be replaced by UploadEventBanner functions
		uc.UploadEventBanner = func(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error) {
			return "banner-key", nil
		}
		uc.UploadEventIcon = func(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error) {
			return "icon-key", nil
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		params := CreateEventParameters{
			Name: "Test Event",
		}

		// Act
		event, _, _, _, err := uc.CreateEvent(ctx, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
		mockS3.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should validate required parameters", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: nil,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfig(),
		}
		uc.UploadEventBanner = func(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error) {
			return "banner-key", nil
		}
		uc.UploadEventIcon = func(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error) {
			return "icon-key", nil
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Test with empty name
		params := CreateEventParameters{
			Name: "",
		}

		// Act
		event, _, _, _, err := uc.CreateEvent(ctx, params, currentUser)

		// Assert
		// The specific validation error depends on implementation
		assert.Error(t, err)
		assert.Nil(t, event)
		mockAuthDg.AssertExpectations(t)
	})
}
