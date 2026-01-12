package worker

import (
	"context"
	"errors"
	"log/slog"
	"math/big"
	"testing"
	"time"

	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	"apps/backend/core-api/internal/entity"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

type mockEventJoinProcessor struct{ mock.Mock }

func (m *mockEventJoinProcessor) ProcessQueuedEventJoin(ctx context.Context, p entity.PendingEventJoin) error {
	return m.Called(ctx, p).Error(0)
}

func (m *mockEventJoinProcessor) IsParticipantJoinedOnChain(ctx context.Context, eventId uuid.UUID, walletAddress string) (bool, error) {
	args := m.Called(ctx, eventId, walletAddress)
	return args.Bool(0), args.Error(1)
}

func (m *mockEventJoinProcessor) HasPendingEventJoin(ctx context.Context, eventId uuid.UUID, credentialId uuid.UUID) (bool, error) {
	args := m.Called(ctx, eventId, credentialId)
	return args.Bool(0), args.Error(1)
}

type mockCertClaimProcessor struct{ mock.Mock }

func (m *mockCertClaimProcessor) ProcessQueuedCertificateClaim(ctx context.Context, p entity.PendingCertificateClaim) error {
	return m.Called(ctx, p).Error(0)
}

func (m *mockCertClaimProcessor) AbortCertificateClaim(ctx context.Context, signatureId uuid.UUID, reason entity.UserSignatureAbortReason) error {
	return m.Called(ctx, signatureId, reason).Error(0)
}

// mockUserSigDg implements offchain_datagateway.UserSignatureDataGateway.
// Only the methods exercised by the worker have expectations wired up.
type mockUserSigDg struct{ mock.Mock }

func (m *mockUserSigDg) GetPendingEventJoinSignatures(ctx context.Context) ([]entity.PendingEventJoin, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.PendingEventJoin), args.Error(1)
}

func (m *mockUserSigDg) GetPendingCertificateClaimSignatures(ctx context.Context) ([]entity.PendingCertificateClaim, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.PendingCertificateClaim), args.Error(1)
}

func (m *mockUserSigDg) UpdateUserSignatureMarkAsExpiredAt(ctx context.Context, id uuid.UUID, at *time.Time) (*entity.UserSignature, error) {
	args := m.Called(ctx, id, at)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSignature), args.Error(1)
}

