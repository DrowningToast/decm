package certificate_share

import (
	"apps/backend/common/customerror"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateCertificateShare(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	certID := uuid.New()
	shareID := uuid.New()

	validUser := &auth.JwtClaims{UserId: userID}

	t.Run("should return ErrUnauthenticated when user is nil", func(t *testing.T) {
		uc := &CertificateShareUsecase{}
		result, err := uc.UpdateCertificateShare(ctx, nil, shareID, nil, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrUnauthenticated)
	})

	t.Run("should return ErrNotFound when share does not exist", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByID", ctx, shareID).Return(nil, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}
		result, err := uc.UpdateCertificateShare(ctx, validUser, shareID, nil, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrNotFound)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrNotFound when certificate does not exist", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 shareID,
			EventCertificateId: certID,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg.On("GetCertificateShareByID", ctx, shareID).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(nil, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.UpdateCertificateShare(ctx, validUser, shareID, nil, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrNotFound)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return ErrForbidden when user does not own the certificate", func(t *testing.T) {
		ownerID := uuid.New()
		share := &entity.CertificateShare{
			Id:                 shareID,
			EventCertificateId: certID,
		}
		cert := &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: &ownerID,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg.On("GetCertificateShareByID", ctx, shareID).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.UpdateCertificateShare(ctx, validUser, shareID, nil, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrForbidden)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should successfully update share with a new password", func(t *testing.T) {
		hashedPw := "$argon2id$hashed"
		share := &entity.CertificateShare{
			Id:                 shareID,
			EventCertificateId: certID,
		}
		cert := &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: &userID,
		}
		updatedShare := &entity.CertificateShare{
			Id:                 shareID,
			EventCertificateId: certID,
			Password:           &hashedPw,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg.On("GetCertificateShareByID", ctx, shareID).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockShareDg.On("UpdateCertificateShare", ctx, shareID, event_datagateway.UpdateCertificateShareParameters{Password: &hashedPw, Active: nil}).Return(updatedShare, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.UpdateCertificateShare(ctx, validUser, shareID, &hashedPw, nil)
		assert.NoError(t, err)
		assert.Equal(t, updatedShare, result)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should successfully remove password protection (nil)", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 shareID,
			EventCertificateId: certID,
		}
		cert := &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: &userID,
		}
		updatedShare := &entity.CertificateShare{
			Id:                 shareID,
			EventCertificateId: certID,
			Password:           nil,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg.On("GetCertificateShareByID", ctx, shareID).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockShareDg.On("UpdateCertificateShare", ctx, shareID, event_datagateway.UpdateCertificateShareParameters{Password: nil, Active: nil}).Return(updatedShare, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.UpdateCertificateShare(ctx, validUser, shareID, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, updatedShare, result)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when GetCertificateShareByID returns error", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByID", ctx, shareID).Return(nil, errors.New("db error"))

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}
		result, err := uc.UpdateCertificateShare(ctx, validUser, shareID, nil, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when certificate datagateway returns error", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 shareID,
			EventCertificateId: certID,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg.On("GetCertificateShareByID", ctx, shareID).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(nil, errors.New("db error"))

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.UpdateCertificateShare(ctx, validUser, shareID, nil, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when update datagateway returns error", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 shareID,
			EventCertificateId: certID,
		}
		cert := &entity.EventCertificate{
			Id:                   certID,
			ReceiverCredentialId: &userID,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockShareDg.On("GetCertificateShareByID", ctx, shareID).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockShareDg.On("UpdateCertificateShare", ctx, shareID, mock.Anything).Return(nil, errors.New("db error"))

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}
		result, err := uc.UpdateCertificateShare(ctx, validUser, shareID, nil, nil)
		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})
}
