package certificate_share

import (
	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"context"
	"encoding/json"
	"log/slog"
	"math/big"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

// CertificateContractInfo holds the DB-stored contract metadata for a claimed certificate.
type CertificateContractInfo struct {
	EventCertificateContractAddress string
	CertificateTokenId              string
}

// CertificateShareData combines the on-chain VC payload with the DB-stored contract info,
// plus the backend-decrypted attendee profile when a BackendPrivateKey is configured.
type CertificateShareData struct {
	Payload                  *entity.CertificatePayload
	Contract                 CertificateContractInfo
	DecryptedUserData        *entity.AttendeeProfileData
	DecryptedCertificateData *entity.CertificateRawData
}

// GetCertificateShareData fetches the on-chain VC data for a certificate identified
// by its share handle. Pass a non-nil password for password-protected shares.
//
// Flow:
//  1. Look up the certificate_share row by handle.
//  2. If password-protected and no password supplied, return ErrForbidden.
//  3. If password-protected and password is wrong, return ErrForbidden.
//  4. Resolve the linked EventCertificate from the off-chain DB.
//  5. Require the certificate to be claimed (CertificateTokenId and EventCertificateAddress must be set).
//  6. Obtain the EventCertificate contract via CertificateContractFactoryDg.
//  7. Call GetTokenData(tokenId) — a free view call, no gas required.
func (uc *CertificateShareUsecase) GetCertificateShareData(ctx context.Context, handle string, password *string) (*CertificateShareData, error) {
	share, err := uc.CertificateShareDg.GetCertificateShareByHandle(ctx, handle)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if share == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("certificate share not found"))
	}

	if !share.Active {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("certificate share not found"))
	}

	if err := uc.CheckSharePassword(share, password); err != nil {
		return nil, err
	}

	return uc.fetchOnChainData(ctx, share.EventCertificateId)
}

// CheckSharePassword verifies that the supplied password satisfies the share's
// password requirement. Returns nil if the share is not password-protected or
// if the password matches. Returns ErrForbidden otherwise.
func (uc *CertificateShareUsecase) CheckSharePassword(share *entity.CertificateShare, password *string) error {
	if share.Password == nil {
		return nil
	}
	if password == nil {
		return customerror.Parse(&customerror.ErrForbidden, errors.New("certificate share is password protected"))
	}
	match, err := hashutils.CompareHash(*password, *share.Password)
	if err != nil || !match {
		return customerror.Parse(&customerror.ErrForbidden, errors.New("incorrect password for certificate share"))
	}
	return nil
}

// fetchOnChainData resolves the EventCertificate from the DB and reads its VC
// data from the blockchain via the contract data gateway.
func (uc *CertificateShareUsecase) fetchOnChainData(ctx context.Context, certID uuid.UUID) (*CertificateShareData, error) {
	cert, err := uc.EventCertificateDataGateway.GetEventCertificateByID(ctx, certID)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get certificate"))
	}
	if cert == nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("certificate not found"))
	}

	// Certificate must be claimed (minted on-chain) before its VC data exists
	if cert.CertificateTokenId == nil || cert.EventCertificateAddress == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("certificate has not been claimed yet — on-chain data is not available"))
	}

	contractAddr := common.HexToAddress(*cert.EventCertificateAddress)
	contractDg, err := uc.CertificateContractFactoryDg.GetContract(contractAddr)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to instantiate certificate contract"))
	}

	tokenId, ok := new(big.Int).SetString(*cert.CertificateTokenId, 10)
	if !ok {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Newf("invalid certificate token ID: %s", *cert.CertificateTokenId))
	}

	payload, err := contractDg.GetTokenData(ctx, tokenId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get on-chain certificate data"))
	}

	decrypted, decryptErr := uc.decryptBackendUserData(payload.Data.BackendEncryptedUserData)
	if decryptErr != nil {
		uc.Logger.WarnContext(ctx, "failed to decrypt backend user data", slog.String("error", decryptErr.Error()))
	}
	decryptedCert, decryptCertErr := uc.decryptCertificateRawData(payload.Proof.EncryptedByBackendRawData)
	if decryptCertErr != nil {
		uc.Logger.WarnContext(ctx, "failed to decrypt certificate raw data", slog.String("error", decryptCertErr.Error()))
	}

	return &CertificateShareData{
		Payload: payload,
		Contract: CertificateContractInfo{
			EventCertificateContractAddress: *cert.EventCertificateAddress,
			CertificateTokenId:              *cert.CertificateTokenId,
		},
		DecryptedUserData:        decrypted,
		DecryptedCertificateData: decryptedCert,
	}, nil
}

// decryptBackendUserData decrypts the ECIES-encrypted attendee profile using the backend
// private key. Returns (nil, nil) when the key is not configured or the ciphertext is empty.
// Returns (nil, err) when decryption or JSON unmarshalling fails.
func (uc *CertificateShareUsecase) decryptBackendUserData(ciphertext string) (*entity.AttendeeProfileData, error) {
	if uc.BackendPrivateKey == nil || ciphertext == "" {
		return nil, nil
	}
	plaintext, err := cyptoutils.DecryptWithPrivateKey(ciphertext, uc.BackendPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt backend user data")
	}
	var profile entity.AttendeeProfileData
	if err := json.Unmarshal([]byte(plaintext), &profile); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal decrypted user data")
	}
	return &profile, nil
}

// decryptCertificateRawData decrypts the ECIES-encrypted PII CSV from
// proof.EncryptedByBackendRawData using the backend private key and parses
// it into a CertificateRawData struct.
// The CSV format is: "{name},{academic_institution},{certificate_title},{certificate_subtitle}"
// Note: fields that contain literal commas will break parsing, as the format has no escaping.
// This is a known limitation tied to the on-chain contract's encoding convention.
// Returns (nil, nil) when the key is not configured or the ciphertext is empty.
// Returns (nil, err) when decryption fails or the CSV does not have exactly 4 fields.
func (uc *CertificateShareUsecase) decryptCertificateRawData(ciphertext string) (*entity.CertificateRawData, error) {
	if uc.BackendPrivateKey == nil || ciphertext == "" {
		return nil, nil
	}
	plaintext, err := cyptoutils.DecryptWithPrivateKey(ciphertext, uc.BackendPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt certificate raw data")
	}
	parts := strings.Split(plaintext, ",")
	if len(parts) != 4 {
		return nil, errors.Newf("expected 4 CSV fields in certificate raw data, got %d", len(parts))
	}
	name := parts[0]
	academicInstitution := parts[1]
	certificateTitle := parts[2]
	certificateSubtitle := parts[3]
	return &entity.CertificateRawData{
		Name:                &name,
		AcademicInstitution: &academicInstitution,
		CertificateTitle:    &certificateTitle,
		CertificateSubtitle: &certificateSubtitle,
	}, nil
}
