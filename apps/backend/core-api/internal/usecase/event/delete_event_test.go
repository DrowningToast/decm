package event

import (
	"apps/backend/common/customerror"
	"apps/backend/common/encryptutils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteEvent(t *testing.T) {
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

		// Act
		params := DeleteEventParameters{HostPassword: password}
		event, err := uc.DeleteEvent(ctx, eventId, currentUser, params)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, event)
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

		// Act
		params := DeleteEventParameters{HostPassword: password}
		event, err := uc.DeleteEvent(ctx, eventId, currentUser, params)

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

		// Act
		params := DeleteEventParameters{HostPassword: password}
		deletedEvent, err := uc.DeleteEvent(ctx, eventId, currentUser, params)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, deletedEvent)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "user is not owner of the event")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should fail when event contract not found", func(t *testing.T) {
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
		}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		mockContractDg := new(MockEventContractDataGateway)
		mockContractDg.On("GetEventContractByEventID", ctx, eventId).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, nil))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventContractDataGateway:   mockContractDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		params := DeleteEventParameters{HostPassword: password}
		deletedEvent, err := uc.DeleteEvent(ctx, eventId, currentUser, params)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, deletedEvent)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockContractDg.AssertExpectations(t)
	})

	t.Run("should return error (not panic) when BYOK user sends only host_password", func(t *testing.T) {
		// Arrange — BYOK user has no encrypted private key
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			EncryptedPrivateKey: nil, // BYOK: no server-side key
			WalletAddress:       "0x1234567890123456789012345678901234567890",
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventId,
			OwnerCredentialId: userId,
		}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		contractAddress := "0x1234567890123456789012345678901234567890"
		mockContractDg := new(MockEventContractDataGateway)
		eventContract := &entity.EventContract{
			EventId:              eventId,
			EventContractAddress: contractAddress,
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventId).Return(eventContract, nil)

		mockBlockchainDg := new(MockBlockchainClientDataGateway)
		mockBlockchainDg.On("GetTransactOpts", ctx).Return(&bind.TransactOpts{}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventContractDataGateway:   mockContractDg,
			BlockchainClientDg:         mockBlockchainDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act — BYOK user mistakenly sends only host_password (no signature/sign_message)
		params := DeleteEventParameters{HostPassword: "some-password"}
		deletedEvent, err := uc.DeleteEvent(ctx, eventId, currentUser, params)

		// Assert — must return a proper error, not panic with a nil dereference
		assert.Error(t, err)
		assert.Nil(t, deletedEvent)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockContractDg.AssertExpectations(t)
		mockBlockchainDg.AssertExpectations(t)
	})

	t.Run("should fail with invalid password", func(t *testing.T) {
		// Arrange - Create valid encrypted private key with one password
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hex.EncodeToString(privateKeyBytes)
		correctPassword := "correct-password"
		encryptedKey, err := encryptutils.EncryptAESGCM(privateKeyHex, correctPassword)
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
		}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		contractAddress := "0x1234567890123456789012345678901234567890"
		mockContractDg := new(MockEventContractDataGateway)
		eventContract := &entity.EventContract{
			EventId:              eventId,
			EventContractAddress: contractAddress,
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventId).Return(eventContract, nil)

		mockBlockchainDg := new(MockBlockchainClientDataGateway)
		transactor := &bind.TransactOpts{}
		mockBlockchainDg.On("GetTransactOpts", ctx).Return(transactor, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventContractDataGateway:   mockContractDg,
			BlockchainClientDg:         mockBlockchainDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		wrongPassword := "wrong-password"

		// Act
		params := DeleteEventParameters{HostPassword: wrongPassword}
		deletedEvent, err := uc.DeleteEvent(ctx, eventId, currentUser, params)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, deletedEvent)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockContractDg.AssertExpectations(t)
		mockBlockchainDg.AssertExpectations(t)
	})
}
