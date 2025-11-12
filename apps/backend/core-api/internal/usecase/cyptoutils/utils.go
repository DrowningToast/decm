package cyptoutils

import (
	"crypto/ecdsa"
	"encoding/hex"

	"apps/backend/common/encryptutils"

	"github.com/ethereum/go-ethereum/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

func DecryptPrivateKey(ciphertext string, password string) (*ecdsa.PrivateKey, *common.Address, error) {
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
