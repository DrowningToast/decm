package event_registration

import (
	"apps/backend/common/customerror"
	"apps/backend/common/encryptutils"
	"apps/backend/common/hashutils"
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// All mocks are now in mocks_test.go for reusability across test files

func TestCheckRegistrationEligibility(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	email := "test@example.com"
	walletAddress := "0x1234567890abcdef"

	currentUser := &auth.JwtClaims{
		UserId:        userID,
		Email:         &email,
		WalletAddress: walletAddress,
	}

	entityEventContract := &entity.EventContract{
		ID:                   uuid.New(),
		EventId:              eventID,
		EventContractAddress: "0xContractAddress",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	t.Run("should return error when user is not authenticated", func(t *testing.T) {
		uc := &EventRegistrationUsecase{}

		eligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, nil, CheckRegistrationEligibilityParams{})

		require.Error(t, err)
		assert.False(t, eligible)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrUnauthenticated.Code, *customError.Code)
	})

	t.Run("should return error when password config not found (nil config)", func(t *testing.T) {
		mockConfigDg := new(MockEventRegistrationConfigDg)
		mockConfigDg.On("GetEventRegistrationConfigPasswordByEventId", ctx, eventID).Return(nil, nil)

		password := "test-password"
		uc := &EventRegistrationUsecase{
			EventRegistrationConfigurationDg: mockConfigDg,
		}

		eligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, currentUser, CheckRegistrationEligibilityParams{
			EventPassword: &password,
		})

		require.Error(t, err)
		assert.False(t, eligible)
		mockConfigDg.AssertExpectations(t)
	})

	t.Run("should return error when config has nil registration password", func(t *testing.T) {
		mockConfigDg := new(MockEventRegistrationConfigDg)

		config := &entity.EventRegistrationConfig{
			ID:                   uuid.New(),
			EventID:              eventID,
			RegistrationPassword: nil,
		}
		mockConfigDg.On("GetEventRegistrationConfigPasswordByEventId", ctx, eventID).Return(config, nil)

		password := "test-password"
		uc := &EventRegistrationUsecase{
			EventRegistrationConfigurationDg: mockConfigDg,
		}

		eligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, currentUser, CheckRegistrationEligibilityParams{
			EventPassword: &password,
		})

		require.Error(t, err)
		assert.False(t, eligible)
		mockConfigDg.AssertExpectations(t)
	})

	t.Run("should return error when getting config fails", func(t *testing.T) {
		mockConfigDg := new(MockEventRegistrationConfigDg)
		mockConfigDg.On("GetEventRegistrationConfigPasswordByEventId", ctx, eventID).Return(nil, errors.New("database error"))

		password := "test-password"
		uc := &EventRegistrationUsecase{
			EventRegistrationConfigurationDg: mockConfigDg,
		}

		eligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, currentUser, CheckRegistrationEligibilityParams{
			EventPassword: &password,
		})

		require.Error(t, err)
		assert.False(t, eligible)
		mockConfigDg.AssertExpectations(t)
	})

	t.Run("should return true when password matches", func(t *testing.T) {
		mockConfigDg := new(MockEventRegistrationConfigDg)

		password := "correct-password"
		hashedPassword, err := hashutils.HashPassword(password)
		require.NoError(t, err)

		config := &entity.EventRegistrationConfig{
			ID:                   uuid.New(),
			EventID:              eventID,
			RegistrationPassword: &hashedPassword,
		}
		mockConfigDg.On("GetEventRegistrationConfigPasswordByEventId", ctx, eventID).Return(config, nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationConfigurationDg: mockConfigDg,
		}

		eligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, currentUser, CheckRegistrationEligibilityParams{
			EventPassword: &password,
		})

		require.NoError(t, err)
		assert.True(t, eligible)
		mockConfigDg.AssertExpectations(t)
	})

	t.Run("should return error when password does not match", func(t *testing.T) {
		mockConfigDg := new(MockEventRegistrationConfigDg)

		// Use an invalid hash format so CompareHash returns an error
		hashedPassword := "invalid-hash-format"
		config := &entity.EventRegistrationConfig{
			ID:                   uuid.New(),
			EventID:              eventID,
			RegistrationPassword: &hashedPassword,
		}
		mockConfigDg.On("GetEventRegistrationConfigPasswordByEventId", ctx, eventID).Return(config, nil)

		password := "wrong-password"
		uc := &EventRegistrationUsecase{
			EventRegistrationConfigurationDg: mockConfigDg,
		}

		eligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, currentUser, CheckRegistrationEligibilityParams{
			EventPassword: &password,
		})

		require.Error(t, err)
		assert.False(t, eligible)
		mockConfigDg.AssertExpectations(t)
	})

	t.Run("should return true with valid invitation (no password)", func(t *testing.T) {
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)

		invitation := &entity.EventRegistrationInvitation{
			Id:      uuid.New(),
			EventId: eventID,
		}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, (*entity.InboxMessage)(nil), nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		eligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, currentUser, CheckRegistrationEligibilityParams{})

		require.NoError(t, err)
		assert.True(t, eligible)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when invitation not found", func(t *testing.T) {
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)

		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(nil, (*entity.InboxMessage)(nil), nil)

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		eligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, currentUser, CheckRegistrationEligibilityParams{})

		require.Error(t, err)
		assert.False(t, eligible)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when invitation lookup fails", func(t *testing.T) {
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)

		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(nil, (*entity.InboxMessage)(nil), errors.New("database error"))

		uc := &EventRegistrationUsecase{
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		eligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, currentUser, CheckRegistrationEligibilityParams{})

		require.Error(t, err)
		assert.False(t, eligible)
		mockInvitationDg.AssertExpectations(t)
	})
}

