import { useInboxMessage } from "@/hooks/inbox/useInboxMessage";
import { type InboxMessageDetail } from "@/services/InboxService/InboxService";
import { EntityInboxMessageType } from "@decm/api";
import { useTranslation } from "react-i18next";

export type InboxContentType = "event-invitation" | "certificate" | "general";

export interface InboxDetail {
    id: string;
    title: string;
    sender: string;
    date: string;
    status: "pending" | "available" | "expired" | "action-required";
    contentType: InboxContentType;
    eventId?: string;
    certificateId?: string;
    description?: string;
    isUserInEvent?: boolean;
    // Action status fields
    acceptedAt?: Date;
    tokenId?: string;
}

interface UseInboxDetailOptions {
    inboxId: string;
}

const deriveStatus = (
    message: InboxMessageDetail,
): "pending" | "available" | "expired" | "action-required" => {
    // Check for expired status based on message type
    if (
        message.entityInboxMessageType ===
        EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation
    ) {
        if (message.cancelledAt) return "expired";
        if (message.validUntil && message.validUntil < new Date()) return "expired";
    }
    if (
        message.entityInboxMessageType ===
        EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation
    ) {
        if (message.revokedAt) return "expired";
    }
    if (!message.isRead) return "action-required";
    return "available";
};

const deriveContentType = (entityInboxMessageType: EntityInboxMessageType): InboxContentType => {
    switch (entityInboxMessageType) {
        case EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation:
            return "event-invitation";
        case EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation:
            return "certificate";
        case EntityInboxMessageType.InboxMessageTypeGeneral:
        default:
            return "general";
    }
};

const formatDate = (date: Date): string => {
    return date.toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
};

const getMessageTitle = (
    entityInboxMessageType: EntityInboxMessageType,
    t: ReturnType<typeof useTranslation>["t"],
): string => {
    switch (entityInboxMessageType) {
        case EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation:
            return t(
                "participant.inbox.messageType.eventRegistrationInvitation",
                "Event Registration Invitation",
            );
        case EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation:
            return t(
                "participant.inbox.messageType.certificateInvitation",
                "Certificate Invitation",
            );
        case EntityInboxMessageType.InboxMessageTypeGeneral:
        default:
            return t("participant.inbox.messageType.general", "Message");
    }
};

const mapInboxMessageDetailToInboxDetail = (
    message: InboxMessageDetail,
    t: ReturnType<typeof useTranslation>["t"],
): InboxDetail => {
    const base = {
        id: message.id,
        title: getMessageTitle(message.entityInboxMessageType, t),
        sender: message.senderCredentialEmail ?? message.senderCredentialWalletAddress ?? "Unknown",
        date: formatDate(message.createdAt),
        status: deriveStatus(message),
        contentType: deriveContentType(message.entityInboxMessageType),
        description: message.messageContent,
    };

    // Add type-specific fields
    if (
        message.entityInboxMessageType ===
        EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation
    ) {
        return { ...base, eventId: message.eventId, acceptedAt: message.acceptedAt };
    }
    if (
        message.entityInboxMessageType ===
        EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation
    ) {
        return {
            ...base,
            certificateId: message.certificateId,
            tokenId: message.tokenId,
            isUserInEvent: message.hasParticipantJoinedEvent,
        };
    }

    return base;
};

export const useInboxDetailUsecase = (options: UseInboxDetailOptions) => {
    const { t } = useTranslation();

    // Fetch inbox message from API
    const {
        data: apiMessage,
        isLoading,
        error,
    } = useInboxMessage({
        messageId: options.inboxId,
    });

    // Transform API response to match InboxDetail interface
    const inboxDetail: InboxDetail | null = apiMessage
        ? mapInboxMessageDetailToInboxDetail(apiMessage, t)
        : null;

    return { inboxDetail, isLoading, error };
};
