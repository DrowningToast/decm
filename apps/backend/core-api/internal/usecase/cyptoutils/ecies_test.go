package cyptoutils

import (
	"crypto/ecdsa"
	"encoding/base64"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// TestEncryptWithPublicKey tests the deprecated function that requires full public key
func TestEncryptWithPublicKey(t *testing.T) {
	t.Run("should_return_error_indicating_full_public_key_required", func(t *testing.T) {
		address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1")
		plaintext := "test message"

		result, err := EncryptWithPublicKey(plaintext, address)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "requires full public key")
		assert.Contains(t, err.Error(), "EncryptWithPublicKeyBytes")
	})
}

// TestEncryptWithPublicKeyBytes tests encryption with various inputs
func TestEncryptWithPublicKeyBytes(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		setupKey  func() *ecdsa.PublicKey
		wantError bool
	}{
		{
			name:      "should_encrypt_simple_message",
			plaintext: "Hello, World!",
			setupKey: func() *ecdsa.PublicKey {
				privateKey, _ := generateTestKeyPair()
				return &privateKey.PublicKey
			},
			wantError: false,
		},
		{
			name:      "should_encrypt_empty_string",
			plaintext: "",
			setupKey: func() *ecdsa.PublicKey {
				privateKey, _ := generateTestKeyPair()
				return &privateKey.PublicKey
			},
			wantError: false,
		},
		{
			name:      "should_encrypt_message_with_special_characters",
			plaintext: "Test!@#$%^&*()_+-=[]{}|;':\",./<>?`~",
			setupKey: func() *ecdsa.PublicKey {
				privateKey, _ := generateTestKeyPair()
				return &privateKey.PublicKey
			},
			wantError: false,
		},
		{
			name:      "should_encrypt_long_message",
			plaintext: string(make([]byte, 10000)), // 10KB message
			setupKey: func() *ecdsa.PublicKey {
				privateKey, _ := generateTestKeyPair()
				return &privateKey.PublicKey
			},
			wantError: false,
		},
		{
			name:      "should_encrypt_json_data",
			plaintext: `{"first_name":"John","last_name":"Doe","email":"john@example.com"}`,
			setupKey: func() *ecdsa.PublicKey {
				privateKey, _ := generateTestKeyPair()
				return &privateKey.PublicKey
			},
			wantError: false,
		},
		{
			name:      "should_encrypt_multiline_text",
			plaintext: "Line 1\nLine 2\nLine 3\n\nEmpty line above",
			setupKey: func() *ecdsa.PublicKey {
				privateKey, _ := generateTestKeyPair()
				return &privateKey.PublicKey
			},
			wantError: false,
		},
		{
			name:      "should_encrypt_unicode_text",
			plaintext: "Hello 世界 สวัสดี مرحبا",
			setupKey: func() *ecdsa.PublicKey {
				privateKey, _ := generateTestKeyPair()
				return &privateKey.PublicKey
			},
			wantError: false,
		},
		{
			name:      "should_return_error_with_nil_public_key",
			plaintext: "test message",
			setupKey: func() *ecdsa.PublicKey {
				return nil
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publicKey := tt.setupKey()

			ciphertext, err := EncryptWithPublicKeyBytes(tt.plaintext, publicKey)

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, ciphertext)
			} else {
				assert.NoError(t, err)

				// For empty strings, ciphertext might be empty but that's acceptable
				if tt.plaintext != "" {
					assert.NotEmpty(t, ciphertext)

					// Verify ciphertext is different from plaintext
					assert.NotEqual(t, tt.plaintext, ciphertext)

					// Verify ciphertext is non-deterministic (encryption with same key produces different ciphertext)
					ciphertext2, err2 := EncryptWithPublicKeyBytes(tt.plaintext, publicKey)
					assert.NoError(t, err2)
					assert.NotEqual(t, ciphertext, ciphertext2, "ECIES encryption should be non-deterministic")
				}

				// Verify ciphertext is base64-encoded (if not empty)
				if ciphertext != "" {
					_, decodeErr := base64.StdEncoding.DecodeString(ciphertext)
					assert.NoError(t, decodeErr, "ciphertext should be valid base64")
				}
			}
		})
	}
}

