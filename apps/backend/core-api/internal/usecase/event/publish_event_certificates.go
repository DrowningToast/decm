package event

// DEPRECATED: This functionality has been merged into eventconfig.ToggleCertificatePublished
// Use the ToggleCertificatePublished endpoint instead: PATCH /api/v1/events/{event_id}/config/certificate/published
//
// This file is kept for backward compatibility but should not be used in new code.
// The merged implementation includes:
// - Signature validation (all issuers must sign)
// - Inbox message creation
// - Published status toggle
//
// TODO: Remove this file and the associated handler after migration is complete

import (
	"apps/backend/common/customerror"
	"apps/backend/common/utils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"decm-database/go/generated"
	"encoding/json"
	"fmt"

	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"

	"github.com/google/uuid"
)

type PublishEventCertificatesResponse struct {
	EventID              uuid.UUID `json:"event_id"`
	PublishedCount       int       `json:"published_count"`
	IsPublished          bool      `json:"is_published"`
	InboxMessagesCreated int       `json:"inbox_messages_created"`
}

func (uc *EventUsecase) PublishEventCertificates(ctx context.Context, eventID uuid.UUID, currentUser *auth.JwtClaims) (*PublishEventCertificatesResponse, error) {
	// 1. Check if current user is authorized (must be verified organizer)
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	if !credential.IsVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user is not a verified organizer"))
	}

	// 2. Check if event exists and user is the owner
	event, err := uc.EventDataGateway.GetEventById(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event not found"))
	}

	if event.OwnerCredentialId != currentUser.UserId {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user is not the event owner"))
	}

	// 3. Get all certificates for the event
	certificates, err := uc.EventCertificateDataGateway.GetEventCertificatesByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if len(certificates) == 0 {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("no certificates found for this event"))
	}

	// 4. Get all event issuers
	eventIssuers, err := uc.EventIssuerDataGateway.GetEventIssuersByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if len(eventIssuers) == 0 {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("no issuers found for this event"))
	}

	// 5. Get certificate config
	certificateConfig, err := uc.EventCertificateConfigDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// 6. Check precondition: all issuers must have signed (signatures are linked to config, not individual certificates)
	signatures, err := uc.EventCertificateSignatureDataGateway.GetEventCertificateSignaturesByEventCertificateConfigID(ctx, certificateConfig.ID)
	if err != nil {
		return nil, err
	}

	// Check if we have signatures from all issuers
	if len(signatures) != len(eventIssuers) {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("not all issuers have signature records for the certificate config"))
	}

	// Check if all issuers have actually signed (issuer_signature is not null)
	for _, signature := range signatures {
		if utils.DerefOrEmpty(signature.IssuerSignature) == "" {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, fmt.Errorf("not all issuers have signed the certificates. Please ensure all issuers have completed signing before publishing"))
		}
	}

	// 7. Mark the event certificate config as published
	_, err = uc.EventCertificateConfigDg.ToggleEventCertificateConfigPublished(ctx, generated.ToggleEventCertificateConfigPublishedParams{
		EventID:     eventID,
		IsPublished: true,
	})
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, fmt.Errorf("failed to mark certificate config as published: %w", err))
	}

	// 7. Create inbox messages for each certificate receiver
	inboxMessagesCreated := 0
	for _, certificate := range certificates {
		// Check if inbox message already exists for this certificate
		if certificate.InboxMessageId != nil {
			// Skip if inbox message already exists
			continue
		}

		// Skip if no receiver email, no credential ID, and no direct wallet address (can't identify the receiver)
		if utils.DerefOrEmpty(certificate.ReceiverEmail) == "" &&
			certificate.ReceiverCredentialId == nil &&
			utils.DerefOrEmpty(certificate.ReceiverWalletAddress) == "" {
			continue
		}

		// Determine receiver email (may be empty for web3-only users)
		var receiverEmail string
		if certificate.ReceiverEmail != nil {
			receiverEmail = *certificate.ReceiverEmail
		}

		// Look up wallet address: prefer credential lookup, fall back to direct wallet address on the certificate
		var receiverWalletAddress *string
		if certificate.ReceiverCredentialId != nil {
			cred, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, *certificate.ReceiverCredentialId)
			if err == nil {
				receiverWalletAddress = &cred.WalletAddress
			}
		} else if utils.DerefOrEmpty(certificate.ReceiverWalletAddress) != "" {
			receiverWalletAddress = certificate.ReceiverWalletAddress
		}

		// Create message content with translations
		messageContent := map[string]string{
			"th": fmt.Sprintf("คุณได้รับเชิญให้รับใบประกาศนียบัตรสำหรับ: %s", event.Title),
			"en": fmt.Sprintf("You have been invited to claim your certificate for: %s", event.Title),
		}
		messageContentJSON, err := json.Marshal(messageContent)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, fmt.Errorf("failed to marshal message content: %w", err))
		}

		fallbackMessage := "You have been invited to claim your certificate"

		// Create inbox message (supports both authenticated and non-authenticated receivers)
		inboxMessage, err := uc.InboxMessageDg.CreateInboxMessage(ctx, offchain_datagateway.CreateInboxMessageParameters{
			SenderCredentialID:     &currentUser.UserId,
			ReceiverCredentialID:   certificate.ReceiverCredentialId,
			ReceiverEmail:          receiverEmail,
			ReceiverWalletAddress:  receiverWalletAddress,
			MessageType:            int(entity.InboxMessageTypeEventCertificateInvitation),
			MessageContent:         string(messageContentJSON),
			FallbackMessageContent: &fallbackMessage,
			IsRead:                 0,
		})
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, fmt.Errorf("failed to create inbox message: %w", err))
		}

		// 8. Update the certificate with the inbox_message_id
		_, err = uc.EventCertificateDataGateway.UpdateEventCertificateInboxMessageID(ctx, certificate.Id, inboxMessage.Id)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, fmt.Errorf("failed to update certificate with inbox message ID: %w", err))
		}

		inboxMessagesCreated++
	}

	return &PublishEventCertificatesResponse{
		EventID:              eventID,
		PublishedCount:       len(certificates),
		IsPublished:          true,
		InboxMessagesCreated: inboxMessagesCreated,
	}, nil
}
