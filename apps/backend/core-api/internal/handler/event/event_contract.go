package event

import (
	"github.com/google/uuid"
)

type EventContractResponse struct {
	ID                           uuid.UUID `json:"id"`
	EventID                      uuid.UUID `json:"event_id"`
	AccessManagerContractAddress string    `json:"access_manager_contract_address"`
	EventContractAddress         string    `json:"event_contract_address"`
	TicketContractAddress        string    `json:"ticket_contract_address"`
	CertificateContractAddress   string    `json:"certificate_contract_address"`
	CreatedAt                    string    `json:"created_at"`
	UpdatedAt                    string    `json:"updated_at"`
}
