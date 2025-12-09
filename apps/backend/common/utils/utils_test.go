package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSecureRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "16 characters",
			length: 16,
		},
		{
			name:   "32 characters",
			length: 32,
		},
		{
			name:   "64 characters",
			length: 64,
		},
		{
			name:   "128 characters",
			length: 128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSecureRandomString(tt.length)
			assert.Len(t, result, tt.length)
			assert.NotEmpty(t, result)

			// Generate another one - should be different
			result2 := GenerateSecureRandomString(tt.length)
			assert.Len(t, result2, tt.length)
			assert.NotEqual(t, result, result2, "Random strings should be different")
		})
	}
}

func TestGenerateSecureRandomString_EmptyLength(t *testing.T) {
	result := GenerateSecureRandomString(0)
	assert.Empty(t, result)
}

func TestGenerateSecureRandomString_ConsistentLength(t *testing.T) {
	// Generate multiple strings and verify they all have the same length
	length := 32
	for i := 0; i < 10; i++ {
		result := GenerateSecureRandomString(length)
		assert.Len(t, result, length, "String %d should have length %d", i, length)
	}
}

func TestGenerateSecureRandomString_Uniqueness(t *testing.T) {
	length := 32
	generated := make(map[string]bool)

	// Generate 100 strings and verify they're all unique
	for i := 0; i < 100; i++ {
		result := GenerateSecureRandomString(length)
		assert.False(t, generated[result], "String should be unique (iteration %d)", i)
		generated[result] = true
	}
}



