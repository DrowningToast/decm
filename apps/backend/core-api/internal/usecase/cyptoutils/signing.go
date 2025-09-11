package cyptoutils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	customerror "apps/backend/common/customerror"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

func HashEthereumMessage(message string) []byte {
	message = fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return crypto.Keccak256([]byte(message))
}

func Sign(message string, privateKey *ecdsa.PrivateKey) (string, error) {
	hashedMessage := HashEthereumMessage(message)

	signature, err := crypto.Sign(hashedMessage, privateKey)
	if err != nil {
		return "", customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	signature[crypto.RecoveryIDOffset] += 27

	return hex.EncodeToString(signature), nil
}

func VerifySignedMessageByAddress(walletAddress ethCommon.Address, message string, signature string) (bool, error) {
	hashedMessage := HashEthereumMessage(message)

	sig := hexutil.MustDecode(signature)

	// Reject non-ethereum signature format
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return false, errors.New("invalid ethereum signature")
	}

	sig[crypto.RecoveryIDOffset] -= 27

	usedPublicKey, err := crypto.SigToPub(hashedMessage, sig)
	if err != nil {
		return false, errors.Wrap(err, "failed to recover public key")
	}

	recoveredAddress := crypto.PubkeyToAddress(*usedPublicKey)

	return recoveredAddress == walletAddress, nil
}

func VerifySignedMessageByPublicKey(publicKey *ecdsa.PublicKey, message string, signature string) (bool, error) {
	recoveredAddress := crypto.PubkeyToAddress(*publicKey)

	return VerifySignedMessageByAddress(recoveredAddress, message, signature)
}
