package worker

import (
	"apps/backend/core-api/internal/worker/process"
	"context"
	"log/slog"
	"time"

	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
)

// BlockchainSubmissionWorker polls pending blockchain submission queues and
// processes them sequentially to avoid nonce conflicts.  On each tick it:
//  1. Checks the current gas price — skips the tick if gas is too high.
//  2. Stops the ticker to prevent tick accumulation during processing.
//  3. Drains the event-join queue, checking gas before every broadcast.
//  4. Then drains the certificate-claim queue, checking gas before every broadcast.
//  5. Resets the ticker for the next interval.
//
// If gas becomes too high mid-queue the current run is aborted immediately and
// the worker resumes at the next normal polling interval.
type BlockchainSubmissionWorker struct {
	joinQueue          *process.JoinQueueProcessor
	certClaimQueue     *process.CertClaimQueueProcessor
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway
	maxGasGwei         float64
	logger             *slog.Logger
	pollInterval       time.Duration
}

// NewBlockchainSubmissionWorker creates a worker with the default poll interval.
func NewBlockchainSubmissionWorker(
	eventJoinProcessor process.EventJoinProcessor,
	certClaimProcessor process.CertificateClaimProcessor,
	userSigDg offchain_datagateway.UserSignatureDataGateway,
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway,
	maxGasGwei float64,
	logger *slog.Logger,
	pollInterval time.Duration,
) *BlockchainSubmissionWorker {
	w := &BlockchainSubmissionWorker{
		blockchainClientDg: blockchainClientDg,
		maxGasGwei:         maxGasGwei,
		logger:             logger,
		pollInterval:       pollInterval,
	}
	w.joinQueue = process.NewJoinQueueProcessor(eventJoinProcessor, userSigDg, blockchainClientDg, logger, w.isGasSuitable)
	w.certClaimQueue = process.NewCertClaimQueueProcessor(certClaimProcessor, eventJoinProcessor, userSigDg, blockchainClientDg, logger, w.isGasSuitable)
	return w
}

// Start runs the polling loop until ctx is cancelled.
func (w *BlockchainSubmissionWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.InfoContext(ctx, "Blockchain submission worker stopped")
			return
		case <-ticker.C:
			if !w.isGasSuitable(ctx) {
				continue
			}
			// Stop the ticker while processing so that slow blockchain
			// transactions don't cause ticks to accumulate.
			ticker.Stop()
			w.processQueues(ctx)
			ticker.Reset(w.pollInterval)
		}
	}
}

// isGasSuitable returns true when the current network gas price is at or below
// the configured maximum, making it safe to submit transactions.
func (w *BlockchainSubmissionWorker) isGasSuitable(ctx context.Context) bool {
	gasInfo, err := w.blockchainClientDg.GetGasPrice(ctx)
	if err != nil {
		w.logger.WarnContext(ctx, "Failed to get gas price, skipping tick",
			slog.String("error", err.Error()),
		)
		return false
	}
	if gasInfo.MaxFeePerGasGwei > w.maxGasGwei {
		w.logger.InfoContext(ctx, "Gas price too high, skipping blockchain submissions",
			slog.Float64("current_gwei", gasInfo.MaxFeePerGasGwei),
			slog.Float64("max_gwei", w.maxGasGwei),
		)
		return false
	}
	return true
}

// processQueues drains the join-event queue first, then the certificate-claim
// queue.  Processing is sequential to avoid nonce conflicts.
func (w *BlockchainSubmissionWorker) processQueues(ctx context.Context) {
	if !w.joinQueue.Run(ctx) {
		return
	}
	w.certClaimQueue.Run(ctx)
}
