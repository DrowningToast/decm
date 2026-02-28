package entity

import (
	"time"

	"github.com/google/uuid"
)

// PendingCertificateClaim carries the joined data needed by the certificate
// claim worker to process a queued claim without additional DB round-trips.
type PendingCertificateClaim struct {
	SignatureId                uuid.UUID
	AuthenticationCredentialId uuid.UUID
	Signature                  string // hex-encoded
	SignMessage                string
	DeadlineBlock              *int32
	EstimatedDeadline          *time.Time
	CreatedAt                  time.Time
	CertificateId              uuid.UUID
	EventId                    uuid.UUID
	WalletAddress              string
}
