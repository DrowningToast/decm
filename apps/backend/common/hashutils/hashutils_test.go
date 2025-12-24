package hashutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCryptographicSalt(t *testing.T) {
	tests := []struct {
		name     string
		saltSize uint32
		hasError bool
	}{
		{
			name:     "16 bytes",
			saltSize: 16,
			hasError: false,
		},
		{
			name:     "32 bytes",
			saltSize: 32,
			hasError: false,
		},
		{
			name:     "64 bytes",
			saltSize: 64,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			salt, err := GenerateCryptographicSalt(tt.saltSize)
			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, salt)
			} else {
				require.NoError(t, err)
				assert.Len(t, salt, int(tt.saltSize))
				// Salt should be different each time
				salt2, err := GenerateCryptographicSalt(tt.saltSize)
				require.NoError(t, err)
				assert.NotEqual(t, salt, salt2)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "test-password-123"

	hashed, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
	assert.Contains(t, hashed, "$argon2id$")

	// Different passwords should produce different hashes
	hashed2, err := HashPassword("different-password")
	require.NoError(t, err)
	assert.NotEqual(t, hashed, hashed2)
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	hashed, err := HashPassword("")
	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
}

func TestHashPasswordWithCustomSalt(t *testing.T) {
	password := "test-password-123"
	salt, err := GenerateCryptographicSalt(16)
	require.NoError(t, err)

	hashed1, err := HashPasswordWithCustomSalt(password, salt)
	require.NoError(t, err)
	assert.NotEmpty(t, hashed1)

	// Same password and salt should produce same hash
	hashed2, err := HashPasswordWithCustomSalt(password, salt)
	require.NoError(t, err)
	assert.Equal(t, hashed1, hashed2)

	// Different salt should produce different hash
	salt2, err := GenerateCryptographicSalt(16)
	require.NoError(t, err)
	hashed3, err := HashPasswordWithCustomSalt(password, salt2)
	require.NoError(t, err)
	assert.NotEqual(t, hashed1, hashed3)
}

func TestCompareHash(t *testing.T) {
	password := "test-password-123"

	hashed, err := HashPassword(password)
	require.NoError(t, err)

	match, err := CompareHash(password, hashed)
	require.NoError(t, err)
	assert.True(t, match)

	// Wrong password should not match
	match, err = CompareHash("wrong-password", hashed)
	require.NoError(t, err)
	assert.False(t, match)
}

func TestCompareHash_InvalidFormat(t *testing.T) {
	tests := []struct {
		name           string
		hashedPassword string
		hasError       bool
	}{
		{
			name:           "invalid format - missing parts",
			hashedPassword: "$argon2id$v=19",
			hasError:       true,
		},
		{
			name:           "invalid format - wrong algorithm",
			hashedPassword: "$bcrypt$v=19$m=65536,t=2,p=4$salt$hash",
			hasError:       true,
		},
		{
			name:           "invalid format - not a hash",
			hashedPassword: "plain-text",
			hasError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CompareHash("password", tt.hashedPassword)
			if tt.hasError {
				assert.Error(t, err)
				assert.False(t, match)
			}
		})
	}
}

func TestCompareHash_RoundTrip(t *testing.T) {
	testCases := []string{
		"test-password-123",
		"",
		"special-chars-!@#$%^&*()",
		"unicode-测试-🚀",
		"very-long-password-" + string(make([]byte, 100)),
	}

	for _, password := range testCases {
		t.Run(password[:min(len(password), 20)], func(t *testing.T) {
			hashed, err := HashPassword(password)
			require.NoError(t, err)

			match, err := CompareHash(password, hashed)
			require.NoError(t, err)
			assert.True(t, match)
		})
	}
}

func TestHashSHA256(t *testing.T) {
	data := "test-data"

	hash, err := HashSHA256(data)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 64) // SHA256 produces 64 hex characters

	// Same input should produce same hash
	hash2, err := HashSHA256(data)
	require.NoError(t, err)
	assert.Equal(t, hash, hash2)

	// Different input should produce different hash
	hash3, err := HashSHA256("different-data")
	require.NoError(t, err)
	assert.NotEqual(t, hash, hash3)
}

func TestHashSHA256_EmptyString(t *testing.T) {
	hash, err := HashSHA256("")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 64)
}

func TestHashSHA256_Unicode(t *testing.T) {
	data := "unicode-测试-🚀"

	hash, err := HashSHA256(data)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 64)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

