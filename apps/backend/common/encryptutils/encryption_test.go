package encryptutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDeterministicAES(t *testing.T) {
	password := "test-password-123"
	plaintext := "sensitive-data@example.com"

	encrypted, err := EncryptDeterministicAES(plaintext, password)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, plaintext, encrypted)

	// Test deterministic - same input should produce same output
	encrypted2, err := EncryptDeterministicAES(plaintext, password)
	require.NoError(t, err)
	assert.Equal(t, encrypted, encrypted2)
}

func TestDecryptDeterministicAES(t *testing.T) {
	password := "test-password-123"
	plaintext := "sensitive-data@example.com"

	encrypted, err := EncryptDeterministicAES(plaintext, password)
	require.NoError(t, err)

	decrypted, err := DecryptDeterministicAES(encrypted, password)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptDeterministicAES_DecryptDeterministicAES_RoundTrip(t *testing.T) {
	password := "test-password-123"
	testCases := []string{
		"test@example.com",
		"John Doe",
		"",
		"special-chars-!@#$%^&*()",
		"unicode-测试-🚀",
		"long-text-" + string(make([]byte, 1000)),
	}

	for _, plaintext := range testCases {
		t.Run(plaintext[:min(len(plaintext), 20)], func(t *testing.T) {
			encrypted, err := EncryptDeterministicAES(plaintext, password)
			require.NoError(t, err)

			decrypted, err := DecryptDeterministicAES(encrypted, password)
			require.NoError(t, err)
			assert.Equal(t, plaintext, decrypted)
		})
	}
}

func TestDecryptDeterministicAES_InvalidKey(t *testing.T) {
	password := "test-password-123"
	wrongPassword := "wrong-password"
	plaintext := "sensitive-data@example.com"

	encrypted, err := EncryptDeterministicAES(plaintext, password)
	require.NoError(t, err)

	_, err = DecryptDeterministicAES(encrypted, wrongPassword)
	assert.Error(t, err)
}

func TestDecryptDeterministicAES_InvalidCiphertext(t *testing.T) {
	password := "test-password-123"
	invalidCiphertext := "invalid-base64-!@#$"

	_, err := DecryptDeterministicAES(invalidCiphertext, password)
	assert.Error(t, err)
}

func TestDecryptDeterministicAES_TooShort(t *testing.T) {
	password := "test-password-123"
	shortCiphertext := "dGVzdA==" // base64 for "test" but too short for GCM

	_, err := DecryptDeterministicAES(shortCiphertext, password)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too short")
}

func TestEncryptAESGCM(t *testing.T) {
	password := "test-password-123"
	plaintext := "sensitive-data@example.com"

	encrypted, err := EncryptAESGCM(plaintext, password)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, plaintext, encrypted)

	// Test non-deterministic - same input should produce different output (due to random nonce)
	encrypted2, err := EncryptAESGCM(plaintext, password)
	require.NoError(t, err)
	assert.NotEqual(t, encrypted, encrypted2) // Should be different due to random nonce
}

func TestDecryptAESGCM(t *testing.T) {
	password := "test-password-123"
	plaintext := "sensitive-data@example.com"

	encrypted, err := EncryptAESGCM(plaintext, password)
	require.NoError(t, err)

	decrypted, err := DecryptAESGCM(encrypted, password)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptAESGCM_DecryptAESGCM_RoundTrip(t *testing.T) {
	password := "test-password-123"
	testCases := []string{
		"test@example.com",
		"John Doe",
		"",
		"special-chars-!@#$%^&*()",
		"unicode-测试-🚀",
	}

	for _, plaintext := range testCases {
		t.Run(plaintext, func(t *testing.T) {
			encrypted, err := EncryptAESGCM(plaintext, password)
			require.NoError(t, err)

			decrypted, err := DecryptAESGCM(encrypted, password)
			require.NoError(t, err)
			assert.Equal(t, plaintext, decrypted)
		})
	}
}

func TestDecryptAESGCM_InvalidKey(t *testing.T) {
	password := "test-password-123"
	wrongPassword := "wrong-password"
	plaintext := "sensitive-data@example.com"

	encrypted, err := EncryptAESGCM(plaintext, password)
	require.NoError(t, err)

	_, err = DecryptAESGCM(encrypted, wrongPassword)
	assert.Error(t, err)
}

func TestEncryptAESWithKey(t *testing.T) {
	testCases := []struct {
		name      string
		keySize   int
		plaintext string
		hasError  bool
	}{
		{
			name:      "AES-128 (16 bytes)",
			keySize:   16,
			plaintext: "test-data",
			hasError:  false,
		},
		{
			name:      "AES-192 (24 bytes)",
			keySize:   24,
			plaintext: "test-data",
			hasError:  false,
		},
		{
			name:      "AES-256 (32 bytes)",
			keySize:   32,
			plaintext: "test-data",
			hasError:  false,
		},
		{
			name:      "invalid key size (15 bytes)",
			keySize:   15,
			plaintext: "test-data",
			hasError:  true,
		},
		{
			name:      "invalid key size (25 bytes)",
			keySize:   25,
			plaintext: "test-data",
			hasError:  true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			key := make([]byte, tt.keySize)
			for i := range key {
				key[i] = byte(i)
			}

			encrypted, err := EncryptAESWithKey(tt.plaintext, key)
			if tt.hasError {
				assert.Error(t, err)
				assert.Empty(t, encrypted)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, encrypted)

				decrypted, err := DecryptAESWithKey(encrypted, key)
				require.NoError(t, err)
				assert.Equal(t, tt.plaintext, decrypted)
			}
		})
	}
}

func TestDecryptAESWithKey(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	plaintext := "test-data"

	encrypted, err := EncryptAESWithKey(plaintext, key)
	require.NoError(t, err)

	decrypted, err := DecryptAESWithKey(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestDecryptAESWithKey_InvalidKey(t *testing.T) {
	key := make([]byte, 32)
	wrongKey := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
		wrongKey[i] = byte(i + 1)
	}
	plaintext := "test-data"

	encrypted, err := EncryptAESWithKey(plaintext, key)
	require.NoError(t, err)

	_, err = DecryptAESWithKey(encrypted, wrongKey)
	assert.Error(t, err)
}

func TestGenerateAESKey(t *testing.T) {
	key1 := GenerateAESKey()
	key2 := GenerateAESKey()

	assert.Len(t, key1, 32)
	assert.Len(t, key2, 32)
	assert.NotEqual(t, key1, key2) // Should be different each time
}

func TestKeyFromPassword(t *testing.T) {
	password := "test-password"
	salt := []byte("test-salt-12345")

	key1 := KeyFromPassword(password, salt)
	key2 := KeyFromPassword(password, salt)

	assert.Len(t, key1, 32)
	assert.Equal(t, key1, key2) // Should be deterministic

	// Different salt should produce different key
	salt2 := []byte("different-salt")
	key3 := KeyFromPassword(password, salt2)
	assert.NotEqual(t, key1, key3)

	// Different password should produce different key
	key4 := KeyFromPassword("different-password", salt)
	assert.NotEqual(t, key1, key4)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
