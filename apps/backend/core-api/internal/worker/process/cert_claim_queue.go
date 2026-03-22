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

// CertificateClaimProcessor is the subset of certificate usecase behaviour
// needed by the cert claim queue.
type CertificateClaimProcessor interface {
	ProcessQueuedCertificateClaim(ctx context.Context, pending entity.PendingCertificateClaim) error
	// AbortCertificateClaim marks the pending signature as aborted with the
	// given reason so it is never retried.
	AbortCertificateClaim(ctx context.Context, signatureId uuid.UUID, reason entity.UserSignatureAbortReason) error
}

// CertClaimQueueProcessor drains the pending certificate-claim signature queue.
type CertClaimQueueProcessor struct {
	certProcessor      CertificateClaimProcessor
	eventJoinProcessor EventJoinProcessor // for on-chain join checks
	userSigDg          offchain_datagateway.UserSignatureDataGateway
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway
	logger             *slog.Logger
	isGasSuitable      func(ctx context.Context) bool
}

func NewCertClaimQueueProcessor(
	certProcessor CertificateClaimProcessor,
	eventJoinProcessor EventJoinProcessor,
	userSigDg offchain_datagateway.UserSignatureDataGateway,
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway,
	logger *slog.Logger,
	isGasSuitable func(ctx context.Context) bool,
) *CertClaimQueueProcessor {
	return &CertClaimQueueProcessor{
		certProcessor:      certProcessor,
		eventJoinProcessor: eventJoinProcessor,
		userSigDg:          userSigDg,
		blockchainClientDg: blockchainClientDg,
		logger:             logger,
		isGasSuitable:      isGasSuitable,
	}
}

// Run processes all pending certificate-claim signatures one by one,
// re-checking gas before each broadcast.
func (p *CertClaimQueueProcessor) Run(ctx context.Context) {
	pending, err := p.userSigDg.GetPendingCertificateClaimSignatures(ctx)
	if err != nil {
		p.logger.ErrorContext(ctx, "Failed to fetch pending certificate claims",
			slog.String("error", err.Error()),
		)
		return
	}

	if len(pending) == 0 {
		return
	}

	p.logger.InfoContext(ctx, "Processing pending certificate claims",
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
			p.logger.WarnContext(ctx, "Certificate claim time deadline expired, marking as expired",
				slog.String("signature_id", item.SignatureId.String()),
				slog.String("certificate_id", item.CertificateId.String()),
			)
			expired = true
		}

		if !expired && blockErr == nil && item.DeadlineBlock != nil && currentBlock >= uint64(*item.DeadlineBlock) {
			p.logger.WarnContext(ctx, "Certificate claim block deadline exceeded, marking as expired",
				slog.String("signature_id", item.SignatureId.String()),
				slog.String("certificate_id", item.CertificateId.String()),
				slog.Uint64("current_block", currentBlock),
				slog.Int("deadline_block", int(*item.DeadlineBlock)),
			)
			expired = true
		}

		if expired {
			expiredAt := now
			if _, markErr := p.userSigDg.UpdateUserSignatureMarkAsExpiredAt(ctx, item.SignatureId, &expiredAt); markErr != nil {
				p.logger.ErrorContext(ctx, "Failed to mark expired claim",
					slog.String("signature_id", item.SignatureId.String()),
					slog.String("error", markErr.Error()),
				)
			}
			continue
		}

		if !p.isGasSuitable(ctx) {
			p.logger.InfoContext(ctx, "Gas too high mid-queue, aborting certificate claim processing")
			return
		}

		joined, joinErr := p.eventJoinProcessor.IsParticipantJoinedOnChain(ctx, item.EventId, item.WalletAddress)
		if joinErr != nil {
			p.logger.WarnContext(ctx, "Failed to check on-chain join status, skipping certificate claim",
				slog.String("signature_id", item.SignatureId.String()),
				slog.String("certificate_id", item.CertificateId.String()),
				slog.String("error", joinErr.Error()),
			)
			continue
		}
		if !joined {
			// The participant has not yet been recorded on-chain. Check
			// whether there is still a pending join tx in the queue.
			hasPending, pendingErr := p.eventJoinProcessor.HasPendingEventJoin(ctx, item.EventId, item.AuthenticationCredentialId)
			if pendingErr != nil {
				p.logger.WarnContext(ctx, "Failed to check pending join queue, deferring certificate claim",
					slog.String("signature_id", item.SignatureId.String()),
					slog.String("certificate_id", item.CertificateId.String()),
					slog.String("error", pendingErr.Error()),
				)
				continue
			}
			if hasPending {
				p.logger.InfoContext(ctx, "Participant not yet joined on chain but join is pending, deferring certificate claim",
					slog.String("signature_id", item.SignatureId.String()),
					slog.String("certificate_id", item.CertificateId.String()),
				)
				continue
			}
			// No pending join and not on-chain — the claim can never succeed.
			p.logger.WarnContext(ctx, "Participant not joined on chain and no pending join found, aborting certificate claim",
				slog.String("signature_id", item.SignatureId.String()),
				slog.String("certificate_id", item.CertificateId.String()),
			)
			if abortErr := p.certProcessor.AbortCertificateClaim(ctx, item.SignatureId, entity.AbortReasonParticipantNotJoined); abortErr != nil {
				p.logger.ErrorContext(ctx, "Failed to abort certificate claim",
					slog.String("signature_id", item.SignatureId.String()),
					slog.String("error", abortErr.Error()),
				)
			}
			continue
		}

		finish := metrics.MeasureWorkerOperation("blockchainSubmission.ProcessCertificateClaim")
		processErr := p.certProcessor.ProcessQueuedCertificateClaim(ctx, item)
		finish(processErr)
		if processErr != nil {
			p.logger.ErrorContext(ctx, "Failed to process certificate claim",
				slog.String("signature_id", item.SignatureId.String()),
				slog.String("certificate_id", item.CertificateId.String()),
				slog.String("error", processErr.Error()),
			)
			continue
		}

		p.logger.InfoContext(ctx, "Certificate claim processed successfully",
			slog.String("signature_id", item.SignatureId.String()),
			slog.String("certificate_id", item.CertificateId.String()),
		)
	}
}
