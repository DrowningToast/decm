package event_registration

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

// IsParticipantJoinedOnChain returns true when the given wallet address is
// recorded as a participant on the event's smart contract.
func (uc *EventRegistrationUsecase) IsParticipantJoinedOnChain(ctx context.Context, eventId uuid.UUID, walletAddress string) (bool, error) {
	entityEventContract, err := uc.EventContractDg.GetEventContractByEventID(ctx, eventId)
	if err != nil {
		return false, errors.Wrap(err, "failed to get event contract")
	}
	if entityEventContract == nil {
		return false, errors.New("event contract not found")
	}

	contractInstance, err := uc.EventContractFactoryDg.GetContract(common.HexToAddress(entityEventContract.EventContractAddress))
	if err != nil {
		return false, errors.Wrap(err, "failed to get event contract instance")
	}

	participants, err := contractInstance.GetParticipants(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get on-chain participants")
	}

	participantAddress := common.HexToAddress(walletAddress)
	for _, p := range participants {
		if p.Cmp(participantAddress) == 0 {
			return true, nil
		}
	}
	return false, nil
}
