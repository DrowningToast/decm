package certificate_share

import (
	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/entity"
	"context"
	"math/big"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

// CertificateShareData is the fully typed representation of the on-chain VC JSON
// returned by the EventCertificate contract's getTokenData(tokenId) view function.
type CertificateShareData = entity.CertificatePayload

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

	if share.Password != nil {
		if password == nil {
			return nil, customerror.Parse(&customerror.ErrForbidden, errors.New("certificate share is password protected"))
		}
		match, err := hashutils.CompareHash(*password, *share.Password)
		if err != nil || !match {
			return nil, customerror.Parse(&customerror.ErrForbidden, errors.New("incorrect password for certificate share"))
		}
	}

	return uc.fetchOnChainData(ctx, share.EventCertificateId)
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

	tokenId := new(big.Int)
	tokenId.SetString(*cert.CertificateTokenId, 10)

	payload, err := contractDg.GetTokenData(ctx, tokenId)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
