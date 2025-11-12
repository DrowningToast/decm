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
		Id:                     result.ID,
		SenderCredentialId:     &result.SenderCredentialID,
		ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:          decryptedEmail,
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
		Id:                     result.ID,
		SenderCredentialId:     &result.SenderCredentialID,
		ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
		ReceiverEmail:          decryptedEmail,
		ReceiverWalletAddress:  pgmapper.PgTextToStringPtr(result.ReceiverWalletAddress),
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

func (r *Repository) GetInboxMessagesByCredentialID(ctx context.Context, params datagateway.GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error) {
	var err error
	var encryptedReceiverEmail pgtype.Text
	receiverEmail := params.ReceiverEmail
	if receiverEmail != nil {
		encryptedReceiverEmail, err = pgmapper.EncryptStringPtrToPgText(receiverEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
	}

	var encryptedReceiverWalletAddress pgtype.Text
	receiverWalletAddress := params.ReceiverWalletAddress
	if receiverWalletAddress != nil {
		encryptedReceiverWalletAddress, err = pgmapper.EncryptStringPtrToPgText(receiverWalletAddress, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
	}

	results, err := r.queries.GetInboxMessagesByCredentialID(ctx, generated.GetInboxMessagesByCredentialIDParams{
		ReceiverCredentialID:  pgmapper.UUIDToPgUUID(&params.CredentialID),
		ReceiverEmail:         encryptedReceiverEmail,
		ReceiverWalletAddress: encryptedReceiverWalletAddress,
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
			MessageType:            int(result.MessageType),
			MessageContent:         string(result.MessageContent),
			FallbackMessageContent: pgmapper.PgTextToStringPtr(result.FallbackMessageContent),
		}
	}

	return messages, nil
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
			Id:                     result.ID,
			SenderCredentialId:     &result.SenderCredentialID,
			ReceiverCredentialId:   pgmapper.PgUUIDToUUIDPtr(result.ReceiverCredentialID),
			ReceiverEmail:          decryptedEmail,
			ReceiverWalletAddress:  pgmapper.PgTextToStringPtr(result.ReceiverWalletAddress),
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

func (r *Repository) GetInboxMessagesByReceiverWalletAddress(ctx context.Context, receiverWalletAddress string) ([]*entity.InboxMessage, error) {
	// Encrypt search term
	encryptedWalletAddress, err := pgmapper.EncryptPII(receiverWalletAddress, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.GetInboxMessagesByReceiverWalletAddress(ctx, pgtype.Text{
		String: encryptedWalletAddress,
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
			MessageType:            int(result.MessageType),
			MessageContent:         string(result.MessageContent),
			FallbackMessageContent: pgmapper.PgTextToStringPtr(result.FallbackMessageContent),
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
			MessageType:            int(result.MessageType),
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

func (r *Repository) UpdateInboxMessageReadStatusAll(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	results, err := r.queries.UpdateInboxMessageReadStatusAll(ctx, pgmapper.UUIDToPgUUID(&credentialID))
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
			MessageType:            int(result.MessageType),
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
