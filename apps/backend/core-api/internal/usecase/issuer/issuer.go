package issuer

import (
	"apps/backend/core-api/internal/datagateway"
	"context"
	"decm-database/go/generated"
)

type IssuerEventViewModel struct {
	ID                     string `json:"id"`
	EventID                string `json:"event_id"`
	EventTitle             string `json:"event_title"`
	EventShortDescription  string `json:"event_short_description"`
	EventStartDate         string `json:"event_start_date"`
	EventEndDate           string `json:"event_end_date"`
	EventLocation          string `json:"event_location"`
	EventOwnerCredentialID string `json:"event_owner_credential_id"`
	OwnerName              string `json:"owner_name"`
	OwnerWalletAddress     string `json:"owner_wallet_address"`
	OwnerEmail             string `json:"owner_email"`
	OwnerGoogleEmail       string `json:"owner_google_email"`
	CertificateCount       int32  `json:"certificate_count"`
	IssuerCredentialID     string `json:"issuer_credential_id"`
	IsSigned               bool   `json:"is_signed"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
}

type IssuerUsecase struct {
	IssuerDg datagateway.IssuerDataGateway
}

func NewIssuerUsecase(issuerDg datagateway.IssuerDataGateway) *IssuerUsecase {
	return &IssuerUsecase{
		IssuerDg: issuerDg,
	}
}

func (u *IssuerUsecase) GetEventsByIssuerCredentialID(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetEventIssuersByCredentialIDRow, error) {
	events, err := u.IssuerDg.GetEventsByIssuerCredentialID(ctx, issuerCredentialID, limitCount, offsetCount)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (u *IssuerUsecase) GetIssuerEventsWithDetails(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetIssuerEventsWithDetailsRow, error) {
	events, err := u.IssuerDg.GetIssuerEventsWithDetails(ctx, issuerCredentialID, limitCount, offsetCount)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (u *IssuerUsecase) GetIssuerEventsViewModel(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]IssuerEventViewModel, error) {
	events, err := u.IssuerDg.GetIssuerEventsWithDetails(ctx, issuerCredentialID, limitCount, offsetCount)
	if err != nil {
		return nil, err
	}

	viewModels := make([]IssuerEventViewModel, len(events))
	for i, event := range events {
		// Combine owner first and last name
		ownerName := ""
		if event.OwnerFirstName.Valid && event.OwnerLastName.Valid {
			ownerName = event.OwnerFirstName.String + " " + event.OwnerLastName.String
		} else if event.OwnerFirstName.Valid {
			ownerName = event.OwnerFirstName.String
		} else if event.OwnerLastName.Valid {
			ownerName = event.OwnerLastName.String
		}

		// Get wallet address (plain text, not encrypted)
		ownerWalletAddress := ""
		if event.OwnerWalletAddress.Valid {
			ownerWalletAddress = event.OwnerWalletAddress.String
		}

		// Get email from profiles (decrypted in repository)
		ownerEmail := ""
		if event.OwnerEmail.Valid {
			ownerEmail = event.OwnerEmail.String
		}

		// Get Google OAuth email from google_connector_ref (decrypted in repository)
		ownerGoogleEmail := ""
		if event.OwnerGoogleConnectorRef.Valid {
			ownerGoogleEmail = event.OwnerGoogleConnectorRef.String
		}

		viewModels[i] = IssuerEventViewModel{
			ID:                     event.ID.String(),
			EventID:                event.EventID.String(),
			EventTitle:             event.EventTitle,
			EventShortDescription:  event.EventShortDescription,
			EventStartDate:         event.EventStartDate.Format("2006-01-02T15:04:05Z07:00"),
			EventEndDate:           event.EventEndDate.Format("2006-01-02T15:04:05Z07:00"),
			EventLocation:          event.EventLocation,
			EventOwnerCredentialID: event.EventOwnerCredentialID.String(),
			OwnerName:              ownerName,
			OwnerWalletAddress:     ownerWalletAddress,
			OwnerEmail:             ownerEmail,
			OwnerGoogleEmail:       ownerGoogleEmail,
			CertificateCount:       event.CertificateCount,
			IssuerCredentialID:     event.IssuerCredentialID.String(),
			IsSigned:               event.IsSigned == 1,
			CreatedAt:              event.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:              event.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return viewModels, nil
}
