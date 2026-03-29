package postgres

import (
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"
	"strings"

	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ offchain_datagateway.InboxMessageDataGateway = (*Repository)(nil)

func (r *Repository) CreateInboxMessage(ctx context.Context, params offchain_datagateway.CreateInboxMessageParameters) (*entity.InboxMessage, error) {
	// Encrypt PII field (receiver_email)
	encryptedEmail, err := pgmapper.EncryptStringPtrToPgText(&params.ReceiverEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	// Normalize wallet address to lowercase plain text (wallet addresses are public, not PII-encrypted)
	var normalizedWalletAddress *string
	if params.ReceiverWalletAddress != nil && *params.ReceiverWalletAddress != "" {
		lower := strings.ToLower(*params.ReceiverWalletAddress)
		normalizedWalletAddress = &lower
	}

	// Convert message_content to pgtype.Text for JSONB field
	messageContent := []byte(params.MessageContent)
	fallbackMessageContent := pgmapper.StringPtrToPgText(params.FallbackMessageContent)

	// Create inbox message
	result, err := r.queries.CreateInboxMessage(ctx, generated.CreateInboxMessageParams{
		SenderCredentialID:     pgmapper.UUIDPtrToPgUUID(params.SenderCredentialID),
		ReceiverCredentialID:   pgmapper.UUIDPtrToPgUUID(params.ReceiverCredentialID),
		ReceiverEmail:          encryptedEmail,
		ReceiverWalletAddress:  pgmapper.StringPtrToPgText(normalizedWalletAddress),
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
		Id:                     result.ID,
		SenderCredentialId:     &result.SenderCredentialID,
		ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:          decryptedEmail,
		MessageType:            entity.InboxMessageType(result.MessageType),
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
		Id:                     result.ID,
		SenderCredentialId:     &result.SenderCredentialID,
		ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:          decryptedEmail,
		ReceiverWalletAddress:  pgmapper.PgTextToStringPtr(result.ReceiverWalletAddress),
		MessageType:            entity.InboxMessageType(result.MessageType),
		MessageContent:         messageContentStr,
		FallbackMessageContent: fallbackMessageContentStr,
		IsRead:                 int(result.IsRead.Int32),
		CreatedAt:              *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:              *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		HiddenAt:               pgmapper.PgTimestampzToTimePtr(result.HiddenAt),
		DeletedAt:              pgmapper.PgTimestampzToTimePtr(result.DeletedAt),
	}, nil
}

func (r *Repository) GetInboxMessagesByCredentialID(ctx context.Context, params offchain_datagateway.GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error) {
	var err error
	var encryptedReceiverEmail pgtype.Text
	receiverEmail := params.ReceiverEmail
	if receiverEmail != nil {
		// Normalize email to lowercase for case-insensitive comparison
		normalizedEmail := strings.ToLower(*receiverEmail)
		encryptedReceiverEmail, err = pgmapper.EncryptStringPtrToPgText(&normalizedEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
	}

	// Wallet address is stored as plain text (not PII-encrypted); normalize to lowercase for matching
	var plainWalletAddress pgtype.Text
	if params.ReceiverWalletAddress != nil && *params.ReceiverWalletAddress != "" {
		normalized := strings.ToLower(*params.ReceiverWalletAddress)
		plainWalletAddress = pgtype.Text{String: normalized, Valid: true}
	}

	results, err := r.queries.GetInboxMessagesByCredentialID(ctx, generated.GetInboxMessagesByCredentialIDParams{
		ReceiverCredentialID:  pgmapper.UUIDToPgUUID(params.CredentialID),
		ReceiverEmail:         encryptedReceiverEmail,
		ReceiverWalletAddress: plainWalletAddress,
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

		messages[i] = &entity.InboxMessage{
			Id:                     result.ID,
			SenderCredentialId:     &result.SenderCredentialID,
			ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:          decryptedEmail,
			ReceiverWalletAddress:  pgmapper.PgTextToStringPtr(result.ReceiverWalletAddress),
			MessageType:            entity.InboxMessageType(result.MessageType),
			MessageContent:         string(result.MessageContent),
			FallbackMessageContent: pgmapper.PgTextToStringPtr(result.FallbackMessageContent),
			IsRead:                 int(result.IsRead.Int32),
			CreatedAt:              *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
			UpdatedAt:              *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
			HiddenAt:               pgmapper.PgTimestampzToTimePtr(result.HiddenAt),
			DeletedAt:              pgmapper.PgTimestampzToTimePtr(result.DeletedAt),
		}
	}

	return messages, nil
}

func (r *Repository) GetUnreadInboxMessageCountByCredentialID(ctx context.Context, params offchain_datagateway.GetInboxMessagesByCredentialIDParameters) (int, error) {
	var err error
	var encryptedReceiverEmail pgtype.Text
	receiverEmail := params.ReceiverEmail
	if receiverEmail != nil {
		// Normalize email to lowercase for case-insensitive comparison
		normalizedEmail := strings.ToLower(*receiverEmail)
		encryptedReceiverEmail, err = pgmapper.EncryptStringPtrToPgText(&normalizedEmail, r.piiEncryptionKey)
		if err != nil {
			return 0, err
		}
	}

	// Wallet address is stored as plain text (not PII-encrypted); normalize to lowercase for matching
	var plainWalletAddress pgtype.Text
	if params.ReceiverWalletAddress != nil && *params.ReceiverWalletAddress != "" {
		normalized := strings.ToLower(*params.ReceiverWalletAddress)
		plainWalletAddress = pgtype.Text{String: normalized, Valid: true}
	}

	count, err := r.queries.GetUnreadInboxMessageCountByCredentialID(ctx, generated.GetUnreadInboxMessageCountByCredentialIDParams{
		ReceiverCredentialID:  pgmapper.UUIDToPgUUID(params.CredentialID),
		ReceiverEmail:         encryptedReceiverEmail,
		ReceiverWalletAddress: plainWalletAddress,
	})
	if err != nil {
		return 0, pgerrutils.ParsePgError(err)
	}

	return int(count), nil
}

func (r *Repository) GetInboxMessagesByReceiverEmail(ctx context.Context, receiverEmail string) ([]*entity.InboxMessage, error) {
	// Normalize email to lowercase for case-insensitive comparison
	normalizedEmail := strings.ToLower(receiverEmail)
	// Encrypt search term
	encryptedEmail, err := pgmapper.EncryptPII(normalizedEmail, r.piiEncryptionKey)
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
			Id:                     result.ID,
			SenderCredentialId:     &result.SenderCredentialID,
			ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:          decryptedEmail,
			ReceiverWalletAddress:  pgmapper.PgTextToStringPtr(result.ReceiverWalletAddress),
			MessageType:            entity.InboxMessageType(result.MessageType),
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

func (r *Repository) GetInboxMessagesByReceiverWalletAddress(ctx context.Context, receiverWalletAddress string) ([]*entity.InboxMessage, error) {
	// Wallet address is plain text; normalize to lowercase for case-insensitive matching
	results, err := r.queries.GetInboxMessagesByReceiverWalletAddress(ctx, pgtype.Text{
		String: strings.ToLower(receiverWalletAddress),
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

		messages[i] = &entity.InboxMessage{
			Id:                     result.ID,
			SenderCredentialId:     &result.SenderCredentialID,
			ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:          decryptedEmail,
			ReceiverWalletAddress:  pgmapper.PgTextToStringPtr(result.ReceiverWalletAddress),
			MessageType:            entity.InboxMessageType(result.MessageType),
			MessageContent:         string(result.MessageContent),
			FallbackMessageContent: pgmapper.PgTextToStringPtr(result.FallbackMessageContent),
			IsRead:                 int(result.IsRead.Int32),
			CreatedAt:              *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
			UpdatedAt:              *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
			HiddenAt:               pgmapper.PgTimestampzToTimePtr(result.HiddenAt),
			DeletedAt:              pgmapper.PgTimestampzToTimePtr(result.DeletedAt),
		}
	}

	return messages, nil
}

func (r *Repository) GetInboxMessagesBySenderCredentialID(ctx context.Context, senderCredentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	results, err := r.queries.GetInboxMessagesBySenderCredentialID(ctx, senderCredentialID)
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

		messages[i] = &entity.InboxMessage{
			Id:                     result.ID,
			SenderCredentialId:     &result.SenderCredentialID,
			ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:          decryptedEmail,
			ReceiverWalletAddress:  pgmapper.PgTextToStringPtr(result.ReceiverWalletAddress),
			MessageType:            entity.InboxMessageType(result.MessageType),
			MessageContent:         string(result.MessageContent),
			FallbackMessageContent: pgmapper.PgTextToStringPtr(result.FallbackMessageContent),
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
		Id:                     result.ID,
		SenderCredentialId:     &result.SenderCredentialID,
		ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:          decryptedEmail,
		ReceiverWalletAddress:  pgmapper.PgTextToStringPtr(result.ReceiverWalletAddress),
		MessageType:            entity.InboxMessageType(result.MessageType),
		MessageContent:         messageContentStr,
		FallbackMessageContent: fallbackMessageContentStr,
		IsRead:                 int(result.IsRead.Int32),
		CreatedAt:              *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:              *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		HiddenAt:               pgmapper.PgTimestampzToTimePtr(result.HiddenAt),
		DeletedAt:              pgmapper.PgTimestampzToTimePtr(result.DeletedAt),
	}, nil
}

func (r *Repository) UpdateInboxMessageReadStatusAll(ctx context.Context, params offchain_datagateway.GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error) {
	var err error

	var encryptedReceiverEmail pgtype.Text
	if params.ReceiverEmail != nil {
		normalizedEmail := strings.ToLower(*params.ReceiverEmail)
		encryptedReceiverEmail, err = pgmapper.EncryptStringPtrToPgText(&normalizedEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
	}

	// Wallet address is plain text; normalize to lowercase for case-insensitive matching
	var plainWalletAddress pgtype.Text
	if params.ReceiverWalletAddress != nil && *params.ReceiverWalletAddress != "" {
		plainWalletAddress = pgtype.Text{String: strings.ToLower(*params.ReceiverWalletAddress), Valid: true}
	}

	results, err := r.queries.UpdateInboxMessageReadStatusAll(ctx, generated.UpdateInboxMessageReadStatusAllParams{
		ReceiverCredentialID:  pgmapper.UUIDToPgUUID(params.CredentialID),
		ReceiverEmail:         encryptedReceiverEmail,
		ReceiverWalletAddress: plainWalletAddress,
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

		messages[i] = &entity.InboxMessage{
			Id:                     result.ID,
			SenderCredentialId:     &result.SenderCredentialID,
			ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:          decryptedEmail,
			ReceiverWalletAddress:  pgmapper.PgTextToStringPtr(result.ReceiverWalletAddress),
			MessageType:            entity.InboxMessageType(result.MessageType),
			MessageContent:         string(result.MessageContent),
			FallbackMessageContent: pgmapper.PgTextToStringPtr(result.FallbackMessageContent),
			IsRead:                 int(result.IsRead.Int32),
			CreatedAt:              *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
			UpdatedAt:              *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
			HiddenAt:               pgmapper.PgTimestampzToTimePtr(result.HiddenAt),
			DeletedAt:              pgmapper.PgTimestampzToTimePtr(result.DeletedAt),
		}
	}

	return messages, nil
}
