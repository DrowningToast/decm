package cyptoutils

import (
	"crypto/ecdsa"
	"encoding/hex"

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
