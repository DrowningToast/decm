package datagateway

import (
	"context"

	"apps/backend/core-api/internal/entity"
	"decm-database/go/generated"
)

type IssuerDataGateway interface {
	ListVerifiedIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error)
	GetEventsByIssuerCredentialID(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetEventIssuersByCredentialIDRow, error)
}
