package event

import (
	"github.com/google/uuid"
)

type EventContractResponse struct {
	ID                           uuid.UUID `json:"id"`
	EventID                      uuid.UUID `json:"event_id"`
	AccessManagerContractAddress string    `json:"access_manager_contract_address"`
	EventContractAddress         string    `json:"event_contract_address"`
	TicketContractAddress        *string   `json:"ticket_contract_address,omitempty"`
	CertificateContractAddress   *string   `json:"certificate_contract_address,omitempty"`
	CreatedAt                    string    `json:"created_at"`
	UpdatedAt                    string    `json:"updated_at"`
}

type CreateEventContractRequest struct {
	AccessManagerContractAddress string `json:"access_manager_contract_address" validate:"required"`
	EventContractAddress         string `json:"event_contract_address" validate:"required"`
	TicketContractAddress        string `json:"ticket_contract_address"`
	CertificateContractAddress   string `json:"certificate_contract_address"`
}

func (r *CreateEventContractRequest) IsValid() error {
	// Add validation logic here if needed
	return nil
}
