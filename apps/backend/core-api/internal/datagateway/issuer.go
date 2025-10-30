package datagateway

import (
	"context"

	"apps/backend/core-api/internal/entity"
)

type IssuerDataGateway interface {
	ListVerifiedIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error)
}