// Remaining interface methods – not called by the worker; satisfy the interface.
func (m *mockUserSigDg) CreateUserSignature(ctx context.Context, p offchain_datagateway.CreateUserSignatureParameters) (*entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetUserSignatureByID(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetUserSignaturesByCredentialID(ctx context.Context, id uuid.UUID) ([]entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) UpdateUserSignatureBroadcastedAt(ctx context.Context, id uuid.UUID, at *time.Time) (*entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetPendingUserSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetStaleUserSignatures(ctx context.Context, before time.Time) ([]entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetBroadcastedUserSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetUserSignaturesByDeadlineBlockRange(ctx context.Context, min, max int32) ([]entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetUserSignaturesExpiringBefore(ctx context.Context, d time.Time) ([]entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetEventJoinSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetCertificateClaimSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetOrphanedUserSignatures(ctx context.Context, before time.Time) ([]entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetUserSignatureWithUsageDetails(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) GetRecentUserSignatureActivity(ctx context.Context, since time.Time) ([]entity.UserSignature, error) {
	return nil, nil
}
func (m *mockUserSigDg) DeleteUserSignature(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockUserSigDg) CountUserSignaturesByCredentialID(ctx context.Context, id uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *mockUserSigDg) CountPendingUserSignatures(ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *mockUserSigDg) CountBroadcastedUserSignatures(ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *mockUserSigDg) UpdateUserSignatureAbortedAt(ctx context.Context, id uuid.UUID, abortedAt time.Time, reason entity.UserSignatureAbortReason) (*entity.UserSignature, error) {
	return nil, nil
}

// mockBlockchainClientDg implements blockchainclient_datagateway.BlockchainClientDataGateway.
// Only GetGasPrice is called by the worker.
type mockBlockchainClientDg struct{ mock.Mock }

func (m *mockBlockchainClientDg) GetGasPrice(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*blockchainclient_datagateway.GasPriceInfo), args.Error(1)
}

func (m *mockBlockchainClientDg) GetCurrentBlockNumber(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}
func (m *mockBlockchainClientDg) GetCalculatedDeadlineBlock(ctx context.Context) (uint64, error) {
	return 0, nil
}
func (m *mockBlockchainClientDg) EstimateDeadlineTime(ctx context.Context, block uint64) (*time.Time, error) {
	return nil, nil
}
func (m *mockBlockchainClientDg) GetTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	return nil, nil
}
func (m *mockBlockchainClientDg) WeiToGwei(wei *big.Int) float64 { return 0 }

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelError + 1}))
}

func gasOk(gwei float64) *blockchainclient_datagateway.GasPriceInfo {
	return &blockchainclient_datagateway.GasPriceInfo{MaxFeePerGasGwei: gwei}
}

func newWorker(
	joinProc EventJoinProcessor,
	certProc CertificateClaimProcessor,
	sigDg offchain_datagateway.UserSignatureDataGateway,
	chainDg blockchainclient_datagateway.BlockchainClientDataGateway,
	maxGas float64,
) *BlockchainSubmissionWorker {
	w := NewBlockchainSubmissionWorker(joinProc, certProc, sigDg, chainDg, maxGas, discardLogger(), 60*time.Minute)
	return w
}

// pendingJoin returns a PendingEventJoin with optional time deadline.
func pendingJoin(deadline *time.Time) entity.PendingEventJoin {
	return entity.PendingEventJoin{
		SignatureId:       uuid.New(),
		EventId:           uuid.New(),
		EstimatedDeadline: deadline,
	}
}

// pendingJoinWithBlock returns a PendingEventJoin with a specific deadline block.
func pendingJoinWithBlock(block int32) entity.PendingEventJoin {
	return entity.PendingEventJoin{
		SignatureId:   uuid.New(),
		EventId:       uuid.New(),
		DeadlineBlock: &block,
	}
}

// pendingClaim returns a PendingCertificateClaim with optional time deadline.
func pendingClaim(deadline *time.Time) entity.PendingCertificateClaim {
	return entity.PendingCertificateClaim{
		SignatureId:       uuid.New(),
		CertificateId:     uuid.New(),
		EstimatedDeadline: deadline,
	}
}

// pendingClaimWithBlock returns a PendingCertificateClaim with a specific deadline block.
func pendingClaimWithBlock(block int32) entity.PendingCertificateClaim {
	return entity.PendingCertificateClaim{
		SignatureId:   uuid.New(),
		CertificateId: uuid.New(),
		DeadlineBlock: &block,
	}
}

func ptr[T any](v T) *T { return &v }

// ---------------------------------------------------------------------------
// isGasSuitable
// ---------------------------------------------------------------------------

func TestIsGasSuitable_ReturnsFalse_OnGetGasPriceError(t *testing.T) {
	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetGasPrice", mock.Anything).Return(nil, errors.New("rpc error"))

	w := newWorker(nil, nil, nil, chainDg, 30)
	assert.False(t, w.isGasSuitable(context.Background()))
	chainDg.AssertExpectations(t)
}

func TestIsGasSuitable_ReturnsFalse_WhenGasExceedsMax(t *testing.T) {
	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(50), nil)

	w := newWorker(nil, nil, nil, chainDg, 30)
	assert.False(t, w.isGasSuitable(context.Background()))
	chainDg.AssertExpectations(t)
}

func TestIsGasSuitable_ReturnsTrue_WhenGasWithinLimit(t *testing.T) {
	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(20), nil)

	w := newWorker(nil, nil, nil, chainDg, 30)
	assert.True(t, w.isGasSuitable(context.Background()))
	chainDg.AssertExpectations(t)
}

func TestIsGasSuitable_ReturnsTrue_WhenGasEqualsMax(t *testing.T) {
	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(30), nil)

	w := newWorker(nil, nil, nil, chainDg, 30)
	assert.True(t, w.isGasSuitable(context.Background()))
	chainDg.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// processJoinQueue
// ---------------------------------------------------------------------------

func TestProcessJoinQueue_ReturnsTrue_WhenQueueEmpty(t *testing.T) {
	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{}, nil)

	w := newWorker(nil, nil, sigDg, nil, 30)
	ok := w.processJoinQueue(context.Background())
	assert.True(t, ok)
	sigDg.AssertExpectations(t)
}

