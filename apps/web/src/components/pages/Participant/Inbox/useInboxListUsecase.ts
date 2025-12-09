import { useSearchNotificationNavStore } from "@/components/BottomNav/stores/notifications";
import { useInboxMessages } from "@/hooks/inbox/useInboxMessages";
import { type InboxMessage, type InboxMessageType } from "@/services/InboxService/InboxService";
import FuzzySearch from "fuzzy-search";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";

export interface InboxItem {
    id: string;
    title: string;
    sender: string;
    date: string;
    status: "pending" | "available" | "expired" | "action-required";
}

const deriveStatus = (isRead: boolean): "pending" | "available" | "expired" | "action-required" => {
    if (!isRead) return "action-required";
    return "available";
};

const formatDate = (date: Date): string => {
    return date.toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
    });
};

const getMessageTitle = (
    messageType: InboxMessageType,
    t: ReturnType<typeof useTranslation>["t"],
): string => {
    switch (messageType) {
        case "event_registration_invitation":
            return t(
                "participant.inbox.messageType.eventRegistrationInvitation",
                "Event Registration Invitation",
            );
        case "event_certificate_invitation":
            return t(
                "participant.inbox.messageType.certificateInvitation",
                "Certificate Invitation",
            );
        case "general":
        default:
            return t("participant.inbox.messageType.general", "Message");
    }
};

const mapInboxMessageToInboxItem = (
    message: InboxMessage,
    t: ReturnType<typeof useTranslation>["t"],
): InboxItem => {
    return {
        id: message.id,
        title: getMessageTitle(message.messageType, t),
        sender: message.senderCredentialEmail ?? message.senderCredentialWalletAddress ?? "Unknown",
        date: formatDate(message.createdAt),
        status: deriveStatus(message.isRead),
    };
};

export const useInboxListUsecase = () => {
    const { t } = useTranslation();
    const { searchQuery } = useSearchNotificationNavStore();

    // Fetch inbox messages from API
    const { data: apiMessages, isLoading, error } = useInboxMessages();

    // Filter out unclaimed certificates from events the user has joined
    const filteredMessages = useMemo(() => {
        return apiMessages.filter((message) => {
            // For certificate invitation messages, filter out unclaimed certificates from events user has joined
            if (message.messageType === "event_certificate_invitation") {
                const hasJoined = message.hasParticipantJoinedEvent === true;
                const isUnclaimed = !message.tokenId; // tokenId is undefined or null for unclaimed certificates

                // Filter out if: user has joined the event AND certificate is unclaimed
                if (hasJoined && isUnclaimed) {
                    return false;
                }
            }
            return true;
        });
    }, [apiMessages]);

    // Transform API response to match InboxItem interface
    const inboxItems = useMemo(() => {
        return filteredMessages.map((message) => mapInboxMessageToInboxItem(message, t));
    }, [filteredMessages, t]);

    // Apply fuzzy search filter
    const searcher = new FuzzySearch(inboxItems, ["title", "sender"], {
        caseSensitive: false,
    });
    const normalizedQuery = searchQuery?.trim() ?? "";
    const filteredItems =
        normalizedQuery.length === 0 ? inboxItems : searcher.search(normalizedQuery);

    return { inboxItems: filteredItems, isLoading, error };
};
