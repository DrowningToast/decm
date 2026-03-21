package contract

import (
	"apps/backend/common/customerror"
	certificateContract "apps/backend/contracts/certificate"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	certificatecontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/certificate_contract"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/log"
	"context"
	"encoding/json"
	"log/slog"
	"math/big"
	"regexp"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// tokenDataPrefix is the data URI prefix prepended by the contract's getTokenData / tokenURI functions.
const tokenDataPrefix = "data:application/json;utf8,"

// issuerAddressesRe matches the broken issuerAddresses value emitted by the contract's
// _addressArrayToString helper, which wraps a JSON array (with per-element quotes) inside
// an outer quoted string:
//
//	"issuerAddresses": "["0xabc...","0xdef..."]"
//
// The capture group holds the content between the square brackets.
var issuerAddressesRe = regexp.MustCompile(`"issuerAddresses":\s*"\[([^\]]*)\]"`)

// fixIssuerAddresses repairs the malformed issuerAddresses JSON produced by the contract.
// The contract embeds a JSON array literal inside a quoted string value, which causes
// the embedded double-quotes to break JSON parsing. This converts it to the expected
// comma-separated string form:
//
//	broken: "issuerAddresses": "["0xabc...","0xdef..."]"
//	fixed:  "issuerAddresses": "0xabc...,0xdef..."
func fixIssuerAddresses(jsonStr string) string {
	return issuerAddressesRe.ReplaceAllStringFunc(jsonStr, func(match string) string {
		sub := issuerAddressesRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		// sub[1] is the content inside the brackets, e.g. `"0xabc...","0xdef..."`
		// Strip the per-element quotes to produce the comma-separated form.
		inner := strings.ReplaceAll(sub[1], `"`, "")
		return `"issuerAddresses": "` + inner + `"`
	})
}

// fixFreeTextField re-encodes a free-text string field value that the contract
// embedded directly into JSON without escaping. The contract never escapes special
// characters (e.g. `"`, `\`, control chars) inside string values, so any of those
// will produce invalid JSON. This function finds the raw value using the known
// surrounding field names as delimiters and replaces it with a properly
// JSON-encoded version.
//
// startMarker is the literal prefix including the opening quote of the value
// (e.g. `"eventName": "`). endMarker is the literal suffix that begins immediately
// after the closing quote of the value (e.g. `","eventDescription":`).
func fixFreeTextField(jsonStr, startMarker, endMarker string) string {
	startIdx := strings.Index(jsonStr, startMarker)
	if startIdx == -1 {
		return jsonStr
	}
	valueStart := startIdx + len(startMarker)

	endIdx := strings.Index(jsonStr[valueStart:], endMarker)
	if endIdx == -1 {
		return jsonStr
	}
	endIdx += valueStart

	rawValue := jsonStr[valueStart:endIdx]

	// Re-encode the raw value as a proper JSON string interior (escapes ", \, control chars).
	encoded, err := json.Marshal(rawValue)
	if err != nil {
		return jsonStr
	}
	// json.Marshal wraps the string in quotes; strip them to get only the interior.
	escapedValue := string(encoded[1 : len(encoded)-1])

	return jsonStr[:valueStart] + escapedValue + jsonStr[endIdx:]
}

// fixFreeTextFields fixes all free-text string fields in the on-chain JSON that may
// contain characters the contract did not escape (most commonly unescaped `"`).
func fixFreeTextFields(jsonStr string) string {
	jsonStr = fixFreeTextField(jsonStr, `"eventName": "`, `","eventDescription":`)
	jsonStr = fixFreeTextField(jsonStr, `"eventDescription": "`, `","certificateTokenId":`)
	jsonStr = fixFreeTextField(jsonStr, `"certificateTitle": "`, `","certificateSubtitle":`)
	jsonStr = fixFreeTextField(jsonStr, `"certificateSubtitle": "`, `","status":`)
	// signMessage is stored as a raw JSON object string by the contract without escaping.
	// e.g. "signMessage": "{"eventContractAddress":"0x...","receivers":["0x..."]}",
	jsonStr = fixFreeTextField(jsonStr, `"signMessage": "`, `","host":`)
	return jsonStr
}

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

	logger := log.FromContext(ctx)

	jsonStr := strings.TrimPrefix(rawData, tokenDataPrefix)
	logger.Debug("GetTokenData: raw JSON from contract", slog.String("raw_json", jsonStr))

	jsonStr = fixIssuerAddresses(jsonStr)
	jsonStr = fixFreeTextFields(jsonStr)
	logger.Debug("GetTokenData: JSON after fixes", slog.String("fixed_json", jsonStr))

	var payload entity.CertificatePayload
	if err := json.Unmarshal([]byte(jsonStr), &payload); err != nil {
		logger.Error("GetTokenData: failed to parse JSON",
			slog.String("error", err.Error()),
			slog.String("raw_json", rawData),
			slog.String("fixed_json", jsonStr),
		)
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
