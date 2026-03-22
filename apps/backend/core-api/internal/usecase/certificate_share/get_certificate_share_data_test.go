package certificate_share

import (
	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
			Active:             true,
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

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		assert.NoError(t, err)
		assert.Equal(t, payload, result.Payload)
		assert.Equal(t, contractAddr, result.Contract.EventCertificateContractAddress)
		assert.Equal(t, tokenID, result.Contract.CertificateTokenId)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should return ErrNotFound when handle does not exist", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(nil, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrNotFound)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrForbidden when share is password-protected and no password given", func(t *testing.T) {
		pw := "secret"
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           &pw,
			Active:             true,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrForbidden)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when share datagateway returns an error", func(t *testing.T) {
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(nil, errors.New("db error"))

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInternalServer when certificate is not found", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Active:             true,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)

		mockShareDg.On("GetCertificateShareByHandle", ctx, handle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(nil, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

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
			Active:             true,
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

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

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
			Active:             true,
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

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

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
			Active:             true,
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

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

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
			Active:             true,
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

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInternalServer)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should return payload when share has no password and password arg is ignored", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Password:           nil,
			Active:             true,
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

		result, err := uc.GetCertificateShareData(ctx, handle, strPtr("any-password-is-ignored"))

		assert.NoError(t, err)
		assert.Equal(t, payload, result.Payload)
		assert.Equal(t, contractAddr, result.Contract.EventCertificateContractAddress)
		assert.Equal(t, tokenID, result.Contract.CertificateTokenId)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should return payload when correct password is supplied for a password-protected share", func(t *testing.T) {
		rawPw := "correct-password"
		hashedPw, err := hashutils.HashPassword(rawPw)
		require.NoError(t, err)

		pwHandle := "pw-handle"
		pwCertID := uuid.New()
		pwTokenID := "7"
		pwContractAddr := "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"

		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: pwCertID,
			Handle:             pwHandle,
			Password:           &hashedPw,
			Active:             true,
		}
		cert := &entity.EventCertificate{
			Id:                      pwCertID,
			CertificateTokenId:      &pwTokenID,
			EventCertificateAddress: &pwContractAddr,
		}
		payload := samplePayload()

		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)
		mockFactory := new(MockCertificateContractFactoryDg)
		mockContract := new(MockCertificateContractDg)

		mockShareDg.On("GetCertificateShareByHandle", ctx, pwHandle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, pwCertID).Return(cert, nil)
		mockFactory.On("GetContract", common.HexToAddress(pwContractAddr)).Return(mockContract, nil)
		mockContract.On("GetTokenData", ctx, big.NewInt(7)).Return(payload, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
		}

		result, err := uc.GetCertificateShareData(ctx, pwHandle, &rawPw)

		assert.NoError(t, err)
		assert.Equal(t, payload, result.Payload)
		assert.Equal(t, pwContractAddr, result.Contract.EventCertificateContractAddress)
		assert.Equal(t, pwTokenID, result.Contract.CertificateTokenId)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should return ErrForbidden when wrong password is supplied", func(t *testing.T) {
		rawPw := "correct-password"
		hashedPw, err := hashutils.HashPassword(rawPw)
		require.NoError(t, err)

		pwHandle := "pw-handle"
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             pwHandle,
			Password:           &hashedPw,
			Active:             true,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockShareDg.On("GetCertificateShareByHandle", ctx, pwHandle).Return(share, nil)

		uc := &CertificateShareUsecase{CertificateShareDg: mockShareDg}

		result, err := uc.GetCertificateShareData(ctx, pwHandle, strPtr("wrong-password"))

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrForbidden)
		mockShareDg.AssertExpectations(t)
	})

	t.Run("should return ErrInvalidArgument when password-protected certificate has not been claimed", func(t *testing.T) {
		rawPw := "correct-password"
		hashedPw, err := hashutils.HashPassword(rawPw)
		require.NoError(t, err)

		pwHandle := "pw-handle"
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             pwHandle,
			Password:           &hashedPw,
			Active:             true,
		}
		cert := &entity.EventCertificate{
			Id:                      certID,
			CertificateTokenId:      nil,
			EventCertificateAddress: nil,
		}
		mockShareDg := new(MockCertificateShareDataGateway)
		mockCertDg := new(MockEventCertificateDataGateway)

		mockShareDg.On("GetCertificateShareByHandle", ctx, pwHandle).Return(share, nil)
		mockCertDg.On("GetEventCertificateByID", ctx, certID).Return(cert, nil)

		uc := &CertificateShareUsecase{
			CertificateShareDg:          mockShareDg,
			EventCertificateDataGateway: mockCertDg,
		}

		result, err := uc.GetCertificateShareData(ctx, pwHandle, &rawPw)

		assert.Nil(t, result)
		assertErrCode(t, err, customerror.ErrInvalidArgument)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
	})
}

// newTestKeyPair generates a throwaway ECDSA key pair for testing decryption.
func newTestKeyPair(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	pk, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	require.NoError(t, err)
	return pk
}

