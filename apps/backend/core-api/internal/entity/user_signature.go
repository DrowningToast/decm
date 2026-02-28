package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserSignature struct {
	Id                         uuid.UUID  `json:"id"`
	AuthenticationCredentialId uuid.UUID  `json:"authentication_credential_id"`
	SignMessage                string     `json:"sign_message"`
	Signature                  string     `json:"signature"`
	DeadlineBlock              *int32     `json:"deadline_block,omitempty"`
	EstimatedDeadline          *time.Time `json:"estimated_deadline,omitempty"`
	BroadcastedAt              *time.Time `json:"broadcasted_at,omitempty"`
	MarkAsExpiredAt            *time.Time `json:"mark_as_expired_at,omitempty"`
	AbortedAt                  *time.Time `json:"aborted_at,omitempty"`
	AbortedReason              *int32     `json:"aborted_reason,omitempty"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at"`
}
