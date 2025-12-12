package pgmapper

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPgTextToStringPtr(t *testing.T) {
	tests := []struct {
		name     string
		input    pgtype.Text
		expected *string
	}{
		{
			name: "valid text",
			input: pgtype.Text{
				String: "test",
				Valid:  true,
			},
			expected: stringPtr("test"),
		},
		{
			name: "invalid text",
			input: pgtype.Text{
				Valid: false,
			},
			expected: nil,
		},
		{
			name: "empty valid text",
			input: pgtype.Text{
				String: "",
				Valid:  true,
			},
			expected: stringPtr(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PgTextToStringPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStringPtrToPgText(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected pgtype.Text
	}{
		{
			name:  "valid string pointer",
			input: stringPtr("test"),
			expected: pgtype.Text{
				String: "test",
				Valid:  true,
			},
		},
		{
			name:  "nil pointer",
			input: nil,
			expected: pgtype.Text{
				Valid: false,
			},
		},
		{
			name:  "empty string",
			input: stringPtr(""),
			expected: pgtype.Text{
				String: "",
				Valid:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringPtrToPgText(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPgTimestampzToTimePtr(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		input    pgtype.Timestamptz
		expected *time.Time
	}{
		{
			name: "valid timestamp",
			input: pgtype.Timestamptz{
				Time:  now,
				Valid: true,
			},
			expected: &now,
		},
		{
			name: "invalid timestamp",
			input: pgtype.Timestamptz{
				Valid: false,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PgTimestampzToTimePtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTimePtrToPgTimestampz(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		input    *time.Time
		expected pgtype.Timestamptz
	}{
		{
			name:  "valid time pointer",
			input: &now,
			expected: pgtype.Timestamptz{
				Time:  now,
				Valid: true,
			},
		},
		{
			name:  "nil pointer",
			input: nil,
			expected: pgtype.Timestamptz{
				Valid: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimePtrToPgTimestampz(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInt32PtrToBoolean(t *testing.T) {
	tests := []struct {
		name     string
		input    *int32
		expected bool
	}{
		{
			name:     "value 1",
			input:    int32Ptr(1),
			expected: true,
		},
		{
			name:     "value 0",
			input:    int32Ptr(0),
			expected: false,
		},
		{
			name:     "nil pointer",
			input:    nil,
			expected: false,
		},
		{
			name:     "value 2",
			input:    int32Ptr(2),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int32PtrToBoolean(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBoolToPgInt4(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected pgtype.Int4
	}{
		{
			name:  "true",
			input: true,
			expected: pgtype.Int4{
				Int32: 1,
				Valid: true,
			},
		},
		{
			name:  "false",
			input: false,
			expected: pgtype.Int4{
				Int32: 0,
				Valid: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolToPgInt4(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBoolToInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected int32
	}{
		{
			name:     "true",
			input:    true,
			expected: 1,
		},
		{
			name:     "false",
			input:    false,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolToInt32(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBoolPtrToPgInt4(t *testing.T) {
	trueVal := true
	falseVal := false
	tests := []struct {
		name     string
		input    *bool
		expected pgtype.Int4
	}{
		{
			name:  "true pointer",
			input: &trueVal,
			expected: pgtype.Int4{
				Int32: 1,
				Valid: true,
			},
		},
		{
			name:  "false pointer",
			input: &falseVal,
			expected: pgtype.Int4{
				Int32: 0,
				Valid: true,
			},
		},
		{
			name:  "nil pointer",
			input: nil,
			expected: pgtype.Int4{
				Valid: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolPtrToPgInt4(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntPtrToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    *int32
		expected bool
	}{
		{
			name:     "value 1",
			input:    int32Ptr(1),
			expected: true,
		},
		{
			name:     "value 0",
			input:    int32Ptr(0),
			expected: false,
		},
		{
			name:     "nil pointer",
			input:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntPtrToBool(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntPtrToPgInt4(t *testing.T) {
	tests := []struct {
		name     string
		input    *int32
		expected pgtype.Int4
	}{
		{
			name:  "valid int32 pointer",
			input: int32Ptr(42),
			expected: pgtype.Int4{
				Int32: 42,
				Valid: true,
			},
		},
		{
			name:  "nil pointer",
			input: nil,
			expected: pgtype.Int4{
				Valid: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntPtrToPgInt4(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInt32ToPgInt4(t *testing.T) {
	tests := []struct {
		name     string
		input    int32
		expected pgtype.Int4
	}{
		{
			name:  "positive value",
			input: 42,
			expected: pgtype.Int4{
				Int32: 42,
				Valid: true,
			},
		},
		{
			name:  "zero",
			input: 0,
			expected: pgtype.Int4{
				Int32: 0,
				Valid: true,
			},
		},
		{
			name:  "negative value",
			input: -10,
			expected: pgtype.Int4{
				Int32: -10,
				Valid: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int32ToPgInt4(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUUIDToPgUUID(t *testing.T) {
	testUUID := uuid.New()
	tests := []struct {
		name     string
		input    *uuid.UUID
		expected pgtype.UUID
	}{
		{
			name:  "valid UUID pointer",
			input: &testUUID,
			expected: pgtype.UUID{
				Bytes: testUUID,
				Valid: true,
			},
		},
		{
			name:  "nil pointer",
			input: nil,
			expected: pgtype.UUID{
				Valid: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UUIDToPgUUID(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPgUUIDToUUIDPtr(t *testing.T) {
	testUUID := uuid.New()
	tests := []struct {
		name     string
		input    pgtype.UUID
		expected *uuid.UUID
	}{
		{
			name: "valid UUID",
			input: pgtype.UUID{
				Bytes: testUUID,
				Valid: true,
			},
			expected: &testUUID,
		},
		{
			name: "invalid UUID",
			input: pgtype.UUID{
				Valid: false,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PgUUIDToUUIDPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPgFloat8ToFloat64Ptr(t *testing.T) {
	tests := []struct {
		name     string
		input    pgtype.Float8
		expected *float64
	}{
		{
			name: "valid float",
			input: pgtype.Float8{
				Float64: 3.14,
				Valid:   true,
			},
			expected: float64Ptr(3.14),
		},
		{
			name: "invalid float",
			input: pgtype.Float8{
				Valid: false,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PgFloat8ToFloat64Ptr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPgInt4ToInt32Ptr(t *testing.T) {
	tests := []struct {
		name     string
		input    pgtype.Int4
		expected *int32
	}{
		{
			name: "valid int",
			input: pgtype.Int4{
				Int32: 42,
				Valid: true,
			},
			expected: int32Ptr(42),
		},
		{
			name: "invalid int",
			input: pgtype.Int4{
				Valid: false,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PgInt4ToInt32Ptr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEncryptPII(t *testing.T) {
	encryptionKey := "test-encryption-key-12345"
	plaintext := "sensitive-data@example.com"

	encrypted, err := EncryptPII(plaintext, encryptionKey)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, plaintext, encrypted)

	// Test deterministic encryption - same input should produce same output
	encrypted2, err := EncryptPII(plaintext, encryptionKey)
	require.NoError(t, err)
	assert.Equal(t, encrypted, encrypted2)
}

func TestDecryptPII(t *testing.T) {
	encryptionKey := "test-encryption-key-12345"
	plaintext := "sensitive-data@example.com"

	encrypted, err := EncryptPII(plaintext, encryptionKey)
	require.NoError(t, err)

	decrypted, err := DecryptPII(encrypted, encryptionKey)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptPII_DecryptPII_RoundTrip(t *testing.T) {
	encryptionKey := "test-encryption-key-12345"
	testCases := []string{
		"test@example.com",
		"John Doe",
		"",
		"special-chars-!@#$%^&*()",
		"unicode-测试-🚀",
	}

	for _, plaintext := range testCases {
		t.Run(plaintext, func(t *testing.T) {
			encrypted, err := EncryptPII(plaintext, encryptionKey)
			require.NoError(t, err)

			decrypted, err := DecryptPII(encrypted, encryptionKey)
			require.NoError(t, err)
			assert.Equal(t, plaintext, decrypted)
		})
	}
}

func TestDecryptPII_InvalidKey(t *testing.T) {
	encryptionKey := "test-encryption-key-12345"
	wrongKey := "wrong-key"
	plaintext := "sensitive-data@example.com"

	encrypted, err := EncryptPII(plaintext, encryptionKey)
	require.NoError(t, err)

	_, err = DecryptPII(encrypted, wrongKey)
	assert.Error(t, err)
}

func TestEncryptStringPtrToPgText(t *testing.T) {
	encryptionKey := "test-encryption-key-12345"
	tests := []struct {
		name     string
		input    *string
		expected pgtype.Text
		hasError bool
	}{
		{
			name:  "valid string",
			input: stringPtr("test@example.com"),
			expected: pgtype.Text{
				Valid: true,
			},
			hasError: false,
		},
		{
			name:  "nil pointer",
			input: nil,
			expected: pgtype.Text{
				Valid: false,
			},
			hasError: false,
		},
		{
			name:  "empty string",
			input: stringPtr(""),
			expected: pgtype.Text{
				Valid: false,
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EncryptStringPtrToPgText(tt.input, encryptionKey)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Valid, result.Valid)
				if tt.input != nil && *tt.input != "" {
					assert.NotEmpty(t, result.String)
					assert.NotEqual(t, *tt.input, result.String)
				}
			}
		})
	}
}

func TestDecryptPgTextToStringPtr(t *testing.T) {
	encryptionKey := "test-encryption-key-12345"
	plaintext := "test@example.com"

	encrypted, err := EncryptPII(plaintext, encryptionKey)
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    pgtype.Text
		expected *string
		hasError bool
	}{
		{
			name: "valid encrypted text",
			input: pgtype.Text{
				String: encrypted,
				Valid:  true,
			},
			expected: stringPtr(plaintext),
			hasError: false,
		},
		{
			name: "invalid text",
			input: pgtype.Text{
				Valid: false,
			},
			expected: nil,
			hasError: false,
		},
		{
			name: "legacy unencrypted data",
			input: pgtype.Text{
				String: "legacy-unencrypted-data",
				Valid:  true,
			},
			expected: stringPtr("legacy-unencrypted-data"),
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DecryptPgTextToStringPtr(tt.input, encryptionKey)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}


