package inbox

import (
	"apps/backend/common/customerror"
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"errors"

	"github.com/google/uuid"
)

func (uc *InboxUsecase) MarkMessageAsRead(ctx context.Context, user auth.JwtClaims, messageID uuid.UUID) (*entity.InboxMessage, error) {
	message, err := uc.InboxMessageDg.GetInboxMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("message not found"))
	}

	// Check if user is authorized to mark this message as read
	if !uc.isAuthorizedToReadMessage(message, user) {
		return nil, customerror.Parse(&customerror.ErrForbidden, errors.New("you are not allowed to mark this message as read"))
	}

	message, err = uc.InboxMessageDg.UpdateInboxMessageReadStatus(ctx, messageID, 1)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (uc *InboxUsecase) MarkAllMessagesAsRead(ctx context.Context, user auth.JwtClaims) ([]*entity.InboxMessage, error) {
	walletAddress := user.WalletAddress
	messages, err := uc.InboxMessageDg.UpdateInboxMessageReadStatusAll(ctx, offchain_datagateway.GetInboxMessagesByCredentialIDParameters{
		CredentialID:          user.UserId,
		ReceiverEmail:         user.Email,
		ReceiverWalletAddress: &walletAddress,
	})
	if err != nil {
		return nil, err
	}
	return messages, nil
}
