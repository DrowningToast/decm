package worker

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"log/slog"
	"time"

	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"

	"github.com/google/uuid"
)

// EventJoinProcessor is the subset of event-registration usecase behaviour
// needed by the worker.
type EventJoinProcessor interface {
	ProcessQueuedEventJoin(ctx context.Context, pending entity.PendingEventJoin) error
	IsParticipantJoinedOnChain(ctx context.Context, eventId uuid.UUID, walletAddress string) (bool, error)
	// HasPendingEventJoin reports whether the given credential has an
	// unbroadcasted, unexpired, unaborted event-join signature waiting in the
	// queue for the specified event.
	HasPendingEventJoin(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (bool, error)
}

// CertificateClaimProcessor is the subset of event usecase behaviour needed
// by the worker.
type CertificateClaimProcessor interface {
	ProcessQueuedCertificateClaim(ctx context.Context, pending entity.PendingCertificateClaim) error
	// AbortCertificateClaim marks the pending signature as aborted with the
	// given reason so it is never retried.
	AbortCertificateClaim(ctx context.Context, signatureId uuid.UUID, reason entity.UserSignatureAbortReason) error
}

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
	eventJoinProcessor        EventJoinProcessor
	certificateClaimProcessor CertificateClaimProcessor
	userSigDg                 offchain_datagateway.UserSignatureDataGateway
	blockchainClientDg        blockchainclient_datagateway.BlockchainClientDataGateway
	maxGasGwei                float64
	logger                    *slog.Logger
	pollInterval              time.Duration
}

// NewBlockchainSubmissionWorker creates a worker with the default poll interval.
func NewBlockchainSubmissionWorker(
	eventJoinProcessor EventJoinProcessor,
	certificateClaimProcessor CertificateClaimProcessor,
	userSigDg offchain_datagateway.UserSignatureDataGateway,
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway,
	maxGasGwei float64,
	logger *slog.Logger,
	pollInterval time.Duration,
) *BlockchainSubmissionWorker {
	return &BlockchainSubmissionWorker{
		eventJoinProcessor:        eventJoinProcessor,
		certificateClaimProcessor: certificateClaimProcessor,
		userSigDg:                 userSigDg,
		blockchainClientDg:        blockchainClientDg,
		maxGasGwei:                maxGasGwei,
		logger:                    logger,
		pollInterval:              pollInterval,
	}
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
	if !w.processJoinQueue(ctx) {
		return
	}
	w.processCertClaimQueue(ctx)
}

// processJoinQueue processes all pending event-join signatures one by one,
// re-checking gas before each broadcast. Returns false if processing was
// aborted due to gas being too high.
func (w *BlockchainSubmissionWorker) processJoinQueue(ctx context.Context) bool {
	pending, err := w.userSigDg.GetPendingEventJoinSignatures(ctx)
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to fetch pending event join signatures",
			slog.String("error", err.Error()),
		)
		return true // DB error doesn't mean gas is bad; proceed to cert queue
	}

	if len(pending) == 0 {
		return true
	}

	w.logger.InfoContext(ctx, "Processing pending event join signatures",
		slog.Int("count", len(pending)),
	)

	currentBlock, blockErr := w.blockchainClientDg.GetCurrentBlockNumber(ctx)
	if blockErr != nil {
		w.logger.WarnContext(ctx, "Failed to get current block number, block deadline check will be skipped",
			slog.String("error", blockErr.Error()),
		)
	}

	now := time.Now()
	for _, p := range pending {
		expired := false

		if p.EstimatedDeadline != nil && now.After(*p.EstimatedDeadline) {
			w.logger.WarnContext(ctx, "Event join signature time deadline expired, marking as expired",
				slog.String("signature_id", p.SignatureId.String()),
				slog.String("event_id", p.EventId.String()),
			)
			expired = true
		}

		if !expired && blockErr == nil && p.DeadlineBlock != nil && currentBlock >= uint64(*p.DeadlineBlock) {
			w.logger.WarnContext(ctx, "Event join signature block deadline exceeded, marking as expired",
				slog.String("signature_id", p.SignatureId.String()),
				slog.String("event_id", p.EventId.String()),
				slog.Uint64("current_block", currentBlock),
				slog.Int("deadline_block", int(*p.DeadlineBlock)),
			)
			expired = true
		}

		if expired {
			expiredAt := now
			if _, markErr := w.userSigDg.UpdateUserSignatureMarkAsExpiredAt(ctx, p.SignatureId, &expiredAt); markErr != nil {
				w.logger.ErrorContext(ctx, "Failed to mark expired join signature",
					slog.String("signature_id", p.SignatureId.String()),
					slog.String("error", markErr.Error()),
				)
			}
			continue
		}

		if !w.isGasSuitable(ctx) {
			w.logger.InfoContext(ctx, "Gas too high mid-queue, aborting event join processing")
			return false
		}

		if err := w.eventJoinProcessor.ProcessQueuedEventJoin(ctx, p); err != nil {
			w.logger.ErrorContext(ctx, "Failed to process event join",
				slog.String("signature_id", p.SignatureId.String()),
				slog.String("event_id", p.EventId.String()),
				slog.String("error", err.Error()),
			)
			continue
		}

		w.logger.InfoContext(ctx, "Event join processed successfully",
			slog.String("signature_id", p.SignatureId.String()),
			slog.String("event_id", p.EventId.String()),
		)
	}
	return true
}

