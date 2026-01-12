package postgres

import (
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ offchain_datagateway.UserSignatureDataGateway = (*Repository)(nil)

func mapGeneratedUserSignatureToEntity(row *generated.UserSignature) *entity.UserSignature {
	return &entity.UserSignature{
		Id:                         row.ID,
		AuthenticationCredentialId: row.AuthenticationCredentialID,
		SignMessage:                row.SignMessage.String,
		Signature:                  row.Signature.String,
		DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
		EstimatedDeadline:          pgmapper.PgTimestampzToTimePtr(row.EstimatedDeadline),
		BroadcastedAt:              pgmapper.PgTimestampzToTimePtr(row.BroadcastedAt),
		MarkAsExpiredAt:            pgmapper.PgTimestampzToTimePtr(row.MarkAsExpiredAt),
		AbortedAt:                  pgmapper.PgTimestampzToTimePtr(row.AbortedAt),
		AbortedReason:              pgmapper.PgInt2ToInt32Ptr(row.AbortedReason),
		CreatedAt:                  row.CreatedAt,
		UpdatedAt:                  row.UpdatedAt,
	}
}

func (r *Repository) UpdateUserSignatureAbortedAt(ctx context.Context, id uuid.UUID, abortedAt time.Time, reason entity.UserSignatureAbortReason) (*entity.UserSignature, error) {
	row, err := r.queries.UpdateUserSignatureAbortedAt(ctx, generated.UpdateUserSignatureAbortedAtParams{
		ID:            id,
		AbortedAt:     pgmapper.TimePtrToPgTimestampz(&abortedAt),
		AbortedReason: pgmapper.Int32ToPgInt2(int32(reason)),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return mapGeneratedUserSignatureToEntity(&row), nil
}

func (r *Repository) CreateUserSignature(ctx context.Context, params offchain_datagateway.CreateUserSignatureParameters) (*entity.UserSignature, error) {
	row, err := r.queries.CreateUserSignature(ctx, generated.CreateUserSignatureParams{
		AuthenticationCredentialID: params.AuthenticationCredentialId,
		SignMessage: pgtype.Text{
			String: params.SignMessage,
			Valid:  true,
		},
		Signature: pgtype.Text{
			String: params.Signature,
			Valid:  true,
		},
		DeadlineBlock:     pgmapper.IntPtrToPgInt4(params.DeadlineBlock),
		EstimatedDeadline: pgmapper.TimePtrToPgTimestampz(params.EstimatedDeadline),
		BroadcastedAt:     pgmapper.TimePtrToPgTimestampz(params.BroadcastedAt),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return mapGeneratedUserSignatureToEntity(&row), nil
}

func (r *Repository) GetUserSignatureByID(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error) {
	row, err := r.queries.GetUserSignatureByID(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return mapGeneratedUserSignatureToEntity(&row), nil
}

func (r *Repository) GetUserSignaturesByCredentialID(ctx context.Context, authenticationCredentialId uuid.UUID) ([]entity.UserSignature, error) {
	rows, err := r.queries.GetUserSignaturesByCredentialID(ctx, authenticationCredentialId)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]entity.UserSignature, 0, len(rows))
	for _, row := range rows {
		signatures = append(signatures, *mapGeneratedUserSignatureToEntity(&row))
	}

	return signatures, nil
}

func (r *Repository) UpdateUserSignatureBroadcastedAt(ctx context.Context, id uuid.UUID, broadcastedAt *time.Time) (*entity.UserSignature, error) {
	row, err := r.queries.UpdateUserSignatureBroadcastedAt(ctx, generated.UpdateUserSignatureBroadcastedAtParams{
		ID:            id,
		BroadcastedAt: pgmapper.TimePtrToPgTimestampz(broadcastedAt),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return mapGeneratedUserSignatureToEntity(&row), nil
}

func (r *Repository) UpdateUserSignatureMarkAsExpiredAt(ctx context.Context, id uuid.UUID, markAsExpiredAt *time.Time) (*entity.UserSignature, error) {
	row, err := r.queries.UpdateUserSignatureMarkAsExpiredAt(ctx, generated.UpdateUserSignatureMarkAsExpiredAtParams{
		ID:              id,
		MarkAsExpiredAt: pgmapper.TimePtrToPgTimestampz(markAsExpiredAt),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return mapGeneratedUserSignatureToEntity(&row), nil
}

func (r *Repository) GetPendingUserSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	rows, err := r.queries.GetPendingUserSignatures(ctx)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]entity.UserSignature, 0, len(rows))
	for _, row := range rows {
		signatures = append(signatures, entity.UserSignature{
			Id:                         row.ID,
			AuthenticationCredentialId: row.AuthenticationCredentialID,
			SignMessage:                row.SignMessage.String,
			Signature:                  row.Signature.String,
			DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
			EstimatedDeadline:          pgmapper.PgTimestampzToTimePtr(row.EstimatedDeadline),
			BroadcastedAt:              pgmapper.PgTimestampzToTimePtr(row.BroadcastedAt),
			CreatedAt:                  row.CreatedAt,
			UpdatedAt:                  row.UpdatedAt,
		})
	}

	return signatures, nil
}

func (r *Repository) GetStaleUserSignatures(ctx context.Context, createdBefore time.Time) ([]entity.UserSignature, error) {
	rows, err := r.queries.GetStaleUserSignatures(ctx, createdBefore)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]entity.UserSignature, 0, len(rows))
	for _, row := range rows {
		signatures = append(signatures, entity.UserSignature{
			Id:                         row.ID,
			AuthenticationCredentialId: row.AuthenticationCredentialID,
			SignMessage:                row.SignMessage.String,
			Signature:                  row.Signature.String,
			DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
			EstimatedDeadline:          pgmapper.PgTimestampzToTimePtr(row.EstimatedDeadline),
			BroadcastedAt:              pgmapper.PgTimestampzToTimePtr(row.BroadcastedAt),
			CreatedAt:                  row.CreatedAt,
			UpdatedAt:                  row.UpdatedAt,
		})
	}

	return signatures, nil
}

func (r *Repository) GetBroadcastedUserSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	rows, err := r.queries.GetBroadcastedUserSignatures(ctx)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]entity.UserSignature, 0, len(rows))
	for _, row := range rows {
		signatures = append(signatures, entity.UserSignature{
			Id:                         row.ID,
			AuthenticationCredentialId: row.AuthenticationCredentialID,
			SignMessage:                row.SignMessage.String,
			Signature:                  row.Signature.String,
			DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
			EstimatedDeadline:          pgmapper.PgTimestampzToTimePtr(row.EstimatedDeadline),
			BroadcastedAt:              pgmapper.PgTimestampzToTimePtr(row.BroadcastedAt),
			CreatedAt:                  row.CreatedAt,
			UpdatedAt:                  row.UpdatedAt,
		})
	}

	return signatures, nil
}

