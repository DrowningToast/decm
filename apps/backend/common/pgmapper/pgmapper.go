package pgmapper

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

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
	return pgtype.Int4{
		Int32: BoolToInt32(*value),
		Valid: value != nil,
	}
}

func IntPtrToBool(value *int32) bool {
	if value != nil {
		return *value == 1
	}
	return false
}

func IntPtrToPgInt4(value *int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: *value,
		Valid: value != nil,
	}
}

func Int32ToPgInt4(value int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: value,
		Valid: true,
	}
}