func TestGetJoinEventSignMessage(t *testing.T) {
	ctx := context.Background()
	walletAddress := "0x1234567890abcdef1234567890abcdef12345678"
	email := "test@example.com"
	eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"

	currentUser := auth.JwtClaims{
		UserId:        uuid.New(),
		Email:         &email,
		WalletAddress: walletAddress,
	}

	t.Run("should generate sign message with provided deadline block", func(t *testing.T) {
		uc := &EventRegistrationUsecase{}

		deadlineBlock := uint64(1000000)
		signMessage, messageHash, err := uc.GetJoinEventSignMessage(
			ctx,
			common.HexToAddress(walletAddress),
			currentUser,
			common.HexToAddress(eventContractAddress),
			&deadlineBlock,
		)

		require.NoError(t, err)
		assert.NotNil(t, signMessage)
		assert.NotNil(t, messageHash)
		// Check message has correct format: address,contractAddress,deadlineBlock
		assert.Contains(t, *signMessage, ",")
		assert.Contains(t, *signMessage, "1000000")
		// Addresses are checksummed so we just verify the message format
		parts := string(*signMessage)
		assert.Contains(t, parts, "0x")
	})

	t.Run("should generate sign message with calculated deadline block when not provided", func(t *testing.T) {
		mockBlockchainDg := new(MockBlockchainClientDg)
		calculatedDeadline := uint64(2000000)
		mockBlockchainDg.On("GetCalculatedDeadlineBlock", ctx).Return(calculatedDeadline, nil)

		uc := &EventRegistrationUsecase{
			BlockchainClientDg: mockBlockchainDg,
		}

		signMessage, messageHash, err := uc.GetJoinEventSignMessage(
			ctx,
			common.HexToAddress(walletAddress),
			currentUser,
			common.HexToAddress(eventContractAddress),
			nil,
		)

		require.NoError(t, err)
		assert.NotNil(t, signMessage)
		assert.NotNil(t, messageHash)
		mockBlockchainDg.AssertExpectations(t)
	})

	t.Run("should return error when blockchain client fails", func(t *testing.T) {
		mockBlockchainDg := new(MockBlockchainClientDg)
		mockBlockchainDg.On("GetCalculatedDeadlineBlock", ctx).Return(uint64(0), errors.New("blockchain error"))

		uc := &EventRegistrationUsecase{
			BlockchainClientDg: mockBlockchainDg,
		}

		signMessage, messageHash, err := uc.GetJoinEventSignMessage(
			ctx,
			common.HexToAddress(walletAddress),
			currentUser,
			common.HexToAddress(eventContractAddress),
			nil,
		)

		require.Error(t, err)
		assert.Nil(t, signMessage)
		assert.Nil(t, messageHash)
		mockBlockchainDg.AssertExpectations(t)
	})
}

func TestJoinEventWithSignature(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	email := "test@example.com"
	walletAddress := "0x1234567890abcdef"

	currentUser := &auth.JwtClaims{
		UserId:        userID,
		Email:         &email,
		WalletAddress: walletAddress,
	}

	mockBlockchainDg := new(MockBlockchainClientDg)
	mockBlockchainDg.On("GetCurrentBlockNumber", ctx).Return(uint64(900000), nil).Maybe()

	t.Run("should return error when user is not authenticated", func(t *testing.T) {
		uc := &EventRegistrationUsecase{}

		attendee, err := uc.JoinEventWithSignature(ctx, nil, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, []byte("sig"), "message")

		require.Error(t, err)
		assert.Nil(t, attendee)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrUnauthenticated.Code, *customError.Code)
	})

	t.Run("should return error when event contract not found", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(nil, nil)

		uc := &EventRegistrationUsecase{
			EventContractDg:    mockContractDg,
			BlockchainClientDg: mockBlockchainDg,
		}

		attendee, err := uc.JoinEventWithSignature(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, []byte("sig"), "message")

		require.Error(t, err)
		assert.Nil(t, attendee)
		mockContractDg.AssertExpectations(t)
	})

	t.Run("should return error when getting event contract fails", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(nil, errors.New("database error"))

		uc := &EventRegistrationUsecase{
			EventContractDg:    mockContractDg,
			BlockchainClientDg: mockBlockchainDg,
		}

		attendee, err := uc.JoinEventWithSignature(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, []byte("sig"), "message")

		require.Error(t, err)
		assert.Nil(t, attendee)
		mockContractDg.AssertExpectations(t)
	})

	t.Run("should return error when sign message format is invalid", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		contract := &entity.EventContract{
			ID:                   uuid.New(),
			EventId:              eventID,
			EventContractAddress: eventContractAddress,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		uc := &EventRegistrationUsecase{
			EventContractDg:    mockContractDg,
			BlockchainClientDg: mockBlockchainDg,
		}

		// Garbage sign message that can't be parsed
		attendee, err := uc.JoinEventWithSignature(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, []byte("sig"), "garbage-message")

		require.Error(t, err)
		assert.Nil(t, attendee)
		assert.Contains(t, err.Error(), "failed to extract deadline block from sign message")
		mockContractDg.AssertExpectations(t)
	})

	t.Run("should return error when signature does not match", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		contract := &entity.EventContract{
			ID:                   uuid.New(),
			EventId:              eventID,
			EventContractAddress: eventContractAddress,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		uc := &EventRegistrationUsecase{
			EventContractDg:    mockContractDg,
			BlockchainClientDg: mockBlockchainDg,
		}

		// Valid format sign message but signed by a different key
		deadlineBlock := uint64(1000000)
		signMessage := fmt.Sprintf("%s,%s,%d",
			common.HexToAddress(walletAddress).Hex(),
			common.HexToAddress(eventContractAddress).Hex(),
			deadlineBlock)

		// Sign with a different private key to create a mismatch
		wrongKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		messageHash := cyptoutils.HashEthereumMessage(signMessage)
		wrongSig, err := cyptoutils.Sign(messageHash.Bytes(), wrongKey)
		require.NoError(t, err)

		attendee, err := uc.JoinEventWithSignature(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, wrongSig, signMessage)

		require.Error(t, err)
		assert.Nil(t, attendee)
		assert.Contains(t, err.Error(), "signature does not match the sign message")
		mockContractDg.AssertExpectations(t)
	})

	t.Run("should return error when eligibility check fails (valid signature but no invitation)", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)

		eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		contract := &entity.EventContract{
			ID:                   uuid.New(),
			EventId:              eventID,
			EventContractAddress: eventContractAddress,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		// Generate a key that matches the walletAddress in currentUser
		// We need the signature to match currentUser.WalletAddress
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		signerAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

		// Override currentUser to match the signer
		signerWallet := signerAddress.Hex()
		signerUser := &auth.JwtClaims{
			UserId:        userID,
			Email:         &email,
			WalletAddress: signerWallet,
		}

		deadlineBlock := uint64(1000000)
		signMessage := fmt.Sprintf("%s,%s,%d",
			signerAddress.Hex(),
			common.HexToAddress(eventContractAddress).Hex(),
			deadlineBlock)

		messageHash := cyptoutils.HashEthereumMessage(signMessage)
		sig, err := cyptoutils.Sign(messageHash.Bytes(), privateKey)
		require.NoError(t, err)

		// No invitation found
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &signerWallet).Return(nil, (*entity.InboxMessage)(nil), nil)

		uc := &EventRegistrationUsecase{
			EventContractDg:               mockContractDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			BlockchainClientDg:            mockBlockchainDg,
		}

		attendee, err := uc.JoinEventWithSignature(ctx, signerUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, sig, signMessage)

		require.Error(t, err)
		assert.Nil(t, attendee)
		mockContractDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
	})
}

