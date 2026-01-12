package event

import (
	"apps/backend/common/customerror"
	"apps/backend/common/encryptutils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateEvent(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	eventId := uuid.New()
	password := "test-password"

	t.Run("should fail when user is not authenticated", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, nil))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		params := UpdateEventParameters{
			HostPassword: password,
		}

		// Act
		event, err := uc.UpdateEvent(ctx, eventId, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when user is not a verified organizer", func(t *testing.T) {
		// Arrange - Create valid encrypted private key
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hex.EncodeToString(privateKeyBytes)
		encryptedKey, err := encryptutils.EncryptAESGCM(privateKeyHex, password)
		require.NoError(t, err)

		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: false, // Not verified
			EncryptedPrivateKey: &encryptedKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		params := UpdateEventParameters{
			HostPassword: password,
		}

		// Act
		event, err := uc.UpdateEvent(ctx, eventId, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when event does not exist", func(t *testing.T) {
		// Arrange - Create valid encrypted private key
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hex.EncodeToString(privateKeyBytes)
		encryptedKey, err := encryptutils.EncryptAESGCM(privateKeyHex, password)
		require.NoError(t, err)

		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			EncryptedPrivateKey: &encryptedKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("GetEventById", ctx, eventId).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, nil))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		params := UpdateEventParameters{
			HostPassword: password,
		}

		// Act
		event, err := uc.UpdateEvent(ctx, eventId, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should fail when user is not the owner", func(t *testing.T) {
		// Arrange - Create valid encrypted private key
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hex.EncodeToString(privateKeyBytes)
		encryptedKey, err := encryptutils.EncryptAESGCM(privateKeyHex, password)
		require.NoError(t, err)

		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			EncryptedPrivateKey: &encryptedKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		differentOwnerId := uuid.New()
		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventId,
			OwnerCredentialId: differentOwnerId, // Different owner
		}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		params := UpdateEventParameters{
			HostPassword: password,
		}

		// Act
		updatedEvent, err := uc.UpdateEvent(ctx, eventId, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedEvent)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "user is not the owner of the event")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should fail when seats count is less than max attendees", func(t *testing.T) {
		// Arrange - Create valid encrypted private key
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hex.EncodeToString(privateKeyBytes)
		encryptedKey, err := encryptutils.EncryptAESGCM(privateKeyHex, password)
		require.NoError(t, err)

		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			EncryptedPrivateKey: &encryptedKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventId,
			OwnerCredentialId: userId,
			MaxAttendees:      100, // Current max is 100
		}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		newSeatsCount := 50 // Trying to reduce to 50, which is less than 100
		params := UpdateEventParameters{
			SeatsCount:   &newSeatsCount,
			HostPassword: password,
		}

		// Act
		updatedEvent, err := uc.UpdateEvent(ctx, eventId, params, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedEvent)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "seats count is less than the max attendees")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})
}
