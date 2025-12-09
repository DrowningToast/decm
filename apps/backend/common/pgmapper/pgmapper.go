package pgmapper

import (
	"strings"
	"time"

	"apps/backend/common/encryptutils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ========== Type Conversion Functions ==========

func PgTextToStringPtr(text pgtype.Text) *string {
	if text.Valid {
		return &text.String
	}
	return nil
}

func StringPtrToPgText(str *string) pgtype.Text {
	if str != nil {
		return pgtype.Text{
			String: *str,
			Valid:  true,
		}
	}
	return pgtype.Text{}
}

func PgTimestampzToTimePtr(timestampz pgtype.Timestamptz) *time.Time {
	if timestampz.Valid {
		return &timestampz.Time
	}
	return nil
}

func TimePtrToPgTimestampz(time *time.Time) pgtype.Timestamptz {
	if time != nil {
		return pgtype.Timestamptz{
			Time:  *time,
			Valid: true,
		}
	}
	return pgtype.Timestamptz{}
}

func Int32PtrToBoolean(value *int32) bool {
	if value != nil {
		return *value == 1
	}
	return false
}

func BoolToPgInt4(value bool) pgtype.Int4 {
	return pgtype.Int4{
		Int32: BoolToInt32(value),
		Valid: true,
	}
}

func BoolToInt32(value bool) int32 {
	var intValue int32 = 0
	if value {
		intValue = 1
	}
	return intValue
}

func BoolToIntPtr(value bool) *int32 {
	var intValue int32 = 0
	if value {
		return &intValue
	}
	return &intValue
}

func BoolPtrToPgInt4(value *bool) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{
			Valid: false,
		}
	}
	return pgtype.Int4{
		Int32: BoolToInt32(*value),
		Valid: true,
	}
}

func IntPtrToBool(value *int32) bool {
	if value != nil {
		return *value == 1
	}
	return false
}

func IntPtrToPgInt4(value *int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{
			Valid: false,
		}
	}
	return pgtype.Int4{
		Int32: *value,
		Valid: true,
	}
}

func Int32ToPgInt4(value int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: value,
		Valid: true,
	}
}

func UUIDToPgUUID(value *uuid.UUID) pgtype.UUID {
	if value == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{
		Bytes: *value,
		Valid: true,
	}
}

func PgUUIDToUUIDPtr(value pgtype.UUID) *uuid.UUID {
	if !value.Valid {
		return nil
	}
	result := uuid.UUID(value.Bytes)
	return &result
}

func PgFloat8ToFloat64Ptr(f8 pgtype.Float8) *float64 {
	if !f8.Valid {
		return nil
	}
	return &f8.Float64
}

func PgInt4ToInt32Ptr(i4 pgtype.Int4) *int32 {
	if !i4.Valid {
		return nil
	}
	return &i4.Int32
}

// ========== PII Encryption Functions ==========

// EncryptPII encrypts plaintext PII data using AES-GCM
func EncryptPII(plaintext string, encryptionKey string) (string, error) {
	return encryptutils.EncryptDeterministicAES(plaintext, encryptionKey)
}

// DecryptPII decrypts ciphertext PII data using AES-GCM
func DecryptPII(ciphertext string, encryptionKey string) (string, error) {
	return encryptutils.DecryptDeterministicAES(ciphertext, encryptionKey)
}

// EncryptStringPtrToPgText encrypts a nullable string pointer and returns pgtype.Text
// Returns invalid pgtype.Text if the input is nil or empty
func EncryptStringPtrToPgText(field *string, encryptionKey string) (pgtype.Text, error) {
	if field == nil || *field == "" {
		return pgtype.Text{Valid: false}, nil
	}
	encrypted, err := EncryptPII(*field, encryptionKey)
	if err != nil {
		return pgtype.Text{}, err
	}
	return pgtype.Text{String: encrypted, Valid: true}, nil
}

// DecryptPgTextToStringPtr decrypts a pgtype.Text field and returns a nullable string pointer
// Returns nil if the pgtype.Text is invalid
// If the data is not encrypted (legacy data), returns the original value
func DecryptPgTextToStringPtr(field pgtype.Text, encryptionKey string) (*string, error) {
	if !field.Valid {
		return nil, nil
	}
	decrypted, err := DecryptPII(field.String, encryptionKey)
	if err != nil {
		// Check if this is an error indicating unencrypted legacy data:
		// - "illegal base64" - data is not base64 encoded
		// - "ciphertext too short" - data is too short to be valid encrypted data
		// - "message authentication failed" - decryption failed (wrong key or corrupted data)
		errMsg := err.Error()
		if strings.Contains(errMsg, "illegal base64") ||
			strings.Contains(errMsg, "base64") ||
			strings.Contains(errMsg, "ciphertext too short") {
			// Return original value for legacy unencrypted data
			return &field.String, nil
		}
		return nil, err
	}
	return &decrypted, nil
}