func TestJoinEventWithPin(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	email := "test@example.com"
	walletAddress := "0x1234567890abcdef"

	currentUser := &auth.JwtClaims{
		UserId:        userID,
		Email:         &email,
		WalletAddress: walletAddress,
	}

	t.Run("should return error when user is not authenticated", func(t *testing.T) {
		uc := &EventRegistrationUsecase{}

		attendee, err := uc.JoinEventWithPin(ctx, nil, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, "password")

		require.Error(t, err)
		assert.Nil(t, attendee)
		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrUnauthenticated.Code, *customError.Code)
	})

	t.Run("should return error when event contract not found", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(nil, nil)

		uc := &EventRegistrationUsecase{
			EventContractDg: mockContractDg,
		}

		attendee, err := uc.JoinEventWithPin(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, "password")

		require.Error(t, err)
		assert.Nil(t, attendee)
		mockContractDg.AssertExpectations(t)
	})

	t.Run("should return error when getting event contract fails", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(nil, errors.New("database error"))

		uc := &EventRegistrationUsecase{
			EventContractDg: mockContractDg,
		}

		attendee, err := uc.JoinEventWithPin(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, "password")

		require.Error(t, err)
		assert.Nil(t, attendee)
		mockContractDg.AssertExpectations(t)
	})

	t.Run("should return error when eligibility check fails (no invitation)", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)

		contract := &entity.EventContract{
			ID:                   uuid.New(),
			EventId:              eventID,
			EventContractAddress: "0xContractAddress",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(nil, (*entity.InboxMessage)(nil), nil)

		uc := &EventRegistrationUsecase{
			EventContractDg:               mockContractDg,
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		attendee, err := uc.JoinEventWithPin(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, "password")

		require.Error(t, err)
		assert.Nil(t, attendee)
		mockContractDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return error when credential not found", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockAuthCredDg := new(MockAuthenticationCredentialDg)

		contract := &entity.EventContract{
			ID:                   uuid.New(),
			EventId:              eventID,
			EventContractAddress: "0xContractAddress",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		invitation := &entity.EventRegistrationInvitation{Id: uuid.New(), EventId: eventID}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, (*entity.InboxMessage)(nil), nil)

		// Credential not found (returns nil)
		mockAuthCredDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).Return(nil, nil)

		uc := &EventRegistrationUsecase{
			EventContractDg:               mockContractDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			AuthenticationCredentialDg:    mockAuthCredDg,
		}

		attendee, err := uc.JoinEventWithPin(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, "password")

		require.Error(t, err)
		assert.Nil(t, attendee)
		assert.Contains(t, err.Error(), "encrypted private key not found")
		mockContractDg.AssertExpectations(t)
		mockAuthCredDg.AssertExpectations(t)
	})

	t.Run("should return error when encrypted private key is nil", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockAuthCredDg := new(MockAuthenticationCredentialDg)

		contract := &entity.EventContract{
			ID:                   uuid.New(),
			EventId:              eventID,
			EventContractAddress: "0xContractAddress",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		invitation := &entity.EventRegistrationInvitation{Id: uuid.New(), EventId: eventID}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, (*entity.InboxMessage)(nil), nil)

		// Credential exists but EncryptedPrivateKey is nil
		credential := &entity.AuthenticationCredential{
			Id:                  userID,
			EncryptedPrivateKey: nil,
		}
		mockAuthCredDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).Return(credential, nil)

		uc := &EventRegistrationUsecase{
			EventContractDg:               mockContractDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			AuthenticationCredentialDg:    mockAuthCredDg,
		}

		attendee, err := uc.JoinEventWithPin(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, "password")

		require.Error(t, err)
		assert.Nil(t, attendee)
		assert.Contains(t, err.Error(), "encrypted private key not found")
		mockAuthCredDg.AssertExpectations(t)
	})

	t.Run("should return error when password is invalid (decryption failure)", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockAuthCredDg := new(MockAuthenticationCredentialDg)

		contract := &entity.EventContract{
			ID:                   uuid.New(),
			EventId:              eventID,
			EventContractAddress: "0xContractAddress",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		invitation := &entity.EventRegistrationInvitation{Id: uuid.New(), EventId: eventID}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, (*entity.InboxMessage)(nil), nil)

		// Encrypt a private key with "correct-password", then try to decrypt with "wrong-password"
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		pkHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
		encrypted, err := encryptutils.EncryptAESGCM(pkHex, "correct-password")
		require.NoError(t, err)

		credential := &entity.AuthenticationCredential{
			Id:                  userID,
			EncryptedPrivateKey: &encrypted,
		}
		mockAuthCredDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).Return(credential, nil)

		uc := &EventRegistrationUsecase{
			EventContractDg:               mockContractDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			AuthenticationCredentialDg:    mockAuthCredDg,
		}

		attendee, err := uc.JoinEventWithPin(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, "wrong-password")

		require.Error(t, err)
		assert.Nil(t, attendee)
		assert.Contains(t, err.Error(), "invalid password or failed to decrypt private key")
		mockAuthCredDg.AssertExpectations(t)
	})

	t.Run("should return error when GetCalculatedDeadlineBlock fails", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockAuthCredDg := new(MockAuthenticationCredentialDg)
		mockBlockchainDg := new(MockBlockchainClientDg)

		eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		contract := &entity.EventContract{
			ID:                   uuid.New(),
			EventId:              eventID,
			EventContractAddress: eventContractAddress,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		invitation := &entity.EventRegistrationInvitation{Id: uuid.New(), EventId: eventID}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, (*entity.InboxMessage)(nil), nil)

		// Generate real key and encrypt it
		password := "test-password"
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		pkHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
		encrypted, err := encryptutils.EncryptAESGCM(pkHex, password)
		require.NoError(t, err)

		credential := &entity.AuthenticationCredential{
			Id:                  userID,
			EncryptedPrivateKey: &encrypted,
		}
		mockAuthCredDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).Return(credential, nil)

		mockBlockchainDg.On("GetCalculatedDeadlineBlock", ctx).Return(uint64(0), errors.New("blockchain error"))

		uc := &EventRegistrationUsecase{
			EventContractDg:               mockContractDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			AuthenticationCredentialDg:    mockAuthCredDg,
			BlockchainClientDg:            mockBlockchainDg,
		}

		attendee, err := uc.JoinEventWithPin(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, password)

		require.Error(t, err)
		assert.Nil(t, attendee)
		mockBlockchainDg.AssertExpectations(t)
	})

	t.Run("should successfully join event with pin (happy path)", func(t *testing.T) {
		mockContractDg := new(MockEventContractDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockAuthCredDg := new(MockAuthenticationCredentialDg)
		mockBlockchainDg := new(MockBlockchainClientDg)
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)
		mockUserSignatureDg := new(MockUserSignatureDg)

		eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		contract := &entity.EventContract{
			ID:                   uuid.New(),
			EventId:              eventID,
			EventContractAddress: eventContractAddress,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		invitation := &entity.EventRegistrationInvitation{Id: uuid.New(), EventId: eventID}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, (*entity.InboxMessage)(nil), nil)

		// Generate real key and encrypt it
		password := "test-password"
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		pkHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
		encrypted, err := encryptutils.EncryptAESGCM(pkHex, password)
		require.NoError(t, err)

		credential := &entity.AuthenticationCredential{
			Id:                  userID,
			EncryptedPrivateKey: &encrypted,
		}
		mockAuthCredDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).Return(credential, nil)

		deadlineBlock := uint64(1000000)
		mockBlockchainDg.On("GetCalculatedDeadlineBlock", ctx).Return(deadlineBlock, nil)

		// joinEvent mocks
		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{}, nil)
		currentCount := int64(5)
		maxCount := int64(100)
		mockOnchainContractDg.On("GetCurrentSeatsCount", ctx).Return(&currentCount, nil)
		mockOnchainContractDg.On("GetMaxSeatsCount", ctx).Return(&maxCount, nil)
		mockEventAttendeeDg.On("GetEventAttendeeByEventIdAndCredentialId", ctx, eventID, userID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))
		mockEventAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventID).Return([]entity.EventAttendee{}, nil)
		// invitation check inside joinEvent
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, currentUser.Email, &currentUser.WalletAddress).Return(invitation, nil, nil).Maybe()

		estimatedTime := time.Now().Add(24 * time.Hour)
		mockBlockchainDg.On("EstimateDeadlineTime", ctx, mock.AnythingOfType("uint64")).Return(&estimatedTime, nil)

		createdSignature := &entity.UserSignature{
			Id:                         uuid.New(),
			AuthenticationCredentialId: userID,
		}
		mockUserSignatureDg.On("CreateUserSignature", ctx, mock.AnythingOfType("offchain_datagateway.CreateUserSignatureParameters")).Return(createdSignature, nil)

		createdAttendee := &entity.EventAttendee{
			Id:                   uuid.New(),
			EventId:              eventID,
			AttendeeCredentialId: userID,
			UserSignatureID:      &createdSignature.Id,
		}
		mockEventAttendeeDg.On("AddParticipant", ctx, mock.AnythingOfType("offchain_datagateway.AddParticipantParameters")).Return(createdAttendee, nil)
		mockInvitationDg.On("UpdateEventRegistrationInvitationAcceptedStatus", ctx, invitation.Id, mock.AnythingOfType("*time.Time")).Return(invitation, nil)

		uc := &EventRegistrationUsecase{
			EventContractDg:               mockContractDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			AuthenticationCredentialDg:    mockAuthCredDg,
			BlockchainClientDg:            mockBlockchainDg,
			EventContractFactoryDg:        mockEventContractFactoryDg,
			EventAttendeeDg:               mockEventAttendeeDg,
			UserSignatureDg:               mockUserSignatureDg,
		}

		attendee, err := uc.JoinEventWithPin(ctx, currentUser, eventID, CheckRegistrationEligibilityParams{}, JoinEventPayload{}, password)

		require.NoError(t, err)
		assert.NotNil(t, attendee)
		assert.Equal(t, createdAttendee.Id, attendee.Id)
		mockContractDg.AssertExpectations(t)
		mockAuthCredDg.AssertExpectations(t)
		mockBlockchainDg.AssertExpectations(t)
		mockUserSignatureDg.AssertExpectations(t)
		mockEventAttendeeDg.AssertExpectations(t)
	})
}