func TestProcessJoinQueue_ReturnsTrue_OnDBError(t *testing.T) {
	// A DB error is not a gas issue — worker should continue to the cert queue.
	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return(nil, errors.New("db error"))

	w := newWorker(nil, nil, sigDg, nil, 30)
	ok := w.processJoinQueue(context.Background())
	assert.True(t, ok)
	sigDg.AssertExpectations(t)
}

func TestProcessJoinQueue_MarksExpiredAndSkips(t *testing.T) {
	past := time.Now().Add(-1 * time.Minute)
	item := pendingJoin(&past)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item}, nil)
	sigDg.On("UpdateUserSignatureMarkAsExpiredAt", mock.Anything, item.SignatureId, mock.AnythingOfType("*time.Time")).
		Return(&entity.UserSignature{}, nil)

	joinProc := new(mockEventJoinProcessor)
	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)

	w := newWorker(joinProc, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.True(t, ok)
	joinProc.AssertNotCalled(t, "ProcessQueuedEventJoin") // must not broadcast expired item
	sigDg.AssertExpectations(t)
}

func TestProcessJoinQueue_ContinuesIfMarkExpiredFails(t *testing.T) {
	// Mark-as-expired DB failure should be logged but not abort the run.
	past := time.Now().Add(-1 * time.Minute)
	item := pendingJoin(&past)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item}, nil)
	sigDg.On("UpdateUserSignatureMarkAsExpiredAt", mock.Anything, item.SignatureId, mock.AnythingOfType("*time.Time")).
		Return(nil, errors.New("db error"))

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)

	w := newWorker(nil, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.True(t, ok)
	sigDg.AssertExpectations(t)
}

func TestProcessJoinQueue_ReturnsFalse_WhenGasTooHighBeforeBroadcast(t *testing.T) {
	item := pendingJoin(nil) // no deadline → not expired

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(99), nil) // exceeds max

	joinProc := new(mockEventJoinProcessor)

	w := newWorker(joinProc, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.False(t, ok)
	joinProc.AssertNotCalled(t, "ProcessQueuedEventJoin")
	sigDg.AssertExpectations(t)
	chainDg.AssertExpectations(t)
}

func TestProcessJoinQueue_ContinuesOnProcessError(t *testing.T) {
	item1 := pendingJoin(nil)
	item2 := pendingJoin(nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item1, item2}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("ProcessQueuedEventJoin", mock.Anything, item1).Return(errors.New("tx failed"))
	joinProc.On("ProcessQueuedEventJoin", mock.Anything, item2).Return(nil)

	w := newWorker(joinProc, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.True(t, ok)
	joinProc.AssertExpectations(t) // item2 must still be attempted
}

func TestProcessJoinQueue_ProcessesAllItems(t *testing.T) {
	item1 := pendingJoin(nil)
	item2 := pendingJoin(nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item1, item2}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("ProcessQueuedEventJoin", mock.Anything, item1).Return(nil)
	joinProc.On("ProcessQueuedEventJoin", mock.Anything, item2).Return(nil)

	w := newWorker(joinProc, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.True(t, ok)
	joinProc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// processCertClaimQueue
// ---------------------------------------------------------------------------

func TestProcessCertClaimQueue_NoOp_WhenQueueEmpty(t *testing.T) {
	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{}, nil)

	w := newWorker(nil, nil, sigDg, nil, 30)
	w.processCertClaimQueue(context.Background()) // must not panic
	sigDg.AssertExpectations(t)
}

func TestProcessCertClaimQueue_NoOp_OnDBError(t *testing.T) {
	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return(nil, errors.New("db error"))

	w := newWorker(nil, nil, sigDg, nil, 30)
	w.processCertClaimQueue(context.Background())
	sigDg.AssertExpectations(t)
}

func TestProcessCertClaimQueue_MarksExpiredAndSkips(t *testing.T) {
	past := time.Now().Add(-1 * time.Minute)
	item := pendingClaim(&past)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)
	sigDg.On("UpdateUserSignatureMarkAsExpiredAt", mock.Anything, item.SignatureId, mock.AnythingOfType("*time.Time")).
		Return(&entity.UserSignature{}, nil)

	certProc := new(mockCertClaimProcessor)
	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)

	w := newWorker(nil, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertNotCalled(t, "ProcessQueuedCertificateClaim")
	sigDg.AssertExpectations(t)
}

func TestProcessCertClaimQueue_Aborts_WhenGasTooHighBeforeBroadcast(t *testing.T) {
	item := pendingClaim(nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(99), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item.EventId, item.WalletAddress).Return(true, nil)

	certProc := new(mockCertClaimProcessor)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertNotCalled(t, "ProcessQueuedCertificateClaim")
	chainDg.AssertExpectations(t)
}

func TestProcessCertClaimQueue_ContinuesOnProcessError(t *testing.T) {
	item1 := pendingClaim(nil)
	item2 := pendingClaim(nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item1, item2}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item1.EventId, item1.WalletAddress).Return(true, nil)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item2.EventId, item2.WalletAddress).Return(true, nil)

	certProc := new(mockCertClaimProcessor)
	certProc.On("ProcessQueuedCertificateClaim", mock.Anything, item1).Return(errors.New("mint failed"))
	certProc.On("ProcessQueuedCertificateClaim", mock.Anything, item2).Return(nil)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertExpectations(t)
}

