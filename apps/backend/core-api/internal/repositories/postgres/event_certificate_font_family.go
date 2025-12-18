package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/pgerrutils"
	datagateway "apps/backend/core-api/internal/datagateway/event"
)

var _ datagateway.EventCertificateFontFamilyDataGateway = (*Repository)(nil)

func (r *Repository) GetAllEventCertificateFontFamilies(ctx context.Context) ([]generated.EventCertificateFontFamily, error) {
	result, err := r.queries.GetAllEventCertificateFontFamilies(ctx)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return result, nil
}

func (r *Repository) GetEventCertificateFontFamilyByID(ctx context.Context, id int32) (*generated.EventCertificateFontFamily, error) {
	result, err := r.queries.GetEventCertificateFontFamilyByID(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return &result, nil
}

func (r *Repository) GetDefaultEventCertificateFontFamily(ctx context.Context) (*generated.EventCertificateFontFamily, error) {
	result, err := r.queries.GetDefaultEventCertificateFontFamily(ctx)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	return &result, nil
}