func (r *Repository) GetUserSignaturesByDeadlineBlockRange(ctx context.Context, minBlock int32, maxBlock int32) ([]entity.UserSignature, error) {
	rows, err := r.queries.GetUserSignaturesByDeadlineBlockRange(ctx, generated.GetUserSignaturesByDeadlineBlockRangeParams{
		MinBlock: pgmapper.Int32ToPgInt4(minBlock),
		MaxBlock: pgmapper.Int32ToPgInt4(maxBlock),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]entity.UserSignature, 0, len(rows))
	for _, row := range rows {
		signatures = append(signatures, *mapGeneratedUserSignatureToEntity(&row))
	}

	return signatures, nil
}

func (r *Repository) GetUserSignaturesExpiringBefore(ctx context.Context, deadline time.Time) ([]entity.UserSignature, error) {
	rows, err := r.queries.GetUserSignaturesExpiringBefore(ctx, pgmapper.TimePtrToPgTimestampz(&deadline))
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]entity.UserSignature, 0, len(rows))
	for _, row := range rows {
		signatures = append(signatures, entity.UserSignature{
			Id:                         row.ID,
			AuthenticationCredentialId: row.AuthenticationCredentialID,
			SignMessage:                row.SignMessage.String,
			Signature:                  row.Signature.String,
			DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
			EstimatedDeadline:          pgmapper.PgTimestampzToTimePtr(row.EstimatedDeadline),
			BroadcastedAt:              pgmapper.PgTimestampzToTimePtr(row.BroadcastedAt),
			CreatedAt:                  row.CreatedAt,
			UpdatedAt:                  row.UpdatedAt,
		})
	}

	return signatures, nil
}

