package certificate_share

import (
	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/entity"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCertificateShareImage(t *testing.T) {
	ctx := context.Background()
	handle := "test-handle"
	certID := uuid.New()
	pngBytes := []byte{0x89, 0x50, 0x4E, 0x47} // minimal PNG magic bytes

	t.Run("should return ErrNotFound when handle does not exist", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(nil, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		result, err := uc.GetCertificateShareImage(ctx, handle, nil)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrNotFound)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when share datagateway returns an error", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(nil, errors.New("db error"))

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		result, err := uc.GetCertificateShareImage(ctx, handle, nil)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrForbidden when share is password-protected and no password given", func(t *testing.T) {
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

		result, err := uc.GetCertificateShareImage(ctx, handle, nil)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrForbidden)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrForbidden when wrong password is supplied", func(t *testing.T) {
		rawPw := "correct"
		hashedPw, err := hashutils.HashPassword(rawPw)
		require.NoError(t, err)

		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           &hashedPw,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		result, err := uc.GetCertificateShareImage(ctx, handle, strPtr("wrong"))

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrForbidden)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should propagate error from image generator", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           nil,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)

		mockImageGen := new(MockCertificateImageGenerator)
		mockImageGen.On("GenerateCertificateImage", ctx, certID).
			Return(nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("chrome failed")))

		uc := &CertificateShareUsecase{
			CertificateShareDg:        mockShareDg,
			CertificateImageGenerator: mockImageGen,
		}

		result, err := uc.GetCertificateShareImage(ctx, handle, nil)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
		mockImageGen.AssertExpectations(t)
	})

	t.Run("should return PNG bytes for a public share", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           nil,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)

		mockImageGen := new(MockCertificateImageGenerator)
		mockImageGen.On("GenerateCertificateImage", ctx, certID).Return(pngBytes, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:        mockShareDg,
			CertificateImageGenerator: mockImageGen,
		}

		result, err := uc.GetCertificateShareImage(ctx, handle, nil)

		assert.NoError(t, err)
		assert.Equal(t, pngBytes, result)
		mockShareDg.AssertExpectations(t)
		mockImageGen.AssertExpectations(t)
	})

	t.Run("should return PNG bytes when correct password is supplied", func(t *testing.T) {
		rawPw := "correct"
		hashedPw, err := hashutils.HashPassword(rawPw)
		require.NoError(t, err)

		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           &hashedPw,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)

		mockImageGen := new(MockCertificateImageGenerator)
		mockImageGen.On("GenerateCertificateImage", ctx, certID).Return(pngBytes, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:        mockShareDg,
			CertificateImageGenerator: mockImageGen,
		}

		result, err := uc.GetCertificateShareImage(ctx, handle, &rawPw)

		assert.NoError(t, err)
		assert.Equal(t, pngBytes, result)
		mockShareDg.AssertExpectations(t)
		mockImageGen.AssertExpectations(t)
	})
}
