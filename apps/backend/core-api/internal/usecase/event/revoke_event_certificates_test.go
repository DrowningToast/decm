package event

import (
	"apps/backend/common/customerror"
	"apps/backend/common/encryptutils"
	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRevokeEventCertificates(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	eventId := uuid.New()
	certificateId1 := uuid.New()
	certificateId2 := uuid.New()
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
		request := RevokeEventCertificatesRequest{
			CertificateIDs: []uuid.UUID{certificateId1},
		}

		// Act
		response, err := uc.RevokeEventCertificates(ctx, eventId, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
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
		request := RevokeEventCertificatesRequest{
			CertificateIDs: []uuid.UUID{certificateId1},
		}

		// Act
		response, err := uc.RevokeEventCertificates(ctx, eventId, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "user is not a verified organizer")
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
		request := RevokeEventCertificatesRequest{
			CertificateIDs: []uuid.UUID{certificateId1},
		}

		// Act
		response, err := uc.RevokeEventCertificates(ctx, eventId, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should fail when certificate config is published", func(t *testing.T) {
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

		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		certConfig := &entity.EventCertificateConfig{
			EventID:     eventId,
			IsPublished: true, // Already published - should fail
		}
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventId).
			Return(certConfig, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventCertificateConfigDg:   mockCertConfigDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		request := RevokeEventCertificatesRequest{
			CertificateIDs: []uuid.UUID{certificateId1},
		}

		// Act
		response, err := uc.RevokeEventCertificates(ctx, eventId, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrForbidden.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "cannot revoke certificates after certificate configuration has been published")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockCertConfigDg.AssertExpectations(t)
	})

	t.Run("should fail when certificate does not exist", func(t *testing.T) {
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

		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		certConfig := &entity.EventCertificateConfig{
			EventID:     eventId,
			IsPublished: false, // Not published yet - OK to revoke
		}
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventId).
			Return(certConfig, nil)

		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certificateId1).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, nil))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventDataGateway:           mockEventDg,
			EventCertificateConfigDg:   mockCertConfigDg,
			EventCertificateDataGateway: mockCertDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		request := RevokeEventCertificatesRequest{
			CertificateIDs: []uuid.UUID{certificateId1},
		}

		// Act
		response, err := uc.RevokeEventCertificates(ctx, eventId, request, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockCertConfigDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should successfully revoke multiple certificates", func(t *testing.T) {
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

		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		certConfig := &entity.EventCertificateConfig{
			EventID:     eventId,
			IsPublished: false,
		}
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventId).
			Return(certConfig, nil)

		// Setup certificates
		mockCertDg := new(MockEventCertificateDataGateway)
		cert1 := &entity.EventCertificate{
			Id:                      certificateId1,
			EventId:                 eventId,
			ReceiverCredentialId:    &userId,
			ReceiverEmail:           stringPtr("user1@example.com"),
			Name:                    stringPtr("Test User 1"),
			AcademicInstitution:     stringPtr("Test University"),
			CertificateTitle:        stringPtr("Test Certificate"),
			CertificateSubtitle:     stringPtr("For Excellence"),
			EventContractAddress:    "0x1234567890123456789012345678901234567890",
			EventCertificateAddress: stringPtr("0x0987654321098765432109876543210987654321"),
			CertificateTokenId:      stringPtr("1"),
		}
		cert2 := &entity.EventCertificate{
			Id:                      certificateId2,
			EventId:                 eventId,
			ReceiverCredentialId:    &userId,
			ReceiverEmail:           stringPtr("user2@example.com"),
			Name:                    stringPtr("Test User 2"),
			AcademicInstitution:     stringPtr("Test University"),
			CertificateTitle:        stringPtr("Test Certificate"),
			CertificateSubtitle:     stringPtr("For Excellence"),
			EventContractAddress:    "0x1234567890123456789012345678901234567890",
			EventCertificateAddress: stringPtr("0x0987654321098765432109876543210987654321"),
			CertificateTokenId:      stringPtr("2"),
		}

		mockCertDg.On("GetEventCertificateByID", ctx, certificateId1).Return(cert1, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certificateId2).Return(cert2, nil)

		// Setup update responses - certificates with RevokedAt set
		now := time.Now()
		revokedCert1 := *cert1
		revokedCert1.RevokedAt = &now
		revokedCert2 := *cert2
		revokedCert2.RevokedAt = &now

		// Use mock.MatchedBy to flexibly match parameters with time.Now()
		mockCertDg.On("UpdateEventCertificate", ctx, certificateId1, mock.MatchedBy(func(params eventdatagateway.UpdateEventCertificateParameters) bool {
			return params.RevokedAt != nil &&
				params.ReceiverCredentialID == cert1.ReceiverCredentialId &&
				params.Name == cert1.Name
		})).Return(&revokedCert1, nil)

		mockCertDg.On("UpdateEventCertificate", ctx, certificateId2, mock.MatchedBy(func(params eventdatagateway.UpdateEventCertificateParameters) bool {
			return params.RevokedAt != nil &&
				params.ReceiverCredentialID == cert2.ReceiverCredentialId &&
				params.Name == cert2.Name
		})).Return(&revokedCert2, nil)

		// Mock issuer reset
		mockIssuerDg := new(MockEventIssuerDataGateway)
		mockIssuerDg.On("ResetAllEventIssuersSigningStatus", ctx, eventId).Return(nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:  mockAuthDg,
			EventDataGateway:            mockEventDg,
			EventCertificateConfigDg:    mockCertConfigDg,
			EventCertificateDataGateway: mockCertDg,
			EventIssuerDataGateway:      mockIssuerDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}
		request := RevokeEventCertificatesRequest{
			CertificateIDs: []uuid.UUID{certificateId1, certificateId2},
		}

		// Act
		response, err := uc.RevokeEventCertificates(ctx, eventId, request, currentUser)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.RevokedCertificates, 2)
		assert.NotNil(t, response.RevokedCertificates[0].RevokedAt)
		assert.NotNil(t, response.RevokedCertificates[1].RevokedAt)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockCertConfigDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockIssuerDg.AssertExpectations(t)
	})
}