func (r *Repository) GetEventJoinSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	rows, err := r.queries.GetEventJoinSignatures(ctx)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]entity.UserSignature, 0, len(rows))
	for _, row := range rows {
		signatures = append(signatures, entity.UserSignature{
			Id:                         row.ID,
			AuthenticationCredentialId: row.AuthenticationCredentialID,
			SignMessage:                row.SignMessage.String,
			Signature:                  row.Signature.String,
			DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
			BroadcastedAt:              pgmapper.PgTimestampzToTimePtr(row.BroadcastedAt),
			CreatedAt:                  row.CreatedAt,
		})
	}

	return signatures, nil
}

func (r *Repository) GetPendingEventJoinSignatures(ctx context.Context) ([]entity.PendingEventJoin, error) {
	rows, err := r.queries.GetEventJoinSignatures(ctx)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	joins := make([]entity.PendingEventJoin, 0, len(rows))
	for _, row := range rows {
		// Skip rows that are already broadcasted, explicitly marked as expired, or aborted.
		if row.BroadcastedAt.Valid || row.MarkAsExpiredAt.Valid || row.AbortedAt.Valid {
			continue
		}
		joins = append(joins, entity.PendingEventJoin{
			SignatureId:                row.ID,
			AuthenticationCredentialId: row.AuthenticationCredentialID,
			Signature:                  row.Signature.String,
			SignMessage:                row.SignMessage.String,
			DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
			EstimatedDeadline:          pgmapper.PgTimestampzToTimePtr(row.EstimatedDeadline),
			CreatedAt:                  row.CreatedAt,
			EventId:                    row.EventID,
			WalletAddress:              row.WalletAddress,
		})
	}

	return joins, nil
}

func (r *Repository) GetCertificateClaimSignatures(ctx context.Context) ([]entity.UserSignature, error) {
	rows, err := r.queries.GetCertificateClaimSignatures(ctx)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]entity.UserSignature, 0, len(rows))
	for _, row := range rows {
		signatures = append(signatures, entity.UserSignature{
			Id:                         row.ID,
			AuthenticationCredentialId: row.AuthenticationCredentialID,
			SignMessage:                row.SignMessage.String,
			Signature:                  row.Signature.String,
			DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
			BroadcastedAt:              pgmapper.PgTimestampzToTimePtr(row.BroadcastedAt),
			CreatedAt:                  row.CreatedAt,
		})
	}

	return signatures, nil
}

func (r *Repository) GetPendingCertificateClaimSignatures(ctx context.Context) ([]entity.PendingCertificateClaim, error) {
	rows, err := r.queries.GetCertificateClaimSignatures(ctx)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	claims := make([]entity.PendingCertificateClaim, 0, len(rows))
	for _, row := range rows {
		// Skip rows that are already broadcasted, explicitly marked as expired, aborted, or already minted.
		if row.BroadcastedAt.Valid || row.MarkAsExpiredAt.Valid || row.AbortedAt.Valid || row.CertificateTokenID.Valid {
			continue
		}
		claims = append(claims, entity.PendingCertificateClaim{
			SignatureId:                row.ID,
			AuthenticationCredentialId: row.AuthenticationCredentialID,
			Signature:                  row.Signature.String,
			SignMessage:                row.SignMessage.String,
			DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
			EstimatedDeadline:          pgmapper.PgTimestampzToTimePtr(row.EstimatedDeadline),
			CreatedAt:                  row.CreatedAt,
			CertificateId:              row.CertificateID,
			EventId:                    row.EventID,
			WalletAddress:              row.WalletAddress,
		})
	}

	return claims, nil
}