func TestProcessCertClaimQueue_ProcessesAllItems(t *testing.T) {
	item1 := pendingClaim(nil)
	item2 := pendingClaim(nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item1, item2}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item1.EventId, item1.WalletAddress).Return(true, nil)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item2.EventId, item2.WalletAddress).Return(true, nil)

	certProc := new(mockCertClaimProcessor)
	certProc.On("ProcessQueuedCertificateClaim", mock.Anything, item1).Return(nil)
	certProc.On("ProcessQueuedCertificateClaim", mock.Anything, item2).Return(nil)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// processQueues
// ---------------------------------------------------------------------------

func TestProcessQueues_SkipsCertQueue_WhenJoinQueueAborts(t *testing.T) {
	// Join queue aborts (gas too high mid-join) → cert queue must not run at all.
	item := pendingJoin(nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(99), nil)

	certProc := new(mockCertClaimProcessor)

	w := newWorker(nil, certProc, sigDg, chainDg, 30)
	w.processQueues(context.Background())

	// GetPendingCertificateClaimSignatures must never be called.
	sigDg.AssertNotCalled(t, "GetPendingCertificateClaimSignatures")
	certProc.AssertNotCalled(t, "ProcessQueuedCertificateClaim")
}

func TestProcessQueues_RunsBothQueues_OnSuccess(t *testing.T) {
	joinItem := pendingJoin(nil)
	certItem := pendingClaim(nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{joinItem}, nil)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{certItem}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("ProcessQueuedEventJoin", mock.Anything, joinItem).Return(nil)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, certItem.EventId, certItem.WalletAddress).Return(true, nil)

	certProc := new(mockCertClaimProcessor)
	certProc.On("ProcessQueuedCertificateClaim", mock.Anything, certItem).Return(nil)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processQueues(context.Background())

	joinProc.AssertExpectations(t)
	certProc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Start
// ---------------------------------------------------------------------------

func TestStart_StopsCleanly_OnContextCancel(t *testing.T) {
	chainDg := new(mockBlockchainClientDg)
	// Gas returns too high so processQueues is never entered — keeps test simple.
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(99), nil).Maybe()

	sigDg := new(mockUserSigDg)

	w := newWorker(nil, nil, sigDg, chainDg, 30)
	w.pollInterval = 5 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		w.Start(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return after context was cancelled")
	}
}

func TestStart_SkipsTick_WhenGasTooHigh(t *testing.T) {
	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(99), nil)

	sigDg := new(mockUserSigDg)
	joinProc := new(mockEventJoinProcessor)
	certProc := new(mockCertClaimProcessor)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.pollInterval = 5 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	w.Start(ctx)

	// No queue fetch or broadcast should have happened.
	sigDg.AssertNotCalled(t, "GetPendingEventJoinSignatures")
	sigDg.AssertNotCalled(t, "GetPendingCertificateClaimSignatures")
}

func TestStart_ProcessesQueues_WhenGasOk(t *testing.T) {
	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{}, nil)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{}, nil)

	w := newWorker(nil, nil, sigDg, chainDg, 30)
	w.pollInterval = 5 * time.Millisecond

	// Let the worker fire at least one tick then stop.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	w.Start(ctx)

	sigDg.AssertCalled(t, "GetPendingEventJoinSignatures", mock.Anything)
	sigDg.AssertCalled(t, "GetPendingCertificateClaimSignatures", mock.Anything)
}