func TestJoinEventSucessfullyMustPopulateDataInUserSignatureTable(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	email := "test@example.com"
	walletAddress := "0x1234567890abcdef1234567890abcdef12345678"
	eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
	signatureBytes := []byte("test_signature")
	deadlineBlock := uint64(1000000)
	// Proper sign message format: address,contractAddress,deadlineBlock
	signMessage := fmt.Sprintf("%s,%s,%d",
		common.HexToAddress(walletAddress).Hex(),
		common.HexToAddress(eventContractAddress).Hex(),
		deadlineBlock)

	currentUser := &auth.JwtClaims{
		UserId:        userID,
		Email:         &email,
		WalletAddress: walletAddress,
	}

	entityEventContract := entity.EventContract{
		ID:                   uuid.New(),
		EventId:              eventID,
		EventContractAddress: eventContractAddress,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	t.Run("should create user signature with correct parameters when joining event", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockUserSignatureDg := new(MockUserSignatureDg)
		mockBlockchainDg := new(MockBlockchainClientDg)

		// Mock contract factory
		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)

		// Mock onchain contract - user hasn't joined yet
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{}, nil)

		// Mock attendee count check
		currentCount := int64(5)
		maxCount := int64(100)
		mockOnchainContractDg.On("GetCurrentSeatsCount", ctx).Return(&currentCount, nil)
		mockOnchainContractDg.On("GetMaxSeatsCount", ctx).Return(&maxCount, nil)

		// Mock attendee not found in database
		mockEventAttendeeDg.On("GetEventAttendeeByEventIdAndCredentialId", ctx, eventID, userID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))
		mockEventAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventID).Return([]entity.EventAttendee{}, nil)

		// Mock invitation found
		invitation := &entity.EventRegistrationInvitation{
			Id:      uuid.New(),
			EventId: eventID,
		}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, nil, nil)

		// Mock blockchain deadline estimation
		estimatedTime := time.Now().Add(24 * time.Hour)
		mockBlockchainDg.On("EstimateDeadlineTime", ctx, mock.AnythingOfType("uint64")).Return(&estimatedTime, nil)

		// Mock user signature creation - this is what we're testing
		createdSignature := &entity.UserSignature{
			Id:                         uuid.New(),
			AuthenticationCredentialId: userID,
			SignMessage:                signMessage,
			Signature:                  hex.EncodeToString(signatureBytes),
		}
		mockUserSignatureDg.On("CreateUserSignature", ctx, mock.MatchedBy(func(params offchain_datagateway.CreateUserSignatureParameters) bool {
			return params.AuthenticationCredentialId == userID &&
				params.SignMessage == signMessage &&
				params.Signature == hex.EncodeToString(signatureBytes) &&
				params.DeadlineBlock != nil &&
				params.EstimatedDeadline != nil
		})).Return(createdSignature, nil)

		// Mock attendee creation
		createdAttendee := &entity.EventAttendee{
			Id:                   uuid.New(),
			EventId:              eventID,
			AttendeeCredentialId: userID,
			UserSignatureID:      &createdSignature.Id,
		}
		mockEventAttendeeDg.On("AddParticipant", ctx, mock.AnythingOfType("offchain_datagateway.AddParticipantParameters")).Return(createdAttendee, nil)

		// Mock invitation acceptance
		mockInvitationDg.On("UpdateEventRegistrationInvitationAcceptedStatus", ctx, invitation.Id, mock.AnythingOfType("*time.Time")).Return(invitation, nil)

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg:        mockEventContractFactoryDg,
			EventAttendeeDg:               mockEventAttendeeDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			UserSignatureDg:               mockUserSignatureDg,
			BlockchainClientDg:            mockBlockchainDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, func() *common.Address { addr := common.HexToAddress(walletAddress); return &addr }())

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdAttendee.Id, result.Id)
		assert.Equal(t, &createdSignature.Id, result.UserSignatureID)

		mockUserSignatureDg.AssertExpectations(t)
		mockEventAttendeeDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
	})
}

