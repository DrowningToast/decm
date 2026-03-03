package certificate_share

import (
	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/entity"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCertificateShareByHandle(t *testing.T) {
	ctx := context.Background()
	handle := "share-handle"
	certID := uuid.New()
	tokenID := "1"
	contractAddr := "0xabcdef1234567890abcdef1234567890abcdef12"

	t.Run("should return Ready status and certificate when share is valid and certificate is claimed", func(t *testing.T) {
		broadcastedAt := time.Now()
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           nil,
		}
		cert := &entity.EventCertificate{
			Id:                      certID,
			CertificateTokenId:      &tokenID,
			EventCertificateAddress: &contractAddr,
			BroadcastedAt:           &broadcastedAt,
		}

		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}

		resultCert, status, err := uc.GetCertificateShareByHandle(ctx, handle)

		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, CertificateShareViewStatusReady, *status)
		assert.Equal(t, cert, resultCert)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return PasswordLocked status when share is password-protected", func(t *testing.T) {
		pw := "secret"
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           &pw,
		}

		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		resultCert, status, err := uc.GetCertificateShareByHandle(ctx, handle)

		assert.NoError(t, err)
		assert.Nil(t, resultCert)
		assert.NotNil(t, status)
		assert.Equal(t, CertificateShareViewStatusPasswordLocked, *status)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ValidButPending status when certificate has not been claimed", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           nil,
		}
		cert := &entity.EventCertificate{
			Id:                      certID,
			CertificateTokenId:      nil,
			EventCertificateAddress: nil,
		}

		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}

		resultCert, status, err := uc.GetCertificateShareByHandle(ctx, handle)

		assert.NoError(t, err)
		assert.Nil(t, resultCert)
		assert.NotNil(t, status)
		assert.Equal(t, CertificateShareViewStatusValidButPending, *status)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return ErrNotFound when handle does not exist", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(nil, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		resultCert, status, err := uc.GetCertificateShareByHandle(ctx, handle)

		assert.Nil(t, resultCert)
		assert.Nil(t, status)
		assertErrCode(t, err, customerror.ErrNotFound)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when share datagateway returns an error", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(nil, errors.New("db connection error"))

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		resultCert, status, err := uc.GetCertificateShareByHandle(ctx, handle)

		assert.Nil(t, resultCert)
		assert.Nil(t, status)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when certificate datagateway returns an error", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           nil,
		}

		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(nil, errors.New("db error"))

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}

		resultCert, status, err := uc.GetCertificateShareByHandle(ctx, handle)

		assert.Nil(t, resultCert)
		assert.Nil(t, status)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})
}

func TestGetCertificateShareByHandleWithPassword(t *testing.T) {
	ctx := context.Background()
	handle := "pw-share-handle"
	certID := uuid.New()
	tokenID := "99"
	contractAddr := "0x9999999999999999999999999999999999999999"

	t.Run("should return Ready status when share has no password (password arg is ignored)", func(t *testing.T) {
		broadcastedAt := time.Now()
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           nil,
		}
		cert := &entity.EventCertificate{
			Id:                      certID,
			CertificateTokenId:      &tokenID,
			EventCertificateAddress: &contractAddr,
			BroadcastedAt:           &broadcastedAt,
		}

		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}

		resultCert, status, err := uc.GetCertificateShareByHandleWithPassword(ctx, handle, "irrelevant")

		assert.NoError(t, err)
		assert.Equal(t, CertificateShareViewStatusReady, *status)
		assert.Equal(t, cert, resultCert)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return Ready status when correct password is supplied", func(t *testing.T) {
		rawPw := "mypassword"
		hashedPw, err := hashutils.HashPassword(rawPw)
		require.NoError(t, err)
		broadcastedAt := time.Now()
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           &hashedPw,
		}
		cert := &entity.EventCertificate{
			Id:                      certID,
			CertificateTokenId:      &tokenID,
			EventCertificateAddress: &contractAddr,
			BroadcastedAt:           &broadcastedAt,
		}

		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}

		resultCert, status, err := uc.GetCertificateShareByHandleWithPassword(ctx, handle, rawPw)

		assert.NoError(t, err)
		assert.Equal(t, CertificateShareViewStatusReady, *status)
		assert.Equal(t, cert, resultCert)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return PasswordLocked status when wrong password is supplied", func(t *testing.T) {
		rawPw := "correct"
		hashedPw, err := hashutils.HashPassword(rawPw)
		require.NoError(t, err)
		pw := hashedPw
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           &pw,
		}

		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		resultCert, status, err := uc.GetCertificateShareByHandleWithPassword(ctx, handle, "wrong")

		assert.NoError(t, err)
		assert.Nil(t, resultCert)
		assert.NotNil(t, status)
		assert.Equal(t, CertificateShareViewStatusPasswordLocked, *status)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrNotFound when handle does not exist", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(nil, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		resultCert, status, err := uc.GetCertificateShareByHandleWithPassword(ctx, handle, "any")

		assert.Nil(t, resultCert)
		assert.Nil(t, status)
		assertErrCode(t, err, customerror.ErrNotFound)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ValidButPending when certificate exists but is not yet claimed", func(t *testing.T) {
		rawPw := "pw"
		hashedPw, err := hashutils.HashPassword(rawPw)
		require.NoError(t, err)
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           &hashedPw,
		}
		cert := &entity.EventCertificate{
			Id:                      certID,
			CertificateTokenId:      nil,
			EventCertificateAddress: nil,
		}

		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}

		resultCert, status, err := uc.GetCertificateShareByHandleWithPassword(ctx, handle, rawPw)

		assert.NoError(t, err)
		assert.Nil(t, resultCert)
		assert.Equal(t, CertificateShareViewStatusValidButPending, *status)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})
}