// ---------------------------------------------------------------------------
// Block-based deadline — processJoinQueue
// ---------------------------------------------------------------------------

func TestProcessJoinQueue_MarksExpired_WhenBlockDeadlineExceeded(t *testing.T) {
	// currentBlock(1000) >= deadlineBlock(500) → must mark as expired, not broadcast.
	item := pendingJoinWithBlock(500)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item}, nil)
	sigDg.On("UpdateUserSignatureMarkAsExpiredAt", mock.Anything, item.SignatureId, mock.AnythingOfType("*time.Time")).
		Return(&entity.UserSignature{}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)

	joinProc := new(mockEventJoinProcessor)

	w := newWorker(joinProc, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.True(t, ok)
	joinProc.AssertNotCalled(t, "ProcessQueuedEventJoin")
	sigDg.AssertExpectations(t)
	chainDg.AssertExpectations(t)
}

func TestProcessJoinQueue_DoesNotMarkExpired_WhenBlockDeadlineNotYetReached(t *testing.T) {
	// currentBlock(500) < deadlineBlock(1000) → not expired by block, should broadcast.
	item := pendingJoinWithBlock(1000)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(500), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("ProcessQueuedEventJoin", mock.Anything, item).Return(nil)

	w := newWorker(joinProc, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.True(t, ok)
	joinProc.AssertExpectations(t)
	sigDg.AssertNotCalled(t, "UpdateUserSignatureMarkAsExpiredAt")
}

func TestProcessJoinQueue_SkipsBlockCheck_WhenGetCurrentBlockFails(t *testing.T) {
	// RPC failure for block number → block deadline check is skipped entirely,
	// item should still be broadcast (time deadline is not set either).
	item := pendingJoinWithBlock(500) // would be expired by block, but check is skipped

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(0), errors.New("rpc error"))
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("ProcessQueuedEventJoin", mock.Anything, item).Return(nil)

	w := newWorker(joinProc, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.True(t, ok)
	joinProc.AssertExpectations(t)
	sigDg.AssertNotCalled(t, "UpdateUserSignatureMarkAsExpiredAt")
}

// ---------------------------------------------------------------------------
// Block-based deadline — processCertClaimQueue
// ---------------------------------------------------------------------------

func TestProcessCertClaimQueue_MarksExpired_WhenBlockDeadlineExceeded(t *testing.T) {
	item := pendingClaimWithBlock(500)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)
	sigDg.On("UpdateUserSignatureMarkAsExpiredAt", mock.Anything, item.SignatureId, mock.AnythingOfType("*time.Time")).
		Return(&entity.UserSignature{}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)

	certProc := new(mockCertClaimProcessor)

	w := newWorker(nil, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertNotCalled(t, "ProcessQueuedCertificateClaim")
	sigDg.AssertExpectations(t)
	chainDg.AssertExpectations(t)
}

func TestProcessCertClaimQueue_DoesNotMarkExpired_WhenBlockDeadlineNotYetReached(t *testing.T) {
	item := pendingClaimWithBlock(1000)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(500), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item.EventId, item.WalletAddress).Return(true, nil)

	certProc := new(mockCertClaimProcessor)
	certProc.On("ProcessQueuedCertificateClaim", mock.Anything, item).Return(nil)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertExpectations(t)
	sigDg.AssertNotCalled(t, "UpdateUserSignatureMarkAsExpiredAt")
}

func TestProcessCertClaimQueue_SkipsBlockCheck_WhenGetCurrentBlockFails(t *testing.T) {
	item := pendingClaimWithBlock(500)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(0), errors.New("rpc error"))
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item.EventId, item.WalletAddress).Return(true, nil)

	certProc := new(mockCertClaimProcessor)
	certProc.On("ProcessQueuedCertificateClaim", mock.Anything, item).Return(nil)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertExpectations(t)
	sigDg.AssertNotCalled(t, "UpdateUserSignatureMarkAsExpiredAt")
}

// ---------------------------------------------------------------------------
// Block deadline expiry – processJoinQueue
// ---------------------------------------------------------------------------

