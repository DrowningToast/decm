package event_test

import (
	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	event_usecase "apps/backend/core-api/internal/usecase/event"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsCertificateClaimed(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	credentialID := uuid.New()
	now := time.Now()

	t.Run("Certificate not claimed", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		mockCertDg.On("GetEventCertificateWithSignature", ctx, eventID, credentialID).Return(nil, nil)

		status, err := uc.IsCertificateClaimed(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.Nil(t, status)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Certificate claimed, signature exists, not broadcasted", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		cert := &eventdatagateway.EventCertificateWithSignature{
			Id:            uuid.New(),
			EventId:       eventID,
			CredentialId:  &credentialID,
			IsBroadcasted: false,
			Signature:     "0xclaimsig",
			SignMessage:   "claim message",
			JoinedAt:      now,
		}

		mockCertDg.On("GetEventCertificateWithSignature", ctx, eventID, credentialID).Return(cert, nil)

		status, err := uc.IsCertificateClaimed(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.False(t, status.IsBroadcasted)
		assert.Equal(t, "0xclaimsig", status.Signature)
		assert.Equal(t, "claim message", status.SignMessage)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Certificate claimed, token exists (broadcasted)", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		tokenID := "1"
		cert := &eventdatagateway.EventCertificateWithSignature{
			Id:            uuid.New(),
			EventId:       eventID,
			CredentialId:  &credentialID,
			IsBroadcasted: false, // Signature might not be broadcasted yet, but token exists
			TokenId:       &tokenID,
			JoinedAt:      now,
		}

		mockCertDg.On("GetEventCertificateWithSignature", ctx, eventID, credentialID).Return(cert, nil)

		status, err := uc.IsCertificateClaimed(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.True(t, status.IsBroadcasted)
		assert.Equal(t, &tokenID, status.CertificateTokenId)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Certificate claimed, IsBroadcasted true but no token", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		cert := &eventdatagateway.EventCertificateWithSignature{
			Id:            uuid.New(),
			EventId:       eventID,
			CredentialId:  &credentialID,
			IsBroadcasted: true,
			TokenId:       nil,
			Signature:     "0xsig",
			SignMessage:   "message",
			JoinedAt:      now,
		}

		mockCertDg.On("GetEventCertificateWithSignature", ctx, eventID, credentialID).Return(cert, nil)

		status, err := uc.IsCertificateClaimed(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.True(t, status.IsBroadcasted)
		assert.Nil(t, status.CertificateTokenId)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Certificate claimed, IsBroadcasted true and token exists", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		tokenID := "42"
		cert := &eventdatagateway.EventCertificateWithSignature{
			Id:            uuid.New(),
			EventId:       eventID,
			CredentialId:  &credentialID,
			IsBroadcasted: true,
			TokenId:       &tokenID,
			Signature:     "0xsig",
			SignMessage:   "message",
			JoinedAt:      now,
		}

		mockCertDg.On("GetEventCertificateWithSignature", ctx, eventID, credentialID).Return(cert, nil)

		status, err := uc.IsCertificateClaimed(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.True(t, status.IsBroadcasted)
		assert.Equal(t, &tokenID, status.CertificateTokenId)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Datagateway returns error", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		expectedErr := errors.New("database connection failed")
		mockCertDg.On("GetEventCertificateWithSignature", ctx, eventID, credentialID).Return(nil, expectedErr)

		status, err := uc.IsCertificateClaimed(ctx, eventID, credentialID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, status)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("BroadcastedAt field is properly returned", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		broadcastTime := now.Add(-time.Hour)
		cert := &eventdatagateway.EventCertificateWithSignature{
			Id:            uuid.New(),
			EventId:       eventID,
			CredentialId:  &credentialID,
			IsBroadcasted: true,
			BroadcastedAt: &broadcastTime,
			Signature:     "0xsig",
			SignMessage:   "message",
			JoinedAt:      now,
		}

		mockCertDg.On("GetEventCertificateWithSignature", ctx, eventID, credentialID).Return(cert, nil)

		status, err := uc.IsCertificateClaimed(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.NotNil(t, status.BroadcastedAt)
		assert.Equal(t, &broadcastTime, status.BroadcastedAt)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("ClaimedAt correctly maps from JoinedAt", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		joinedAt := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
		cert := &eventdatagateway.EventCertificateWithSignature{
			Id:            uuid.New(),
			EventId:       eventID,
			CredentialId:  &credentialID,
			IsBroadcasted: false,
			Signature:     "0xsig",
			SignMessage:   "message",
			JoinedAt:      joinedAt,
		}

		mockCertDg.On("GetEventCertificateWithSignature", ctx, eventID, credentialID).Return(cert, nil)

		status, err := uc.IsCertificateClaimed(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, joinedAt, status.ClaimedAt)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Empty signature and sign message", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		cert := &eventdatagateway.EventCertificateWithSignature{
			Id:            uuid.New(),
			EventId:       eventID,
			CredentialId:  &credentialID,
			IsBroadcasted: false,
			Signature:     "",
			SignMessage:   "",
			JoinedAt:      now,
		}

		mockCertDg.On("GetEventCertificateWithSignature", ctx, eventID, credentialID).Return(cert, nil)

		status, err := uc.IsCertificateClaimed(ctx, eventID, credentialID)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, "", status.Signature)
		assert.Equal(t, "", status.SignMessage)
		mockCertDg.AssertExpectations(t)
	})
}
