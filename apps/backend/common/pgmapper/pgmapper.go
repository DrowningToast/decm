package pgmapper

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type mapper struct {
	piiEncryptionKey string
}

func NewMapper(piiEncryptionKey string) *mapper {
	return &mapper{
		piiEncryptionKey: piiEncryptionKey,
	}
}

func (m *mapper) PgTextToStringPtr(text pgtype.Text) *string {
	if text.Valid {
		return &text.String
	}
	return nil
}

func (m *mapper) StringPtrToPgText(str *string) pgtype.Text {
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

func BoolToInt(value bool) int32 {
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

func IntPtrToBool(value *int32) bool {
	if value != nil {
		return *value == 1
	}
	return false
}