func (r *Repository) GetOrphanedUserSignatures(ctx context.Context, createdBefore time.Time) ([]entity.UserSignature, error) {
	rows, err := r.queries.GetOrphanedUserSignatures(ctx, createdBefore)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]entity.UserSignature, 0, len(rows))
	for _, row := range rows {
		signatures = append(signatures, entity.UserSignature{
			Id:                         row.ID,
			AuthenticationCredentialId: row.AuthenticationCredentialID,
			CreatedAt:                  row.CreatedAt,
			BroadcastedAt:              pgmapper.PgTimestampzToTimePtr(row.BroadcastedAt),
		})
	}

	return signatures, nil
}

func (r *Repository) GetUserSignatureWithUsageDetails(ctx context.Context, id uuid.UUID) (*entity.UserSignature, error) {
	row, err := r.queries.GetUserSignatureWithUsageDetails(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.UserSignature{
		Id:                         row.ID,
		AuthenticationCredentialId: row.AuthenticationCredentialID,
		SignMessage:                row.SignMessage.String,
		Signature:                  row.Signature.String,
		DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
		EstimatedDeadline:          pgmapper.PgTimestampzToTimePtr(row.EstimatedDeadline),
		BroadcastedAt:              pgmapper.PgTimestampzToTimePtr(row.BroadcastedAt),
		CreatedAt:                  row.CreatedAt,
		UpdatedAt:                  row.UpdatedAt,
	}, nil
}

func (r *Repository) GetRecentUserSignatureActivity(ctx context.Context, since time.Time) ([]entity.UserSignature, error) {
	rows, err := r.queries.GetRecentUserSignatureActivity(ctx, since)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	signatures := make([]entity.UserSignature, 0, len(rows))
	for _, row := range rows {
		signatures = append(signatures, entity.UserSignature{
			Id:                         row.ID,
			AuthenticationCredentialId: row.AuthenticationCredentialID,
			CreatedAt:                  row.CreatedAt,
			BroadcastedAt:              pgmapper.PgTimestampzToTimePtr(row.BroadcastedAt),
			DeadlineBlock:              pgmapper.PgInt4ToInt32Ptr(row.DeadlineBlock),
			EstimatedDeadline:          pgmapper.PgTimestampzToTimePtr(row.EstimatedDeadline),
		})
	}

	return signatures, nil
}

func (r *Repository) DeleteUserSignature(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteUserSignature(ctx, id)
	if err != nil {
		return pgerrutils.ParsePgError(err)
	}
	return nil
}

func (r *Repository) CountUserSignaturesByCredentialID(ctx context.Context, authenticationCredentialId uuid.UUID) (int64, error) {
	count, err := r.queries.CountUserSignaturesByCredentialID(ctx, authenticationCredentialId)
	if err != nil {
		return 0, pgerrutils.ParsePgError(err)
	}
	return count, nil
}

func (r *Repository) CountPendingUserSignatures(ctx context.Context) (int64, error) {
	count, err := r.queries.CountPendingUserSignatures(ctx)
	if err != nil {
		return 0, pgerrutils.ParsePgError(err)
	}
	return count, nil
}

func (r *Repository) CountBroadcastedUserSignatures(ctx context.Context) (int64, error) {
	count, err := r.queries.CountBroadcastedUserSignatures(ctx)
	if err != nil {
		return 0, pgerrutils.ParsePgError(err)
	}
	return count, nil
}