// encryptProfile JSON-marshals the profile and ECIES-encrypts it with the given public key.
func encryptProfile(t *testing.T, profile *entity.AttendeeProfileData, pk *ecdsa.PrivateKey) string {
	t.Helper()
	b, err := json.Marshal(profile)
	require.NoError(t, err)
	ciphertext, err := cyptoutils.EncryptWithPublicKeyBytes(string(b), &pk.PublicKey)
	require.NoError(t, err)
	return ciphertext
}

func TestDecryptBackendUserData(t *testing.T) {
	pk := newTestKeyPair(t)

	profile := &entity.AttendeeProfileData{
		FirstName: strPtr("Alice"),
		LastName:  strPtr("Smith"),
		Email:     strPtr("alice@example.com"),
	}
	ciphertext := encryptProfile(t, profile, pk)

	t.Run("should decrypt valid ciphertext and return AttendeeProfileData", func(t *testing.T) {
		uc := &CertificateShareUsecase{BackendPrivateKey: pk}
		result, err := uc.decryptBackendUserData(ciphertext)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Alice", *result.FirstName)
		assert.Equal(t, "Smith", *result.LastName)
		assert.Equal(t, "alice@example.com", *result.Email)
	})

	t.Run("should return nil without error when ciphertext is empty", func(t *testing.T) {
		uc := &CertificateShareUsecase{BackendPrivateKey: pk}
		result, err := uc.decryptBackendUserData("")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return nil without error when private key is not configured", func(t *testing.T) {
		uc := &CertificateShareUsecase{}
		result, err := uc.decryptBackendUserData(ciphertext)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when ciphertext is corrupt", func(t *testing.T) {
		uc := &CertificateShareUsecase{BackendPrivateKey: pk}
		result, err := uc.decryptBackendUserData("not-valid-ciphertext!")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// encryptCsv ECIES-encrypts a raw CSV string with the given public key.
func encryptCsv(t *testing.T, csv string, pk *ecdsa.PrivateKey) string {
	t.Helper()
	ciphertext, err := cyptoutils.EncryptWithPublicKeyBytes(csv, &pk.PublicKey)
	require.NoError(t, err)
	return ciphertext
}

func TestDecryptCertificateRawData(t *testing.T) {
	pk := newTestKeyPair(t)

	t.Run("should decrypt valid CSV and return CertificateRawData with all fields", func(t *testing.T) {
		csv := "Alice Smith,MIT,My Cert,My Subtitle"
		uc := &CertificateShareUsecase{BackendPrivateKey: pk}
		result, err := uc.decryptCertificateRawData(encryptCsv(t, csv, pk))
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Alice Smith", *result.Name)
		assert.Equal(t, "MIT", *result.AcademicInstitution)
		assert.Equal(t, "My Cert", *result.CertificateTitle)
		assert.Equal(t, "My Subtitle", *result.CertificateSubtitle)
	})

	t.Run("should decrypt valid CSV with empty fields", func(t *testing.T) {
		csv := "Alice Smith,,My Cert,"
		uc := &CertificateShareUsecase{BackendPrivateKey: pk}
		result, err := uc.decryptCertificateRawData(encryptCsv(t, csv, pk))
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Alice Smith", *result.Name)
		assert.Equal(t, "", *result.AcademicInstitution)
		assert.Equal(t, "My Cert", *result.CertificateTitle)
		assert.Equal(t, "", *result.CertificateSubtitle)
	})

	t.Run("should return nil without error when ciphertext is empty", func(t *testing.T) {
		uc := &CertificateShareUsecase{BackendPrivateKey: pk}
		result, err := uc.decryptCertificateRawData("")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return nil without error when private key is not configured", func(t *testing.T) {
		csv := "Alice Smith,MIT,My Cert,My Subtitle"
		uc := &CertificateShareUsecase{}
		result, err := uc.decryptCertificateRawData(encryptCsv(t, csv, pk))
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when ciphertext is corrupt", func(t *testing.T) {
		uc := &CertificateShareUsecase{BackendPrivateKey: pk}
		result, err := uc.decryptCertificateRawData("not-valid-ciphertext!")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when CSV does not have exactly 4 fields", func(t *testing.T) {
		csv := "one,two,three"
		uc := &CertificateShareUsecase{BackendPrivateKey: pk}
		result, err := uc.decryptCertificateRawData(encryptCsv(t, csv, pk))
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetCertificateShareDataDecryption(t *testing.T) {
	ctx := context.Background()
	handle := "test-handle"
	certID := uuid.New()
	tokenID := "42"
	contractAddr := "0x1234567890123456789012345678901234567890"

	setupMocks := func(payload *entity.CertificatePayload) (
		*MockCertificateShareDataGateway,
		*MockEventCertificateDataGateway,
		*MockCertificateContractFactoryDg,
		*MockCertificateContractDg,
	) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: certID,
			Handle:             handle,
			Active:             true,
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
		mockContract.On("GetTokenData", ctx, big.NewInt(42)).Return(payload, nil)
		return mockShareDg, mockCertDg, mockFactory, mockContract
	}

	t.Run("should populate DecryptedUserData when private key is configured", func(t *testing.T) {
		pk := newTestKeyPair(t)
		profile := &entity.AttendeeProfileData{
			FirstName: strPtr("Bob"),
			LastName:  strPtr("Jones"),
			Email:     strPtr("bob@example.com"),
		}
		payload := samplePayload()
		payload.Data.BackendEncryptedUserData = encryptProfile(t, profile, pk)

		mockShareDg, mockCertDg, mockFactory, mockContract := setupMocks(payload)
		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
			BackendPrivateKey:            pk,
		}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		require.NoError(t, err)
		require.NotNil(t, result.DecryptedUserData)
		assert.Equal(t, "Bob", *result.DecryptedUserData.FirstName)
		assert.Equal(t, "Jones", *result.DecryptedUserData.LastName)
		assert.Equal(t, "bob@example.com", *result.DecryptedUserData.Email)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should set DecryptedUserData to nil when no private key is configured", func(t *testing.T) {
		payload := samplePayload()
		mockShareDg, mockCertDg, mockFactory, mockContract := setupMocks(payload)
		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
		}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		require.NoError(t, err)
		assert.Nil(t, result.DecryptedUserData)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should set DecryptedUserData to nil when backendEncryptedUserData is empty", func(t *testing.T) {
		pk := newTestKeyPair(t)
		payload := samplePayload() // BackendEncryptedUserData is "" by default
		mockShareDg, mockCertDg, mockFactory, mockContract := setupMocks(payload)
		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
			BackendPrivateKey:            pk,
		}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		require.NoError(t, err)
		assert.Nil(t, result.DecryptedUserData)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should set DecryptedUserData to nil and not fail when ciphertext is corrupt", func(t *testing.T) {
		pk := newTestKeyPair(t)
		payload := samplePayload()
		payload.Data.BackendEncryptedUserData = "this-is-not-valid-ciphertext"
		mockShareDg, mockCertDg, mockFactory, mockContract := setupMocks(payload)
		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
			BackendPrivateKey:            pk,
		}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		require.NoError(t, err)
		assert.Nil(t, result.DecryptedUserData)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should populate DecryptedCertificateData when private key is configured", func(t *testing.T) {
		pk := newTestKeyPair(t)
		csv := "Alice Smith,MIT,My Cert,My Subtitle"
		payload := samplePayload()
		payload.Proof.EncryptedByBackendRawData = encryptCsv(t, csv, pk)

		mockShareDg, mockCertDg, mockFactory, mockContract := setupMocks(payload)
		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
			BackendPrivateKey:            pk,
		}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		require.NoError(t, err)
		require.NotNil(t, result.DecryptedCertificateData)
		assert.Equal(t, "Alice Smith", *result.DecryptedCertificateData.Name)
		assert.Equal(t, "MIT", *result.DecryptedCertificateData.AcademicInstitution)
		assert.Equal(t, "My Cert", *result.DecryptedCertificateData.CertificateTitle)
		assert.Equal(t, "My Subtitle", *result.DecryptedCertificateData.CertificateSubtitle)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should set DecryptedCertificateData to nil when no private key is configured", func(t *testing.T) {
		payload := samplePayload()
		mockShareDg, mockCertDg, mockFactory, mockContract := setupMocks(payload)
		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
		}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		require.NoError(t, err)
		assert.Nil(t, result.DecryptedCertificateData)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should set DecryptedCertificateData to nil when EncryptedByBackendRawData is empty", func(t *testing.T) {
		pk := newTestKeyPair(t)
		payload := samplePayload() // EncryptedByBackendRawData is "" by default
		mockShareDg, mockCertDg, mockFactory, mockContract := setupMocks(payload)
		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
			BackendPrivateKey:            pk,
		}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		require.NoError(t, err)
		assert.Nil(t, result.DecryptedCertificateData)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})

	t.Run("should set DecryptedCertificateData to nil and not fail when ciphertext is corrupt", func(t *testing.T) {
		pk := newTestKeyPair(t)
		payload := samplePayload()
		payload.Proof.EncryptedByBackendRawData = "this-is-not-valid-ciphertext"
		mockShareDg, mockCertDg, mockFactory, mockContract := setupMocks(payload)
		uc := &CertificateShareUsecase{
			CertificateShareDg:           mockShareDg,
			EventCertificateDataGateway:  mockCertDg,
			CertificateContractFactoryDg: mockFactory,
			BackendPrivateKey:            pk,
		}

		result, err := uc.GetCertificateShareData(ctx, handle, nil)

		require.NoError(t, err)
		assert.Nil(t, result.DecryptedCertificateData)
		mockShareDg.AssertExpectations(t)
		mockCertDg.AssertExpectations(t)
		mockFactory.AssertExpectations(t)
		mockContract.AssertExpectations(t)
	})
}
