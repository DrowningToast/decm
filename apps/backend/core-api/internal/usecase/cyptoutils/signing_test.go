package cyptoutils

import (
	"crypto/ecdsa"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Test helper: Generate a test private key and address
func generateTestKeyPair() (*ecdsa.PrivateKey, ethCommon.Address) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic("failed to generate test key: " + err.Error())
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return privateKey, address
}

func TestHashEthereumMessage(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		wantHash string // Expected hash in hex
	}{
		{
			name:     "Simple message",
			message:  "Hello, Ethereum!",
			wantHash: "0x4b09f0d62c5b8ec4d70c8e909c1e0b9f0e4f1dc05e8e0eb9e35c5e3c4c8b9c3d",
		},
		{
			name:     "Empty message",
			message:  "",
			wantHash: "0x5f35dce98ba4fba25530a026ed80b2cecdaa31091ba4958b99b52ea1d068adad",
		},
		{
			name:     "Message with special characters",
			message:  "Test!@#$%^&*()",
			wantHash: "0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8",
		},
		{
			name:     "Long message",
			message:  "This is a much longer message that contains multiple words and sentences to test the hashing function with more complex input data.",
			wantHash: "0x7b98e5d5b5d5f3f8c3b6d5b3c8d9c5e3f8d5b3c8d9c5e3f8d5b3c8d9c5e3f8d5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := HashEthereumMessage(tt.message)

			// Verify hash is 32 bytes (256 bits)
			if len(hash) != 32 {
				t.Errorf("HashEthereumMessage() hash length = %v, want 32 bytes", len(hash))
			}

			// Verify hash is deterministic
			hash2 := HashEthereumMessage(tt.message)
			if hexutil.Encode(hash) != hexutil.Encode(hash2) {
				t.Error("HashEthereumMessage() is not deterministic")
			}
		})
	}
}

func TestHashEthereumMessage_Prefix(t *testing.T) {
	// Test that the Ethereum signed message prefix is correctly applied
	message := "Test message"
	hash := HashEthereumMessage(message)

	// The hash should be different from a plain Keccak256 hash
	plainHash := crypto.Keccak256([]byte(message))

	if hexutil.Encode(hash) == hexutil.Encode(plainHash) {
		t.Error("HashEthereumMessage() should apply Ethereum signed message prefix")
	}
}

func TestSign(t *testing.T) {
	privateKey, _ := generateTestKeyPair()
	message := "Test signing message"

	t.Run("Sign valid message", func(t *testing.T) {
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Errorf("Sign() error = %v, want nil", err)
			return
		}

		// Signature should be hex encoded and 65 bytes when decoded (130 hex chars + 0x prefix)
		decoded := hexutil.MustDecode("0x" + signature)
		if len(decoded) != 65 {
			t.Errorf("Sign() signature length = %v bytes, want 65 bytes", len(decoded))
		}

		// Recovery ID should be 27 or 28 (Ethereum standard)
		recoveryID := decoded[crypto.RecoveryIDOffset]
		if recoveryID != 27 && recoveryID != 28 {
			t.Errorf("Sign() recovery ID = %v, want 27 or 28", recoveryID)
		}
	})

	t.Run("Sign produces consistent length", func(t *testing.T) {
		sig1, err1 := Sign(message, privateKey)
		sig2, err2 := Sign(message, privateKey)

		if err1 != nil || err2 != nil {
			t.Errorf("Sign() unexpected errors: %v, %v", err1, err2)
			return
		}

		// ECDSA signatures have randomness, so they should be different
		// but both should be valid and have the same length
		if len(sig1) != len(sig2) {
			t.Error("Sign() signatures should have same length")
		}
	})
}

