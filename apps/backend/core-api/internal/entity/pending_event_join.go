package entity

import (
	"time"

	"github.com/google/uuid"
)

// PendingEventJoin carries the joined data needed by the blockchain submission
// worker to add a participant on-chain without extra DB round-trips.
type PendingEventJoin struct {
	SignatureId                uuid.UUID
	AuthenticationCredentialId uuid.UUID
	Signature                  string // hex-encoded
	SignMessage                string
	DeadlineBlock              *int32
	EstimatedDeadline          *time.Time
	CreatedAt                  time.Time
	EventId                    uuid.UUID
	WalletAddress              string
}
