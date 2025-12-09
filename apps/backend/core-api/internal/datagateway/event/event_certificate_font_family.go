package datagateway

import (
	"context"
	"decm-database/go/generated"
)

type EventCertificateFontFamilyDataGateway interface {
	GetAllEventCertificateFontFamilies(ctx context.Context) ([]generated.EventCertificateFontFamily, error)
	GetEventCertificateFontFamilyByID(ctx context.Context, id int32) (*generated.EventCertificateFontFamily, error)
	GetDefaultEventCertificateFontFamily(ctx context.Context) (*generated.EventCertificateFontFamily, error)
}


