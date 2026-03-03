package certificate_share

import (
	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/entity"
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// samplePayload returns a minimal valid CertificatePayload for happy-path tests.
func samplePayload() *entity.CertificatePayload {
	return &entity.CertificatePayload{
		Header: entity.CertificatePayloadHeader{
			Context:      []string{"https://www.w3.org/2018/credentials/v1"},
			Type:         []string{"VerifiableCredential", "EventCertificate"},
			Id:           "cert-id-1",
			Issuer:       "0xIssuer",
			IssuanceDate: "2024-01-01T00:00:00Z",
		},
		Data: entity.CertificatePayloadData{
			EventName:     "Test Event",
			CertificateId: "cert-id-1",
			UserId:        "user-1",
			Status:        "VALID",
		},
	}
}

// assertErrCode asserts that err is a *customerror.Err with the given ErrSignature's code.
func assertErrCode(t *testing.T, err error, sig customerror.ErrSignature) {
	t.Helper()
	customErr := customerror.TryParseAsCustomErr(err)
	require.NotNil(t, customErr, "expected a customerror.Err but got: %v", err)
	assert.Equal(t, sig.Code, *customErr.Code)
}

func TestGetCertificateShareData(t *testing.T) {
	ctx := context.Background()
	handle := "test-handle"
	certID := uuid.New()
	tokenID := "42"
	contractAddr := "0x1234567890123456789012345678901234567890"

	t.Run("should return payload for a valid public share with a claimed certificate", func(t *testing.T) {
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
		}
		payload := samplePayload()

		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockFactory := new(MockCertificateContractFactoryDg)
		mockContract := new(MockCertificateContractDg)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockFactory.On("GetContract", common.HexToAddress(contractAddr)).Return(mockContract, nil)
		mockContract.On("GetTokenData", ctx, big.NewInt(42)).Return(payload, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
		}

		result, err := uc.GetCertificateShareData(ctx, handle)

		assert.NoError(t, err)
		assert.Equal(t, payload, result)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should return ErrNotFound when handle does not exist", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(nil, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		result, err := uc.GetCertificateShareData(ctx, handle)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrNotFound)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrForbidden when share is password-protected", func(t *testing.T) {
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

		result, err := uc.GetCertificateShareData(ctx, handle)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrForbidden)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when share datagateway returns an error", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(nil, errors.New("db error"))

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		result, err := uc.GetCertificateShareData(ctx, handle)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when certificate is not found", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(nil, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}

		result, err := uc.GetCertificateShareData(ctx, handle)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return ErrInvalidArgument when certificate has no token ID (not claimed)", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
		}
		cert := &entity.EventCertificate{
			Id:                      certID,
			CertificateTokenId:      nil,
			EventCertificateAddress: &contractAddr,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}

		result, err := uc.GetCertificateShareData(ctx, handle)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInvalidArgument)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return ErrInvalidArgument when certificate has no contract address (not claimed)", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
		}
		cert := &entity.EventCertificate{
			Id:                      certID,
			CertificateTokenId:      &tokenID,
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

		result, err := uc.GetCertificateShareData(ctx, handle)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInvalidArgument)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})

	t.Run("should return error when GetContract fails", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
		}
		cert := &entity.EventCertificate{
			Id:                      certID,
			CertificateTokenId:      &tokenID,
			EventCertificateAddress: &contractAddr,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockFactory := new(MockCertificateContractFactoryDg)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockFactory.On("GetContract", common.HexToAddress(contractAddr)).
			Return(nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("rpc connection refused")))

		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
		}

		result, err := uc.GetCertificateShareData(ctx, handle)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
	})

	t.Run("should propagate error when GetTokenData fails", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
		}
		cert := &entity.EventCertificate{
			Id:                      certID,
			CertificateTokenId:      &tokenID,
			EventCertificateAddress: &contractAddr,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockFactory := new(MockCertificateContractFactoryDg)
		mockContract := new(MockCertificateContractDg)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockFactory.On("GetContract", common.HexToAddress(contractAddr)).Return(mockContract, nil)
		mockContract.On("GetTokenData", ctx, big.NewInt(42)).
			Return(nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("contract call failed")))

		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
		}

		result, err := uc.GetCertificateShareData(ctx, handle)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})
}

func TestGetCertificateShareDataWithPassword(t *testing.T) {
	ctx := context.Background()
	handle := "pw-handle"
	certID := uuid.New()
	tokenID := "7"
	contractAddr := "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"

	t.Run("should return payload when share has no password (password arg is ignored)", func(t *testing.T) {
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
		}
		payload := samplePayload()

		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockFactory := new(MockCertificateContractFactoryDg)
		mockContract := new(MockCertificateContractDg)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockFactory.On("GetContract", common.HexToAddress(contractAddr)).Return(mockContract, nil)
		mockContract.On("GetTokenData", ctx, big.NewInt(7)).Return(payload, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
		}

		result, err := uc.GetCertificateShareDataWithPassword(ctx, handle, "any-password-is-ignored")

		assert.NoError(t, err)
		assert.Equal(t, payload, result)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should return payload when correct password is supplied", func(t *testing.T) {
		rawPw := "correct-password"
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
			CertificateTokenId:      &tokenID,
			EventCertificateAddress: &contractAddr,
		}
		payload := samplePayload()

		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockFactory := new(MockCertificateContractFactoryDg)
		mockContract := new(MockCertificateContractDg)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)
		mockFactory.On("GetContract", common.HexToAddress(contractAddr)).Return(mockContract, nil)
		mockContract.On("GetTokenData", ctx, big.NewInt(7)).Return(payload, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
		}

		result, err := uc.GetCertificateShareDataWithPassword(ctx, handle, rawPw)

		assert.NoError(t, err)
		assert.Equal(t, payload, result)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should return ErrForbidden when wrong password is supplied", func(t *testing.T) {
		rawPw := "correct-password"
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

		result, err := uc.GetCertificateShareDataWithPassword(ctx, handle, "wrong-password")

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrForbidden)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrNotFound when handle does not exist", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(nil, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		result, err := uc.GetCertificateShareDataWithPassword(ctx, handle, "any")

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrNotFound)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInvalidArgument when certificate has not been claimed", func(t *testing.T) {
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

		result, err := uc.GetCertificateShareDataWithPassword(ctx, handle, "")

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInvalidArgument)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})
}
