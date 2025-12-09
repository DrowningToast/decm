package cyptoutils

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"strings"

	"apps/backend/common/encryptutils"

	"github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func GetAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) (*ethCommon.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &address, nil
}

func GetPublicKeyFromPrivateKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key")
	}
	return publicKeyECDSA, nil
}

// RecoverPublicKeyFromSignature recovers the public key from a signature and message hash
// This is used in wallet extension flow where user signs a message
func RecoverPublicKeyFromSignature(messageHash ethCommon.Hash, signature []byte) (*ecdsa.PublicKey, error) {
	// Make a copy to avoid modifying the original signature
	sig := make([]byte, len(signature))
	copy(sig, signature)

	// Adjust recovery ID from Ethereum format (27/28) to go-ethereum format (0/1) if needed
	if len(sig) == 65 && (sig[64] == 27 || sig[64] == 28) {
		sig[64] -= 27
	}

	// Recover public key from signature
	publicKey, err := crypto.SigToPub(messageHash.Bytes(), sig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to recover public key from signature")
	}

	return publicKey, nil
}

// PublicKeyToAddress converts an ECDSA public key to an Ethereum address
func PublicKeyToAddress(publicKey *ecdsa.PublicKey) ethCommon.Address {
	return crypto.PubkeyToAddress(*publicKey)
}

func ParsePrivateKey(privateKey string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, errors.New("invalid private key")
	}
	return crypto.ToECDSA(privateKeyBytes)
}

func ParseAddress(address string) (*ethCommon.Address, error) {
	if !ethCommon.IsHexAddress(address) {
		return nil, errors.New("invalid address")
	}

	parsedAddress := ethCommon.HexToAddress(address)
	return &parsedAddress, nil
}

func DecryptPrivateKey(ciphertext string, password string) (*ecdsa.PrivateKey, *ethCommon.Address, error) {
	hostDecryptedPk, err := encryptutils.DecryptAESGCM(ciphertext, password)
	if err != nil {
		return nil, nil, err
	}

	hostParsedPk, err := ParsePrivateKey(hostDecryptedPk)
	if err != nil {
		return nil, nil, err
	}

	hostAddress, err := GetAddressFromPrivateKey(hostParsedPk)
	if err != nil {
		return nil, nil, err
	}

	parsedAddress, err := ParseAddress(hostAddress.Hex())
	if err != nil {
		return nil, nil, err
	}

	return hostParsedPk, parsedAddress, nil
}

// GetRevertReason attempts to extract the revert reason from a failed transaction
func GetRevertReason(ctx context.Context, client *ethclient.Client, tx *types.Transaction, receipt *types.Receipt) (string, error) {
	if receipt.Status == types.ReceiptStatusSuccessful {
		return "", nil
	}

	// Get the transaction sender
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get transaction sender")
	}

	// Call the contract to get the revert reason
	msg := ethereum.CallMsg{
		From:      from,
		To:        tx.To(),
		Gas:       tx.Gas(),
		GasPrice:  tx.GasPrice(),
		GasFeeCap: tx.GasFeeCap(),
		GasTipCap: tx.GasTipCap(),
		Value:     tx.Value(),
		Data:      tx.Data(),
	}

	// Try to call at the block where transaction failed
	_, err = client.CallContract(ctx, msg, receipt.BlockNumber)
	if err != nil {
		// Parse the error message
		errMsg := err.Error()
		// Common patterns for revert reasons
		if strings.Contains(errMsg, "execution reverted:") {
			parts := strings.SplitN(errMsg, "execution reverted:", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
		return errMsg, nil
	}

	return "unknown revert reason", nil
}

// DecodeRevertReason attempts to decode the revert reason from error data
func DecodeRevertReason(errData []byte) string {
	if len(errData) < 4 {
		return "insufficient data"
	}

	// Check for Error(string) selector (0x08c379a0)
	errorSelector := []byte{0x08, 0xc3, 0x79, 0xa0}
	if len(errData) >= 4 && hex.EncodeToString(errData[:4]) == hex.EncodeToString(errorSelector) {
		// Decode the string parameter
		if len(errData) < 68 {
			return "invalid error data length"
		}
		// Skip selector (4 bytes) + offset (32 bytes) + length (32 bytes)
		lengthBytes := errData[36:68]
		length := new(big.Int).SetBytes(lengthBytes).Uint64()
		if len(errData) < int(68+length) {
			return "truncated error data"
		}
		reason := string(errData[68 : 68+length])
		return reason
	}

	// Check for custom error (4 bytes selector)
	selector := hex.EncodeToString(errData[:4])

	// Map known custom errors from Event contract
	customErrors := map[string]string{
		"8e4a23d6": "Event__InvalidEventName()",
		"4af0bf99": "Event__CannotReduceSeatsCount()",
		"85ab58b1": "Event__SeatsCountReached()",
		"d773224f": "Event__ParticipantIsNotJoined()",
		"de3f615d": "Event__ParticipantIsAlreadyJoined()",
		"d92e233d": "Event__AddressCannotBeZero()",
		"8baa579f": "Event__InvalidSignature()",
	}

	if errorName, ok := customErrors[selector]; ok {
		return errorName
	}

	return "unknown error (0x" + selector + ")"
}
