package eventcontract_datagateway

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type CreateContractParams struct {
	AccessManagerContractAddress common.Address
	HostAddress                  common.Address
	EventName                    string
	EventDescription             string
	SeatsCount                   int64
}

type CreateContractResponse struct {
	EventContractAddress         common.Address
	AccessManagerContractAddress common.Address
}

type EventContractFactoryDataGateway interface {
	GetContract(address common.Address) (EventContractDataGateway, error)
	CreateContract(ctx context.Context, params CreateContractParams) (*CreateContractResponse, error)
}
