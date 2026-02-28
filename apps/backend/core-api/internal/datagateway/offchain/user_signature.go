package offchain_datagateway

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateUserSignatureParameters struct {
	AuthenticationCredentialId uuid.UUID
	SignMessage                string
	Signature                  string
	DeadlineBlock              *int32
	EstimatedDeadline          *time.Time
	BroadcastedAt              *time.Time
}

type UserSignatureDataGateway interface {
	CreateUserSignature(ctx context.Context, params CreateUserSignatureParameters) (*entity.UserSignature, error)
	GetUserSignatureByID(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error)
	GetUserSignaturesByCredentialID(ctx context.Context, authenticationCredentialId uuid.UUID) ([]entity.UserSignature, error)
	UpdateUserSignatureBroadcastedAt(ctx context.Context, id uuid.UUID, broadcastedAt *time.Time) (*entity.UserSignature, error)
	UpdateUserSignatureMarkAsExpiredAt(ctx context.Context, id uuid.UUID, markAsExpiredAt *time.Time) (*entity.UserSignature, error)
	UpdateUserSignatureAbortedAt(ctx context.Context, id uuid.UUID, abortedAt time.Time, reason entity.UserSignatureAbortReason) (*entity.UserSignature, error)
	GetPendingUserSignatures(ctx context.Context) ([]entity.UserSignature, error)
	GetStaleUserSignatures(ctx context.Context, createdBefore time.Time) ([]entity.UserSignature, error)
	GetBroadcastedUserSignatures(ctx context.Context) ([]entity.UserSignature, error)
	GetUserSignaturesByDeadlineBlockRange(ctx context.Context, minBlock int32, maxBlock int32) ([]entity.UserSignature, error)
	GetUserSignaturesExpiringBefore(ctx context.Context, deadline time.Time) ([]entity.UserSignature, error)
	GetEventJoinSignatures(ctx context.Context) ([]entity.UserSignature, error)
	GetPendingEventJoinSignatures(ctx context.Context) ([]entity.PendingEventJoin, error)
	GetCertificateClaimSignatures(ctx context.Context) ([]entity.UserSignature, error)
	GetPendingCertificateClaimSignatures(ctx context.Context) ([]entity.PendingCertificateClaim, error)
	GetOrphanedUserSignatures(ctx context.Context, createdBefore time.Time) ([]entity.UserSignature, error)
	GetUserSignatureWithUsageDetails(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error)
	GetRecentUserSignatureActivity(ctx context.Context, since time.Time) ([]entity.UserSignature, error)
	DeleteUserSignature(ctx context.Context, id uuid.UUID) error
	CountUserSignaturesByCredentialID(ctx context.Context, authenticationCredentialId uuid.UUID) (int64, error)
	CountPendingUserSignatures(ctx context.Context) (int64, error)
	CountBroadcastedUserSignatures(ctx context.Context) (int64, error)
}
