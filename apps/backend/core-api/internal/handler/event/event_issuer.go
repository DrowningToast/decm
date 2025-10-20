package event

import (
	"github.com/google/uuid"
)

type EventIssuerResponse struct {
	ID                 uuid.UUID `json:"id"`
	EventID            uuid.UUID `json:"event_id"`
	IssuerCredentialID uuid.UUID `json:"issuer_credential_id"`
	IsSigned           int32     `json:"is_signed"`
	Signature          string    `json:"signature"`
	SignMessage        string    `json:"sign_message"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
}
