package entity

import (
	"github.com/google/uuid"
)

type EventCertificateSignature struct {
	ID                 uuid.UUID `json:"id"`
	EventCertificateID uuid.UUID `json:"event_certificate_id"`
	IssuerCredentialID uuid.UUID `json:"issuer_credential_id"`
	IssuerSignature    *string   `json:"issuer_signature"`
	HostSignature      string    `json:"host_signature"`
	SignMessage        *string   `json:"sign_message"`
	SignMessageDigest  *string   `json:"sign_message_digest"`
}
