package entity

import (
	"time"

	"github.com/google/uuid"
)

// EventContract represents the contract addresses for an event
type EventContract struct {
	ID                           uuid.UUID `json:"id"`
	EventID                      uuid.UUID `json:"event_id"`
	AccessManagerContractAddress string    `json:"access_manager_contract_address"`
	EventContractAddress         string    `json:"event_contract_address"`
	TicketContractAddress        *string   `json:"ticket_contract_address,omitempty"`
	CertificateContractAddress   *string   `json:"certificate_contract_address,omitempty"`
	CreatedAt                    time.Time `json:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at"`
}
