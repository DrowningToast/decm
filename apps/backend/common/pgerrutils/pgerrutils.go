package pgerrutils

import (
	"apps/backend/common/customerror"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func ParsePgError(err error) *customerror.Err {
	var pgErr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) {
		return customerror.Parse(&customerror.ErrNotFound, err)
	} else if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return customerror.Parse(&customerror.ErrDuplicateEntry, err)
		case pgerrcode.ForeignKeyViolation:
			return customerror.Parse(&customerror.ErrInvalidArgument, err)
		case pgerrcode.NotNullViolation:
			return customerror.Parse(&customerror.ErrInvalidArgument, err)
		case pgerrcode.InvalidTextRepresentation:
			return customerror.Parse(&customerror.ErrInvalidArgument, err)
		default:
			return customerror.Parse(&customerror.ErrInternalServer, err)
		}
	}
	return customerror.Parse(&customerror.ErrInternalServer, err)
}
