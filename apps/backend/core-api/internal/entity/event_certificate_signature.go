package entity

import (
	"github.com/google/uuid"
)

type EventCertificateSignature struct {
	Id                       uuid.UUID `json:"id"`
	EventCertificateConfigId uuid.UUID `json:"event_certificate_config_id"`
	IssuerCredentialId       uuid.UUID `json:"issuer_credential_id"`
	IssuerSignature          *string   `json:"issuer_signature,omitempty"`
	HostSignature            string    `json:"host_signature"`
	SignMessage              *string   `json:"sign_message,omitempty"`
	SignMessageDigest        *string   `json:"sign_message_digest,omitempty"`
}