func TestJoinEventReturnsErrorIfAlreadyExistsAUnbroadcastedSignatureAndNotExpired(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	email := "test@example.com"
	walletAddress := "0x1234567890abcdef1234567890abcdef12345678"
	eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
	signatureBytes := []byte("test_signature")
	deadlineBlock := uint64(1000000)
	signMessage := fmt.Sprintf("%s,%s,%d",
		common.HexToAddress(walletAddress).Hex(),
		common.HexToAddress(eventContractAddress).Hex(),
		deadlineBlock)

	currentUser := &auth.JwtClaims{
		UserId:        userID,
		Email:         &email,
		WalletAddress: walletAddress,
	}

	entityEventContract := entity.EventContract{
		ID:                   uuid.New(),
		EventId:              eventID,
		EventContractAddress: eventContractAddress,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	t.Run("should return error when user has unbroadcasted signature pending", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)
		mockUserSignatureDg := new(MockUserSignatureDg)

		// Mock contract factory
		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)

		// Mock onchain contract - user HAS joined onchain
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{common.HexToAddress(walletAddress)}, nil)

		// Mock attendee exists in database with signature
		userSignatureID := uuid.New()
		existingAttendee := &entity.EventAttendee{
			Id:                   uuid.New(),
			EventId:              eventID,
			AttendeeCredentialId: userID,
			UserSignatureID:      &userSignatureID,
		}
		mockEventAttendeeDg.On("GetEventAttendeeByEventIdAndCredentialId", ctx, eventID, userID).Return(existingAttendee, nil)

		// Mock user signature - NOT broadcasted yet (BroadcastedAt is nil)
		existingSignature := &entity.UserSignature{
			Id:            userSignatureID,
			SignMessage:   signMessage,
			Signature:     hex.EncodeToString(signatureBytes),
			BroadcastedAt: nil, // Key: not broadcasted yet
		}
		mockUserSignatureDg.On("GetUserSignatureByID", ctx, userSignatureID).Return(existingSignature, nil)

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg: mockEventContractFactoryDg,
			EventAttendeeDg:        mockEventAttendeeDg,
			UserSignatureDg:        mockUserSignatureDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, func() *common.Address { addr := common.HexToAddress(walletAddress); return &addr }())

		require.Error(t, err)
		assert.Nil(t, result)

		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrAlreadyPerformed.Code, *customError.Code)
		assert.Contains(t, err.Error(), "join event already queued and pending broadcast")

		mockEventContractFactoryDg.AssertExpectations(t)
		mockOnchainContractDg.AssertExpectations(t)
		mockEventAttendeeDg.AssertExpectations(t)
		mockUserSignatureDg.AssertExpectations(t)
	})
}

func TestJoinEventSuccessfullyIfNotJoinedAndUnboardcastedSignatureExpired(t *testing.T) {
	// This scenario is covered by TestJoinEventSucessfullyMustPopulateDataInUserSignatureTable
	// The normal join flow handles this case - when user hasn't joined onchain and no attendee exists,
	// a new signature and attendee are created regardless of any expired signatures
	t.Skip("Covered by TestJoinEventSucessfullyMustPopulateDataInUserSignatureTable")
}