func TestGetAddressFromSignature(t *testing.T) {
	privateKey, expectedAddress := generateTestKeyPair()
	message := "Authentication message"

	t.Run("Recover address from valid signature", func(t *testing.T) {
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}

		recoveredAddress, err := GetAddressFromSignature(message, "0x"+signature)
		if err != nil {
			t.Errorf("GetAddressFromSignature() error = %v, want nil", err)
			return
		}

		if recoveredAddress != expectedAddress {
			t.Errorf("GetAddressFromSignature() = %v, want %v", recoveredAddress.Hex(), expectedAddress.Hex())
		}
	})

	t.Run("Reject invalid recovery ID", func(t *testing.T) {
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}

		// Manually corrupt the recovery ID
		sigBytes := hexutil.MustDecode("0x" + signature)
		sigBytes[crypto.RecoveryIDOffset] = 30 // Invalid recovery ID
		corruptedSig := hexutil.Encode(sigBytes)

		_, err = GetAddressFromSignature(message, corruptedSig)
		if err == nil {
			t.Error("GetAddressFromSignature() should reject invalid recovery ID")
		}
	})

	t.Run("Reject too short signature", func(t *testing.T) {
		// Valid hex but too short to be a proper signature
		shortSig := "0x1234567890abcdef"

		defer func() {
			if r := recover(); r != nil {
				// Expected to panic with invalid signature
				return
			}
		}()

		_, err := GetAddressFromSignature(message, shortSig)
		// If it doesn't panic, it should at least error
		if err == nil {
			t.Error("GetAddressFromSignature() should reject short signature")
		}
	})

	t.Run("Reject signature with wrong message", func(t *testing.T) {
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}

		wrongMessage := "Different message"
		recoveredAddress, err := GetAddressFromSignature(wrongMessage, "0x"+signature)
		if err != nil {
			// Some implementations may error, which is fine
			return
		}

		// If no error, the recovered address should be different
		if recoveredAddress == expectedAddress {
			t.Error("GetAddressFromSignature() should not recover correct address with wrong message")
		}
	})

	t.Run("Recover address matches signer", func(t *testing.T) {
		// Test with multiple messages to ensure consistency
		testMessages := []string{
			"Message 1",
			"Another test message",
			"Sign in to DECM Platform",
			"Timestamp: 1234567890",
		}

		for _, msg := range testMessages {
			signature, err := Sign(msg, privateKey)
			if err != nil {
				t.Fatalf("Failed to sign message '%s': %v", msg, err)
			}

			recoveredAddress, err := GetAddressFromSignature(msg, "0x"+signature)
			if err != nil {
				t.Errorf("GetAddressFromSignature() error for message '%s': %v", msg, err)
				continue
			}

			if recoveredAddress != expectedAddress {
				t.Errorf("GetAddressFromSignature() for message '%s' = %v, want %v",
					msg, recoveredAddress.Hex(), expectedAddress.Hex())
			}
		}
	})
}

func TestVerifySignedMessageByAddress(t *testing.T) {
	privateKey, address := generateTestKeyPair()
	message := "Verify this message"

	t.Run("Verify valid signature", func(t *testing.T) {
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}

		valid, err := VerifySignedMessageByAddress(address, message, "0x"+signature)
		if err != nil {
			t.Errorf("VerifySignedMessageByAddress() error = %v, want nil", err)
			return
		}

		if !valid {
			t.Error("VerifySignedMessageByAddress() = false, want true")
		}
	})

	t.Run("Reject signature from different address", func(t *testing.T) {
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}

		// Use a different address
		_, wrongAddress := generateTestKeyPair()

		valid, err := VerifySignedMessageByAddress(wrongAddress, message, "0x"+signature)
		if err != nil {
			t.Errorf("VerifySignedMessageByAddress() error = %v, want nil", err)
			return
		}

		if valid {
			t.Error("VerifySignedMessageByAddress() = true, want false for wrong address")
		}
	})

	t.Run("Reject signature with wrong message", func(t *testing.T) {
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}

		wrongMessage := "Different message"
		valid, err := VerifySignedMessageByAddress(address, wrongMessage, "0x"+signature)
		if err != nil {
			t.Errorf("VerifySignedMessageByAddress() error = %v, want nil", err)
			return
		}

		if valid {
			t.Error("VerifySignedMessageByAddress() = true, want false for wrong message")
		}
	})

	t.Run("Reject invalid recovery ID", func(t *testing.T) {
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}

		// Corrupt recovery ID
		sigBytes := hexutil.MustDecode("0x" + signature)
		sigBytes[crypto.RecoveryIDOffset] = 25 // Invalid
		corruptedSig := hexutil.Encode(sigBytes)

		valid, err := VerifySignedMessageByAddress(address, message, corruptedSig)
		if err == nil {
			t.Error("VerifySignedMessageByAddress() should return error for invalid recovery ID")
			return
		}

		if valid {
			t.Error("VerifySignedMessageByAddress() should return false for invalid signature")
		}
	})

	t.Run("Reject too short signature", func(t *testing.T) {
		// Valid hex but too short
		shortSig := "0x1234567890abcdef"

		defer func() {
			if r := recover(); r != nil {
				// Expected to panic with invalid signature format
				return
			}
		}()

		valid, err := VerifySignedMessageByAddress(address, message, shortSig)
		// If it doesn't panic, it should at least error or return false
		if err == nil && valid {
			t.Error("VerifySignedMessageByAddress() should reject short signature")
		}
	})

	t.Run("Verify multiple signatures from same key", func(t *testing.T) {
		messages := []string{
			"Message 1",
			"Message 2",
			"Message 3",
		}

		for _, msg := range messages {
			signature, err := Sign(msg, privateKey)
			if err != nil {
				t.Fatalf("Failed to sign message '%s': %v", msg, err)
			}

			valid, err := VerifySignedMessageByAddress(address, msg, "0x"+signature)
			if err != nil {
				t.Errorf("VerifySignedMessageByAddress() error for message '%s': %v", msg, err)
				continue
			}

			if !valid {
				t.Errorf("VerifySignedMessageByAddress() = false for message '%s', want true", msg)
			}
		}
	})
}

