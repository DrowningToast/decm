package event_test

import (
	"apps/backend/core-api/internal/entity"
	event_usecase "apps/backend/core-api/internal/usecase/event"
	"apps/backend/services/auth"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetMyCertificatesListViewModel(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	email := "user@example.com"
	currentUser := &auth.JwtClaims{
		UserId: userID,
		Email:  &email,
	}

	t.Run("Returns claimed and unclaimed certificates", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		broadcastedAt := time.Now().Add(-time.Hour)
		claimed := []*entity.EventCertificate{
			{
				Id:            uuid.New(),
				EventId:       uuid.New(),
				BroadcastedAt: &broadcastedAt,
			},
		}
		unclaimed := []*entity.EventCertificate{
			{
				Id:      uuid.New(),
				EventId: uuid.New(),
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return(unclaimed, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.ClaimedCertificates, 1)
		assert.Len(t, result.UnclaimedCertificates, 1)
		assert.Equal(t, 1, result.TotalClaimed)
		assert.Equal(t, 1, result.TotalUnclaimed)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Certificate with token_id has CLAIMED status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		tokenId := "42"
		claimed := []*entity.EventCertificate{
			{
				Id:                 uuid.New(),
				EventId:            uuid.New(),
				CertificateTokenId: &tokenId,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusClaimed, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Certificate with token_id but no BroadcastedAt still has CLAIMED status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		tokenId := "1"
		claimed := []*entity.EventCertificate{
			{
				Id:                 uuid.New(),
				EventId:            uuid.New(),
				CertificateTokenId: &tokenId,
				BroadcastedAt:      nil,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusClaimed, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate with BroadcastedAt has CLAIMED status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		broadcastedAt := time.Now().Add(-time.Hour)
		claimed := []*entity.EventCertificate{
			{
				Id:            uuid.New(),
				EventId:       uuid.New(),
				BroadcastedAt: &broadcastedAt,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusClaimed, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate without BroadcastedAt and future deadline has PENDING status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		futureDeadline := time.Now().Add(24 * time.Hour)
		claimed := []*entity.EventCertificate{
			{
				Id:                uuid.New(),
				EventId:           uuid.New(),
				BroadcastedAt:     nil,
				EstimatedDeadline: &futureDeadline,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusPending, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate without BroadcastedAt and past deadline has EXPIRED status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		pastDeadline := time.Now().Add(-time.Hour)
		claimed := []*entity.EventCertificate{
			{
				Id:                uuid.New(),
				EventId:           uuid.New(),
				BroadcastedAt:     nil,
				EstimatedDeadline: &pastDeadline,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusExpired, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate without BroadcastedAt and nil deadline has PENDING status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		claimed := []*entity.EventCertificate{
			{
				Id:                uuid.New(),
				EventId:           uuid.New(),
				BroadcastedAt:     nil,
				EstimatedDeadline: nil,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusPending, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Empty certificates lists", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.ClaimedCertificates, 0)
		assert.Len(t, result.UnclaimedCertificates, 0)
		assert.Equal(t, 0, result.TotalClaimed)
		assert.Equal(t, 0, result.TotalUnclaimed)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("GetClaimedCertificatesByCredentialID returns error", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		expectedErr := errors.New("database error")
		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(nil, expectedErr)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get claimed certificates")
		mockCertDg.AssertExpectations(t)
	})

	t.Run("GetUnclaimedReadyCertificatesByCredentialID returns error", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		expectedErr := errors.New("database error")
		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return(nil, expectedErr)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get unclaimed ready certificates")
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Multiple claimed certificates with mixed statuses", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		broadcastedAt := time.Now().Add(-time.Hour)
		pastDeadline := time.Now().Add(-time.Hour)
		futureDeadline := time.Now().Add(24 * time.Hour)
		tokenId := "99"

		claimed := []*entity.EventCertificate{
			{
				Id:            uuid.New(),
				EventId:       uuid.New(),
				BroadcastedAt: &broadcastedAt,
			},
			{
				Id:                 uuid.New(),
				EventId:            uuid.New(),
				CertificateTokenId: &tokenId,
			},
			{
				Id:                uuid.New(),
				EventId:           uuid.New(),
				BroadcastedAt:     nil,
				EstimatedDeadline: &pastDeadline,
			},
			{
				Id:                uuid.New(),
				EventId:           uuid.New(),
				BroadcastedAt:     nil,
				EstimatedDeadline: &futureDeadline,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Len(t, result.ClaimedCertificates, 4)
		assert.Equal(t, event_usecase.ClaimCertificateStatusClaimed, result.ClaimedCertificates[0].Status)  // broadcasted
		assert.Equal(t, event_usecase.ClaimCertificateStatusClaimed, result.ClaimedCertificates[1].Status)  // has token_id
		assert.Equal(t, event_usecase.ClaimCertificateStatusExpired, result.ClaimedCertificates[2].Status)  // past deadline
		assert.Equal(t, event_usecase.ClaimCertificateStatusPending, result.ClaimedCertificates[3].Status)  // future deadline
		mockCertDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate embeds EventCertificate fields", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
		}

		certID := uuid.New()
		eventID := uuid.New()
		eventName := "Test Event"
		broadcastedAt := time.Now().Add(-time.Hour)
		signatureCreatedAt := time.Now().Add(-2 * time.Hour)
		claimed := []*entity.EventCertificate{
			{
				Id:                 certID,
				EventId:            eventID,
				EventName:          &eventName,
				BroadcastedAt:      &broadcastedAt,
				SignatureCreatedAt: &signatureCreatedAt,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		vm := result.ClaimedCertificates[0]
		assert.Equal(t, certID, vm.Id)
		assert.Equal(t, eventID, vm.EventId)
		assert.Equal(t, &eventName, vm.EventName)
		assert.Equal(t, &broadcastedAt, vm.BroadcastedAt)
		assert.Equal(t, &signatureCreatedAt, vm.SignatureCreatedAt)
		mockCertDg.AssertExpectations(t)
	})
}
