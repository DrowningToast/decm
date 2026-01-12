package event

import (
	"apps/backend/common/customerror"
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQueueCertificateClaim(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	certificateId := uuid.New()
	eventId := uuid.New()
	contractAddress := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"

	// Create test private key and address
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	participantAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	participantPublicKey := &privateKey.PublicKey

	// Create test signature and sign message
	deadlineBlock := uint64(12345)
	signMessage, err := cyptoutils.GetSignMessage(participantAddress, common.HexToAddress(contractAddress), deadlineBlock)
	assert.NoError(t, err)
	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	signature, err := cyptoutils.Sign(messageHash.Bytes(), privateKey)
	assert.NoError(t, err)

	t.Run("should queue certificate claim in user signature and edit the certificate to link", func(t *testing.T) {
		// Arrange
		currentUser := &auth.JwtClaims{UserId: userId}

		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
			CertificateTokenId:      nil, // Not minted yet
			UserClaimSignatureId:    nil, // No signature yet
		}

		mockUserSigDg := new(MockUserSignatureDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockBlockchainDg := new(MockBlockchainClientDataGateway)

		signatureId := uuid.New()
		userSignature := &entity.UserSignature{
			Id:                         signatureId,
			AuthenticationCredentialId: userId,
			Signature:                  hex.EncodeToString(signature),
			SignMessage:                signMessage,
		}

		estimatedTime := time.Now().Add(5 * time.Minute)
		mockBlockchainDg.On("EstimateDeadlineTime", ctx, deadlineBlock).
			Return(&estimatedTime, nil)

		mockUserSigDg.On("CreateUserSignature", ctx, mock.MatchedBy(func(params offchain_datagateway.CreateUserSignatureParameters) bool {
			return params.Signature == hex.EncodeToString(signature) &&
				params.SignMessage == signMessage &&
				params.AuthenticationCredentialId == userId
		})).Return(userSignature, nil)

		updatedCertificate := *certificate
		updatedCertificate.UserClaimSignatureId = &signatureId

		mockCertDg.On("UpdateEventCertificate", ctx, certificateId, mock.MatchedBy(func(params event_datagateway.UpdateEventCertificateParameters) bool {
			return params.UserClaimSignatureId != nil && *params.UserClaimSignatureId == signatureId
		})).Return(&updatedCertificate, nil)

		uc := &EventUsecase{
			UserSignatureDg:             mockUserSigDg,
			EventCertificateDataGateway: mockCertDg,
			BlockchainClientDg:          mockBlockchainDg,
		}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, &participantAddress, participantPublicKey)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resultCert)
		assert.NotNil(t, resultSig)
		assert.Equal(t, signatureId, *resultCert.UserClaimSignatureId)
		assert.Equal(t, hex.EncodeToString(signature), resultSig.Signature)
		mockUserSigDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockBlockchainDg.AssertExpectations(t)
	})

	t.Run("should return error if the certificate is already minted", func(t *testing.T) {
		// Arrange
		currentUser := &auth.JwtClaims{UserId: userId}

		tokenId := "123"
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
			CertificateTokenId:      &tokenId, // Already minted
			UserClaimSignatureId:    nil,
		}

		uc := &EventUsecase{}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, &participantAddress, participantPublicKey)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultCert)
		assert.Nil(t, resultSig)
		assert.Contains(t, err.Error(), "certificate already minted")
	})

	t.Run("should return error if the certificate is already claimed in the database", func(t *testing.T) {
		// Arrange - Same as above, testing the same condition
		currentUser := &auth.JwtClaims{UserId: userId}

		tokenId := "123"
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
			CertificateTokenId:      &tokenId, // Already has token ID = claimed
			UserClaimSignatureId:    nil,
		}

		uc := &EventUsecase{}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, &participantAddress, participantPublicKey)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultCert)
		assert.Nil(t, resultSig)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "certificate already minted")
	})

	t.Run("should queue certificate again successfully if the signature is expired", func(t *testing.T) {
		// Arrange
		currentUser := &auth.JwtClaims{UserId: userId}

		existingSignatureId := uuid.New()
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
			CertificateTokenId:      nil,
			UserClaimSignatureId:    &existingSignatureId, // Already has signature
		}

		// Create expired signature (deadline in the past)
		pastTime := time.Now().Add(-10 * time.Minute)
		existingSignature := &entity.UserSignature{
			Id:                         existingSignatureId,
			AuthenticationCredentialId: userId,
			Signature:                  "old_signature",
			SignMessage:                "old_message",
			BroadcastedAt:              nil,       // Not broadcasted
			EstimatedDeadline:          &pastTime, // Expired
		}

		mockUserSigDg := new(MockUserSignatureDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockBlockchainDg := new(MockBlockchainClientDataGateway)

		mockUserSigDg.On("GetUserSignatureByID", ctx, existingSignatureId).
			Return(existingSignature, nil)

		newSignatureId := uuid.New()
		newUserSignature := &entity.UserSignature{
			Id:                         newSignatureId,
			AuthenticationCredentialId: userId,
			Signature:                  hex.EncodeToString(signature),
			SignMessage:                signMessage,
		}

		estimatedTime := time.Now().Add(5 * time.Minute)
		mockBlockchainDg.On("EstimateDeadlineTime", ctx, deadlineBlock).
			Return(&estimatedTime, nil)

		mockUserSigDg.On("CreateUserSignature", ctx, mock.MatchedBy(func(params offchain_datagateway.CreateUserSignatureParameters) bool {
			return params.Signature == hex.EncodeToString(signature)
		})).Return(newUserSignature, nil)

		updatedCertificate := *certificate
		updatedCertificate.UserClaimSignatureId = &newSignatureId

		mockCertDg.On("UpdateEventCertificate", ctx, certificateId, mock.MatchedBy(func(params event_datagateway.UpdateEventCertificateParameters) bool {
			return params.UserClaimSignatureId != nil && *params.UserClaimSignatureId == newSignatureId
		})).Return(&updatedCertificate, nil)

		uc := &EventUsecase{
			UserSignatureDg:             mockUserSigDg,
			EventCertificateDataGateway: mockCertDg,
			BlockchainClientDg:          mockBlockchainDg,
		}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, &participantAddress, participantPublicKey)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resultCert)
		assert.NotNil(t, resultSig)
		assert.Equal(t, newSignatureId, *resultCert.UserClaimSignatureId)
		assert.NotEqual(t, existingSignatureId, *resultCert.UserClaimSignatureId)
		mockUserSigDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockBlockchainDg.AssertExpectations(t)
	})

	t.Run("should return error if the signature isn't broadcasted and not expired", func(t *testing.T) {
		// Arrange
		currentUser := &auth.JwtClaims{UserId: userId}

		existingSignatureId := uuid.New()
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
			CertificateTokenId:      nil,
			UserClaimSignatureId:    &existingSignatureId,
		}

		// Create non-expired signature (deadline in the future)
		futureTime := time.Now().Add(10 * time.Minute)
		existingSignature := &entity.UserSignature{
			Id:                         existingSignatureId,
			AuthenticationCredentialId: userId,
			Signature:                  "pending_signature",
			SignMessage:                "pending_message",
			BroadcastedAt:              nil,         // Not broadcasted yet
			EstimatedDeadline:          &futureTime, // Not expired yet
		}

		mockUserSigDg := new(MockUserSignatureDataGateway)
		mockUserSigDg.On("GetUserSignatureByID", ctx, existingSignatureId).
			Return(existingSignature, nil)

		uc := &EventUsecase{
			UserSignatureDg: mockUserSigDg,
		}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, &participantAddress, participantPublicKey)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultCert)
		assert.Nil(t, resultSig)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrAlreadyPerformed.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "signature not expired")
		mockUserSigDg.AssertExpectations(t)
	})

	t.Run("should return error if signature is already broadcasted", func(t *testing.T) {
		// Arrange
		currentUser := &auth.JwtClaims{UserId: userId}

		existingSignatureId := uuid.New()
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
			CertificateTokenId:      nil,
			UserClaimSignatureId:    &existingSignatureId,
		}

		broadcastedTime := time.Now().Add(-5 * time.Minute)
		existingSignature := &entity.UserSignature{
			Id:                         existingSignatureId,
			AuthenticationCredentialId: userId,
			Signature:                  "broadcasted_signature",
			SignMessage:                "broadcasted_message",
			BroadcastedAt:              &broadcastedTime, // Already broadcasted
		}

		mockUserSigDg := new(MockUserSignatureDataGateway)
		mockUserSigDg.On("GetUserSignatureByID", ctx, existingSignatureId).
			Return(existingSignature, nil)

		uc := &EventUsecase{
			UserSignatureDg: mockUserSigDg,
		}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, &participantAddress, participantPublicKey)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultCert)
		assert.Nil(t, resultSig)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrInvalidArgument.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "signature already broadcasted")
		mockUserSigDg.AssertExpectations(t)
	})

	t.Run("should allow requeue when worker has explicitly marked signature as expired", func(t *testing.T) {
		// Arrange — MarkAsExpiredAt is set even though EstimatedDeadline is still in the future.
		// The worker sets this when it determines the signature can no longer be broadcast.
		currentUser := &auth.JwtClaims{UserId: userId}

		existingSignatureId := uuid.New()
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
			CertificateTokenId:      nil,
			UserClaimSignatureId:    &existingSignatureId,
		}

		futureTime := time.Now().Add(10 * time.Minute)
		workerExpiredAt := time.Now().Add(-1 * time.Minute)
		existingSignature := &entity.UserSignature{
			Id:                         existingSignatureId,
			AuthenticationCredentialId: userId,
			Signature:                  "old_signature",
			SignMessage:                "old_message",
			BroadcastedAt:              nil,
			EstimatedDeadline:          &futureTime,      // deadline not yet passed
			MarkAsExpiredAt:            &workerExpiredAt, // but worker already marked it
		}

		mockUserSigDg := new(MockUserSignatureDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockBlockchainDg := new(MockBlockchainClientDataGateway)

		mockUserSigDg.On("GetUserSignatureByID", ctx, existingSignatureId).
			Return(existingSignature, nil)

		newSignatureId := uuid.New()
		newUserSignature := &entity.UserSignature{
			Id:                         newSignatureId,
			AuthenticationCredentialId: userId,
			Signature:                  hex.EncodeToString(signature),
			SignMessage:                signMessage,
		}

		estimatedTime := time.Now().Add(5 * time.Minute)
		mockBlockchainDg.On("EstimateDeadlineTime", ctx, deadlineBlock).
			Return(&estimatedTime, nil)

		mockUserSigDg.On("CreateUserSignature", ctx, mock.MatchedBy(func(params offchain_datagateway.CreateUserSignatureParameters) bool {
			return params.Signature == hex.EncodeToString(signature)
		})).Return(newUserSignature, nil)

		updatedCertificate := *certificate
		updatedCertificate.UserClaimSignatureId = &newSignatureId

		mockCertDg.On("UpdateEventCertificate", ctx, certificateId, mock.MatchedBy(func(params event_datagateway.UpdateEventCertificateParameters) bool {
			return params.UserClaimSignatureId != nil && *params.UserClaimSignatureId == newSignatureId
		})).Return(&updatedCertificate, nil)

		uc := &EventUsecase{
			UserSignatureDg:             mockUserSigDg,
			EventCertificateDataGateway: mockCertDg,
			BlockchainClientDg:          mockBlockchainDg,
		}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, &participantAddress, participantPublicKey)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resultCert)
		assert.NotNil(t, resultSig)
		assert.Equal(t, newSignatureId, *resultCert.UserClaimSignatureId)
		mockUserSigDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockBlockchainDg.AssertExpectations(t)
	})

	t.Run("should return error when signature has no deadline info and is not expired", func(t *testing.T) {
		// Arrange — EstimatedDeadline == nil and MarkAsExpiredAt == nil.
		// Without any expiry signal we cannot determine the signature is safe to
		// overwrite, so re-queue must be blocked (Gap 1 fix).
		currentUser := &auth.JwtClaims{UserId: userId}

		existingSignatureId := uuid.New()
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
			CertificateTokenId:      nil,
			UserClaimSignatureId:    &existingSignatureId,
		}

		existingSignature := &entity.UserSignature{
			Id:                         existingSignatureId,
			AuthenticationCredentialId: userId,
			Signature:                  "pending_signature",
			SignMessage:                "pending_message",
			BroadcastedAt:              nil, // Not broadcasted
			EstimatedDeadline:          nil, // No deadline info
			MarkAsExpiredAt:            nil, // Not explicitly expired
		}

		mockUserSigDg := new(MockUserSignatureDataGateway)
		mockUserSigDg.On("GetUserSignatureByID", ctx, existingSignatureId).
			Return(existingSignature, nil)

		uc := &EventUsecase{
			UserSignatureDg: mockUserSigDg,
		}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, &participantAddress, participantPublicKey)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultCert)
		assert.Nil(t, resultSig)
		customErr := customerror.TryParseAsCustomErr(err)
		assert.NotNil(t, customErr)
		assert.Equal(t, customerror.ErrAlreadyPerformed.Code, *customErr.Code)
		assert.Contains(t, err.Error(), "signature not expired")
		mockUserSigDg.AssertExpectations(t)
	})

	t.Run("should return error when user is not authenticated", func(t *testing.T) {
		// Arrange
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
		}

		uc := &EventUsecase{}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, nil, certificate, signature, signMessage, &participantAddress, participantPublicKey)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultCert)
		assert.Nil(t, resultSig)
		assert.Contains(t, err.Error(), "not authenticated")
	})

	t.Run("should return error when participant address is nil", func(t *testing.T) {
		// Arrange
		currentUser := &auth.JwtClaims{UserId: userId}
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
		}

		uc := &EventUsecase{}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, nil, participantPublicKey)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultCert)
		assert.Nil(t, resultSig)
		assert.Contains(t, err.Error(), "participant address is required")
	})

	t.Run("should return error when participant public key is nil", func(t *testing.T) {
		// Arrange
		currentUser := &auth.JwtClaims{UserId: userId}
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: &contractAddress,
		}

		uc := &EventUsecase{}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, &participantAddress, nil)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultCert)
		assert.Nil(t, resultSig)
		assert.Contains(t, err.Error(), "participant public key is required")
	})

	t.Run("should return error when certificate contract not deployed", func(t *testing.T) {
		// Arrange
		currentUser := &auth.JwtClaims{UserId: userId}
		certificate := &entity.EventCertificate{
			Id:                      certificateId,
			EventId:                 eventId,
			EventCertificateAddress: nil, // Not deployed
		}

		uc := &EventUsecase{}

		// Act
		resultCert, resultSig, err := uc.queueCertificateClaim(ctx, currentUser, certificate, signature, signMessage, &participantAddress, participantPublicKey)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultCert)
		assert.Nil(t, resultSig)
		assert.Contains(t, err.Error(), "certificate contract not deployed")
	})
}
