package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CertificateClaimStatus struct {
	IsBroadcasted      bool       `json:"is_broadcasted"`
	BroadcastedAt      *time.Time `json:"broadcasted_at,omitempty"`
	Signature          string     `json:"signature"`
	SignMessage        string     `json:"sign_message"`
	CertificateTokenId *string    `json:"certificate_token_id,omitempty"`
	ClaimedAt          time.Time  `json:"claimed_at"`
}

func (uc *EventUsecase) IsCertificateClaimed(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*CertificateClaimStatus, error) {
	// 1. Get certificate details and signature in one query
	cert, err := uc.EventCertificateDataGateway.GetEventCertificateWithSignature(ctx, eventID, credentialID)
	if err != nil {
		return nil, err
	}

	// Returns nil if not found
	if cert == nil {
		return nil, nil
	}

	return &CertificateClaimStatus{
		IsBroadcasted:      cert.IsBroadcasted || cert.TokenId != nil,
		BroadcastedAt:      cert.BroadcastedAt,
		Signature:          cert.Signature,
		SignMessage:        cert.SignMessage,
		CertificateTokenId: cert.TokenId,
		ClaimedAt:          cert.JoinedAt,
	}, nil
}
