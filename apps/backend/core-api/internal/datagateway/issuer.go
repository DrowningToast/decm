package datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"
)

type IssuerDataGateway interface {
	ListVerifiedIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error)
}
