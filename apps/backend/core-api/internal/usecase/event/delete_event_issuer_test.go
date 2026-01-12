package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"decm-database/go/generated"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteEventIssuer(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	issuerId := uuid.New()
	eventId := uuid.New()
	issuerCredentialId := uuid.New()
	certConfigId := uuid.New()

	t.Run("should fail when user is not authenticated", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, userId).
			Return(nil, errors.New("not found"))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		err := uc.DeleteEventIssuer(ctx, issuerId, currentUser)

		// Assert
		assert.Error(t, err)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrInternalServer.Code, *customErr.Code)
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when user is not a verified organizer", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: false, // Not verified
		}
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, userId).
			Return(credential, nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		err := uc.DeleteEventIssuer(ctx, issuerId, currentUser)

		// Assert
		assert.Error(t, err)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrUnauthorized.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "user is not a verified organizer")
		mockAuthDg.AssertExpectations(t)
	})

	t.Run("should fail when issuer does not exist", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
		}
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, userId).
			Return(credential, nil)

		mockIssuerDg := new(MockEventIssuerDataGateway)
		mockIssuerDg.On("GetEventIssuerByID", ctx, issuerId).
			Return(generated.EventIssuer{}, errors.New("not found"))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventIssuerDataGateway:     mockIssuerDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		err := uc.DeleteEventIssuer(ctx, issuerId, currentUser)

		// Assert
		assert.Error(t, err)
		mockAuthDg.AssertExpectations(t)
		mockIssuerDg.AssertExpectations(t)
	})

	t.Run("should successfully delete issuer without certificate config", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
		}
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, userId).
			Return(credential, nil)

		mockIssuerDg := new(MockEventIssuerDataGateway)
		issuer := generated.EventIssuer{
			ID:                 issuerId,
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
		}
		mockIssuerDg.On("GetEventIssuerByID", ctx, issuerId).
			Return(issuer, nil)
		mockIssuerDg.On("DeleteEventIssuer", ctx, issuerId).Return(nil)

		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventId).
			Return(nil, errors.New("not found"))

		uc := &EventUsecase{
			AuthenticationCredentialDg: mockAuthDg,
			EventIssuerDataGateway:     mockIssuerDg,
			EventCertificateConfigDg:   mockCertConfigDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		err := uc.DeleteEventIssuer(ctx, issuerId, currentUser)

		// Assert
		assert.NoError(t, err)
		mockAuthDg.AssertExpectations(t)
		mockIssuerDg.AssertExpectations(t)
		mockCertConfigDg.AssertExpectations(t)
	})

	t.Run("should successfully delete issuer with certificate signature", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
		}
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, userId).
			Return(credential, nil)

		mockIssuerDg := new(MockEventIssuerDataGateway)
		issuer := generated.EventIssuer{
			ID:                 issuerId,
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
		}
		mockIssuerDg.On("GetEventIssuerByID", ctx, issuerId).
			Return(issuer, nil)
		mockIssuerDg.On("DeleteEventIssuer", ctx, issuerId).Return(nil)

		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		certConfig := &entity.EventCertificateConfig{
			ID:      certConfigId,
			EventID: eventId,
		}
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventId).
			Return(certConfig, nil)

		// Certificate signature for this issuer
		signatureId := uuid.New()
		signatures := []*entity.EventCertificateSignature{
			{
				Id:                 signatureId,
				IssuerCredentialId: issuerCredentialId, // Matches issuer
			},
		}

		mockSigDg := new(MockEventCertificateSignatureDataGateway)
		mockSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, certConfigId).
			Return(signatures, nil)
		mockSigDg.On("DeleteEventCertificateSignature", ctx, signatureId).Return(nil)

		uc := &EventUsecase{
			AuthenticationCredentialDg:             mockAuthDg,
			EventIssuerDataGateway:                 mockIssuerDg,
			EventCertificateConfigDg:               mockCertConfigDg,
			EventCertificateSignatureDataGateway:   mockSigDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		err := uc.DeleteEventIssuer(ctx, issuerId, currentUser)

		// Assert
		assert.NoError(t, err)
		mockAuthDg.AssertExpectations(t)
		mockIssuerDg.AssertExpectations(t)
		mockCertConfigDg.AssertExpectations(t)
		mockSigDg.AssertExpectations(t)
	})

	t.Run("should fail when signature deletion fails", func(t *testing.T) {
		// Arrange
		mockAuthDg := new(MockAuthenticationCredentialDg)
		credential := &entity.AuthenticationCredential{
			Id:                  userId,
			IsVerifiedOrganizer: true,
		}
		mockAuthDg.On("GetAuthenticationCredentialById", ctx, userId).
			Return(credential, nil)

		mockIssuerDg := new(MockEventIssuerDataGateway)
		issuer := generated.EventIssuer{
			ID:                 issuerId,
			EventID:            eventId,
			IssuerCredentialID: issuerCredentialId,
		}
		mockIssuerDg.On("GetEventIssuerByID", ctx, issuerId).
			Return(issuer, nil)

		mockCertConfigDg := new(MockEventCertificateConfigDataGateway)
		certConfig := &entity.EventCertificateConfig{
			ID:      certConfigId,
			EventID: eventId,
		}
		mockCertConfigDg.On("GetEventCertificateConfigByEventID", ctx, eventId).
			Return(certConfig, nil)

		signatureId := uuid.New()
		signatures := []*entity.EventCertificateSignature{
			{
				Id:                 signatureId,
				IssuerCredentialId: issuerCredentialId,
			},
		}

		mockSigDg := new(MockEventCertificateSignatureDataGateway)
		mockSigDg.On("GetEventCertificateSignaturesByEventCertificateConfigID", ctx, certConfigId).
			Return(signatures, nil)
		mockSigDg.On("DeleteEventCertificateSignature", ctx, signatureId).
			Return(errors.New("delete failed"))

		uc := &EventUsecase{
			AuthenticationCredentialDg:           mockAuthDg,
			EventIssuerDataGateway:               mockIssuerDg,
			EventCertificateConfigDg:             mockCertConfigDg,
			EventCertificateSignatureDataGateway: mockSigDg,
		}

		currentUser := &auth.JwtClaims{UserId: userId}

		// Act
		err := uc.DeleteEventIssuer(ctx, issuerId, currentUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete certificate signature for issuer")
		mockAuthDg.AssertExpectations(t)
		mockIssuerDg.AssertExpectations(t)
		mockCertConfigDg.AssertExpectations(t)
		mockSigDg.AssertExpectations(t)
	})
}
