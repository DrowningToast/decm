package entity

import (
	"time"

	"github.com/google/uuid"
)

type EventCertificate struct {
	Id                      uuid.UUID  `json:"id"`
	EventId                 uuid.UUID  `json:"event_id"`
	EventName               *string    `json:"event_name,omitempty"`
	ReceiverCredentialId    *uuid.UUID `json:"receiver_credential_id,omitempty"`
	ReceiverEmail           *string    `json:"receiver_email,omitempty"`
	Name                    *string    `json:"name,omitempty"`
	AcademicInstitution     *string    `json:"academic_institution,omitempty"`
	CertificateTitle        *string    `json:"certificate_title,omitempty"`
	CertificateSubtitle     *string    `json:"certificate_subtitle,omitempty"`
	EventContractAddress    string     `json:"event_contract_address"`
	EventCertificateAddress *string    `json:"event_certificate_address,omitempty"`
	CertificateTokenId      *string    `json:"certificate_token_id,omitempty"`
	CertificateDigest       *string    `json:"certificate_digest,omitempty"`
	UserClaimSignatureId    *uuid.UUID `json:"user_claim_signature_id,omitempty"`
	InboxMessageId          *uuid.UUID `json:"inbox_message_id,omitempty"`
	BroadcastedAt           *time.Time `json:"broadcasted_at,omitempty"`
	SignatureCreatedAt      *time.Time `json:"signature_created_at,omitempty"`
	EstimatedDeadline       *time.Time `json:"estimated_deadline,omitempty"`
	AbortedAt               *time.Time `json:"aborted_at,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	RevokedAt               *time.Time `json:"revoked_at,omitempty"`
}
