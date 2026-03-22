package entity

import (
	"time"

	"github.com/google/uuid"
)

type CertificateShare struct {
	Id                 uuid.UUID `json:"id"`
	EventCertificateId uuid.UUID `json:"event_certificate_id"`
	Active             bool      `json:"active"`
	Handle             string    `json:"handle"`
	Password           *string   `json:"password,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
