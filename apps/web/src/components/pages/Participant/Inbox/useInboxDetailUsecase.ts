import { useInboxMessage } from "@/hooks/inbox/useInboxMessage";
import { EntityInboxMessageType } from "@decm/api";

export type InboxContentType = "event-invitation" | "certificate";

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
}

interface UseInboxDetailOptions {
    inboxId: string;
}

const deriveStatus = (message: {
    is_read?: number;
    cancelled_at?: string;
    deleted_at?: string;
    hidden_at?: string;
    valid_until?: string;
}): "pending" | "available" | "expired" | "action-required" => {
    if (message.cancelled_at) return "expired";
    if (message.valid_until && new Date(message.valid_until) < new Date()) return "expired";
    if (message.is_read === 0) return "action-required";
    return "available";
};

const deriveContentType = (messageType?: EntityInboxMessageType): InboxContentType => {
    if (messageType === EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation) {
        return "event-invitation";
    }
    if (messageType === EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation) {
        return "certificate";
    }
    return "event-invitation";
};

export const useInboxDetailUsecase = (options: UseInboxDetailOptions) => {
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
        ? {
              id: apiMessage.id ?? "",
              title:
                  apiMessage.message_type ===
                  EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation
                      ? "Event Registration Invitation"
                      : apiMessage.message_type ===
                          EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation
                        ? "Certificate Invitation"
                        : "Message",
              sender:
                  apiMessage.sender_credential_email ??
                  apiMessage.sender_credential_wallet_address ??
                  "Unknown",
              date: apiMessage.created_at ?? "",
              status: deriveStatus(apiMessage),
              contentType: deriveContentType(apiMessage.message_type),
              eventId: apiMessage.event_id,
              description: apiMessage.message_content,
          }
        : null;

    return { inboxDetail, isLoading, error };
};
