import { EntityInboxMessageType, type EntityInboxMessage } from "@decm/api";
import { InboxMessageContent, type InboxMessage, type InboxMessageType } from "./InboxService";

export const mapEntityInboxMessageTypeToInboxMessageType = (
    entityInboxMessageType: EntityInboxMessageType,
): InboxMessageType => {
    switch (entityInboxMessageType) {
        case EntityInboxMessageType.InboxMessageTypeGeneral:
            return "general";
        case EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation:
            return "event_registration_invitation";
        case EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation:
            return "event_certificate_invitation";
    }
    return "general";
};

export const parseInboxMessageContent = (content: string): InboxMessageContent => {
    return InboxMessageContent.parse(JSON.parse(content));
};

export const mapEntityInboxMessageToInboxMessage = (
    entityInboxMessage: EntityInboxMessage,
): InboxMessage => {
    return {
        id: entityInboxMessage.id,
        messageType: mapEntityInboxMessageTypeToInboxMessageType(entityInboxMessage.message_type),
        messageContent: parseInboxMessageContent(entityInboxMessage.message_content),
        isRead: entityInboxMessage.is_read === 1,
        fallbackMessageContent: entityInboxMessage.fallback_message_content,
        receiverCredentialId: entityInboxMessage.receiver_credential_id,
        receiverEmail: entityInboxMessage.receiver_email,
        receiverWalletAddress: entityInboxMessage.receiver_wallet_address,
        senderCredentialId: entityInboxMessage.sender_credential_id,
        createdAt: entityInboxMessage.created_at,
        updatedAt: entityInboxMessage.updated_at,
        hiddenAt: entityInboxMessage.hidden_at,
        deletedAt: entityInboxMessage.deleted_at,
    };
};
