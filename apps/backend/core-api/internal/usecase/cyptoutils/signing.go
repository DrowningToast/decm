package cyptoutils

import (
	"crypto/ecdsa"
	"fmt"
	"strconv"
	"strings"

	customerror "apps/backend/common/customerror"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

func HashEthereumMessage(message string) ethCommon.Hash {
	message = fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return crypto.Keccak256Hash([]byte(message))
}

func HashMessage(message string) []byte {
	return crypto.Keccak256([]byte(message))
}

// SignTransaction for sign a transaction with private key, which is compatible with Solidity's ECDSA.recover
func Sign(digest []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	signature, err := crypto.Sign(digest, privateKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Adjust the recovery ID (v) to be compatible with Solidity's ECDSA.recover
	// Go-ethereum returns 0/1, but Solidity expects 27/28
	// https://github.com/ethereum/go-ethereum/issues/19751#issuecomment-50490073
	if signature[64] == 0 || signature[64] == 1 {
		signature[64] += 27
	}

	return signature, nil
}

// SignMessage for sign a message w/o Adjust the recovery ID (v)

func GetAddressFromSignature(message string, signature string) (ethCommon.Address, error) {
	hashedMessage := HashEthereumMessage(message)

	// check if signature is valid
	if len(signature)%2 == 1 {
		return ethCommon.Address{}, errors.Wrap(customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid ethereum signature")), "failed to get address from signature")
	}
	sig := hexutil.MustDecode(signature)

	// Reject non-ethereum signature format
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return ethCommon.Address{}, errors.Wrap(customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid ethereum signature")), "failed to get address from signature")
	}

	// Adjust recovery ID from Ethereum format (27/28) to go-ethereum format (0/1)
	sig[crypto.RecoveryIDOffset] -= 27

	usedPublicKey, err := crypto.SigToPub(hashedMessage.Bytes(), sig)
	if err != nil {
		return ethCommon.Address{}, errors.Wrap(customerror.Parse(&customerror.ErrInvalidArgument, err), "failed to recover public key")
	}
	return crypto.PubkeyToAddress(*usedPublicKey), nil
}

func VerifySignedMessageByAddress(walletAddress ethCommon.Address, message string, signature string) (bool, error) {
	hashedMessage := HashEthereumMessage(message)

	sig := hexutil.MustDecode(signature)

	// Reject non-ethereum signature format
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return false, errors.Wrap(customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid ethereum signature")), "failed to verify signed message by address")
	}

	sig[crypto.RecoveryIDOffset] -= 27

	usedPublicKey, err := crypto.SigToPub(hashedMessage.Bytes(), sig)
	if err != nil {
		return false, errors.Wrap(err, "failed to recover public key")
	}

	recoveredAddress := crypto.PubkeyToAddress(*usedPublicKey)

	return recoveredAddress == walletAddress, nil
}

func VerifySignatureByDigest(walletAddress ethCommon.Address, messageDigest ethCommon.Hash, signature []byte) (bool, error) {
	// Make a copy of the signature to avoid mutating the original
	// Slices are reference types, so modifying signature[64] would modify the original array
	signatureCopy := make([]byte, len(signature))
	copy(signatureCopy, signature)

	// Adjust recovery ID from Ethereum format (27/28) to go-ethereum format (0/1)
	if len(signatureCopy) == 65 && (signatureCopy[64] == 27 || signatureCopy[64] == 28) {
		signatureCopy[64] -= 27
	}

	usedPublicKey, err := crypto.SigToPub(messageDigest.Bytes(), signatureCopy)
	if err != nil {
		return false, errors.Wrap(err, "failed to recover public key")
	}

	signatureWallet := crypto.PubkeyToAddress(*usedPublicKey)
	return signatureWallet.Cmp(walletAddress) == 0, nil
}

func VerifySignedMessageByPublicKey(publicKey *ecdsa.PublicKey, message string, signature string) (bool, error) {
	recoveredAddress := crypto.PubkeyToAddress(*publicKey)

	return VerifySignedMessageByAddress(recoveredAddress, message, signature)
}

func GetSignMessage(signerAddress ethCommon.Address, contractAddress ethCommon.Address, deadlineBlock uint64) (string, error) {
	rawSignMessage := fmt.Sprintf("%s,%s,%d", signerAddress.Hex(), contractAddress.Hex(), deadlineBlock)
	return rawSignMessage, nil
}

func ValidateSignMessage(signMessage string, signerAddress ethCommon.Address, contractAddress ethCommon.Address, deadlineBlock *uint64) (bool, error) {
	// split strings into 3 parts
	parts := strings.Split(signMessage, ",")
	if len(parts) != 3 {
		return false, customerror.NewWithPreset(&customerror.ErrInvalidArgument, errors.New("invalid sign message"))
	}

	if parts[0] != signerAddress.Hex() {
		return false, customerror.NewWithPreset(&customerror.ErrInvalidArgument, errors.New("invalid signer address"))
	}

	if parts[1] != contractAddress.Hex() {
		return false, customerror.NewWithPreset(&customerror.ErrInvalidArgument, errors.New("invalid contract address"))
	}

	if deadlineBlock != nil && parts[2] != strconv.FormatUint(*deadlineBlock, 10) {
		return false, customerror.NewWithPreset(&customerror.ErrInvalidArgument, errors.New("invalid deadline block"))
	} else {
		// validate if the value is int or not
		value, err := strconv.ParseUint(parts[2], 10, 64)
		if err != nil {
			return false, customerror.NewWithPreset(&customerror.ErrInvalidArgument, err)
		}
		*deadlineBlock = value
		return true, nil
	}
}

func ExtractDeadlineBlockFromSignMessage(signMessage string) (*uint64, error) {
	parts := strings.Split(signMessage, ",")
	if len(parts) != 3 {
		return nil, customerror.NewWithPreset(&customerror.ErrInvalidArgument, errors.New("invalid sign message"))
	}

	value, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return nil, customerror.NewWithPreset(&customerror.ErrInvalidArgument, err)
	}
	return &value, nil
}
