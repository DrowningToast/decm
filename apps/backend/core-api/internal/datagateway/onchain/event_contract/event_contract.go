package eventcontract_datagateway

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type EventContractDataGateway interface {
	GetCurrentSeatsCount(ctx context.Context) (*int64, error)
	GetMaxSeatsCount(ctx context.Context) (*int64, error)
	GetParticipants(ctx context.Context) ([]common.Address, error)
	AddParticipant(ctx context.Context, participantAddress common.Address, signMessage string, signature []byte) error
	UpdateEvent(ctx context.Context, eventName string, eventDescription string, seatsCount int64, signMessage string, signature []byte) error
}
