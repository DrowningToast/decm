package event

import (
	"apps/backend/common/customerror"
	"apps/backend/common/encryptutils"
	"apps/backend/core-api/internal/entity"
	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestImportCertificateReceivers(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	eventID := uuid.New()
	hostPin := "test-pin"

	// Generate a real encrypted private key for tests that reach the early PIN decryption step.
	rawKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	realEncryptedKey, err := encryptutils.EncryptAESGCM(hex.EncodeToString(crypto.FromECDSA(rawKey)), hostPin)
	require.NoError(t, err)

	t.Run("should fail when receiver has neither email nor wallet address", func(t *testing.T) {
		uc := &EventUsecase{
			cfg: createMockConfig(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{FirstName: strPtr("John")}, // Missing both email and wallet
		}

		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "must have exactly one of email or wallet_address")
	})

	t.Run("should fail when receiver has both email and wallet address", func(t *testing.T) {
		uc := &EventUsecase{
			cfg: createMockConfig(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("a@b.com"), WalletAddress: strPtr("0x123")}, // Has both
		}

		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "must have exactly one of email or wallet_address")
	})

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
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("test@example.com")},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
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
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("test@example.com")},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when event does not exist", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: strPtr(realEncryptedKey),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(nil, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			cfg:                        createMockConfig(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("test@example.com")},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrNotFound.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should fail when event contract does not exist", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: strPtr(realEncryptedKey),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userId,
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		mockEventContractDg.On("GetEventContractByEventID", ctx, eventID).
			Return(nil, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventContractDataGateway:   mockEventContractDg,
			cfg:                        createMockConfig(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("test@example.com")},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrNotFound.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockEventContractDg.AssertExpectations(t)
	})

	t.Run("should delete existing certificates and signatures before importing new ones", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: strPtr(realEncryptedKey),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userId,
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		certificateAddress := "0xCertificateContractAddress"
		eventContract := &entity.EventContract{
			EventId:                    eventID,
			EventContractAddress:       "0xEventContractAddress",
			CertificateContractAddress: &certificateAddress,
		}
		mockEventContractDg.On("GetEventContractByEventID", ctx, eventID).
			Return(eventContract, nil)

		mockEventIssuerDg := new(MockEventIssuerDataGateway)
		mockEventIssuerDg.On("ResetAllEventIssuersSigningStatus", ctx, eventID).
			Return(nil)

		// Setup existing certificates
		oldCertID1 := uuid.New()
		oldCertID2 := uuid.New()
		oldCertificates := []*entity.EventCertificate{
			{Id: oldCertID1, EventId: eventID},
			{Id: oldCertID2, EventId: eventID},
		}

		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return(oldCertificates, nil)

		// Setup certificate config
		configID := uuid.New()
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		// Setup existing signatures (now linked to config, not individual certificates)
		oldSigID1 := uuid.New()
		oldSigID2 := uuid.New()
		oldSignatures := []*entity.EventCertificateSignature{
			{Id: oldSigID1, EventCertificateConfigId: configID},
			{Id: oldSigID2, EventCertificateConfigId: configID},
		}

		mockCertSigDg := new(MockEventCertificateSignatureDataGateway)
		mockCertSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return(oldSignatures, nil)
		mockCertSigDg.On("DeleteEventCertificateSignature", ctx, oldSigID1).
			Return(nil)
		mockCertSigDg.On("DeleteEventCertificateSignature", ctx, oldSigID2).
			Return(nil)

		mockCertDg.On("DeleteEventCertificate", ctx, oldCertID1).
			Return(nil)
		mockCertDg.On("DeleteEventCertificate", ctx, oldCertID2).
			Return(nil)
		// Stop execution after deletion — we're only testing the deletion logic here
		mockCertDg.On("CreateEventCertificate", ctx, mock.Anything).
			Return(nil, errors.New("stop here"))

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventDataGateway:                     mockEventDg,
			EventContractDataGateway:             mockEventContractDg,
			EventIssuerDataGateway:               mockEventIssuerDg,
			EventCertificateDataGateway:          mockCertDg,
			EventCertificateSignatureDataGateway: mockCertSigDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			cfg:                                  createMockConfig(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("new@example.com")},
		}

		// Act
		// Note: This will fail at blockchain operations, but we're testing the deletion logic
		// The deletion should happen before blockchain operations
		_, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert
		// We expect an error due to blockchain operations, but deletion should have been called
		assert.Error(t, err)
		mockCertDg.AssertExpectations(t)
		mockCertSigDg.AssertExpectations(t)
	})

	t.Run("should handle error when getting existing certificates", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: strPtr(realEncryptedKey),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userId,
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		certificateAddress := "0xCertificateContractAddress"
		eventContract := &entity.EventContract{
			EventId:                    eventID,
			EventContractAddress:       "0xEventContractAddress",
			CertificateContractAddress: &certificateAddress,
		}
		mockEventContractDg.On("GetEventContractByEventID", ctx, eventID).
			Return(eventContract, nil)

		mockEventIssuerDg := new(MockEventIssuerDataGateway)
		mockEventIssuerDg.On("ResetAllEventIssuersSigningStatus", ctx, eventID).
			Return(nil)

		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return(nil, errors.New("database error"))

		uc := &EventUsecase{
			AuthenticationCredentialDg:  mockAuthDg,
			EventDataGateway:            mockEventDg,
			EventContractDataGateway:    mockEventContractDg,
			EventIssuerDataGateway:      mockEventIssuerDg,
			EventCertificateDataGateway: mockCertDg,
			cfg:                         createMockConfig(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("test@example.com")},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should handle error when deleting certificate signatures", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: strPtr(realEncryptedKey),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userId,
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		certificateAddress := "0xCertificateContractAddress"
		eventContract := &entity.EventContract{
			EventId:                    eventID,
			EventContractAddress:       "0xEventContractAddress",
			CertificateContractAddress: &certificateAddress,
		}
		mockEventContractDg.On("GetEventContractByEventID", ctx, eventID).
			Return(eventContract, nil)

		mockEventIssuerDg := new(MockEventIssuerDataGateway)
		mockEventIssuerDg.On("ResetAllEventIssuersSigningStatus", ctx, eventID).
			Return(nil)

		oldCertID := uuid.New()
		oldCertificates := []*entity.EventCertificate{
			{Id: oldCertID, EventId: eventID},
		}

		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return(oldCertificates, nil)

		// Setup certificate config
		configID := uuid.New()
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		oldSigID := uuid.New()
		oldSignatures := []*entity.EventCertificateSignature{
			{Id: oldSigID, EventCertificateConfigId: configID},
		}

		mockCertSigDg := new(MockEventCertificateSignatureDataGateway)
		mockCertSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return(oldSignatures, nil)
		mockCertSigDg.On("DeleteEventCertificateSignature", ctx, oldSigID).
			Return(errors.New("delete signature error"))

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventDataGateway:                     mockEventDg,
			EventContractDataGateway:             mockEventContractDg,
			EventIssuerDataGateway:               mockEventIssuerDg,
			EventCertificateDataGateway:          mockCertDg,
			EventCertificateSignatureDataGateway: mockCertSigDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			cfg:                                  createMockConfig(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("test@example.com")},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockCertSigDg.AssertExpectations(t)
	})

	t.Run("should handle error when deleting certificates", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: strPtr(realEncryptedKey),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userId,
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		certificateAddress := "0xCertificateContractAddress"
		eventContract := &entity.EventContract{
			EventId:                    eventID,
			EventContractAddress:       "0xEventContractAddress",
			CertificateContractAddress: &certificateAddress,
		}
		mockEventContractDg.On("GetEventContractByEventID", ctx, eventID).
			Return(eventContract, nil)

		mockEventIssuerDg := new(MockEventIssuerDataGateway)
		mockEventIssuerDg.On("ResetAllEventIssuersSigningStatus", ctx, eventID).
			Return(nil)

		oldCertID := uuid.New()
		oldCertificates := []*entity.EventCertificate{
			{Id: oldCertID, EventId: eventID},
		}

		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return(oldCertificates, nil)
		mockCertDg.On("DeleteEventCertificate", ctx, oldCertID).
			Return(errors.New("delete certificate error"))

		// Setup certificate config
		configID := uuid.New()
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		mockCertSigDg := new(MockEventCertificateSignatureDataGateway)
		mockCertSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return([]*entity.EventCertificateSignature{}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventDataGateway:                     mockEventDg,
			EventContractDataGateway:             mockEventContractDg,
			EventIssuerDataGateway:               mockEventIssuerDg,
			EventCertificateDataGateway:          mockCertDg,
			EventCertificateSignatureDataGateway: mockCertSigDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			cfg:                                  createMockConfig(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("test@example.com")},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should fail when user does not have encrypted private key", func(t *testing.T) {
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

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("test@example.com")},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when certificate contract deployment reverts", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: strPtr(realEncryptedKey),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userId,
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		eventContractNoAddr := &entity.EventContract{
			EventId:                      eventID,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
			CertificateContractAddress:   nil, // No existing certificate contract
		}
		mockEventContractDg.On("GetEventContractByEventID", ctx, eventID).
			Return(eventContractNoAddr, nil)

		mockEventIssuerDg := new(MockEventIssuerDataGateway)
		mockEventIssuerDg.On("ResetAllEventIssuersSigningStatus", ctx, eventID).
			Return(nil)

		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return([]*entity.EventCertificate{}, nil)

		configID := uuid.New()
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		mockCertSigDg := new(MockEventCertificateSignatureDataGateway)
		mockCertSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return([]*entity.EventCertificateSignature{}, nil)

		mockBlockchainDg := new(MockBlockchainClientDataGateway)
		mockBlockchainDg.On("GetTransactOpts", ctx).
			Return(&bind.TransactOpts{}, nil)

		deployErr := customerror.Parse(&customerror.ErrInternalServer, fmt.Errorf("transaction reverted"))

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventDataGateway:                     mockEventDg,
			EventContractDataGateway:             mockEventContractDg,
			EventIssuerDataGateway:               mockEventIssuerDg,
			EventCertificateDataGateway:          mockCertDg,
			EventCertificateSignatureDataGateway: mockCertSigDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			BlockchainClientDg:                   mockBlockchainDg,
			cfg:                                  createMockConfig(),
			deployCertificateContract: func(_ context.Context, _ *bind.TransactOpts, _, _ common.Address) (common.Address, error) {
				return common.Address{}, deployErr
			},
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("test@example.com")},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrInternalServer.Code, *customErr.Code)
		mockEventContractDg.AssertExpectations(t) // UpdateEventContract must NOT have been called
		mockBlockchainDg.AssertExpectations(t)
	})

	t.Run("should not call deployCertificateContract when certificate contract address already exists", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			IsVerifiedIssuer:    false,
			EncryptedPrivateKey: strPtr(realEncryptedKey),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:                eventID,
			Title:             "Test Event",
			OwnerCredentialId: userId,
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		existingCertAddr := "0xExistingCertificateContract"
		eventContractWithAddr := &entity.EventContract{
			EventId:                      eventID,
			AccessManagerContractAddress: "0xAccessManager",
			EventContractAddress:         "0xEventContract",
			CertificateContractAddress:   &existingCertAddr,
		}
		mockEventContractDg.On("GetEventContractByEventID", ctx, eventID).
			Return(eventContractWithAddr, nil)

		mockEventIssuerDg := new(MockEventIssuerDataGateway)
		mockEventIssuerDg.On("ResetAllEventIssuersSigningStatus", ctx, eventID).
			Return(nil)

		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return([]*entity.EventCertificate{}, nil)
		// CreateEventCertificate is called after contract address check — return error to stop execution.
		// The key assertion is that deployCertificateContract panics if called (proving it was skipped).
		mockCertDg.On("CreateEventCertificate", ctx, mock.Anything).
			Return(nil, errors.New("stop here"))

		configID := uuid.New()
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		mockCertSigDg := new(MockEventCertificateSignatureDataGateway)
		mockCertSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return([]*entity.EventCertificateSignature{}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventDataGateway:                     mockEventDg,
			EventContractDataGateway:             mockEventContractDg,
			EventIssuerDataGateway:               mockEventIssuerDg,
			EventCertificateDataGateway:          mockCertDg,
			EventCertificateSignatureDataGateway: mockCertSigDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			cfg:                                  createMockConfig(),
			// If called, this panics — proving it is NOT called when address already exists
			deployCertificateContract: func(_ context.Context, _ *bind.TransactOpts, _, _ common.Address) (common.Address, error) {
				panic("deployCertificateContract must not be called when certificate address already exists")
			},
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: strPtr("test@example.com")},
		}

		// Act — deployCertificateContract must NOT be called (would panic), but CreateEventCertificate will error to stop execution
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, ImportCertificateReceiversOptions{HostPin: &hostPin}, currentUser)

		// Assert — error came from CreateEventCertificate mock, not from deployCertificateContract panic
		assert.Error(t, err)
		assert.Nil(t, result)
		// No panic = deployCertificateContract was never called
	})

	// M1: empty receiver list must be rejected before any DB access
	t.Run("should fail when receiver list is empty", func(t *testing.T) {
		uc := &EventUsecase{cfg: createMockConfig()}
		currentUser := &auth.JwtClaims{UserId: userId}

		result, err := uc.ImportCertificateReceivers(
			ctx, eventID,
			[]ImportCertificateReceiversRequest{},
			ImportCertificateReceiversOptions{HostPin: &hostPin},
			currentUser,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		customErr := customerror.TryParseAsCustomErr(err)
		require.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "at least one receiver")
	})

	// C2: wrong PIN must fail BEFORE any destructive database operations
	t.Run("should fail with wrong PIN before performing destructive operations", func(t *testing.T) {
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		encryptedKey, err := encryptutils.EncryptAESGCM(hex.EncodeToString(crypto.FromECDSA(privateKey)), "correct-pin")
		require.NoError(t, err)

		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			EncryptedPrivateKey: &encryptedKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{Id: eventID, OwnerCredentialId: userId}
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		certAddr := "0xCertContractAddress"
		mockContractDg := new(MockEventContractDataGateway)
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(&entity.EventContract{
			EventId:                    eventID,
			EventContractAddress:       "0xEventContractAddress",
			CertificateContractAddress: &certAddr,
		}, nil)

		// Intentionally provide NO expectations on the issuer gateway.
		// Any call to ResetAllEventIssuersSigningStatus means auth ran too late.
		mockIssuerDg := new(MockEventIssuerDataGateway)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventContractDataGateway:   mockContractDg,
			EventIssuerDataGateway:     mockIssuerDg,
			cfg:                        createMockConfig(),
		}

		wrongPin := "wrong-pin"
		result, err := uc.ImportCertificateReceivers(
			ctx, eventID,
			[]ImportCertificateReceiversRequest{{Email: strPtr("test@example.com")}},
			ImportCertificateReceiversOptions{HostPin: &wrongPin},
			&auth.JwtClaims{UserId: userId},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockIssuerDg.AssertNotCalled(t, "ResetAllEventIssuersSigningStatus")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockContractDg.AssertExpectations(t)
	})

	// C2: invalid BYOK signature must fail BEFORE any destructive database operations
	t.Run("should fail with invalid BYOK signature before performing destructive operations", func(t *testing.T) {
		walletKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		walletAddr := crypto.PubkeyToAddress(walletKey.PublicKey)

		// Sign the wrong message so the signature is valid-format but wrong
		wrongHash := cyptoutils.HashEthereumMessage("wrong message content")
		wrongSig, err := cyptoutils.Sign(wrongHash.Bytes(), walletKey)
		require.NoError(t, err)
		invalidSig := hexutil.Encode(wrongSig)

		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
			EncryptedPrivateKey: nil, // BYOK
			WalletAddress:       walletAddr.Hex(),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{Id: eventID, OwnerCredentialId: userId}
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		certAddr := "0xCertContractAddress"
		mockContractDg := new(MockEventContractDataGateway)
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(&entity.EventContract{
			EventId:                    eventID,
			EventContractAddress:       "0xEventContractAddress",
			CertificateContractAddress: &certAddr,
		}, nil)

		mockIssuerDg := new(MockEventIssuerDataGateway)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventContractDataGateway:   mockContractDg,
			EventIssuerDataGateway:     mockIssuerDg,
			cfg:                        createMockConfig(),
			logger:                     slog.Default(),
		}

		hostSignMessage := "ignored"
		result, err := uc.ImportCertificateReceivers(
			ctx, eventID,
			[]ImportCertificateReceiversRequest{{Email: strPtr("test@example.com")}},
			ImportCertificateReceiversOptions{
				HostSignature:   &invalidSig,
				HostSignMessage: &hostSignMessage,
			},
			&auth.JwtClaims{UserId: userId, WalletAddress: walletAddr.Hex()},
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		customErr := customerror.TryParseAsCustomErr(err)
		require.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockIssuerDg.AssertNotCalled(t, "ResetAllEventIssuersSigningStatus")
	})
}
