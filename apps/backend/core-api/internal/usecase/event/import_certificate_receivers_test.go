package event

import (
	"context"
	"decm-database/go/generated"
	"errors"
	"testing"

	"apps/backend/common/customerror"
	"apps/backend/core-api/config"
	"apps/backend/core-api/config/blockchain"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for testing
type MockEventCertificateDataGateway struct {
	mock.Mock
}

func (m *MockEventCertificateDataGateway) CreateEventCertificate(ctx context.Context, params datagateway.CreateEventCertificateParameters) (*entity.EventCertificate, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetEventCertificateByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetEventCertificateByInboxMessageID(ctx context.Context, inboxMessageID uuid.UUID) (*entity.EventCertificate, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateDataGateway) GetEventCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificate), args.Error(1)
}

func (m *MockEventCertificateDataGateway) GetAllEventCertificateIDsByEventID(ctx context.Context, eventID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

func (m *MockEventCertificateDataGateway) UpdateEventCertificate(ctx context.Context, id uuid.UUID, params datagateway.UpdateEventCertificateParameters) (*entity.EventCertificate, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateDataGateway) UpdateEventCertificateInboxMessageID(ctx context.Context, id uuid.UUID, inboxMessageID uuid.UUID) (*entity.EventCertificate, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateDataGateway) DeleteEventCertificate(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventCertificateDataGateway) GetClaimedCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateDataGateway) GetUnclaimedReadyCertificatesByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventCertificate, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateDataGateway) GetClaimedCertificatesByCredentialID(ctx context.Context, credentialID uuid.UUID, email *string) ([]*entity.EventCertificate, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateDataGateway) GetUnclaimedReadyCertificatesByCredentialID(ctx context.Context, credentialID uuid.UUID, email *string) ([]*entity.EventCertificate, error) {
	return nil, errors.New("not implemented")
}

type MockEventCertificateSignatureDataGateway struct {
	mock.Mock
}

func (m *MockEventCertificateSignatureDataGateway) CreateEventCertificateSignature(ctx context.Context, params datagateway.CreateEventCertificateSignatureParameters) (*entity.EventCertificateSignature, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateSignature), args.Error(1)
}

func (m *MockEventCertificateSignatureDataGateway) GetEventCertificateSignatureByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificateSignature, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateSignatureDataGateway) GetEventCertificateSignaturesByEventCertificateConfigID(ctx context.Context, eventCertificateConfigID uuid.UUID) ([]*entity.EventCertificateSignature, error) {
	args := m.Called(ctx, eventCertificateConfigID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EventCertificateSignature), args.Error(1)
}

func (m *MockEventCertificateSignatureDataGateway) UpdateEventCertificateSignature(ctx context.Context, id uuid.UUID, params datagateway.UpdateEventCertificateSignatureParameters) (*entity.EventCertificateSignature, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateSignatureDataGateway) UpdateEventCertificateIssuerSignature(ctx context.Context, id uuid.UUID, issuerSignature *string) (*entity.EventCertificateSignature, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateSignatureDataGateway) DeleteEventCertificateSignature(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockEventIssuerDataGateway struct {
	mock.Mock
}

func (m *MockEventIssuerDataGateway) CreateEventIssuer(ctx context.Context, params generated.CreateEventIssuerParams) (*generated.EventIssuer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventIssuerDataGateway) GetEventIssuerByID(ctx context.Context, id uuid.UUID) (generated.EventIssuer, error) {
	return generated.EventIssuer{}, errors.New("not implemented")
}

func (m *MockEventIssuerDataGateway) GetEventIssuersByEventID(ctx context.Context, eventID uuid.UUID) ([]generated.EventIssuer, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]generated.EventIssuer), args.Error(1)
}

func (m *MockEventIssuerDataGateway) UpdateEventIssuer(ctx context.Context, params generated.UpdateEventIssuerParams) (*generated.EventIssuer, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventIssuerDataGateway) DeleteEventIssuer(ctx context.Context, eventID uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *MockEventIssuerDataGateway) GetEventIssuerByEventIDAndIssuerCredentialID(ctx context.Context, eventID uuid.UUID, issuerCredentialID uuid.UUID) (generated.EventIssuer, error) {
	return generated.EventIssuer{}, errors.New("not implemented")
}

func (m *MockEventIssuerDataGateway) UpdateEventIssuerSigningStatus(ctx context.Context, eventID uuid.UUID, issuerCredentialID uuid.UUID, isSigned int32) error {
	return errors.New("not implemented")
}

