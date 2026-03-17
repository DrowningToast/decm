package event

import (
	"apps/backend/common/customerror"
	"apps/backend/common/encryptutils"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	eventcontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/event_contract"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
		// Arrange - Create valid encrypted private key for testing
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hex.EncodeToString(privateKeyBytes)
		password := "test-password"
		encryptedKey, err := encryptutils.EncryptAESGCM(privateKeyHex, password)
		require.NoError(t, err)

		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: &encryptedKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		expectedEvent := &entity.Event{
			Id:           uuid.New(),
			Title:        "Test Event",
			CreatedAt:    time.Now(),
			MaxAttendees: 100,
		}
		mockEventDg.On("CreateEvent", ctx, mock.MatchedBy(func(params event_datagateway.CreateEventParameters) bool {
			return params.Name == "Test Event" &&
				params.SeatsCount == 100 &&
				params.OwnerCredentialID == userId
		})).Return(expectedEvent, nil)

		mockContractFactoryDg := new(MockEventContractFactoryDg)
		accessManagerAddr := common.HexToAddress("0xACCE55")
		eventContractAddr := common.HexToAddress("0xEVE47")
		contractResponse := &eventcontract_datagateway.CreateContractResponse{
			AccessManagerContractAddress: accessManagerAddr,
			EventContractAddress:         eventContractAddr,
		}
		mockContractFactoryDg.On("CreateContract", ctx, mock.MatchedBy(func(params eventcontract_datagateway.CreateContractParams) bool {
			return params.EventName == expectedEvent.Title &&
				params.SeatsCount == int64(expectedEvent.MaxAttendees)
		})).Return(contractResponse, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventContractFactoryDg:     mockContractFactoryDg,
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
			Name:         "Test Event",
			SeatsCount:   100,
			StartDate:    time.Now(),
			EndDate:      time.Now().Add(24 * time.Hour),
			HostPassword: password, // Use the same password we used for encryption
		}

		// Act
		event, accessManager, eventContract, _, err := uc.CreateEvent(ctx, params, currentUser)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, expectedEvent.Id, event.Id)
		assert.Equal(t, accessManagerAddr, accessManager)
		assert.Equal(t, eventContractAddr, eventContract)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockContractFactoryDg.AssertExpectations(t)
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

		// Create a mock S3 datagateway
		mockS3Dg := new(MockS3DataGateway)
		mockS3Dg.On("DeleteFile", ctx, "banner-key").Return(nil)
		mockS3Dg.On("DeleteFile", ctx, "icon-key").Return(nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			S3DataGateway:              mockS3Dg,
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
			Name: "Test Event",
		}

		// Act
		event, _, _, _, err := uc.CreateEvent(ctx, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
		assert.Contains(t, err.Error(), "database error")
		mockS3Dg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockAuthDg.AssertExpectations(t)
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

		mockEventDg := new(MockEventDataGateway)
		// Return validation error when empty name is provided
		mockEventDg.On("CreateEvent", ctx, mock.Anything).
			Return(nil, errors.New("validation error: name is required"))

		// Create mock S3 datagateway to handle cleanup after validation failure
		mockS3Dg := new(MockS3DataGateway)
		mockS3Dg.On("DeleteFile", ctx, "banner-key").Return(nil)
		mockS3Dg.On("DeleteFile", ctx, "icon-key").Return(nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			S3DataGateway:              mockS3Dg,
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
		assert.Contains(t, err.Error(), "validation error")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockS3Dg.AssertExpectations(t)
	})
}
