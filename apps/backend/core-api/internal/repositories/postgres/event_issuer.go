package postgres

import (
	"context"
	"decm-database/go/generated"

	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"

	"github.com/google/uuid"
)

var _ event_datagateway.EventIssuerDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventIssuer(ctx context.Context, params generated.CreateEventIssuerParams) (*generated.EventIssuer, error) {
	result, err := r.queries.CreateEventIssuer(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) GetEventIssuersByEventID(ctx context.Context, eventID uuid.UUID) ([]generated.EventIssuer, error) {
	return r.queries.GetEventIssuersByEventID(ctx, eventID)
}

func (r *Repository) GetEventIssuerByID(ctx context.Context, id uuid.UUID) (generated.EventIssuer, error) {
	return r.queries.GetEventIssuerByID(ctx, id)
}

func (r *Repository) UpdateEventIssuer(ctx context.Context, params generated.UpdateEventIssuerParams) (*generated.EventIssuer, error) {
	result, err := r.queries.UpdateEventIssuer(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) DeleteEventIssuer(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEventIssuer(ctx, id)
}

func (r *Repository) GetEventIssuerByEventIDAndIssuerCredentialID(ctx context.Context, eventID uuid.UUID, issuerCredentialID uuid.UUID) (generated.EventIssuer, error) {
	return r.queries.GetEventIssuerByEventIDAndIssuerCredentialID(ctx, generated.GetEventIssuerByEventIDAndIssuerCredentialIDParams{
		EventID:            eventID,
		IssuerCredentialID: issuerCredentialID,
	})
}

// UpdateEventIssuerSigningStatus updates is_signed field for an event issuer
func (r *Repository) UpdateEventIssuerSigningStatus(ctx context.Context, eventID uuid.UUID, issuerCredentialID uuid.UUID, isSigned int32) error {
	// First get issuer record
	issuer, err := r.queries.GetEventIssuerByEventIDAndIssuerCredentialID(ctx, generated.GetEventIssuerByEventIDAndIssuerCredentialIDParams{
		EventID:            eventID,
		IssuerCredentialID: issuerCredentialID,
	})
	if err != nil {
		return err
	}

	// Update is_signed field
	_, err = r.queries.UpdateEventIssuer(ctx, generated.UpdateEventIssuerParams{
		ID:       issuer.ID,
		IsSigned: isSigned,
	})

	return err
}

// ResetAllEventIssuersSigningStatus resets is_signed to 0 for all issuers of an event
func (r *Repository) ResetAllEventIssuersSigningStatus(ctx context.Context, eventID uuid.UUID) error {
	return r.queries.ResetAllEventIssuersSigningStatus(ctx, eventID)
}

// HasSignedIssuers checks if the event has at least one signed issuer
func (r *Repository) HasSignedIssuers(ctx context.Context, eventID uuid.UUID) (bool, error) {
	return r.queries.HasSignedIssuers(ctx, eventID)
}

// GetSignedIssuersCount returns the count of signed issuers for an event
func (r *Repository) GetSignedIssuersCount(ctx context.Context, eventID uuid.UUID) (int64, error) {
	return r.queries.GetSignedIssuersCount(ctx, eventID)
}

// GetTotalIssuersCount returns the total count of issuers for an event
func (r *Repository) GetTotalIssuersCount(ctx context.Context, eventID uuid.UUID) (int64, error) {
	return r.queries.GetTotalIssuersCount(ctx, eventID)
}

// AllIssuersHaveSigned checks if all assigned issuers have signed
func (r *Repository) AllIssuersHaveSigned(ctx context.Context, eventID uuid.UUID) (bool, error) {
	return r.queries.AllIssuersHaveSigned(ctx, eventID)
}
