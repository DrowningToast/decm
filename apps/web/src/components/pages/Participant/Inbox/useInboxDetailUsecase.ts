import { useInboxMessage } from "@/hooks/inbox/useInboxMessage";

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
              id: apiMessage.id,
              title: apiMessage.title,
              sender: apiMessage.sender,
              date: apiMessage.date,
              status: apiMessage.status,
              contentType: apiMessage.contentType,
              eventId: apiMessage.eventId,
              certificateId: apiMessage.certificateId,
              description: apiMessage.description,
              isUserInEvent: apiMessage.isUserInEvent,
          }
        : null;

    return { inboxDetail, isLoading, error };
};
