package certificate_share

import (
	"apps/backend/common/customerror"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCertificateShare(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	certID := uuid.New()
	shareID := uuid.New()

	validUser := &auth.JwtClaims{UserId: userID}

	// pendingCert returns an EventCertificate in PENDING status (no CertificateTokenId / BroadcastedAt).
	pendingCert := func() *entity.EventCertificate {
		return &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: &userID,
		}
	}

	// claimedCert returns an EventCertificate in CLAIMED status.
	claimedCert := func() *entity.EventCertificate {
		tokenID := "99"
		now := time.Now()
		return &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: &userID,
			CertificateTokenId:   &tokenID,
			BroadcastedAt:        &now,
		}
	}

	// abortedCert returns an EventCertificate in ABORTED status.
	abortedCert := func() *entity.EventCertificate {
		now := time.Now()
		return &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: &userID,
			AbortedAt:            &now,
		}
	}

	returnedShare := func() *entity.CertificateShare {
		return &entity.CertificateShare{
			Id:                 shareID,
			EventCertificateId: certID,
			Handle:             "somegeneratedhandle1234",
			Active:             false,
		}
	}

	t.Run("should return ErrUnauthenticated when currentUser is nil", func(t *testing.T) {
		uc := &CertificateShareUsecase{}
		result, err := uc.CreateCertificateShare(ctx, nil, certID, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrUnauthenticated)
	})

	t.Run("should return ErrInternalServer when GetEventCertificateByID fails", func(t *testing.T) {
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(nil, errors.New("db error"))

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.CreateCertificateShare(ctx, validUser, certID, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return ErrNotFound when certificate does not exist", func(t *testing.T) {
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(nil, nil)

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.CreateCertificateShare(ctx, validUser, certID, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrNotFound)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return ErrForbidden when user is not the receiver", func(t *testing.T) {
		otherID := uuid.New()
		cert := &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: &otherID,
		}
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.CreateCertificateShare(ctx, validUser, certID, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrForbidden)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return ErrForbidden when certificate has no receiver credential and no email", func(t *testing.T) {
		cert := &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: nil,
			ReceiverEmail:        nil,
		}
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.CreateCertificateShare(ctx, validUser, certID, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrForbidden)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return ErrForbidden when certificate email does not match user email", func(t *testing.T) {
		otherEmail := "other@example.com"
		cert := &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: nil,
			ReceiverEmail:        &otherEmail,
		}
		userEmail := "user@example.com"
		userWithEmail := &auth.JwtClaims{UserId: userID, Email: &userEmail}
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.CreateCertificateShare(ctx, userWithEmail, certID, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrForbidden)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should allow share when certificate has no credential ID but email matches user", func(t *testing.T) {
		userEmail := "user@example.com"
		cert := &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: nil,
			ReceiverEmail:        &userEmail,
		}
		userWithEmail := &auth.JwtClaims{UserId: userID, Email: &userEmail}
		newShare := returnedShare()
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)
		mockShareDg.On("CreateCertificateShare", ctx, mock.MatchedBy(func(p event_datagateway.CreateCertificateShareParameters) bool {
			return p.EventCertificateId == certID && p.Active && p.Password == nil && len(p.Handle) == 32
		})).Return(newShare, nil)

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDg:          mockShareDg,
		}
		result, err := uc.CreateCertificateShare(ctx, userWithEmail, certID, nil)
		assert.NoError(t, err)
		assert.Equal(t, newShare, result)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInvalidArgument when certificate status is not PENDING or CLAIMED (e.g. aborted)", func(t *testing.T) {
		mockCertDg := new(MockEventCertificateDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(abortedCert(), nil)

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.CreateCertificateShare(ctx, validUser, certID, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInvalidArgument)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return existing share when one already exists (idempotency)", func(t *testing.T) {
		existing := returnedShare()
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(pendingCert(), nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(existing, nil)

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDg:          mockShareDg,
		}
		result, err := uc.CreateCertificateShare(ctx, validUser, certID, nil)
		assert.NoError(t, err)
		assert.Equal(t, existing, result)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when GetCertificateShareByEventCertificateID fails", func(t *testing.T) {
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(pendingCert(), nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, errors.New("db error"))

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDg:          mockShareDg,
		}
		result, err := uc.CreateCertificateShare(ctx, validUser, certID, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should create and return new share when no existing share (PENDING status)", func(t *testing.T) {
		newShare := returnedShare()
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(pendingCert(), nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)
		mockShareDg.On("CreateCertificateShare", ctx, mock.MatchedBy(func(p event_datagateway.CreateCertificateShareParameters) bool {
			return p.EventCertificateId == certID && p.Active && p.Password == nil && len(p.Handle) == 32
		})).Return(newShare, nil)

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDg:          mockShareDg,
		}
		result, err := uc.CreateCertificateShare(ctx, validUser, certID, nil)
		assert.NoError(t, err)
		assert.Equal(t, newShare, result)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should create and return new share when certificate is CLAIMED", func(t *testing.T) {
		newShare := returnedShare()
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(claimedCert(), nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)
		mockShareDg.On("CreateCertificateShare", ctx, mock.MatchedBy(func(p event_datagateway.CreateCertificateShareParameters) bool {
			return p.EventCertificateId == certID && p.Active && p.Password == nil && len(p.Handle) == 32
		})).Return(newShare, nil)

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDg:          mockShareDg,
		}
		result, err := uc.CreateCertificateShare(ctx, validUser, certID, nil)
		assert.NoError(t, err)
		assert.Equal(t, newShare, result)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when CreateCertificateShare datagateway fails", func(t *testing.T) {
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(pendingCert(), nil)
		mockShareDg.On("GetCertificateShareByEventCertificateID", ctx, certID).Return(nil, nil)
		mockShareDg.On("CreateCertificateShare", ctx, mock.Anything).Return(nil, errors.New("db error"))

		uc := &CertificateShareUsecase{
			EventCertificateDataGateway: mockCertDg,
			CertificateShareDg:          mockShareDg,
		}
		result, err := uc.CreateCertificateShare(ctx, validUser, certID, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockCertDg.AssertExpectations(t)
		mockShareDg.AssertExpectations(t)
	})
}
