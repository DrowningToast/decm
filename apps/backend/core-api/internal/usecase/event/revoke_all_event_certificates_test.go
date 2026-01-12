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

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRevokeAllEventCertificates(t *testing.T) {
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
		response, err := uc.RevokeAllEventCertificates(ctx, eventId, currentUser)

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

		// Act
		response, err := uc.RevokeAllEventCertificates(ctx, eventId, currentUser)

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

		// Act
		response, err := uc.RevokeAllEventCertificates(ctx, eventId, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should fail when user is not the event owner", func(t *testing.T) {
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
		response, err := uc.RevokeAllEventCertificates(ctx, eventId, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "user is not owner of the event")
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
	})

	t.Run("should return empty array when no certificates exist", func(t *testing.T) {
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
			OwnerCredentialId: userId, // Correct owner
		}
		mockEventDg.On("GetEventById", ctx, eventId).Return(event, nil)

		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetAllEventCertificateIDsByEventID", ctx, eventId).
			Return([]uuid.UUID{}, nil) // No certificates

		uc := &EventUsecase{
			AuthenticationCredentialDg:  mockAuthDg,
			EventDataGateway:            mockEventDg,
			EventCertificateDataGateway: mockCertDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		response, err := uc.RevokeAllEventCertificates(ctx, eventId, currentUser)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response.RevokedCertificates)
		assert.Len(t, response.RevokedCertificates, 0)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should successfully revoke all certificates", func(t *testing.T) {
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

		// Setup certificate IDs
		certificateId1 := uuid.New()
		certificateId2 := uuid.New()
		certificateIds := []uuid.UUID{certificateId1, certificateId2}

		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetAllEventCertificateIDsByEventID", ctx, eventId).
			Return(certificateIds, nil)

		// Setup certificates
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

		// Use mock.MatchedBy to flexibly match parameters with time.Now()
		mockCertDg.On("UpdateEventCertificate", ctx, certificateId1, mock.MatchedBy(func(params eventdatagateway.UpdateEventCertificateParameters) bool {
			return params.RevokedAt != nil &&
				params.ReceiverCredentialID == cert1.ReceiverCredentialId &&
				params.Name == cert1.Name
		})).Return(cert1, nil) // Return same cert (will have RevokedAt set by implementation)

		mockCertDg.On("UpdateEventCertificate", ctx, certificateId2, mock.MatchedBy(func(params eventdatagateway.UpdateEventCertificateParameters) bool {
			return params.RevokedAt != nil &&
				params.ReceiverCredentialID == cert2.ReceiverCredentialId &&
				params.Name == cert2.Name
		})).Return(cert2, nil) // Return same cert (will have RevokedAt set by implementation)

		uc := &EventUsecase{
			AuthenticationCredentialDg:  mockAuthDg,
			EventDataGateway:            mockEventDg,
			EventCertificateDataGateway: mockCertDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		response, err := uc.RevokeAllEventCertificates(ctx, eventId, currentUser)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.RevokedCertificates, 2)
		mockAuthDg.AssertExpectations(t)
		mockEventDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})
}
