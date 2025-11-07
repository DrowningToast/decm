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

type InboxMessageDataGateway interface {
	CreateInboxMessage(ctx context.Context, params CreateInboxMessageParameters) (*entity.InboxMessage, error)
	GetInboxMessageByID(ctx context.Context, id uuid.UUID) (*entity.InboxMessage, error)
	GetInboxMessagesByReceiverEmail(ctx context.Context, receiverEmail string) ([]*entity.InboxMessage, error)
	UpdateInboxMessageReadStatus(ctx context.Context, id uuid.UUID, isRead int) (*entity.InboxMessage, error)
}
