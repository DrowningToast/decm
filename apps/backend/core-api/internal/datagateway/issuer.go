package datagateway

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/core-api/internal/entity"
)

type IssuerDataGateway interface {
	ListVerifiedIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error)
	GetEventsByIssuerCredentialID(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetEventIssuersByCredentialIDRow, error)
}