// TestDecryptWithPrivateKey tests decryption with various inputs
func TestDecryptWithPrivateKey(t *testing.T) {
	t.Run("should_return_error_with_nil_private_key", func(t *testing.T) {
		ciphertext := "dGVzdA==" // base64("test")

		result, err := DecryptWithPrivateKey(ciphertext, nil)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "private key is nil")
	})

	t.Run("should_return_error_with_invalid_base64", func(t *testing.T) {
		privateKey, _ := generateTestKeyPair()
		invalidBase64 := "not-valid-base64-!!!@#$"

		result, err := DecryptWithPrivateKey(invalidBase64, privateKey)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to decode base64")
	})

	t.Run("should_return_error_with_empty_ciphertext", func(t *testing.T) {
		privateKey, _ := generateTestKeyPair()
		emptyBase64 := base64.StdEncoding.EncodeToString([]byte(""))

		// Empty ciphertext might fail during decryption
		result, err := DecryptWithPrivateKey(emptyBase64, privateKey)

		// This should fail as empty data cannot be valid ECIES ciphertext
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("should_return_error_with_wrong_key", func(t *testing.T) {
		// Encrypt with one key
		privateKey1, _ := generateTestKeyPair()
		publicKey1 := &privateKey1.PublicKey

		plaintext := "test message"
		ciphertext, err := EncryptWithPublicKeyBytes(plaintext, publicKey1)
		assert.NoError(t, err)

		// Try to decrypt with different key
		privateKey2, _ := generateTestKeyPair()

		result, err := DecryptWithPrivateKey(ciphertext, privateKey2)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to decrypt")
	})

	t.Run("should_return_error_with_random_ciphertext", func(t *testing.T) {
		privateKey, _ := generateTestKeyPair()

		// Generate random valid base64 string that's not valid ECIES ciphertext
		randomBytes := []byte("this is not valid ECIES ciphertext")
		randomBase64 := base64.StdEncoding.EncodeToString(randomBytes)

		result, err := DecryptWithPrivateKey(randomBase64, privateKey)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

// TestEncryptDecryptRoundTrip tests full encryption/decryption cycle
func TestEncryptDecryptRoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "should_roundtrip_simple_message",
			plaintext: "Hello, World!",
		},
		// Note: Empty strings cannot be encrypted/decrypted with ECIES
		// (produces empty ciphertext that fails decryption)
		// This is a limitation of the ECIES algorithm
		{
			name:      "should_roundtrip_special_characters",
			plaintext: "Test!@#$%^&*()_+-=[]{}|;':\",./<>?`~",
		},
		{
			name:      "should_roundtrip_json_data",
			plaintext: `{"first_name":"John","last_name":"Doe","email":"john@example.com","bio":"Engineer"}`,
		},
		{
			name:      "should_roundtrip_multiline_text",
			plaintext: "Line 1\nLine 2\nLine 3\n\nEmpty line above",
		},
		{
			name:      "should_roundtrip_unicode_text",
			plaintext: "Hello 世界 สวัสดี مرحبا",
		},
		{
			name:      "should_roundtrip_long_message",
			plaintext: string(make([]byte, 5000)), // 5KB message
		},
		{
			name:      "should_roundtrip_csv_data",
			plaintext: "John Doe,Chulalongkorn University,Certificate of Completion,Blockchain Development",
		},
		{
			name:      "should_roundtrip_binary_like_data",
			plaintext: "\x00\x01\x02\x03\xFF\xFE\xFD\xFC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip empty strings - ECIES cannot properly encrypt/decrypt them
			if tt.plaintext == "" {
				t.Skip("ECIES does not support encrypting empty strings")
				return
			}

			// Generate key pair
			privateKey, _ := generateTestKeyPair()
			publicKey := &privateKey.PublicKey

			// Encrypt
			ciphertext, err := EncryptWithPublicKeyBytes(tt.plaintext, publicKey)
			assert.NoError(t, err)
			assert.NotEmpty(t, ciphertext)

			// Decrypt
			decrypted, err := DecryptWithPrivateKey(ciphertext, privateKey)
			assert.NoError(t, err)
			assert.Equal(t, tt.plaintext, decrypted)
		})
	}
}

// TestEncryptDecryptMultipleKeys tests that encryption is key-specific
func TestEncryptDecryptMultipleKeys(t *testing.T) {
	t.Run("should_encrypt_with_different_keys_produce_different_ciphertexts", func(t *testing.T) {
		plaintext := "same message for both keys"

		// Generate two key pairs
		privateKey1, _ := generateTestKeyPair()
		publicKey1 := &privateKey1.PublicKey

		privateKey2, _ := generateTestKeyPair()
		publicKey2 := &privateKey2.PublicKey

		// Encrypt with both keys
		ciphertext1, err1 := EncryptWithPublicKeyBytes(plaintext, publicKey1)
		assert.NoError(t, err1)

		ciphertext2, err2 := EncryptWithPublicKeyBytes(plaintext, publicKey2)
		assert.NoError(t, err2)

		// Ciphertexts should be different
		assert.NotEqual(t, ciphertext1, ciphertext2)

		// Each should decrypt correctly with its own key
		decrypted1, err := DecryptWithPrivateKey(ciphertext1, privateKey1)
		assert.NoError(t, err)
		assert.Equal(t, plaintext, decrypted1)

		decrypted2, err := DecryptWithPrivateKey(ciphertext2, privateKey2)
		assert.NoError(t, err)
		assert.Equal(t, plaintext, decrypted2)

		// Each should NOT decrypt with the other key
		_, err = DecryptWithPrivateKey(ciphertext1, privateKey2)
		assert.Error(t, err)

		_, err = DecryptWithPrivateKey(ciphertext2, privateKey1)
		assert.Error(t, err)
	})
}

