package process

import (
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/handler/metrics"
	"context"
	"log/slog"
	"time"

	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"

	"github.com/google/uuid"
)

// EventJoinProcessor is the subset of event-registration usecase behaviour
// needed by the join queue.
type EventJoinProcessor interface {
	ProcessQueuedEventJoin(ctx context.Context, pending entity.PendingEventJoin) error
	IsParticipantJoinedOnChain(ctx context.Context, eventId uuid.UUID, walletAddress string) (bool, error)
	// HasPendingEventJoin reports whether the given credential has an
	// unbroadcasted, unexpired, unaborted event-join signature waiting in the
	// queue for the specified event.
	HasPendingEventJoin(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (bool, error)
}

// JoinQueueProcessor drains the pending event-join signature queue.
type JoinQueueProcessor struct {
	processor          EventJoinProcessor
	userSigDg          offchain_datagateway.UserSignatureDataGateway
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway
	logger             *slog.Logger
	isGasSuitable      func(ctx context.Context) bool
}

func NewJoinQueueProcessor(
	processor EventJoinProcessor,
	userSigDg offchain_datagateway.UserSignatureDataGateway,
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway,
	logger *slog.Logger,
	isGasSuitable func(ctx context.Context) bool,
) *JoinQueueProcessor {
	return &JoinQueueProcessor{
		processor:          processor,
		userSigDg:          userSigDg,
		blockchainClientDg: blockchainClientDg,
		logger:             logger,
		isGasSuitable:      isGasSuitable,
	}
}

// Run processes all pending event-join signatures one by one, re-checking gas
// before each broadcast. Returns false if processing was aborted due to gas
// being too high.
func (p *JoinQueueProcessor) Run(ctx context.Context) bool {
	pending, err := p.userSigDg.GetPendingEventJoinSignatures(ctx)
	if err != nil {
		p.logger.ErrorContext(ctx, "Failed to fetch pending event join signatures",
			slog.String("error", err.Error()),
		)
		return true // DB error doesn't mean gas is bad; proceed to cert queue
	}

	if len(pending) == 0 {
		return true
	}

	p.logger.InfoContext(ctx, "Processing pending event join signatures",
		slog.Int("count", len(pending)),
	)

	currentBlock, blockErr := p.blockchainClientDg.GetCurrentBlockNumber(ctx)
	if blockErr != nil {
		p.logger.WarnContext(ctx, "Failed to get current block number, block deadline check will be skipped",
			slog.String("error", blockErr.Error()),
		)
	}

	now := time.Now()
	for _, item := range pending {
		expired := false

		if item.EstimatedDeadline != nil && now.After(*item.EstimatedDeadline) {
			p.logger.WarnContext(ctx, "Event join signature time deadline expired, marking as expired",
				slog.String("signature_id", item.SignatureId.String()),
				slog.String("event_id", item.EventId.String()),
			)
			expired = true
		}

		if !expired && blockErr == nil && item.DeadlineBlock != nil && currentBlock >= uint64(*item.DeadlineBlock) {
			p.logger.WarnContext(ctx, "Event join signature block deadline exceeded, marking as expired",
				slog.String("signature_id", item.SignatureId.String()),
				slog.String("event_id", item.EventId.String()),
				slog.Uint64("current_block", currentBlock),
				slog.Int("deadline_block", int(*item.DeadlineBlock)),
			)
			expired = true
		}

		if expired {
			expiredAt := now
			if _, markErr := p.userSigDg.UpdateUserSignatureMarkAsExpiredAt(ctx, item.SignatureId, &expiredAt); markErr != nil {
				p.logger.ErrorContext(ctx, "Failed to mark expired join signature",
					slog.String("signature_id", item.SignatureId.String()),
					slog.String("error", markErr.Error()),
				)
			}
			continue
		}

		if !p.isGasSuitable(ctx) {
			p.logger.InfoContext(ctx, "Gas too high mid-queue, aborting event join processing")
			return false
		}

		finish := metrics.MeasureWorkerOperation("blockchainSubmission.ProcessEventJoin")
		processErr := p.processor.ProcessQueuedEventJoin(ctx, item)
		finish(processErr)
		if processErr != nil {
			p.logger.ErrorContext(ctx, "Failed to process event join",
				slog.String("signature_id", item.SignatureId.String()),
				slog.String("event_id", item.EventId.String()),
				slog.String("error", processErr.Error()),
			)
			continue
		}

		p.logger.InfoContext(ctx, "Event join processed successfully",
			slog.String("signature_id", item.SignatureId.String()),
			slog.String("event_id", item.EventId.String()),
		)
	}
	return true
}
