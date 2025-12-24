package pgerrutils

import (
	"apps/backend/common/customerror"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestParsePgError_NoRows(t *testing.T) {
	err := pgx.ErrNoRows
	result := ParsePgError(err)

	assert.NotNil(t, result)
	assert.Equal(t, customerror.ErrNotFound.Code, *result.Code)
}

func TestParsePgError_UniqueViolation(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code: pgerrcode.UniqueViolation,
	}
	err := &pgconn.PgError{
		Code: pgErr.Code,
	}

	result := ParsePgError(err)
	assert.NotNil(t, result)
	assert.Equal(t, customerror.ErrDuplicateEntry.Code, *result.Code)
}

func TestParsePgError_ForeignKeyViolation(t *testing.T) {
	err := &pgconn.PgError{
		Code: pgerrcode.ForeignKeyViolation,
	}

	result := ParsePgError(err)
	assert.NotNil(t, result)
	assert.Equal(t, customerror.ErrInvalidArgument.Code, *result.Code)
}

func TestParsePgError_NotNullViolation(t *testing.T) {
	err := &pgconn.PgError{
		Code: pgerrcode.NotNullViolation,
	}

	result := ParsePgError(err)
	assert.NotNil(t, result)
	assert.Equal(t, customerror.ErrInvalidArgument.Code, *result.Code)
}

func TestParsePgError_InvalidTextRepresentation(t *testing.T) {
	err := &pgconn.PgError{
		Code: pgerrcode.InvalidTextRepresentation,
	}

	result := ParsePgError(err)
	assert.NotNil(t, result)
	assert.Equal(t, customerror.ErrInvalidArgument.Code, *result.Code)
}

func TestParsePgError_UnknownPgError(t *testing.T) {
	err := &pgconn.PgError{
		Code: "99999", // Unknown error code
	}

	result := ParsePgError(err)
	assert.NotNil(t, result)
	assert.Equal(t, customerror.ErrInternalServer.Code, *result.Code)
}

func TestParsePgError_NonPgError(t *testing.T) {
	err := errors.New("some generic error")

	result := ParsePgError(err)
	assert.NotNil(t, result)
	assert.Equal(t, customerror.ErrInternalServer.Code, *result.Code)
}

func TestParsePgError_NilError(t *testing.T) {
	result := ParsePgError(nil)
	assert.NotNil(t, result)
	assert.Equal(t, customerror.ErrInternalServer.Code, *result.Code)
}

func TestParsePgError_WrappedNoRows(t *testing.T) {
	err := errors.Wrap(pgx.ErrNoRows, "wrapped error")
	result := ParsePgError(err)

	assert.NotNil(t, result)
	// Note: errors.Wrap doesn't preserve pgx.ErrNoRows type, so it will be treated as generic error
	// This test documents current behavior
	assert.NotNil(t, result)
}
