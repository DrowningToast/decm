package pgclient

import (
	"context"
	"decm-database/go/generated"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	_ PGTX = (*pgx.Conn)(nil)
	_ PGTX = (*pgxpool.Conn)(nil)
)

type PGTX interface {
	generated.DBTX

	Begin(context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Client interface {
	PGTX
	generated.DBTX

	SendBatch(ctx context.Context, b *pgx.Batch) (br pgx.BatchResults)
	Ping(ctx context.Context) error
}
