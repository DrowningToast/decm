package event

import (
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type EventIssuerResponse struct {
	ID                 uuid.UUID       `json:"id"`
	EventID            uuid.UUID       `json:"event_id"`
	IssuerCredentialID uuid.UUID       `json:"issuer_credential_id"`
	IssuerProfile      *entity.Profile `json:"issuer_profile"`
	IsSigned           int32           `json:"is_signed"`
	Signature          string          `json:"signature"`
	SignMessage        string          `json:"sign_message"`
	CreatedAt          string          `json:"created_at"`
	UpdatedAt          string          `json:"updated_at"`
}