func (m *MockEventIssuerDataGateway) ResetAllEventIssuersSigningStatus(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *MockEventIssuerDataGateway) AllIssuersHaveSigned(ctx context.Context, eventID uuid.UUID) (bool, error) {
	args := m.Called(ctx, eventID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEventIssuerDataGateway) GetSignedIssuersCount(ctx context.Context, eventID uuid.UUID) (int64, error) {
	args := m.Called(ctx, eventID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockEventIssuerDataGateway) GetTotalIssuersCount(ctx context.Context, eventID uuid.UUID) (int64, error) {
	args := m.Called(ctx, eventID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockEventIssuerDataGateway) HasSignedIssuers(ctx context.Context, eventID uuid.UUID) (bool, error) {
	args := m.Called(ctx, eventID)
	return args.Bool(0), args.Error(1)
}

type MockEventCertificateConfigDataGateway struct {
	mock.Mock
}

func (m *MockEventCertificateConfigDataGateway) CreateEventCertificateConfig(ctx context.Context, params generated.CreateEventCertificateConfigParams) (*entity.EventCertificateConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateConfigDataGateway) GetEventCertificateConfigByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificateConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateConfigDataGateway) GetEventCertificateConfigByEventID(ctx context.Context, eventId uuid.UUID) (*entity.EventCertificateConfig, error) {
	args := m.Called(ctx, eventId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventCertificateConfig), args.Error(1)
}

func (m *MockEventCertificateConfigDataGateway) UpdateEventCertificateConfig(ctx context.Context, params generated.UpdateEventCertificateConfigParams) (*entity.EventCertificateConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateConfigDataGateway) UpdateEventCertificateTextConfig(ctx context.Context, params generated.UpdateEventCertificateTextConfigParams) (*entity.EventCertificateConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateConfigDataGateway) ToggleEventCertificateConfigPublished(ctx context.Context, params generated.ToggleEventCertificateConfigPublishedParams) (*entity.EventCertificateConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventCertificateConfigDataGateway) DeleteEventCertificateConfig(ctx context.Context, eventID uuid.UUID) error {
	return errors.New("not implemented")
}

type MockEventContractDataGateway struct {
	mock.Mock
}

func (m *MockEventContractDataGateway) CreateEventContract(ctx context.Context, params generated.CreateEventContractParams) (*entity.EventContract, error) {
	return nil, errors.New("not implemented")
}

func (m *MockEventContractDataGateway) GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*entity.EventContract, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventContract), args.Error(1)
}

func (m *MockEventContractDataGateway) UpdateEventContract(ctx context.Context, params generated.UpdateEventContractParams) (*entity.EventContract, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EventContract), args.Error(1)
}

func (m *MockEventContractDataGateway) DeleteEventContract(ctx context.Context, eventID uuid.UUID) error {
	return errors.New("not implemented")
}

// Helper function to create a mock config for testing
func createMockConfigForImport() *config.Config {
	return &config.Config{
		Blockchain: blockchain.BlockchainConfig{
			ChainID:                  1,
			DecmAccessManagerAddress: "0x1234567890123456789012345678901234567890",
		},
	}
}

func TestImportCertificateReceivers(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	eventID := uuid.New()
	hostPin := "test-pin"

	t.Run("should fail when user is not authenticated", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(nil, errors.New("not found"))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfigForImport(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: "test@example.com", HostPin: hostPin},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, currentUser)

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
			cfg:                        createMockConfigForImport(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: "test@example.com", HostPin: hostPin},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, currentUser)

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
			EncryptedPrivateKey: stringPtr("encrypted-key"),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(nil, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			cfg:                        createMockConfigForImport(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: "test@example.com", HostPin: hostPin},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, currentUser)

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
			EncryptedPrivateKey: stringPtr("encrypted-key"),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:    eventID,
			Title: "Test Event",
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
			cfg:                        createMockConfigForImport(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: "test@example.com", HostPin: hostPin},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, currentUser)

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
			EncryptedPrivateKey: stringPtr("encrypted-key"),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:    eventID,
			Title: "Test Event",
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		certificateAddress := "0xCertificateContractAddress"
		eventContract := &entity.EventContract{
			EventID:                    eventID,
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

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventDataGateway:                     mockEventDg,
			EventContractDataGateway:             mockEventContractDg,
			EventIssuerDataGateway:               mockEventIssuerDg,
			EventCertificateDataGateway:          mockCertDg,
			EventCertificateSignatureDataGateway: mockCertSigDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			cfg:                                  createMockConfigForImport(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: "new@example.com", HostPin: hostPin},
		}

		// Act
		// Note: This will fail at blockchain operations, but we're testing the deletion logic
		// The deletion should happen before blockchain operations
		_, err := uc.ImportCertificateReceivers(ctx, eventID, requests, currentUser)

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
			EncryptedPrivateKey: stringPtr("encrypted-key"),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:    eventID,
			Title: "Test Event",
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		certificateAddress := "0xCertificateContractAddress"
		eventContract := &entity.EventContract{
			EventID:                    eventID,
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
			cfg:                         createMockConfigForImport(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: "test@example.com", HostPin: hostPin},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, currentUser)

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
			EncryptedPrivateKey: stringPtr("encrypted-key"),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:    eventID,
			Title: "Test Event",
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		certificateAddress := "0xCertificateContractAddress"
		eventContract := &entity.EventContract{
			EventID:                    eventID,
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
			cfg:                                  createMockConfigForImport(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: "test@example.com", HostPin: hostPin},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, currentUser)

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
			EncryptedPrivateKey: stringPtr("encrypted-key"),
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id:    eventID,
			Title: "Test Event",
		}
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(event, nil)

		mockEventContractDg := new(MockEventContractDataGateway)
		certificateAddress := "0xCertificateContractAddress"
		eventContract := &entity.EventContract{
			EventID:                    eventID,
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
			cfg:                                  createMockConfigForImport(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: "test@example.com", HostPin: hostPin},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, currentUser)

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
			cfg:                        createMockConfigForImport(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		requests := []ImportCertificateReceiversRequest{
			{Email: "test@example.com", HostPin: hostPin},
		}

		// Act
		result, err := uc.ImportCertificateReceivers(ctx, eventID, requests, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
	})
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
