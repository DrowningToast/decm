import {
    EntityInboxMessageType,
    type EntityInboxMessage,
    type InboxInboxMessagesViewModel,
    type InboxmessagesGetInboxMessageResponse,
} from "@decm/api";
import type { InboxMessage, InboxMessageDetail, InboxMessageType } from "./InboxService";

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
        default:
            return "general";
    }
};

export const mapInboxMessagesViewModelToInboxMessage = (
    viewModel: InboxInboxMessagesViewModel,
): InboxMessage => {
    return {
        id: viewModel.id,
        messageType: mapEntityInboxMessageTypeToInboxMessageType(viewModel.message_type),
        messageContent: viewModel.message_content ?? "",
        isRead: viewModel.is_read === 1,
        receiverEmail: viewModel.receiver_email ?? undefined,
        receiverWalletAddress: viewModel.receiver_wallet_address ?? undefined,
        senderCredentialEmail: viewModel.sender_credential_email ?? undefined,
        senderCredentialWalletAddress: viewModel.sender_credential_wallet_address ?? undefined,
        createdAt: new Date(viewModel.created_at),
        updatedAt: new Date(viewModel.updated_at),
        hiddenAt: viewModel.hidden_at ? new Date(viewModel.hidden_at) : undefined,
        deletedAt: viewModel.deleted_at ? new Date(viewModel.deleted_at) : undefined,
        // Certificate-specific fields (only populated for certificate invitation messages)
        eventId: viewModel.event_id ?? undefined,
        certificateId: viewModel.certificate_id ?? undefined,
        certificateTitle: viewModel.certificate_title ?? undefined,
        tokenId: viewModel.token_id ?? undefined,
        hasParticipantJoinedEvent: viewModel.has_participant_joined_event ?? undefined,
        eventName: viewModel.event_name ?? undefined,
    };
};

export const mapInboxMessagesDetailViewModelToInboxMessageDetail = (
    response: InboxmessagesGetInboxMessageResponse,
): InboxMessageDetail => {
    const mapped = {
        id: response.inbox_message.id,
        messageType: mapEntityInboxMessageTypeToInboxMessageType(
            response.inbox_message.message_type,
        ),
        messageContent: response.inbox_message.message_content,
        isRead: response.inbox_message.is_read === 1,
        receiverEmail: response.inbox_message.receiver_email,
        receiverWalletAddress: response.inbox_message.receiver_wallet_address,
        senderCredentialEmail: response.inbox_message.sender_credential_email,
        senderCredentialWalletAddress: response.inbox_message.sender_credential_wallet_address,
        createdAt: new Date(response.inbox_message.created_at),
        updatedAt: new Date(response.inbox_message.updated_at),
        hiddenAt: response.inbox_message.hidden_at
            ? new Date(response.inbox_message.hidden_at)
            : undefined,
        deletedAt: response.inbox_message.deleted_at
            ? new Date(response.inbox_message.deleted_at)
            : undefined,
        code: response.event_registration_invitation?.code,
        eventId: response.event_registration_invitation?.event_id,
        email: response.event_registration_invitation?.email,
        firstName: response.event_registration_invitation?.first_name,
        lastName: response.event_registration_invitation?.last_name,
        phoneNumber: response.event_registration_invitation?.phone_number,
        academicInstitution: response.event_registration_invitation?.academic_institution,
        validUntil: response.event_registration_invitation?.valid_until
            ? new Date(response.event_registration_invitation?.valid_until)
            : undefined,
        cancelledAt: response.event_registration_invitation?.cancelled_at
            ? new Date(response.event_registration_invitation?.cancelled_at)
            : undefined,
        acceptedAt: response.event_registration_invitation?.accepted_at
            ? new Date(response.event_registration_invitation?.accepted_at)
            : undefined,
    } as const;
    if (
        response.inbox_message_type ===
        EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation
    ) {
        return {
            ...mapped,
            entityInboxMessageType:
                EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
            eventId: response.event_registration_invitation?.event_id,
            code: response.event_registration_invitation?.code,
            email: response.event_registration_invitation?.email,
            firstName: response.event_registration_invitation?.first_name,
        };
    }
    if (
        response.inbox_message_type ===
        EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation
    ) {
        if (!response.event_certificate) {
            throw new Error(
                "inbox message type is event certificate invitation but event certificate is not found",
            );
        }
        // Cast to include has_participant_joined_event field (regenerate API types after backend update)
        const eventCertificate = response.event_certificate as typeof response.event_certificate & {
            has_participant_joined_event?: boolean;
        };
        return {
            ...mapped,
            entityInboxMessageType:
                EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
            eventName: eventCertificate.event_name,
            certificateId: eventCertificate.certificate_id,
            certificateTitle: eventCertificate.certificate_title,
            tokenId: eventCertificate.token_id,
            revokedAt: eventCertificate.revoked_at
                ? new Date(eventCertificate.revoked_at)
                : undefined,
            hasParticipantJoinedEvent: eventCertificate.has_participant_joined_event ?? false,
        };
    }
    return {
        ...mapped,
        entityInboxMessageType: EntityInboxMessageType.InboxMessageTypeGeneral,
    };
};

/**
 * Maps EntityInboxMessage (used in EventRegistration responses) to InboxMessage
 * This is a simpler entity type compared to InboxInboxMessagesViewModel
 */
export const mapEntityInboxMessageToInboxMessage = (entity: EntityInboxMessage): InboxMessage => {
    return {
        id: entity.id,
        messageType: mapEntityInboxMessageTypeToInboxMessageType(entity.message_type),
        messageContent: entity.message_content,
        isRead: entity.is_read === 1,
        receiverEmail: entity.receiver_email,
        receiverWalletAddress: entity.receiver_wallet_address,
        senderCredentialEmail: entity.sender_credential_id, // Entity uses sender_credential_id
        senderCredentialWalletAddress: undefined, // Not available in EntityInboxMessage
        createdAt: new Date(entity.created_at),
        updatedAt: new Date(entity.updated_at),
        hiddenAt: entity.hidden_at ? new Date(entity.hidden_at) : undefined,
        deletedAt: entity.deleted_at ? new Date(entity.deleted_at) : undefined,
    };
};
