package datagateway

import (
	"context"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type CreateEventCertificateSignatureParameters struct {
	EventCertificateID uuid.UUID
	IssuerCredentialID uuid.UUID
	IssuerSignature    *string
	HostSignature      string
	SignMessage        *string
	SignMessageDigest  *string
}

type UpdateEventCertificateSignatureParameters struct {
	IssuerSignature   *string
	HostSignature     string
	SignMessage       *string
	SignMessageDigest *string
}

type EventCertificateSignatureDataGateway interface {
	CreateEventCertificateSignature(ctx context.Context, params CreateEventCertificateSignatureParameters) (*entity.EventCertificateSignature, error)
	GetEventCertificateSignatureByID(ctx context.Context, id uuid.UUID) (*entity.EventCertificateSignature, error)
	GetEventCertificateSignaturesByEventCertificateID(ctx context.Context, eventCertificateID uuid.UUID) ([]*entity.EventCertificateSignature, error)
	UpdateEventCertificateSignature(ctx context.Context, id uuid.UUID, params UpdateEventCertificateSignatureParameters) (*entity.EventCertificateSignature, error)
	UpdateEventCertificateIssuerSignature(ctx context.Context, id uuid.UUID, issuerSignature *string) (*entity.EventCertificateSignature, error)
	DeleteEventCertificateSignature(ctx context.Context, id uuid.UUID) error
}
