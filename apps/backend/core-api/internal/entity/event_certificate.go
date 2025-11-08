package entity

import (
	"time"

	"github.com/google/uuid"
)

type EventCertificate struct {
	ID                      uuid.UUID  `json:"id"`
	EventID                 uuid.UUID  `json:"event_id"`
	ReceiverCredentialID    *uuid.UUID `json:"receiver_credential_id"`
	ReceiverEmail           *string    `json:"receiver_email"`
	Name                    *string    `json:"name"`
	AcademicInstitution     *string    `json:"academic_institution"`
	CertificateTitle        *string    `json:"certificate_title"`
	CertificateSubtitle     *string    `json:"certificate_subtitle"`
	EventContractAddress    string     `json:"event_contract_address"`
	EventCertificateAddress *string    `json:"event_certificate_address"`
	CertificateTokenID      *string    `json:"certificate_token_id"`
	CreatedAt               time.Time  `json:"created_at"`
	RevokedAt               *time.Time `json:"revoked_at"`
}
