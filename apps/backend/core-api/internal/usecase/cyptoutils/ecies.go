package cyptoutils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// EncryptWithPublicKey encrypts data using ECIES with the provided Ethereum public key
// Returns base64-encoded ciphertext
func EncryptWithPublicKey(plaintext string, walletAddress common.Address) (string, error) {
	// For Ethereum addresses, we need to get the public key from the address
	// Since we can't derive public key from address alone, we'll use a hybrid approach
	// This function expects the full public key to be derivable from context

	// NOTE: In practice, you should pass the actual public key, not just the address
	// For now, we'll document that this requires the full public key
	return "", fmt.Errorf("EncryptWithPublicKey requires full public key, not just address - use EncryptWithPublicKeyBytes instead")
}

// EncryptWithPublicKeyBytes encrypts data using ECIES with the provided ECDSA public key
// Returns base64-encoded ciphertext that can only be decrypted by the corresponding private key
func EncryptWithPublicKeyBytes(plaintext string, publicKey *ecdsa.PublicKey) (string, error) {
	if publicKey == nil {
		return "", fmt.Errorf("public key is nil")
	}

	// Convert ECDSA public key to ECIES public key
	eciesPublicKey := ecies.ImportECDSAPublic(publicKey)
	if eciesPublicKey == nil {
		return "", fmt.Errorf("failed to import ECDSA public key to ECIES")
	}

	// Encrypt the plaintext
	ciphertext, err := ecies.Encrypt(rand.Reader, eciesPublicKey, []byte(plaintext), nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt with ECIES: %w", err)
	}

	// Encode to base64 for storage/transmission
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptWithPrivateKey decrypts ECIES-encrypted data using the provided private key
// Accepts base64-encoded ciphertext
func DecryptWithPrivateKey(ciphertext string, privateKey *ecdsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", fmt.Errorf("private key is nil")
	}

	// Decode from base64
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	// Convert ECDSA private key to ECIES private key
	eciesPrivateKey := ecies.ImportECDSA(privateKey)
	if eciesPrivateKey == nil {
		return "", fmt.Errorf("failed to import ECDSA private key to ECIES")
	}

	// Decrypt the ciphertext
	plaintext, err := eciesPrivateKey.Decrypt(ciphertextBytes, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt with ECIES: %w", err)
	}

	return string(plaintext), nil
}

// GetPublicKeyFromAddress retrieves the public key for a given wallet address
// NOTE: This requires looking up the public key from a credential store or database
// as public keys cannot be derived from Ethereum addresses alone
func GetPublicKeyFromAddress(walletAddress common.Address) (*ecdsa.PublicKey, error) {
	// This function should be implemented to fetch the public key from storage
	// For now, return an error indicating this needs to be implemented
	return nil, fmt.Errorf("GetPublicKeyFromAddress must fetch public key from credential storage")
}
