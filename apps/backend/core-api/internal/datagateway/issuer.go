package datagateway

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/core-api/internal/entity"
)

type IssuerDataGateway interface {
	ListVerifiedIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error)
	ListIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error)
	SearchIssuerCredentialsByWalletAddress(ctx context.Context, searchQuery string, limitCount int, offsetCount int) ([]entity.AuthenticationCredential, error)
	ListAllIssuerCredentials(ctx context.Context, limitCount int) ([]entity.AuthenticationCredential, error)
	GetEventsByIssuerCredentialID(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetEventIssuersByCredentialIDRow, error)
	GetIssuerEventsWithDetails(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetIssuerEventsWithDetailsRow, error)
}
