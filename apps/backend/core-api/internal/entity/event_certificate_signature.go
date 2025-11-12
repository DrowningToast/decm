package entity

import (
	"github.com/google/uuid"
)

type EventCertificateSignature struct {
	Id                 uuid.UUID `json:"id"`
	EventCertificateId uuid.UUID `json:"event_certificate_id"`
	IssuerCredentialId uuid.UUID `json:"issuer_credential_id"`
	IssuerSignature    *string   `json:"issuer_signature"`
	HostSignature      string    `json:"host_signature"`
	SignMessage        *string   `json:"sign_message"`
	SignMessageDigest  *string   `json:"sign_message_digest"`
}
