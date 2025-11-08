package inbox

import (
	"context"

	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
)

func (uc *InboxUsecase) GetMyInboxMessages(ctx context.Context, user auth.JwtClaims) ([]*entity.InboxMessage, error) {
	return uc.GetInboxMessagesByCredentailID(ctx, user.UserId)
}

func (uc *InboxUsecase) GetInboxMessagesByCredentailID(ctx context.Context, credentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	authenticationCredential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, credentialID)
	if err != nil {
		return nil, err
	}

	return uc.InboxMessageDg.GetInboxMessagesByCredentialID(ctx, datagateway.GetInboxMessagesByCredentialIDParameters{
		CredentialID:          credentialID,
		ReceiverEmail:         authenticationCredential.GoogleConnectorRef,
		ReceiverWalletAddress: &authenticationCredential.WalletAddress,
	})
}

func (uc *InboxUsecase) GetInboxMessagesByReceiverEmail(ctx context.Context, receiverEmail string) ([]*entity.InboxMessage, error) {
	return uc.InboxMessageDg.GetInboxMessagesByReceiverEmail(ctx, receiverEmail)
}

func (uc *InboxUsecase) GetInboxMessagesByReceiverWalletAddress(ctx context.Context, receiverWalletAddress string) ([]*entity.InboxMessage, error) {
	return uc.InboxMessageDg.GetInboxMessagesByReceiverWalletAddress(ctx, receiverWalletAddress)
}

func (uc *InboxUsecase) GetInboxMessagesBySenderCredentialID(ctx context.Context, senderCredentialID uuid.UUID) ([]*entity.InboxMessage, error) {
	return uc.InboxMessageDg.GetInboxMessagesBySenderCredentialID(ctx, senderCredentialID)
}