// processCertClaimQueue processes all pending certificate-claim signatures one
// by one, re-checking gas before each broadcast.
func (w *BlockchainSubmissionWorker) processCertClaimQueue(ctx context.Context) {
	pending, err := w.userSigDg.GetPendingCertificateClaimSignatures(ctx)
	if err != nil {
		w.logger.ErrorContext(ctx, "Failed to fetch pending certificate claims",
			slog.String("error", err.Error()),
		)
		return
	}

	if len(pending) == 0 {
		return
	}

	w.logger.InfoContext(ctx, "Processing pending certificate claims",
		slog.Int("count", len(pending)),
	)

	currentBlock, blockErr := w.blockchainClientDg.GetCurrentBlockNumber(ctx)
	if blockErr != nil {
		w.logger.WarnContext(ctx, "Failed to get current block number, block deadline check will be skipped",
			slog.String("error", blockErr.Error()),
		)
	}

	now := time.Now()
	for _, p := range pending {
		expired := false

		if p.EstimatedDeadline != nil && now.After(*p.EstimatedDeadline) {
			w.logger.WarnContext(ctx, "Certificate claim time deadline expired, marking as expired",
				slog.String("signature_id", p.SignatureId.String()),
				slog.String("certificate_id", p.CertificateId.String()),
			)
			expired = true
		}

		if !expired && blockErr == nil && p.DeadlineBlock != nil && currentBlock >= uint64(*p.DeadlineBlock) {
			w.logger.WarnContext(ctx, "Certificate claim block deadline exceeded, marking as expired",
				slog.String("signature_id", p.SignatureId.String()),
				slog.String("certificate_id", p.CertificateId.String()),
				slog.Uint64("current_block", currentBlock),
				slog.Int("deadline_block", int(*p.DeadlineBlock)),
			)
			expired = true
		}

		if expired {
			expiredAt := now
			if _, markErr := w.userSigDg.UpdateUserSignatureMarkAsExpiredAt(ctx, p.SignatureId, &expiredAt); markErr != nil {
				w.logger.ErrorContext(ctx, "Failed to mark expired claim",
					slog.String("signature_id", p.SignatureId.String()),
					slog.String("error", markErr.Error()),
				)
			}
			continue
		}

		if !w.isGasSuitable(ctx) {
			w.logger.InfoContext(ctx, "Gas too high mid-queue, aborting certificate claim processing")
			return
		}

		joined, joinErr := w.eventJoinProcessor.IsParticipantJoinedOnChain(ctx, p.EventId, p.WalletAddress)
		if joinErr != nil {
			w.logger.WarnContext(ctx, "Failed to check on-chain join status, skipping certificate claim",
				slog.String("signature_id", p.SignatureId.String()),
				slog.String("certificate_id", p.CertificateId.String()),
				slog.String("error", joinErr.Error()),
			)
			continue
		}
		if !joined {
			// The participant has not yet been recorded on-chain.  Check
			// whether there is still a pending join tx in the queue.
			hasPending, pendingErr := w.eventJoinProcessor.HasPendingEventJoin(ctx, p.EventId, p.AuthenticationCredentialId)
			if pendingErr != nil {
				w.logger.WarnContext(ctx, "Failed to check pending join queue, deferring certificate claim",
					slog.String("signature_id", p.SignatureId.String()),
					slog.String("certificate_id", p.CertificateId.String()),
					slog.String("error", pendingErr.Error()),
				)
				continue
			}
			if hasPending {
				w.logger.InfoContext(ctx, "Participant not yet joined on chain but join is pending, deferring certificate claim",
					slog.String("signature_id", p.SignatureId.String()),
					slog.String("certificate_id", p.CertificateId.String()),
				)
				continue
			}
			// No pending join and not on-chain — the claim can never succeed.
			w.logger.WarnContext(ctx, "Participant not joined on chain and no pending join found, aborting certificate claim",
				slog.String("signature_id", p.SignatureId.String()),
				slog.String("certificate_id", p.CertificateId.String()),
			)
			if abortErr := w.certificateClaimProcessor.AbortCertificateClaim(ctx, p.SignatureId, entity.AbortReasonParticipantNotJoined); abortErr != nil {
				w.logger.ErrorContext(ctx, "Failed to abort certificate claim",
					slog.String("signature_id", p.SignatureId.String()),
					slog.String("error", abortErr.Error()),
				)
			}
			continue
		}

		if err := w.certificateClaimProcessor.ProcessQueuedCertificateClaim(ctx, p); err != nil {
			w.logger.ErrorContext(ctx, "Failed to process certificate claim",
				slog.String("signature_id", p.SignatureId.String()),
				slog.String("certificate_id", p.CertificateId.String()),
				slog.String("error", err.Error()),
			)
			continue
		}

		w.logger.InfoContext(ctx, "Certificate claim processed successfully",
			slog.String("signature_id", p.SignatureId.String()),
			slog.String("certificate_id", p.CertificateId.String()),
		)
	}
}