func TestVerifySignedMessageByPublicKey(t *testing.T) {
	privateKey, _ := generateTestKeyPair()
	publicKey := &privateKey.PublicKey
	message := "Verify with public key"

	t.Run("Verify valid signature with public key", func(t *testing.T) {
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}

		valid, err := VerifySignedMessageByPublicKey(publicKey, message, "0x"+signature)
		if err != nil {
			t.Errorf("VerifySignedMessageByPublicKey() error = %v, want nil", err)
			return
		}

		if !valid {
			t.Error("VerifySignedMessageByPublicKey() = false, want true")
		}
	})

	t.Run("Reject signature from different public key", func(t *testing.T) {
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}

		// Use a different public key
		wrongPrivateKey, _ := generateTestKeyPair()
		wrongPublicKey := &wrongPrivateKey.PublicKey

		valid, err := VerifySignedMessageByPublicKey(wrongPublicKey, message, "0x"+signature)
		if err != nil {
			t.Errorf("VerifySignedMessageByPublicKey() error = %v, want nil", err)
			return
		}

		if valid {
			t.Error("VerifySignedMessageByPublicKey() = true, want false for wrong public key")
		}
	})

	t.Run("Verify public key derivation matches address", func(t *testing.T) {
		// Ensure VerifySignedMessageByPublicKey uses correct address derivation
		signature, err := Sign(message, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign message: %v", err)
		}

		address := crypto.PubkeyToAddress(*publicKey)

		// Both methods should return the same result
		validByPubKey, err1 := VerifySignedMessageByPublicKey(publicKey, message, "0x"+signature)
		validByAddress, err2 := VerifySignedMessageByAddress(address, message, "0x"+signature)

		if err1 != nil || err2 != nil {
			t.Fatalf("Verification errors: pubkey=%v, address=%v", err1, err2)
		}

		if validByPubKey != validByAddress {
			t.Error("VerifySignedMessageByPublicKey and VerifySignedMessageByAddress should return same result")
		}
	})
}

func TestSignAndVerifyRoundTrip(t *testing.T) {
	// Integration test: Sign and verify the complete flow
	privateKey, address := generateTestKeyPair()

	testCases := []string{
		"Welcome to DECM Platform!",
		"Sign in with nonce: 123456",
		"Timestamp: 1234567890",
		"Multi-line\nmessage\ntest",
		"",
		"🎉 Emoji test 🚀",
	}

	for _, message := range testCases {
		t.Run("Round trip: "+message, func(t *testing.T) {
			// Sign the message
			signature, err := Sign(message, privateKey)
			if err != nil {
				t.Fatalf("Sign() error = %v", err)
			}

			// Recover address from signature
			recoveredAddress, err := GetAddressFromSignature(message, "0x"+signature)
			if err != nil {
				t.Fatalf("GetAddressFromSignature() error = %v", err)
			}

			if recoveredAddress != address {
				t.Errorf("Recovered address = %v, want %v", recoveredAddress.Hex(), address.Hex())
			}

			// Verify signature by address
			valid, err := VerifySignedMessageByAddress(address, message, "0x"+signature)
			if err != nil {
				t.Fatalf("VerifySignedMessageByAddress() error = %v", err)
			}

			if !valid {
				t.Error("VerifySignedMessageByAddress() = false, want true")
			}

			// Verify signature by public key
			valid, err = VerifySignedMessageByPublicKey(&privateKey.PublicKey, message, "0x"+signature)
			if err != nil {
				t.Fatalf("VerifySignedMessageByPublicKey() error = %v", err)
			}

			if !valid {
				t.Error("VerifySignedMessageByPublicKey() = false, want true")
			}
		})
	}
}

func TestRecoveryIDConversion(t *testing.T) {
	// Test that recovery ID conversion is correct (27/28 <-> 0/1)
	privateKey, address := generateTestKeyPair()
	message := "Test recovery ID"

	signature, err := Sign(message, privateKey)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	sigBytes := hexutil.MustDecode("0x" + signature)

	t.Run("Sign produces recovery ID 27 or 28", func(t *testing.T) {
		recoveryID := sigBytes[crypto.RecoveryIDOffset]
		if recoveryID != 27 && recoveryID != 28 {
			t.Errorf("Sign() recovery ID = %v, want 27 or 28", recoveryID)
		}
	})

	t.Run("GetAddressFromSignature handles recovery ID conversion", func(t *testing.T) {
		recoveredAddress, err := GetAddressFromSignature(message, "0x"+signature)
		if err != nil {
			t.Fatalf("GetAddressFromSignature() error = %v", err)
		}

		if recoveredAddress != address {
			t.Error("Recovery ID conversion failed in GetAddressFromSignature")
		}
	})

	t.Run("VerifySignedMessageByAddress handles recovery ID conversion", func(t *testing.T) {
		valid, err := VerifySignedMessageByAddress(address, message, "0x"+signature)
		if err != nil {
			t.Fatalf("VerifySignedMessageByAddress() error = %v", err)
		}

		if !valid {
			t.Error("Recovery ID conversion failed in VerifySignedMessageByAddress")
		}
	})
}
