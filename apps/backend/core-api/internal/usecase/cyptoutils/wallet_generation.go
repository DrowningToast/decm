package cyptoutils

import (
	"apps/backend/common/customerror"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func GenerateMnemonic(wordsCount *int) (*string, error) {
	wordCount := 12
	if wordsCount != nil {
		wordCount = *wordsCount
	}

	// Validate word count
	validCounts := []int{12, 15, 18, 21, 24}
	valid := false
	for _, count := range validCounts {
		if wordCount == count {
			valid = true
			break
		}
	}
	if !valid {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("invalid word count: %d (must be 12, 15, 18, 21, or 24)", wordCount))
	}

	// Generate entropy and mnemonic
	var entropyBits int
	switch wordCount {
	case 12:
		entropyBits = 128
	case 15:
		entropyBits = 160
	case 18:
		entropyBits = 192
	case 21:
		entropyBits = 224
	case 24:
		entropyBits = 256
	}

	entropy, _ := bip39.NewEntropy(entropyBits)

	// Generate mnemonic
	mnemonic, _ := bip39.NewMnemonic(entropy)

	return &mnemonic, nil
}

// GenerateSeedPhrase generates a seed phrase for a given word count
// Default is 12 words
func GenerateSeed(wordsCount *int) ([]byte, error) {
	mnemonic, err := GenerateMnemonic(wordsCount)
	if err != nil {
		return nil, err
	}

	seed := bip39.NewSeed(*mnemonic, "")
	return seed, nil
}

func GenerateSeedFromMnemonic(mnemonic *string) ([]byte, error) {
	seed := bip39.NewSeed(*mnemonic, "")
	return seed, nil
}

func GeneratePrivateKeyFromSeed(seed []byte) (*ecdsa.PrivateKey, error) {
	// Generate master key from seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err).Extend("failed to create master key")
	}

	// Parse derivation path (simplified - only supports m/44'/60'/0'/0/0 format)
	// For Ethereum: m/44'/60'/0'/0/0
	purposeKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err).Extend("failed to create purpose key")
	}

	coinKey, err := purposeKey.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err).Extend("failed to create coin key")
	}

	accountKey, err := coinKey.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err).Extend("failed to create account key")
	}

	changeKey, err := accountKey.NewChildKey(0)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err).Extend("failed to create change key")
	}

	addressKey, err := changeKey.NewChildKey(0)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err).Extend("failed to create address key")
	}

	// Convert to ECDSA private key
	privateKey, err := crypto.ToECDSA(addressKey.Key)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err).Extend("failed to create private key")
	}

	return privateKey, nil
}
