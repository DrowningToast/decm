package datagateway

import (
	"context"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

type CreateInboxMessageParameters struct {
	SenderCredentialID     *uuid.UUID
	ReceiverCredentialID   *uuid.UUID
	ReceiverEmail          string
	MessageType            int
	MessageContent         string
	FallbackMessageContent *string
	IsRead                 int
}

type GetInboxMessagesByCredentialIDParameters struct {
	CredentialID          uuid.UUID
	ReceiverEmail         *string
	ReceiverWalletAddress *string
}

type InboxMessageDataGateway interface {
	CreateInboxMessage(ctx context.Context, params CreateInboxMessageParameters) (*entity.InboxMessage, error)
	GetInboxMessageByID(ctx context.Context, id uuid.UUID) (*entity.InboxMessage, error)
	// Search with receiver authentication credential ID
	GetInboxMessagesByCredentialID(ctx context.Context, params GetInboxMessagesByCredentialIDParameters) ([]*entity.InboxMessage, error)
	// Search with receiver email (doesn't search through receiver wallet address column, not authentication credential)
	GetInboxMessagesByReceiverEmail(ctx context.Context, receiverEmail string) ([]*entity.InboxMessage, error)
	// Search with receiver wallet address (doesn't search through receiver email column, not authentication credential)
	GetInboxMessagesByReceiverWalletAddress(ctx context.Context, walletAddress string) ([]*entity.InboxMessage, error)
	// Search with sender authentication credential ID
	GetInboxMessagesBySenderCredentialID(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error)
	UpdateInboxMessageReadStatus(ctx context.Context, id uuid.UUID, isRead int) (*entity.InboxMessage, error)
	UpdateInboxMessageReadStatusAll(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error)
}
