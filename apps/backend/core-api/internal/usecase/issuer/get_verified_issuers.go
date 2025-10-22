package issuer

import (
	"apps/backend/core-api/internal/entity"
	"context"
)

func (u *IssuerUsecase) GetVerifiedIssuers(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error) {
	issuers, err := u.IssuerDg.ListVerifiedIssuerProfiles(ctx, limitCount, offsetCount)

	if err != nil {
		return nil, err
	}

	return issuers, nil
}