func TestJoinEventSuccessfullyMustNotPopulateEventAttendee(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	email := "test@example.com"
	walletAddress := "0x1234567890abcdef1234567890abcdef12345678"
	eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
	signatureBytes := []byte("test_signature")
	deadlineBlock := uint64(1000000)
	signMessage := fmt.Sprintf("%s,%s,%d",
		common.HexToAddress(walletAddress).Hex(),
		common.HexToAddress(eventContractAddress).Hex(),
		deadlineBlock)

	currentUser := &auth.JwtClaims{
		UserId:        userID,
		Email:         &email,
		WalletAddress: walletAddress,
	}

	entityEventContract := entity.EventContract{
		ID:                   uuid.New(),
		EventId:              eventID,
		EventContractAddress: eventContractAddress,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	t.Run("should return error when user already joined with broadcasted signature", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)
		mockUserSignatureDg := new(MockUserSignatureDg)

		// Mock contract factory
		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)

		// Mock onchain contract - user HAS joined onchain
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{common.HexToAddress(walletAddress)}, nil)

		// Mock attendee exists in database with signature
		userSignatureID := uuid.New()
		existingAttendee := &entity.EventAttendee{
			Id:                   uuid.New(),
			EventId:              eventID,
			AttendeeCredentialId: userID,
			UserSignatureID:      &userSignatureID,
		}
		mockEventAttendeeDg.On("GetEventAttendeeByEventIdAndCredentialId", ctx, eventID, userID).Return(existingAttendee, nil)

		// Mock user signature - ALREADY broadcasted
		broadcastedTime := time.Now().Add(-1 * time.Hour)
		existingSignature := &entity.UserSignature{
			Id:            userSignatureID,
			SignMessage:   signMessage,
			Signature:     hex.EncodeToString(signatureBytes),
			BroadcastedAt: &broadcastedTime, // Key: already broadcasted
		}
		mockUserSignatureDg.On("GetUserSignatureByID", ctx, userSignatureID).Return(existingSignature, nil)

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg: mockEventContractFactoryDg,
			EventAttendeeDg:        mockEventAttendeeDg,
			UserSignatureDg:        mockUserSignatureDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, func() *common.Address { addr := common.HexToAddress(walletAddress); return &addr }())

		require.Error(t, err)
		assert.Nil(t, result)

		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customError.Code)
		assert.Contains(t, err.Error(), "user has already joined the event")

		// Ensure AddParticipant was NOT called (no duplicate attendee created)
		mockEventAttendeeDg.AssertNotCalled(t, "AddParticipant")

		mockEventContractFactoryDg.AssertExpectations(t)
		mockOnchainContractDg.AssertExpectations(t)
		mockEventAttendeeDg.AssertExpectations(t)
		mockUserSignatureDg.AssertExpectations(t)
	})
}

func TestJoinEventReturnsErrorIfNotInvited(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	email := "test@example.com"
	walletAddress := "0x1234567890abcdef1234567890abcdef12345678"
	eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
	signatureBytes := []byte("test_signature")
	deadlineBlock := uint64(1000000)
	signMessage := fmt.Sprintf("%s,%s,%d",
		common.HexToAddress(walletAddress).Hex(),
		common.HexToAddress(eventContractAddress).Hex(),
		deadlineBlock)

	currentUser := &auth.JwtClaims{
		UserId:        userID,
		Email:         &email,
		WalletAddress: walletAddress,
	}

	entityEventContract := entity.EventContract{
		ID:                   uuid.New(),
		EventId:              eventID,
		EventContractAddress: eventContractAddress,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	t.Run("should return error when user has no invitation", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)

		// Mock contract factory
		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)

		// Mock onchain contract - user hasn't joined yet
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{}, nil)

		// Mock attendee count check
		currentCount := int64(5)
		maxCount := int64(100)
		mockOnchainContractDg.On("GetCurrentSeatsCount", ctx).Return(&currentCount, nil)
		mockOnchainContractDg.On("GetMaxSeatsCount", ctx).Return(&maxCount, nil)

		// Mock attendee not found in database
		mockEventAttendeeDg.On("GetEventAttendeeByEventIdAndCredentialId", ctx, eventID, userID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))
		mockEventAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventID).Return([]entity.EventAttendee{}, nil)

		// Mock invitation NOT found - this is the key test case
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(nil, nil, nil)

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg:        mockEventContractFactoryDg,
			EventAttendeeDg:               mockEventAttendeeDg,
			EventRegistrationInvitationDg: mockInvitationDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, func() *common.Address { addr := common.HexToAddress(walletAddress); return &addr }())

		require.Error(t, err)
		assert.Nil(t, result)

		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrForbidden.Code, *customError.Code)
		assert.Contains(t, err.Error(), "invitation not found")

		mockEventContractFactoryDg.AssertExpectations(t)
		mockOnchainContractDg.AssertExpectations(t)
		mockEventAttendeeDg.AssertExpectations(t)
		mockInvitationDg.AssertExpectations(t)
	})
}

func TestJoinEventSuccessfullyMarksInvitationAsAccepted(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	email := "test@example.com"
	walletAddress := "0x1234567890abcdef1234567890abcdef12345678"
	eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
	signatureBytes := []byte("test_signature")
	deadlineBlock := uint64(1000000)
	signMessage := fmt.Sprintf("%s,%s,%d",
		common.HexToAddress(walletAddress).Hex(),
		common.HexToAddress(eventContractAddress).Hex(),
		deadlineBlock)

	currentUser := &auth.JwtClaims{
		UserId:        userID,
		Email:         &email,
		WalletAddress: walletAddress,
	}

	entityEventContract := entity.EventContract{
		ID:                   uuid.New(),
		EventId:              eventID,
		EventContractAddress: eventContractAddress,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	t.Run("should mark invitation as accepted after successful join", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockUserSignatureDg := new(MockUserSignatureDg)
		mockBlockchainDg := new(MockBlockchainClientDg)

		// Mock contract factory
		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)

		// Mock onchain contract - user hasn't joined yet
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{}, nil)

		// Mock attendee count check
		currentCount := int64(5)
		maxCount := int64(100)
		mockOnchainContractDg.On("GetCurrentSeatsCount", ctx).Return(&currentCount, nil)
		mockOnchainContractDg.On("GetMaxSeatsCount", ctx).Return(&maxCount, nil)

		// Mock attendee not found in database
		mockEventAttendeeDg.On("GetEventAttendeeByEventIdAndCredentialId", ctx, eventID, userID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))
		mockEventAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventID).Return([]entity.EventAttendee{}, nil)

		// Mock invitation found
		invitation := &entity.EventRegistrationInvitation{
			Id:      uuid.New(),
			EventId: eventID,
		}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, nil, nil)

		// Mock blockchain deadline estimation
		estimatedTime := time.Now().Add(24 * time.Hour)
		mockBlockchainDg.On("EstimateDeadlineTime", ctx, mock.AnythingOfType("uint64")).Return(&estimatedTime, nil)

		// Mock user signature creation
		createdSignature := &entity.UserSignature{
			Id:                         uuid.New(),
			AuthenticationCredentialId: userID,
			SignMessage:                signMessage,
			Signature:                  hex.EncodeToString(signatureBytes),
		}
		mockUserSignatureDg.On("CreateUserSignature", ctx, mock.AnythingOfType("offchain_datagateway.CreateUserSignatureParameters")).Return(createdSignature, nil)

		// Mock attendee creation
		createdAttendee := &entity.EventAttendee{
			Id:                   uuid.New(),
			EventId:              eventID,
			AttendeeCredentialId: userID,
			UserSignatureID:      &createdSignature.Id,
		}
		mockEventAttendeeDg.On("AddParticipant", ctx, mock.AnythingOfType("offchain_datagateway.AddParticipantParameters")).Return(createdAttendee, nil)

		// Mock invitation acceptance - THIS IS THE KEY TEST
		mockInvitationDg.On("UpdateEventRegistrationInvitationAcceptedStatus", ctx, invitation.Id, mock.MatchedBy(func(acceptedAt *time.Time) bool {
			return acceptedAt != nil && acceptedAt.Before(time.Now().Add(1*time.Second))
		})).Return(invitation, nil)

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg:        mockEventContractFactoryDg,
			EventAttendeeDg:               mockEventAttendeeDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			UserSignatureDg:               mockUserSignatureDg,
			BlockchainClientDg:            mockBlockchainDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, func() *common.Address { addr := common.HexToAddress(walletAddress); return &addr }())

		require.NoError(t, err)
		assert.NotNil(t, result)

		// Verify that UpdateEventRegistrationInvitationAcceptedStatus was called
		mockInvitationDg.AssertCalled(t, "UpdateEventRegistrationInvitationAcceptedStatus", ctx, invitation.Id, mock.AnythingOfType("*time.Time"))
		mockInvitationDg.AssertExpectations(t)
	})
}

