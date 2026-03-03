package certificatecontract_datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type IssuerProof struct {
	IssuerSignature string
	IssuerPublicKey string
}

type MintNftParams struct {
	ReceiverAddress          common.Address
	UserId                   string
	CertificateId            string
	IssuerId                 string
	EncryptedUserData        string
	BackendEncryptedUserData string
	IssuerAddresses          []common.Address
	SignedMessageDigest      string // participant's sign message (contract EIP-191 hashes this)
	Signature                []byte // participant's 65-byte ECDSA sig (v = 27 or 28)
	HostSignature            string // host ECDSA sig hex (0x-prefixed)
	HostPublicKey            string
	SignMessage              string // original JSON from event_certificate_signatures
	UserEncryptedProof       string // PII CSV ECIES-encrypted with user's public key
	BackendEncryptedProof    string // PII CSV ECIES-encrypted with backend's public key
	CertificateTitle         string
	CertificateSubtitle      string
	Hash                     string // Keccak256 of PII CSV (= certificate_digest in DB)
	IssuerProofs             []IssuerProof
}

type CertificateMintedEvent struct {
	TokenId         *big.Int
	ReceiverAddress common.Address
	CertificateId   string
	UserId          string
	IssuerId        string
	Raw             types.Log
}

type CertificateContractDataGateway interface {
	// View calls (free, no gas)
	GetTokenData(ctx context.Context, tokenId *big.Int) (*entity.CertificatePayload, error)
	UsedSignatures(ctx context.Context, signature []byte) (bool, error)

	// Write calls (gas required, blocks until mined)
	MintNft(ctx context.Context, params MintNftParams) (*big.Int, error) // returns minted tokenId
	RevokeCertificate(ctx context.Context, tokenId *big.Int, signedMessageDigest string, signature []byte) error

	// Event log operations
	FilterCertificateMinted(ctx context.Context, fromBlock uint64, toBlock uint64, receiverAddresses []common.Address) ([]*CertificateMintedEvent, error)
	ParseCertificateMinted(log types.Log) (*CertificateMintedEvent, error)
}
