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
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
		}

		broadcastedAt := time.Now().Add(-time.Hour)
		certID1 := uuid.New()
		certID2 := uuid.New()
		claimed := []*entity.EventCertificate{
			{
				Id:            certID1,
				EventId:       uuid.New(),
				BroadcastedAt: &broadcastedAt,
			},
		}
		unclaimed := []*entity.EventCertificate{
			{
				Id:      certID2,
				EventId: uuid.New(),
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return(unclaimed, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID1).Return(nil, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID2).Return(nil, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.ClaimedCertificates, 1)
		assert.Len(t, result.UnclaimedCertificates, 1)
		assert.Equal(t, 1, result.TotalClaimed)
		assert.Equal(t, 1, result.TotalUnclaimed)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("Certificate with token_id has CLAIMED status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
		}

		tokenId := "42"
		certID := uuid.New()
		claimed := []*entity.EventCertificate{
			{
				Id:                 certID,
				EventId:            uuid.New(),
				CertificateTokenId: &tokenId,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusClaimed, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("Certificate with token_id but no BroadcastedAt still has CLAIMED status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
		}

		tokenId := "1"
		certID := uuid.New()
		claimed := []*entity.EventCertificate{
			{
				Id:                 certID,
				EventId:            uuid.New(),
				CertificateTokenId: &tokenId,
				BroadcastedAt:      nil,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusClaimed, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate with BroadcastedAt has CLAIMED status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
		}

		broadcastedAt := time.Now().Add(-time.Hour)
		certID := uuid.New()
		claimed := []*entity.EventCertificate{
			{
				Id:            certID,
				EventId:       uuid.New(),
				BroadcastedAt: &broadcastedAt,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusClaimed, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate without BroadcastedAt and future deadline has PENDING status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
		}

		futureDeadline := time.Now().Add(24 * time.Hour)
		certID := uuid.New()
		claimed := []*entity.EventCertificate{
			{
				Id:                certID,
				EventId:           uuid.New(),
				BroadcastedAt:     nil,
				EstimatedDeadline: &futureDeadline,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusPending, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate without BroadcastedAt and past deadline has EXPIRED status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
		}

		pastDeadline := time.Now().Add(-time.Hour)
		certID := uuid.New()
		claimed := []*entity.EventCertificate{
			{
				Id:                certID,
				EventId:           uuid.New(),
				BroadcastedAt:     nil,
				EstimatedDeadline: &pastDeadline,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusExpired, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate without BroadcastedAt and nil deadline has PENDING status", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
		}

		certID := uuid.New()
		claimed := []*entity.EventCertificate{
			{
				Id:                certID,
				EventId:           uuid.New(),
				BroadcastedAt:     nil,
				EstimatedDeadline: nil,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Equal(t, event_usecase.ClaimCertificateStatusPending, result.ClaimedCertificates[0].Status)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("Empty certificates lists", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
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
		mockShareDg.AssertExpectations(t)
	})

	t.Run("GetClaimedCertificatesByCredentialID returns error", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
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
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
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
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
		}

		broadcastedAt := time.Now().Add(-time.Hour)
		pastDeadline := time.Now().Add(-time.Hour)
		futureDeadline := time.Now().Add(24 * time.Hour)
		tokenId := "99"
		certID1 := uuid.New()
		certID2 := uuid.New()
		certID3 := uuid.New()
		certID4 := uuid.New()

		claimed := []*entity.EventCertificate{
			{
				Id:            certID1,
				EventId:       uuid.New(),
				BroadcastedAt: &broadcastedAt,
			},
			{
				Id:                 certID2,
				EventId:            uuid.New(),
				CertificateTokenId: &tokenId,
			},
			{
				Id:                certID3,
				EventId:           uuid.New(),
				BroadcastedAt:     nil,
				EstimatedDeadline: &pastDeadline,
			},
			{
				Id:                certID4,
				EventId:           uuid.New(),
				BroadcastedAt:     nil,
				EstimatedDeadline: &futureDeadline,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID1).Return(nil, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID2).Return(nil, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID3).Return(nil, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID4).Return(nil, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		assert.Len(t, result.ClaimedCertificates, 4)
		assert.Equal(t, event_usecase.ClaimCertificateStatusClaimed, result.ClaimedCertificates[0].Status) // broadcasted
		assert.Equal(t, event_usecase.ClaimCertificateStatusClaimed, result.ClaimedCertificates[1].Status) // has token_id
		assert.Equal(t, event_usecase.ClaimCertificateStatusExpired, result.ClaimedCertificates[2].Status) // past deadline
		assert.Equal(t, event_usecase.ClaimCertificateStatusPending, result.ClaimedCertificates[3].Status) // future deadline
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate embeds EventCertificate fields", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
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
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		vm := result.ClaimedCertificates[0]
		assert.Equal(t, certID, vm.Id)
		assert.Equal(t, eventID, vm.EventId)
		assert.Equal(t, &eventName, vm.EventName)
		assert.Equal(t, &broadcastedAt, vm.BroadcastedAt)
		assert.Equal(t, &signatureCreatedAt, vm.SignatureCreatedAt)
		assert.Nil(t, vm.CertificateShare)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("Claimed certificate with share returns share viewmodel", func(t *testing.T) {
		mockCertDg := new(mockEventCertificateDg)
		mockShareDg := new(mockCertificateShareDg)
		uc := &event_usecase.EventUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDataGateway: mockShareDg,
		}

		certID := uuid.New()
		broadcastedAt := time.Now().Add(-time.Hour)
		shareHandle := "abc123"
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Active:             true,
			Handle:             shareHandle,
			Password:           nil,
		}
		claimed := []*entity.EventCertificate{
			{
				Id:            certID,
				EventId:       uuid.New(),
				BroadcastedAt: &broadcastedAt,
			},
		}

		mockCertDg.On("GetClaimedCertificatesByCredentialID", ctx, userID, &email).Return(claimed, nil)
		mockCertDg.On("GetUnclaimedReadyCertificatesByCredentialID", ctx, userID, &email).Return([]*entity.EventCertificate{}, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(share, nil)

		result, err := uc.GetMyCertificatesListViewModel(ctx, currentUser)

		assert.NoError(t, err)
		vm := result.ClaimedCertificates[0]
		assert.NotNil(t, vm.CertificateShare)
		assert.Equal(t, shareHandle, vm.CertificateShare.Handle)
		assert.Equal(t, true, vm.CertificateShare.Active)
		assert.Equal(t, false, vm.CertificateShare.HasPassword)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})
}