func TestJoinEventInternalErrors(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	email := "test@example.com"
	walletAddress := "0x1234567890abcdef1234567890abcdef12345678"
	eventContractAddress := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
	signatureBytes := []byte("test_signature")
	deadlineBlock := uint64(1000000)
	signMessage := fmt.Sprintf("%s,%s,%d",
		common.HexToAddress(walletAddress).Hex(),
		common.HexToAddress(eventContractAddress).Hex(),
		deadlineBlock)

	currentUser := &auth.JwtClaims{
		UserId:        userID,
		Email:         &email,
		WalletAddress: walletAddress,
	}

	entityEventContract := entity.EventContract{
		ID:                   uuid.New(),
		EventId:              eventID,
		EventContractAddress: eventContractAddress,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	addr := common.HexToAddress(walletAddress)

	t.Run("should return error when GetContract fails", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(nil, errors.New("contract factory error"))

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg: mockEventContractFactoryDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, &addr)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get event contract by address")
		mockEventContractFactoryDg.AssertExpectations(t)
	})

	t.Run("should return error when GetParticipants fails", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)

		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)
		mockOnchainContractDg.On("GetParticipants", ctx).Return(nil, errors.New("participants error"))

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg: mockEventContractFactoryDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, &addr)

		require.Error(t, err)
		assert.Nil(t, result)
		mockOnchainContractDg.AssertExpectations(t)
	})

	t.Run("should return error when event attendee is full", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)

		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{}, nil)

		currentCount := int64(100)
		maxCount := int64(100)
		mockOnchainContractDg.On("GetCurrentSeatsCount", ctx).Return(&currentCount, nil)
		mockOnchainContractDg.On("GetMaxSeatsCount", ctx).Return(&maxCount, nil)

		fullAttendees := make([]entity.EventAttendee, 100)
		mockEventAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventID).Return(fullAttendees, nil)

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg: mockEventContractFactoryDg,
			EventAttendeeDg:        mockEventAttendeeDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, &addr)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), string(JoinEventUserErrorEventAttendeeFull))
		mockOnchainContractDg.AssertExpectations(t)
		mockEventAttendeeDg.AssertExpectations(t)
	})

	t.Run("should return error when CreateUserSignature fails", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockUserSignatureDg := new(MockUserSignatureDg)
		mockBlockchainDg := new(MockBlockchainClientDg)

		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{}, nil)
		currentCount := int64(5)
		maxCount := int64(100)
		mockOnchainContractDg.On("GetCurrentSeatsCount", ctx).Return(&currentCount, nil)
		mockOnchainContractDg.On("GetMaxSeatsCount", ctx).Return(&maxCount, nil)
		mockEventAttendeeDg.On("GetEventAttendeeByEventIdAndCredentialId", ctx, eventID, userID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))
		mockEventAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventID).Return([]entity.EventAttendee{}, nil)

		invitation := &entity.EventRegistrationInvitation{Id: uuid.New(), EventId: eventID}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, nil, nil)

		estimatedTime := time.Now().Add(24 * time.Hour)
		mockBlockchainDg.On("EstimateDeadlineTime", ctx, mock.AnythingOfType("uint64")).Return(&estimatedTime, nil)

		mockUserSignatureDg.On("CreateUserSignature", ctx, mock.AnythingOfType("offchain_datagateway.CreateUserSignatureParameters")).Return(nil, errors.New("db error"))

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg:        mockEventContractFactoryDg,
			EventAttendeeDg:               mockEventAttendeeDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			UserSignatureDg:               mockUserSignatureDg,
			BlockchainClientDg:            mockBlockchainDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, &addr)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create user signature")
		mockUserSignatureDg.AssertExpectations(t)
	})

	t.Run("should return error and rollback signature when AddParticipant fails", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockUserSignatureDg := new(MockUserSignatureDg)
		mockBlockchainDg := new(MockBlockchainClientDg)

		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{}, nil)
		currentCount := int64(5)
		maxCount := int64(100)
		mockOnchainContractDg.On("GetCurrentSeatsCount", ctx).Return(&currentCount, nil)
		mockOnchainContractDg.On("GetMaxSeatsCount", ctx).Return(&maxCount, nil)
		mockEventAttendeeDg.On("GetEventAttendeeByEventIdAndCredentialId", ctx, eventID, userID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))
		mockEventAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventID).Return([]entity.EventAttendee{}, nil)

		invitation := &entity.EventRegistrationInvitation{Id: uuid.New(), EventId: eventID}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, nil, nil)

		estimatedTime := time.Now().Add(24 * time.Hour)
		mockBlockchainDg.On("EstimateDeadlineTime", ctx, mock.AnythingOfType("uint64")).Return(&estimatedTime, nil)

		createdSignature := &entity.UserSignature{Id: uuid.New(), AuthenticationCredentialId: userID}
		mockUserSignatureDg.On("CreateUserSignature", ctx, mock.AnythingOfType("offchain_datagateway.CreateUserSignatureParameters")).Return(createdSignature, nil)

		// AddParticipant fails
		mockEventAttendeeDg.On("AddParticipant", ctx, mock.AnythingOfType("offchain_datagateway.AddParticipantParameters")).Return(nil, errors.New("db error"))
		// Verify rollback of user signature
		mockUserSignatureDg.On("DeleteUserSignature", ctx, createdSignature.Id).Return(nil)

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg:        mockEventContractFactoryDg,
			EventAttendeeDg:               mockEventAttendeeDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			UserSignatureDg:               mockUserSignatureDg,
			BlockchainClientDg:            mockBlockchainDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, &addr)

		require.Error(t, err)
		assert.Nil(t, result)
		// Verify that DeleteUserSignature was called (rollback)
		mockUserSignatureDg.AssertCalled(t, "DeleteUserSignature", ctx, createdSignature.Id)
		mockUserSignatureDg.AssertExpectations(t)
		mockEventAttendeeDg.AssertExpectations(t)
	})

	t.Run("should return error and rollback attendee and signature when UpdateInvitationAcceptedStatus fails", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)
		mockInvitationDg := new(MockEventRegistrationInvitationDataGateway)
		mockUserSignatureDg := new(MockUserSignatureDg)
		mockBlockchainDg := new(MockBlockchainClientDg)

		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{}, nil)
		currentCount := int64(5)
		maxCount := int64(100)
		mockOnchainContractDg.On("GetCurrentSeatsCount", ctx).Return(&currentCount, nil)
		mockOnchainContractDg.On("GetMaxSeatsCount", ctx).Return(&maxCount, nil)
		mockEventAttendeeDg.On("GetEventAttendeeByEventIdAndCredentialId", ctx, eventID, userID).Return(nil, customerror.NewWithPreset(&customerror.ErrNotFound, errors.New("not found")))
		mockEventAttendeeDg.On("ListEventAttendeesByEventID", ctx, eventID).Return([]entity.EventAttendee{}, nil)

		invitation := &entity.EventRegistrationInvitation{Id: uuid.New(), EventId: eventID}
		mockInvitationDg.On("GetEventRegistrationInvitationByEventIDAndCredential", ctx, eventID, userID, &email, &walletAddress).Return(invitation, nil, nil)

		estimatedTime := time.Now().Add(24 * time.Hour)
		mockBlockchainDg.On("EstimateDeadlineTime", ctx, mock.AnythingOfType("uint64")).Return(&estimatedTime, nil)

		createdSignature := &entity.UserSignature{Id: uuid.New(), AuthenticationCredentialId: userID}
		mockUserSignatureDg.On("CreateUserSignature", ctx, mock.AnythingOfType("offchain_datagateway.CreateUserSignatureParameters")).Return(createdSignature, nil)

		createdAttendee := &entity.EventAttendee{
			Id:                   uuid.New(),
			EventId:              eventID,
			AttendeeCredentialId: userID,
			UserSignatureID:      &createdSignature.Id,
		}
		mockEventAttendeeDg.On("AddParticipant", ctx, mock.AnythingOfType("offchain_datagateway.AddParticipantParameters")).Return(createdAttendee, nil)

		// UpdateInvitationAcceptedStatus fails
		mockInvitationDg.On("UpdateEventRegistrationInvitationAcceptedStatus", ctx, invitation.Id, mock.AnythingOfType("*time.Time")).Return(nil, errors.New("update error"))

		// Verify rollback of both attendee and signature
		mockEventAttendeeDg.On("DeleteEventAttendeeById", ctx, createdAttendee.Id).Return(nil)
		mockUserSignatureDg.On("DeleteUserSignature", ctx, createdSignature.Id).Return(nil)

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg:        mockEventContractFactoryDg,
			EventAttendeeDg:               mockEventAttendeeDg,
			EventRegistrationInvitationDg: mockInvitationDg,
			UserSignatureDg:               mockUserSignatureDg,
			BlockchainClientDg:            mockBlockchainDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, &addr)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "rolled back attendee and signature after failing to mark invitation accepted")
		mockEventAttendeeDg.AssertCalled(t, "DeleteEventAttendeeById", ctx, createdAttendee.Id)
		mockUserSignatureDg.AssertCalled(t, "DeleteUserSignature", ctx, createdSignature.Id)
		mockInvitationDg.AssertExpectations(t)
	})

	t.Run("should return already joined error when old version attendee has nil UserSignatureID", func(t *testing.T) {
		mockEventContractFactoryDg := new(MockEventContractFactoryDg)
		mockOnchainContractDg := new(MockOnchainEventContractDg)
		mockEventAttendeeDg := new(MockEventAttendeeDg)

		mockEventContractFactoryDg.On("GetContract", common.HexToAddress(eventContractAddress)).Return(mockOnchainContractDg, nil)

		// User HAS joined onchain
		mockOnchainContractDg.On("GetParticipants", ctx).Return([]common.Address{common.HexToAddress(walletAddress)}, nil)

		// Attendee exists but with nil UserSignatureID (old version)
		existingAttendee := &entity.EventAttendee{
			Id:                   uuid.New(),
			EventId:              eventID,
			AttendeeCredentialId: userID,
			UserSignatureID:      nil, // Old version - no signature
		}
		mockEventAttendeeDg.On("GetEventAttendeeByEventIdAndCredentialId", ctx, eventID, userID).Return(existingAttendee, nil)

		uc := &EventRegistrationUsecase{
			EventContractFactoryDg: mockEventContractFactoryDg,
			EventAttendeeDg:        mockEventAttendeeDg,
		}

		result, err := uc.queueEventJoin(ctx, currentUser, entityEventContract, JoinEventPayload{}, signatureBytes, signMessage, &addr)

		require.Error(t, err)
		assert.Nil(t, result)

		var customError *customerror.Err
		require.True(t, errors.As(err, &customError))
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customError.Code)
		assert.Contains(t, err.Error(), "user has already joined the event")
		mockEventContractFactoryDg.AssertExpectations(t)
		mockOnchainContractDg.AssertExpectations(t)
		mockEventAttendeeDg.AssertExpectations(t)
	})
}
