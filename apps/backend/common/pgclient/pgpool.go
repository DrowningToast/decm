package pgclient

import (
	"apps/backend/core-api/config/postgres"
	"apps/backend/services/log"
	"context"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

func NewPool(ctx context.Context, pgConfig *postgres.Config) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(pgConfig.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pgxpool")
	}

	connConfig.MaxConns = lo.CoalesceOrEmpty(pgConfig.MaxConns, int32(postgres.MaxConnections))
	connConfig.MinConns = lo.CoalesceOrEmpty(pgConfig.MinConns, int32(postgres.MinConnections))
	logger := log.NewLogger()
	connConfig.ConnConfig.Tracer = pgConfig.QueryTracer(*logger)

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pgxpool")
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return pool, nil
}
