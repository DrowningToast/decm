package postgres

import (
	"decm-database/go/generated"

	"apps/backend/common/pgclient"
)

type Repository struct {
	db               pgclient.Client
	queries          *generated.Queries
	piiEncryptionKey string
}

func NewRepository(db pgclient.Client, piiEncryptionKey string) *Repository {
	return &Repository{
		db:               db,
		queries:          generated.New(db),
		piiEncryptionKey: piiEncryptionKey,
	}
}
