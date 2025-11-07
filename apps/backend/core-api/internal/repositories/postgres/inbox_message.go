package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	datagateway "apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ datagateway.InboxMessageDataGateway = (*Repository)(nil)

func (r *Repository) CreateInboxMessage(ctx context.Context, params datagateway.CreateInboxMessageParameters) (*entity.InboxMessage, error) {
	// Encrypt PII field (receiver_email)
	encryptedEmail, err := pgmapper.EncryptStringPtrToPgText(&params.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	// Convert message_content to pgtype.Text for JSONB field
	messageContent := []byte(params.MessageContent)
	fallbackMessageContent := pgmapper.StringPtrToPgText(params.FallbackMessageContent)

	// Create inbox message
	result, err := r.queries.CreateInboxMessage(ctx, generated.CreateInboxMessageParams{
		SenderCredentialID:     pgmapper.UUIDToPgUUID(params.SenderCredentialID),
		ReceiverCredentialID:   pgmapper.UUIDToPgUUID(params.ReceiverCredentialID),
		ReceiverEmail:          encryptedEmail,
		MessageType:            int32(params.MessageType),
		MessageContent:         messageContent,
		FallbackMessageContent: fallbackMessageContent,
		IsRead:                 pgmapper.BoolToPgInt4(params.IsRead == 1),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII field for return
	decryptedEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	// Convert message_content from pgtype.Text to string
	messageContentStr := string(result.MessageContent)
	fallbackMessageContentStr := pgmapper.PgTextToStringPtr(result.FallbackMessageContent)

	return &entity.InboxMessage{
		ID:                     result.ID,
		SenderCredentialID:     &result.SenderCredentialID,
		ReceiverCredentialID:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:          *decryptedEmail,
		MessageType:            int(result.MessageType),
		MessageContent:         messageContentStr,
		FallbackMessageContent: fallbackMessageContentStr,
		IsRead:                 int(result.IsRead.Int32),
		CreatedAt:              *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:              *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		HiddenAt:               pgmapper.PgTimestampzToTimePtr(result.HiddenAt),
		DeletedAt:              pgmapper.PgTimestampzToTimePtr(result.DeletedAt),
	}, nil
}

func (r *Repository) GetInboxMessageByID(ctx context.Context, id uuid.UUID) (*entity.InboxMessage, error) {
	result, err := r.queries.GetInboxMessageByID(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII field
	decryptedEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	// Convert message_content from []byte to string
	messageContentStr := string(result.MessageContent)
	fallbackMessageContentStr := pgmapper.PgTextToStringPtr(result.FallbackMessageContent)

	return &entity.InboxMessage{
		ID:                     result.ID,
		SenderCredentialID:     &result.SenderCredentialID,
		ReceiverCredentialID:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:          *decryptedEmail,
		MessageType:            int(result.MessageType),
		MessageContent:         messageContentStr,
		FallbackMessageContent: fallbackMessageContentStr,
		IsRead:                 int(result.IsRead.Int32),
		CreatedAt:              *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:              *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		HiddenAt:               pgmapper.PgTimestampzToTimePtr(result.HiddenAt),
		DeletedAt:              pgmapper.PgTimestampzToTimePtr(result.DeletedAt),
	}, nil
}

func (r *Repository) GetInboxMessagesByReceiverEmail(ctx context.Context, receiverEmail string) ([]*entity.InboxMessage, error) {
	// Encrypt search term
	encryptedEmail, err := pgmapper.EncryptPII(receiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.GetInboxMessagesByReceiverEmail(ctx, pgtype.Text{
		String: encryptedEmail,
		Valid:  true,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	messages := make([]*entity.InboxMessage, len(results))
	for i, result := range results {
		// Decrypt PII field
		decryptedEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		// Convert message_content from pgtype.Text to string
		messageContentStr := string(result.MessageContent)
		fallbackMessageContentStr := pgmapper.PgTextToStringPtr(result.FallbackMessageContent)

		messages[i] = &entity.InboxMessage{
			ID:                     result.ID,
			SenderCredentialID:     &result.SenderCredentialID,
			ReceiverCredentialID:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:          *decryptedEmail,
			MessageType:            int(result.MessageType),
			MessageContent:         messageContentStr,
			FallbackMessageContent: fallbackMessageContentStr,
			IsRead:                 int(result.IsRead.Int32),
			CreatedAt:              *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
			UpdatedAt:              *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
			HiddenAt:               pgmapper.PgTimestampzToTimePtr(result.HiddenAt),
			DeletedAt:              pgmapper.PgTimestampzToTimePtr(result.DeletedAt),
		}
	}

	return messages, nil
}

func (r *Repository) UpdateInboxMessageReadStatus(ctx context.Context, id uuid.UUID, isRead int) (*entity.InboxMessage, error) {
	result, err := r.queries.UpdateInboxMessageReadStatus(ctx, generated.UpdateInboxMessageReadStatusParams{
		ID:     id,
		IsRead: pgmapper.Int32ToPgInt4(int32(isRead)),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII field
	decryptedEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	// Convert message_content from []byte to string
	messageContentStr := string(result.MessageContent)
	fallbackMessageContentStr := pgmapper.PgTextToStringPtr(result.FallbackMessageContent)

	return &entity.InboxMessage{
		ID:                     result.ID,
		SenderCredentialID:     &result.SenderCredentialID,
		ReceiverCredentialID:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:          *decryptedEmail,
		MessageType:            int(result.MessageType),
		MessageContent:         messageContentStr,
		FallbackMessageContent: fallbackMessageContentStr,
		IsRead:                 int(result.IsRead.Int32),
		CreatedAt:              *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:              *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		HiddenAt:               pgmapper.PgTimestampzToTimePtr(result.HiddenAt),
		DeletedAt:              pgmapper.PgTimestampzToTimePtr(result.DeletedAt),
	}, nil
}
