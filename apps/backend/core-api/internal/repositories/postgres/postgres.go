package postgres

import (
	"decm-database/go/generated"

	"apps/backend/common/pgclient"
)

type Repository struct {
	db      pgclient.Client
	queries *generated.Queries
}

func NewRepository(db pgclient.Client) *Repository {
	return &Repository{
		db:      db,
		queries: generated.New(db),
	}
}