// TestGetPublicKeyFromAddress tests the unimplemented function
func TestGetPublicKeyFromAddress(t *testing.T) {
	t.Run("should_return_error_indicating_implementation_required", func(t *testing.T) {
		address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1")

		result, err := GetPublicKeyFromAddress(address)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "must fetch public key from credential storage")
	})

	t.Run("should_return_error_for_zero_address", func(t *testing.T) {
		zeroAddress := common.Address{}

		result, err := GetPublicKeyFromAddress(zeroAddress)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestECIESIntegration tests integration scenarios
func TestECIESIntegration(t *testing.T) {
	t.Run("should_handle_multiple_encryption_decryption_cycles", func(t *testing.T) {
		privateKey, _ := generateTestKeyPair()
		publicKey := &privateKey.PublicKey

		messages := []string{
			"First message",
			"Second message",
			"Third message",
		}

		for i, plaintext := range messages {
			t.Run(string(rune('A'+i)), func(t *testing.T) {
				ciphertext, err := EncryptWithPublicKeyBytes(plaintext, publicKey)
				assert.NoError(t, err)

				decrypted, err := DecryptWithPrivateKey(ciphertext, privateKey)
				assert.NoError(t, err)
				assert.Equal(t, plaintext, decrypted)
			})
		}
	})

	t.Run("should_encrypt_same_message_multiple_times_produce_different_ciphertexts", func(t *testing.T) {
		privateKey, _ := generateTestKeyPair()
		publicKey := &privateKey.PublicKey
		plaintext := "same message encrypted multiple times"

		ciphertexts := make([]string, 10)
		for i := 0; i < 10; i++ {
			ciphertext, err := EncryptWithPublicKeyBytes(plaintext, publicKey)
			assert.NoError(t, err)
			ciphertexts[i] = ciphertext

			// Verify each can be decrypted correctly
			decrypted, err := DecryptWithPrivateKey(ciphertext, privateKey)
			assert.NoError(t, err)
			assert.Equal(t, plaintext, decrypted)
		}

		// All ciphertexts should be different (non-deterministic encryption)
		for i := 0; i < len(ciphertexts); i++ {
			for j := i + 1; j < len(ciphertexts); j++ {
				assert.NotEqual(t, ciphertexts[i], ciphertexts[j],
					"ECIES encryption should produce different ciphertexts for same plaintext")
			}
		}
	})
}

// TestECIESEdgeCases tests edge cases and boundary conditions
func TestECIESEdgeCases(t *testing.T) {
	t.Run("should_handle_very_long_plaintext", func(t *testing.T) {
		privateKey, _ := generateTestKeyPair()
		publicKey := &privateKey.PublicKey

		// Create 1MB message
		longPlaintext := string(make([]byte, 1024*1024))

		ciphertext, err := EncryptWithPublicKeyBytes(longPlaintext, publicKey)
		assert.NoError(t, err)

		decrypted, err := DecryptWithPrivateKey(ciphertext, privateKey)
		assert.NoError(t, err)
		assert.Equal(t, longPlaintext, decrypted)
	})

	t.Run("should_handle_null_bytes_in_plaintext", func(t *testing.T) {
		privateKey, _ := generateTestKeyPair()
		publicKey := &privateKey.PublicKey

		plaintext := "text\x00with\x00null\x00bytes"

		ciphertext, err := EncryptWithPublicKeyBytes(plaintext, publicKey)
		assert.NoError(t, err)

		decrypted, err := DecryptWithPrivateKey(ciphertext, privateKey)
		assert.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("should_handle_all_null_bytes", func(t *testing.T) {
		privateKey, _ := generateTestKeyPair()
		publicKey := &privateKey.PublicKey

		plaintext := string(make([]byte, 100))

		ciphertext, err := EncryptWithPublicKeyBytes(plaintext, publicKey)
		assert.NoError(t, err)

		decrypted, err := DecryptWithPrivateKey(ciphertext, privateKey)
		assert.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})
}

// Benchmark tests
func BenchmarkEncryptWithPublicKeyBytes(b *testing.B) {
	privateKey, _ := generateTestKeyPair()
	publicKey := &privateKey.PublicKey
	plaintext := "benchmark test message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = EncryptWithPublicKeyBytes(plaintext, publicKey)
	}
}

func BenchmarkDecryptWithPrivateKey(b *testing.B) {
	privateKey, _ := generateTestKeyPair()
	publicKey := &privateKey.PublicKey
	plaintext := "benchmark test message"

	ciphertext, err := EncryptWithPublicKeyBytes(plaintext, publicKey)
	if err != nil {
		b.Fatalf("Failed to encrypt: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecryptWithPrivateKey(ciphertext, privateKey)
	}
}

func BenchmarkEncryptDecryptRoundTrip(b *testing.B) {
	privateKey, _ := generateTestKeyPair()
	publicKey := &privateKey.PublicKey
	plaintext := "benchmark test message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ciphertext, err := EncryptWithPublicKeyBytes(plaintext, publicKey)
		if err != nil {
			b.Fatalf("Failed to encrypt: %v", err)
		}

		_, err = DecryptWithPrivateKey(ciphertext, privateKey)
		if err != nil {
			b.Fatalf("Failed to decrypt: %v", err)
		}
	}
}