func TestProcessJoinQueue_MarksExpiredByBlockDeadline(t *testing.T) {
	// DeadlineBlock is 500; current block is 1000 → expired.
	item := pendingJoinWithBlock(500)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item}, nil)
	sigDg.On("UpdateUserSignatureMarkAsExpiredAt", mock.Anything, item.SignatureId, mock.AnythingOfType("*time.Time")).
		Return(&entity.UserSignature{}, nil)

	joinProc := new(mockEventJoinProcessor)
	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)

	w := newWorker(joinProc, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.True(t, ok)
	joinProc.AssertNotCalled(t, "ProcessQueuedEventJoin")
	sigDg.AssertExpectations(t)
	chainDg.AssertExpectations(t)
}

func TestProcessJoinQueue_NotExpired_WhenBlockBelowDeadline(t *testing.T) {
	// DeadlineBlock is 2000; current block is 1000 → not expired.
	item := pendingJoinWithBlock(2000)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("ProcessQueuedEventJoin", mock.Anything, item).Return(nil)

	w := newWorker(joinProc, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.True(t, ok)
	joinProc.AssertExpectations(t)
	sigDg.AssertNotCalled(t, "UpdateUserSignatureMarkAsExpiredAt")
}

func TestProcessJoinQueue_SkipsBlockCheck_WhenGetCurrentBlockNumberFails(t *testing.T) {
	// Even when block fetch fails the item should still be processed normally.
	item := pendingJoinWithBlock(500) // would be expired if block check ran

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingEventJoinSignatures", mock.Anything).Return([]entity.PendingEventJoin{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(0), errors.New("rpc error"))
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("ProcessQueuedEventJoin", mock.Anything, item).Return(nil)

	w := newWorker(joinProc, nil, sigDg, chainDg, 30)
	ok := w.processJoinQueue(context.Background())

	assert.True(t, ok)
	joinProc.AssertExpectations(t)
	sigDg.AssertNotCalled(t, "UpdateUserSignatureMarkAsExpiredAt")
}

// ---------------------------------------------------------------------------
// Block deadline expiry – processCertClaimQueue
// ---------------------------------------------------------------------------

func TestProcessCertClaimQueue_MarksExpiredByBlockDeadline(t *testing.T) {
	// DeadlineBlock is 500; current block is 1000 → expired.
	item := pendingClaimWithBlock(500)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)
	sigDg.On("UpdateUserSignatureMarkAsExpiredAt", mock.Anything, item.SignatureId, mock.AnythingOfType("*time.Time")).
		Return(&entity.UserSignature{}, nil)

	certProc := new(mockCertClaimProcessor)
	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)

	w := newWorker(nil, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertNotCalled(t, "ProcessQueuedCertificateClaim")
	sigDg.AssertExpectations(t)
	chainDg.AssertExpectations(t)
}

func TestProcessCertClaimQueue_NotExpired_WhenBlockBelowDeadline(t *testing.T) {
	// DeadlineBlock is 2000; current block is 1000 → not expired.
	item := pendingClaimWithBlock(2000)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item.EventId, item.WalletAddress).Return(true, nil)

	certProc := new(mockCertClaimProcessor)
	certProc.On("ProcessQueuedCertificateClaim", mock.Anything, item).Return(nil)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertExpectations(t)
	sigDg.AssertNotCalled(t, "UpdateUserSignatureMarkAsExpiredAt")
}

func TestProcessCertClaimQueue_SkipsBlockCheck_WhenGetCurrentBlockNumberFails(t *testing.T) {
	// Even when block fetch fails the item should still be processed normally.
	item := pendingClaimWithBlock(500) // would be expired if block check ran

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(0), errors.New("rpc error"))
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item.EventId, item.WalletAddress).Return(true, nil)

	certProc := new(mockCertClaimProcessor)
	certProc.On("ProcessQueuedCertificateClaim", mock.Anything, item).Return(nil)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertExpectations(t)
	sigDg.AssertNotCalled(t, "UpdateUserSignatureMarkAsExpiredAt")
}

// ---------------------------------------------------------------------------
// Not-joined on-chain → pending-join queue check → abort logic
// ---------------------------------------------------------------------------

