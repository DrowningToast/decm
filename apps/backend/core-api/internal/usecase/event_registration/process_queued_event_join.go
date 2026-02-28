package event_registration

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"encoding/hex"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
)

// ProcessQueuedEventJoin is the exported bridge called by the blockchain
// submission worker. It submits the stored participant signature on-chain via
// AddParticipant, then marks the signature as broadcasted on success.
func (uc *EventRegistrationUsecase) ProcessQueuedEventJoin(ctx context.Context, pending entity.PendingEventJoin) error {
	// 1. Look up the event contract to obtain the on-chain contract address.
	entityEventContract, err := uc.EventContractDg.GetEventContractByEventID(ctx, pending.EventId)
	if err != nil {
		return errors.Wrap(err, "failed to get event contract")
	}
	if entityEventContract == nil {
		return errors.New("event contract not found")
	}

	// 2. Get the live contract instance from the factory.
	contractInstance, err := uc.EventContractFactoryDg.GetContract(common.HexToAddress(entityEventContract.EventContractAddress))
	if err != nil {
		return errors.Wrap(err, "failed to get event contract instance")
	}

	// 3. Derive the participant's address from the stored wallet address
	//    (already validated during the original join flow).
	participantAddress := common.HexToAddress(pending.WalletAddress)

	// 4. Idempotency guard: if the participant is already recorded on-chain
	//    (e.g. from a previous worker run that crashed after the tx mined),
	//    just mark the signature as broadcasted and return.
	participants, err := contractInstance.GetParticipants(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get on-chain participants")
	}
	for _, p := range participants {
		if p.Cmp(participantAddress) == 0 {
			now := time.Now()
			_, _ = uc.UserSignatureDg.UpdateUserSignatureBroadcastedAt(ctx, pending.SignatureId, &now)
			return nil
		}
	}

	// 5. Decode the hex-encoded signature back to raw bytes.
	sigBytes, err := hex.DecodeString(pending.Signature)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex signature")
	}

	// 6. Submit the on-chain transaction. AddParticipant internally creates a
	//    transactor, broadcasts, and waits for the receipt.
	if err := contractInstance.AddParticipant(ctx, participantAddress, pending.SignMessage, sigBytes); err != nil {
		return errors.Wrap(err, "failed to add participant to blockchain")
	}

	// 7. Mark the signature as processed so the worker does not retry it.
	now := time.Now()
	if _, err = uc.UserSignatureDg.UpdateUserSignatureBroadcastedAt(ctx, pending.SignatureId, &now); err != nil {
		return errors.Wrap(err, "failed to update signature broadcasted_at")
	}

	return nil
}
