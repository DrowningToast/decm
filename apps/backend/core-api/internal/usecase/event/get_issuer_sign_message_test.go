package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIssuerSignMessage(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	userID := uuid.New()
	configID := uuid.New()

	t.Run("should fail when user is not found", func(t *testing.T) {
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, nil))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfigForSign(),
		}

		_, err := uc.GetIssuerSignMessage(ctx, eventID, &auth.JwtClaims{UserId: userID})
		assert.Error(t, err)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when user is not a verified issuer", func(t *testing.T) {
		encKey := "some_key"
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).
			Return(&entity.AuthenticationCredential{
				Id:                  userID,
				IsVerifiedIssuer:    false,
				EncryptedPrivateKey: &encKey,
			}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfigForSign(),
		}

		_, err := uc.GetIssuerSignMessage(ctx, eventID, &auth.JwtClaims{UserId: userID})
		assert.Error(t, err)
		customErr := customerror.TryParseAsCustomErr(err)
		require.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when user is not BYOK (has encrypted private key)", func(t *testing.T) {
		encKey := "encrypted_key"
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).
			Return(&entity.AuthenticationCredential{
				Id:                  userID,
				IsVerifiedIssuer:    true,
				EncryptedPrivateKey: &encKey, // Not BYOK
			}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			cfg:                        createMockConfigForSign(),
		}

		_, err := uc.GetIssuerSignMessage(ctx, eventID, &auth.JwtClaims{UserId: userID})
		assert.Error(t, err)
		customErr := customerror.TryParseAsCustomErr(err)
		require.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when certificate config not found", func(t *testing.T) {
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).
			Return(&entity.AuthenticationCredential{
				Id:                  userID,
				IsVerifiedIssuer:    true,
				EncryptedPrivateKey: nil, // BYOK
			}, nil)

		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(nil, customerror.Parse(&customerror.ErrNotFound, nil))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventCertificateConfigDg:   mockCertConfigDg,
			cfg:                        createMockConfigForSign(),
		}

		_, err := uc.GetIssuerSignMessage(ctx, eventID, &auth.JwtClaims{UserId: userID})
		assert.Error(t, err)
		mockAuthDg.AssertExpectations(t)
		mockCertConfigDg.AssertExpectations(t)
	})

	t.Run("should fail when no signature record found for issuer", func(t *testing.T) {
		otherIssuerID := uuid.New()

		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).
			Return(&entity.AuthenticationCredential{
				Id:                  userID,
				IsVerifiedIssuer:    true,
				EncryptedPrivateKey: nil,
			}, nil)

		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		mockSigDg := new(MockEventCertificateSignatureDataGateway)
		signMsg := "some message"
		mockSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return([]*entity.EventCertificateSignature{{
				Id:                       uuid.New(),
				EventCertificateConfigId: configID,
				IssuerCredentialId:       otherIssuerID, // Different issuer
				SignMessage:              &signMsg,
			}}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			EventCertificateSignatureDataGateway: mockSigDg,
			cfg:                                  createMockConfigForSign(),
		}

		_, err := uc.GetIssuerSignMessage(ctx, eventID, &auth.JwtClaims{UserId: userID})
		assert.Error(t, err)
		customErr := customerror.TryParseAsCustomErr(err)
		require.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrNotFound.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
		mockSigDg.AssertExpectations(t)
	})

	t.Run("should fail when sign message is nil", func(t *testing.T) {
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).
			Return(&entity.AuthenticationCredential{
				Id:                  userID,
				IsVerifiedIssuer:    true,
				EncryptedPrivateKey: nil,
			}, nil)

		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		mockSigDg := new(MockEventCertificateSignatureDataGateway)
		mockSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return([]*entity.EventCertificateSignature{{
				Id:                       uuid.New(),
				EventCertificateConfigId: configID,
				IssuerCredentialId:       userID,
				SignMessage:              nil, // Missing sign message
			}}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			EventCertificateSignatureDataGateway: mockSigDg,
			cfg:                                  createMockConfigForSign(),
		}

		_, err := uc.GetIssuerSignMessage(ctx, eventID, &auth.JwtClaims{UserId: userID})
		assert.Error(t, err)
		customErr := customerror.TryParseAsCustomErr(err)
		require.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrNotFound.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
		mockSigDg.AssertExpectations(t)
	})

	t.Run("should return sign message for BYOK issuer", func(t *testing.T) {
		expectedSignMessage := `{"eventContractAddress":"0xABC","receivers":["0x123"]}`

		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialByIdWithEncryptedPrivateKey", ctx, userID).
			Return(&entity.AuthenticationCredential{
				Id:                  userID,
				IsVerifiedIssuer:    true,
				EncryptedPrivateKey: nil,
			}, nil)

		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventID).
			Return(&entity.EventCertificateConfig{ID: configID, EventID: eventID}, nil)

		mockSigDg := new(MockEventCertificateSignatureDataGateway)
		mockSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, configID).
			Return([]*entity.EventCertificateSignature{{
				Id:                       uuid.New(),
				EventCertificateConfigId: configID,
				IssuerCredentialId:       userID,
				SignMessage:              &expectedSignMessage,
			}}, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			EventCertificateSignatureDataGateway: mockSigDg,
			cfg:                                  createMockConfigForSign(),
		}

		result, err := uc.GetIssuerSignMessage(ctx, eventID, &auth.JwtClaims{UserId: userID})
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedSignMessage, result.SignMessage)
		mockAuthDg.AssertExpectations(t)
		mockSigDg.AssertExpectations(t)
	})
}
