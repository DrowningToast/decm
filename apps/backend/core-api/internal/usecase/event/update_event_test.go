package event

import (
	"apps/backend/common/customerror"
	"apps/backend/common/encryptutils"
	"apps/backend/core-api/internal/entity"
	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
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

		seatsCount := 100
		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventId,
			OwnerCredentialId: userId,
			MaxAttendees:      seatsCount,
		}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		contractAddress := "0x1234567890123456789012345678901234567890"
		mockContractDg := new(MockEventContractDataGateway)
		eventContract := &entity.EventContract{
			EventId:              eventId,
			EventContractAddress: contractAddress,
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventId).Return(eventContract, nil)

		// Auth fails before UpdateEvent and GetTransactOpts — do NOT register those expectations.
		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventContractDataGateway:   mockContractDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act — BYOK user mistakenly sends only host_password (no signature/sign_message)
		params := UpdateEventParameters{
			SeatsCount:   &seatsCount,
			HostPassword: "some-password",
		}
		updatedEvent, err := uc.UpdateEvent(ctx, eventId, params, currentUser)

		// Assert — must return a proper ErrUnauthorized before any DB/S3 side-effects
		assert.Error(t, err)
		assert.Nil(t, updatedEvent)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockEventDg.AssertNotCalled(t, "UpdateEvent")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockContractDg.AssertExpectations(t)
	})

	// C3: wrong password must fail BEFORE any S3 upload or DB update
	t.Run("should fail with wrong password before S3 upload or DB update", func(t *testing.T) {
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		encryptedKey, err := encryptutils.EncryptAESGCM(hex.EncodeToString(crypto.FromECDSA(privateKey)), "correct-password")
		require.NoError(t, err)

		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			EncryptedPrivateKey: &encryptedKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).Return(credential, nil)

		seatsCount := 100
		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{Id: eventId, OwnerCredentialId: userId, MaxAttendees: seatsCount}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		contractAddress := "0x1234567890123456789012345678901234567890"
		mockContractDg := new(MockEventContractDataGateway)
		mockContractDg.On("GetEventContractByEventID", ctx, eventId).Return(&entity.EventContract{
			EventId:              eventId,
			EventContractAddress: contractAddress,
		}, nil)

		// Do NOT register S3DataGateway or UpdateEvent expectations.
		// Any call to them means auth ran too late.
		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventContractDataGateway:   mockContractDg,
		}

		params := UpdateEventParameters{SeatsCount: &seatsCount, HostPassword: "wrong-password"}
		updatedEvent, err := uc.UpdateEvent(ctx, eventId, params, &auth.JwtClaims{UserId: userId})

		assert.Error(t, err)
		assert.Nil(t, updatedEvent)
		mockEventDg.AssertNotCalled(t, "UpdateEvent")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockContractDg.AssertExpectations(t)
	})

	// C3: invalid wallet signature must fail BEFORE any S3 upload or DB update
	t.Run("should fail with invalid wallet signature before S3 upload or DB update", func(t *testing.T) {
		keyA, err := crypto.GenerateKey()
		require.NoError(t, err)
		addrA := crypto.PubkeyToAddress(keyA.PublicKey)

		// Sign with a different key so recovered address != credential address
		keyB, err := crypto.GenerateKey()
		require.NoError(t, err)
		wrongHash := cyptoutils.HashEthereumMessage("some message")
		wrongSig, err := cyptoutils.Sign(wrongHash.Bytes(), keyB)
		require.NoError(t, err)
		invalidSig := hexutil.Encode(wrongSig)

		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			EncryptedPrivateKey: nil,         // BYOK
			WalletAddress:       addrA.Hex(), // credential owns keyA's address
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).Return(credential, nil)

		seatsCount := 100
		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{Id: eventId, OwnerCredentialId: userId, MaxAttendees: seatsCount}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		contractAddress := "0x1234567890123456789012345678901234567890"
		mockContractDg := new(MockEventContractDataGateway)
		mockContractDg.On("GetEventContractByEventID", ctx, eventId).Return(&entity.EventContract{
			EventId:              eventId,
			EventContractAddress: contractAddress,
		}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventContractDataGateway:   mockContractDg,
		}

		// invalidSig was signed by keyB — recovered address won't match addrA
		params := UpdateEventParameters{
			SeatsCount:  &seatsCount,
			Signature:   invalidSig,
			SignMessage: "some message",
		}
		updatedEvent, err := uc.UpdateEvent(ctx, eventId, params, &auth.JwtClaims{UserId: userId})

		assert.Error(t, err)
		assert.Nil(t, updatedEvent)
		customErr := customerror.TryParseAsCustomErr(err)
		require.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockEventDg.AssertNotCalled(t, "UpdateEvent")
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
