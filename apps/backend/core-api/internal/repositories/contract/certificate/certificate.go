package contract

import (
	"apps/backend/common/customerror"
	certificateContract "apps/backend/contracts/certificate"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	certificatecontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/certificate_contract"
	"apps/backend/core-api/internal/entity"
	"context"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// tokenDataPrefix is the data URI prefix prepended by the contract's getTokenData / tokenURI functions.
const tokenDataPrefix = "data:application/json;utf8,"

type CertificateContractRepository struct {
	client             *ethclient.Client
	contract           *certificateContract.EventCertificate
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway
}

func NewCertificateContractRepository(client *ethclient.Client, contract *certificateContract.EventCertificate, blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway) *CertificateContractRepository {
	return &CertificateContractRepository{
		client:             client,
		contract:           contract,
		blockchainClientDg: blockchainClientDg,
	}
}

var _ certificatecontract_datagateway.CertificateContractDataGateway = (*CertificateContractRepository)(nil)

func (r *CertificateContractRepository) GetTokenData(ctx context.Context, tokenId *big.Int) (*entity.CertificatePayload, error) {
	rawData, err := r.contract.GetTokenData(&bind.CallOpts{Context: ctx}, tokenId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get token data from contract"))
	}

	jsonStr := strings.TrimPrefix(rawData, tokenDataPrefix)

	var payload entity.CertificatePayload
	if err := json.Unmarshal([]byte(jsonStr), &payload); err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to parse on-chain certificate data"))
	}

	return &payload, nil
}

func (r *CertificateContractRepository) UsedSignatures(ctx context.Context, signature []byte) (bool, error) {
	used, err := r.contract.UsedSignatures(&bind.CallOpts{Context: ctx}, signature)
	if err != nil {
		return false, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to check used signatures"))
	}
	return used, nil
}

func (r *CertificateContractRepository) MintNft(ctx context.Context, params certificatecontract_datagateway.MintNftParams) (*big.Int, error) {
	transactor, err := r.blockchainClientDg.GetTransactOpts(ctx)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get transaction options"))
	}

	issuerProofs := make([]certificateContract.CertificateVCStructsIssuerProof, len(params.IssuerProofs))
	for i, p := range params.IssuerProofs {
		issuerProofs[i] = certificateContract.CertificateVCStructsIssuerProof{
			IssuerSignature: p.IssuerSignature,
			IssuerPublicKey: p.IssuerPublicKey,
		}
	}

	tx, err := r.contract.MintNft(
		transactor,
		params.ReceiverAddress,
		params.UserId,
		params.CertificateId,
		params.IssuerId,
		params.EncryptedUserData,
		params.BackendEncryptedUserData,
		params.IssuerAddresses,
		params.SignedMessageDigest,
		params.Signature,
		params.HostSignature,
		params.HostPublicKey,
		params.SignMessage,
		params.UserEncryptedProof,
		params.BackendEncryptedProof,
		params.CertificateTitle,
		params.CertificateSubtitle,
		params.Hash,
		issuerProofs,
	)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to mint NFT"))
	}

	receipt, err := bind.WaitMined(ctx, r.client, tx)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to wait for mint transaction to be mined"))
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("mint transaction reverted"))
	}

	for _, log := range receipt.Logs {
		event, err := r.contract.ParseCertificateMinted(*log)
		if err == nil {
			return event.TokenId, nil
		}
	}

	return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("CertificateMinted event not found in mint receipt"))
}

func (r *CertificateContractRepository) RevokeCertificate(ctx context.Context, tokenId *big.Int, signedMessageDigest string, signature []byte) error {
	transactor, err := r.blockchainClientDg.GetTransactOpts(ctx)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get transaction options"))
	}

	tx, err := r.contract.RevokeCertificate(transactor, tokenId, signedMessageDigest, signature)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to revoke certificate"))
	}

	receipt, err := bind.WaitMined(ctx, r.client, tx)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to wait for revoke transaction to be mined"))
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return customerror.Parse(&customerror.ErrInternalServer, errors.New("revoke transaction reverted"))
	}

	return nil
}

func (r *CertificateContractRepository) FilterCertificateMinted(ctx context.Context, fromBlock uint64, toBlock uint64, receiverAddresses []common.Address) ([]*certificatecontract_datagateway.CertificateMintedEvent, error) {
	iter, err := r.contract.FilterCertificateMinted(
		&bind.FilterOpts{Start: fromBlock, End: &toBlock, Context: ctx},
		nil,
		receiverAddresses,
	)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to filter CertificateMinted events"))
	}
	defer iter.Close()

	var events []*certificatecontract_datagateway.CertificateMintedEvent
	for iter.Next() {
		e := iter.Event
		events = append(events, &certificatecontract_datagateway.CertificateMintedEvent{
			TokenId:         e.TokenId,
			ReceiverAddress: e.ReceiverAddress,
			CertificateId:   e.CertificateId,
			UserId:          e.UserId,
			IssuerId:        e.IssuerId,
			Raw:             e.Raw,
		})
	}
	if err := iter.Error(); err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "error iterating CertificateMinted events"))
	}

	return events, nil
}

func (r *CertificateContractRepository) ParseCertificateMinted(log types.Log) (*certificatecontract_datagateway.CertificateMintedEvent, error) {
	e, err := r.contract.ParseCertificateMinted(log)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to parse CertificateMinted event"))
	}

	return &certificatecontract_datagateway.CertificateMintedEvent{
		TokenId:         e.TokenId,
		ReceiverAddress: e.ReceiverAddress,
		CertificateId:   e.CertificateId,
		UserId:          e.UserId,
		IssuerId:        e.IssuerId,
		Raw:             e.Raw,
	}, nil
}