func TestProcessCertClaimQueue_Skips_WhenNotJoined_AndPendingJoinExists(t *testing.T) {
	// User is not joined on-chain yet, but there IS a pending join tx in the
	// queue → worker must defer (skip) and NOT abort.
	item := pendingClaim(nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item.EventId, item.WalletAddress).Return(false, nil)
	joinProc.On("HasPendingEventJoin", mock.Anything, item.EventId, item.AuthenticationCredentialId).Return(true, nil)

	certProc := new(mockCertClaimProcessor)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertNotCalled(t, "ProcessQueuedCertificateClaim")
	certProc.AssertNotCalled(t, "AbortCertificateClaim")
	joinProc.AssertExpectations(t)
}

func TestProcessCertClaimQueue_Aborts_WhenNotJoined_AndNoPendingJoin(t *testing.T) {
	// User is not joined on-chain AND there is no pending join in the queue →
	// the claim can never succeed, so the worker must abort it.
	item := pendingClaim(nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item.EventId, item.WalletAddress).Return(false, nil)
	joinProc.On("HasPendingEventJoin", mock.Anything, item.EventId, item.AuthenticationCredentialId).Return(false, nil)

	certProc := new(mockCertClaimProcessor)
	certProc.On("AbortCertificateClaim", mock.Anything, item.SignatureId, entity.AbortReasonParticipantNotJoined).Return(nil)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertNotCalled(t, "ProcessQueuedCertificateClaim")
	certProc.AssertExpectations(t)
}

func TestProcessCertClaimQueue_Skips_WhenNotJoined_AndPendingJoinCheckFails(t *testing.T) {
	// If the pending-join check itself errors, we must defer safely (skip) and
	// NOT abort — it is better to leave the claim alone than to abort it when
	// we are uncertain of the state.
	item := pendingClaim(nil)

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, item.EventId, item.WalletAddress).Return(false, nil)
	joinProc.On("HasPendingEventJoin", mock.Anything, item.EventId, item.AuthenticationCredentialId).Return(false, errors.New("db error"))

	certProc := new(mockCertClaimProcessor)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertNotCalled(t, "ProcessQueuedCertificateClaim")
	certProc.AssertNotCalled(t, "AbortCertificateClaim")
}

func TestProcessCertClaimQueue_ContinuesProcessingNext_AfterAbort(t *testing.T) {
	// After aborting one claim (not joined + no pending join), the worker must
	// continue processing the next item in the queue.
	//
	// item1 and item2 must have distinct EventId/WalletAddress so that the
	// IsParticipantJoinedOnChain mock can return different values for each.
	event1 := uuid.New()
	cred1 := uuid.New()
	item1 := entity.PendingCertificateClaim{
		SignatureId:                uuid.New(),
		CertificateId:              uuid.New(),
		EventId:                    event1,
		WalletAddress:              "0xAAA",
		AuthenticationCredentialId: cred1,
	}

	event2 := uuid.New()
	item2 := entity.PendingCertificateClaim{
		SignatureId:                uuid.New(),
		CertificateId:              uuid.New(),
		EventId:                    event2,
		WalletAddress:              "0xBBB",
		AuthenticationCredentialId: uuid.New(),
	}

	sigDg := new(mockUserSigDg)
	sigDg.On("GetPendingCertificateClaimSignatures", mock.Anything).Return([]entity.PendingCertificateClaim{item1, item2}, nil)

	chainDg := new(mockBlockchainClientDg)
	chainDg.On("GetCurrentBlockNumber", mock.Anything).Return(uint64(1000), nil)
	chainDg.On("GetGasPrice", mock.Anything).Return(gasOk(10), nil)

	joinProc := new(mockEventJoinProcessor)
	// item1: not joined, no pending join → abort
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, event1, "0xAAA").Return(false, nil)
	joinProc.On("HasPendingEventJoin", mock.Anything, event1, cred1).Return(false, nil)
	// item2: joined → process normally
	joinProc.On("IsParticipantJoinedOnChain", mock.Anything, event2, "0xBBB").Return(true, nil)

	certProc := new(mockCertClaimProcessor)
	certProc.On("AbortCertificateClaim", mock.Anything, item1.SignatureId, entity.AbortReasonParticipantNotJoined).Return(nil)
	certProc.On("ProcessQueuedCertificateClaim", mock.Anything, item2).Return(nil)

	w := newWorker(joinProc, certProc, sigDg, chainDg, 30)
	w.processCertClaimQueue(context.Background())

	certProc.AssertExpectations(t)
	joinProc.AssertExpectations(t)
}
