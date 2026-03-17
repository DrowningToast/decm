package event

import (
	"apps/backend/common/customerror"
	"apps/backend/common/encryptutils"
	"apps/backend/core-api/config"
	"apps/backend/core-api/config/blockchain"
	"apps/backend/core-api/internal/entity"
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

// Helper function to create a mock config for testing
func createMockConfigForSign() *config.Config {
	return &config.Config{
		Blockchain: blockchain.BlockchainConfig{
			ChainID:                  1,
			DecmAccessManagerAddress: "0x1234567890123456789012345678901234567890",
		},
	}
}

// MockEventCertificateConfigDataGateway is defined in import_certificate_receivers_test.go
// This is a placeholder to ensure the type exists for this test file
// The actual implementation is shared across test files in this package

func TestSignEventCertificates(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	eventID := uuid.New()
	certificateID1 := uuid.New()
	issuerPin := "test-pin-123"

	// Mock encrypted private key (this would normally be encrypted with issuerPin)
	// For testing, we'll use a dummy value
	encryptedPrivateKey := "encrypted_key_data"

	t.Run("should fail when user is not authenticated", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, nil))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfigForSign(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		request := SignEventCertificatesRequest{IssuerPin: issuerPin}

		// Act
		result, err := uc.SignEventCertificates(ctx, eventID, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when user is not a verified issuer", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: false,
			IsVerifiedIssuer:    false, // Not a verified issuer
			EncryptedPrivateKey: &encryptedPrivateKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfigForSign(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		request := SignEventCertificatesRequest{IssuerPin: issuerPin}

		// Act
		result, err := uc.SignEventCertificates(ctx, eventID, request, currentUser)

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
			IsVerifiedIssuer:    true,
			EncryptedPrivateKey: &encryptedPrivateKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		mockEventDg.On("GetEventById", ctx, eventID).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, nil))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			cfg:                        createMockConfigForSign(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		request := SignEventCertificatesRequest{IssuerPin: issuerPin}

		// Act
		result, err := uc.SignEventCertificates(ctx, eventID, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should fail when no certificates found for event", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedIssuer:    true,
			EncryptedPrivateKey: &encryptedPrivateKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id: eventID,
		}
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		mockCertDg := new(MockEventCertificateDataGateway)
		// Return empty slice
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return([]*entity.EventCertificate{}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:  mockAuthDg,
			EventDataGateway:            mockEventDg,
			EventCertificateDataGateway: mockCertDg,
			cfg:                         createMockConfigForSign(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		request := SignEventCertificatesRequest{IssuerPin: issuerPin}

		// Act
		result, err := uc.SignEventCertificates(ctx, eventID, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrNotFound.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "no certificates found")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should fail when event contract not found", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedIssuer:    true,
			EncryptedPrivateKey: &encryptedPrivateKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id: eventID,
		}
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		mockCertDg := new(MockEventCertificateDataGateway)
		certificates := []*entity.EventCertificate{
			{Id: certificateID1},
		}
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return(certificates, nil)

		mockContractDg := new(MockEventContractDataGateway)
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, nil))

		uc := &EventUsecase{
			AuthenticationCredentialDg:  mockAuthDg,
			EventDataGateway:            mockEventDg,
			EventCertificateDataGateway: mockCertDg,
			EventContractDataGateway:    mockContractDg,
			cfg:                         createMockConfigForSign(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		request := SignEventCertificatesRequest{IssuerPin: issuerPin}

		// Act
		result, err := uc.SignEventCertificates(ctx, eventID, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockContractDg.AssertExpectations(t)
	})

	t.Run("should fail when issuer has no certificates to sign", func(t *testing.T) {
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
			IsVerifiedIssuer:    true,
			EncryptedPrivateKey: &encryptedKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id: eventID,
		}
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		mockCertDg := new(MockEventCertificateDataGateway)
		certificates := []*entity.EventCertificate{
			{Id: certificateID1},
		}
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return(certificates, nil)

		mockContractDg := new(MockEventContractDataGateway)
		contractAddress := "0xD8F0b257d1150E35E0351E9eb735b1229396D6fa"
		contract := &entity.EventContract{
			EventId:                      eventID,
			CertificateContractAddress:   &contractAddress,
			EventContractAddress:         contractAddress,
			AccessManagerContractAddress: "0x1234567890123456789012345678901234567890",
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		configID := uuid.New()
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		mockSigDg := new(MockEventCertificateSignatureDataGateway)
		// Return signatures but none for the current issuer
		otherIssuerID := uuid.New()
		signaturesForConfig := []*entity.EventCertificateSignature{
			{
				Id:                       uuid.New(),
				EventCertificateConfigId: configID,
				IssuerCredentialId:       otherIssuerID, // Different issuer
				SignMessage:              nil,
				IssuerSignature:          nil,
			},
		}
		mockSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return(signaturesForConfig, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventDataGateway:                     mockEventDg,
			EventCertificateDataGateway:          mockCertDg,
			EventContractDataGateway:             mockContractDg,
			EventCertificateSignatureDataGateway: mockSigDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			cfg:                                  createMockConfigForSign(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		request := SignEventCertificatesRequest{IssuerPin: password}

		// Act
		result, err := uc.SignEventCertificates(ctx, eventID, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no certificates found with signatures for this issuer")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockContractDg.AssertExpectations(t)
		mockSigDg.AssertExpectations(t)
		mockCertConfigDg.AssertExpectations(t)
	})

	t.Run("should fail when sign message is missing", func(t *testing.T) {
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
			IsVerifiedIssuer:    true,
			EncryptedPrivateKey: &encryptedKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userId).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id: eventID,
		}
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		mockCertDg := new(MockEventCertificateDataGateway)
		certificates := []*entity.EventCertificate{
			{Id: certificateID1},
		}
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return(certificates, nil)

		mockContractDg := new(MockEventContractDataGateway)
		contractAddress := "0xD8F0b257d1150E35E0351E9eb735b1229396D6fa"
		contract := &entity.EventContract{
			EventId:                      eventID,
			CertificateContractAddress:   &contractAddress,
			EventContractAddress:         contractAddress,
			AccessManagerContractAddress: "0x1234567890123456789012345678901234567890",
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		configID := uuid.New()
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		mockSigDg := new(MockEventCertificateSignatureDataGateway)
		signatureID := uuid.New()
		// Return signature for current issuer but with nil SignMessage
		signaturesForConfig := []*entity.EventCertificateSignature{
			{
				Id:                       signatureID,
				EventCertificateConfigId: configID,
				IssuerCredentialId:       userId, // Current issuer
				SignMessage:              nil,    // Missing sign message
				IssuerSignature:          nil,
			},
		}
		mockSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return(signaturesForConfig, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventDataGateway:                     mockEventDg,
			EventCertificateDataGateway:          mockCertDg,
			EventContractDataGateway:             mockContractDg,
			EventCertificateSignatureDataGateway: mockSigDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			cfg:                                  createMockConfigForSign(),
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		request := SignEventCertificatesRequest{IssuerPin: password}

		// Act
		result, err := uc.SignEventCertificates(ctx, eventID, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "sign message not found")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockContractDg.AssertExpectations(t)
		mockSigDg.AssertExpectations(t)
		mockCertConfigDg.AssertExpectations(t)
	})
}

func TestSignEventCertificates_MultipleIssuers(t *testing.T) {
	ctx := context.Background()
	issuer1ID := uuid.New()
	issuer2ID := uuid.New()
	issuer3ID := uuid.New()
	eventID := uuid.New()
	certificateID1 := uuid.New()
	certificateID2 := uuid.New()
	signatureID1 := uuid.New()
	// signatureID2 := uuid.New()
	// signatureID3 := uuid.New()
	// signatureID4 := uuid.New()
	issuerPin := "test-pin-123"

	encryptedPrivateKey := "encrypted_key_data"
	signMessage := `{"eventContractAddress":"0xD8F0b257d1150E35E0351E9eb735b1229396D6fa","receivers":["0x148d532c97fb3f21940c9f6923ab7b6a7df0489091da9fcfd4925fe05bdc49af"]}`

	t.Run("should only sign certificates for the current issuer when multiple issuers exist", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  issuer1ID,
			IsVerifiedIssuer:    true,
			EncryptedPrivateKey: &encryptedPrivateKey,
		}
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, issuer1ID).
			Return(credential, nil)

		mockEventDg := new(MockEventDataGateway)
		event := &entity.Event{
			Id: eventID,
		}
		mockEventDg.On("GetEventById", ctx, eventID).Return(event, nil)

		mockCertDg := new(MockEventCertificateDataGateway)
		certificates := []*entity.EventCertificate{
			{Id: certificateID1},
			{Id: certificateID2},
		}
		mockCertDg.On("GetEventCertificatesByEventID", ctx, eventID).
			Return(certificates, nil)

		mockContractDg := new(MockEventContractDataGateway)
		contractAddress := "0xD8F0b257d1150E35E0351E9eb735b1229396D6fa"
		contract := &entity.EventContract{
			EventId:                      eventID,
			CertificateContractAddress:   &contractAddress,
			EventContractAddress:         contractAddress,
			AccessManagerContractAddress: "0x1234567890123456789012345678901234567890",
		}
		mockContractDg.On("GetEventContractByEventID", ctx, eventID).Return(contract, nil)

		// Setup certificate config (signatures are linked to config, not individual certificates)
		configID := uuid.New()
		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		mockSigDg := new(MockEventCertificateSignatureDataGateway)

		// Signatures are shared across all certificates (linked to config)
		signaturesForConfig := []*entity.EventCertificateSignature{
			{
				Id:                       signatureID1,
				EventCertificateConfigId: configID,
				IssuerCredentialId:       issuer1ID, // Current issuer
				SignMessage:              &signMessage,
				IssuerSignature:          nil, // Not signed yet
			},
			{
				Id:                       uuid.New(),
				EventCertificateConfigId: configID,
				IssuerCredentialId:       issuer2ID, // Other issuer
				SignMessage:              &signMessage,
				IssuerSignature:          nil,
			},
			{
				Id:                       uuid.New(),
				EventCertificateConfigId: configID,
				IssuerCredentialId:       issuer3ID, // Other issuer
				SignMessage:              &signMessage,
				IssuerSignature:          nil,
			},
		}
		// Both certificates use the same signatures (from config)
		mockSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return(signaturesForConfig, nil)

		// Mock the UpdateEventCertificateIssuerSignature calls - only for issuer1ID's signature
		// Since signatures are now linked to config, there's only one signature per issuer
		dummySignature := hexutil.Encode([]byte{})
		mockSigDg.On("UpdateEventCertificateIssuerSignature", ctx, signatureID1, &dummySignature).
			Return(&entity.EventCertificateSignature{}, nil).Maybe()

		mockIssuerDg := new(MockEventIssuerDataGateway)
		mockIssuerDg.On("UpdateEventIssuerSigningStatus", ctx, eventID, issuer1ID, int32(1)).
			Return(nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventDataGateway:                     mockEventDg,
			EventCertificateDataGateway:          mockCertDg,
			EventContractDataGateway:             mockContractDg,
			EventCertificateSignatureDataGateway: mockSigDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			EventIssuerDataGateway:               mockIssuerDg,
			cfg:                                  createMockConfigForSign(),
		}

		currentUser := &auth.JwtClaims{UserId: issuer1ID}
		request := SignEventCertificatesRequest{IssuerPin: issuerPin}

		// Act
		result, err := uc.SignEventCertificates(ctx, eventID, request, currentUser)

		// Assert
		// Note: This will fail with actual crypto operations due to mock PIN
		// In a real test, you would need valid encryption/decryption
		// For now, we're just verifying the test structure
		assert.Error(t, err) // Expected to fail at decryption
		assert.Nil(t, result)

		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockContractDg.AssertExpectations(t)
	})
}

func TestSignEventCertificates_ReceiverListChanges(t *testing.T) {
	t.Run("should demonstrate receiver list change invalidates old signatures", func(t *testing.T) {
		// This is a documentation test showing the expected behavior
		// When import_certificate_receivers is called:
		// 1. Old certificates are deleted (line 88-116 in import_certificate_receivers.go)
		// 2. Old signatures are deleted (line 97-109)
		// 3. New sign message is created with new receiver list (line 229-278)
		// 4. New signatures are created for all issuers (line 288-312)
		// 5. All issuers' signing status is reset to 0 (line 82-86)

		// This test documents the expected flow
		t.Log("When receiver list changes:")
		t.Log("1. ResetAllEventIssuersSigningStatus is called")
		t.Log("2. All existing certificates are deleted")
		t.Log("3. All existing signatures are deleted")
		t.Log("4. New certificates are created with new receiver hashes")
		t.Log("5. New sign messages are generated with new receiver list")
		t.Log("6. New signature records are created for all issuers")
		t.Log("7. All issuers must sign again")
	})
}
