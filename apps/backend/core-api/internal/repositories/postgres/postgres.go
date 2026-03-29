package postgres

import (
	"apps/backend/common/pgclient"
	"decm-database/go/generated"
)

type Repository struct {
	db               pgclient.Client
	queries          generated.Querier
	piiEncryptionKey string
}

func NewRepository(db pgclient.Client, piiEncryptionKey string) *Repository {
	return &Repository{
		db:               db,
		queries:          generated.New(db),
		piiEncryptionKey: piiEncryptionKey,
	}
}
