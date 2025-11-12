import { useSearchNotificationNavStore } from "@/components/BottomNav/stores/notifications";
import { useInboxMessages } from "@/hooks/inbox/useInboxMessages";
import { EntityInboxMessageType } from "@decm/api";
import FuzzySearch from "fuzzy-search";
import { useMemo } from "react";

export interface InboxItem {
    id: string;
    title: string;
    sender: string;
    date: string;
    status: "pending" | "available" | "expired" | "action-required";
}

const deriveStatus = (message: {
    is_read?: number;
    deleted_at?: string;
    hidden_at?: string;
}): "pending" | "available" | "expired" | "action-required" => {
    if (message.is_read === 0) return "action-required";
    return "available";
};

export const useInboxListUsecase = () => {
    const { searchQuery } = useSearchNotificationNavStore();

    // Fetch inbox messages from API
    const { data: apiMessages = [], isLoading, error } = useInboxMessages();

    // Transform API response to match InboxItem interface
    const inboxItems = useMemo(() => {
        return apiMessages.map((message) => ({
            id: message.id ?? "",
            title:
                message.message_type ===
                EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation
                    ? "Event Registration Invitation"
                    : message.message_type ===
                        EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation
                      ? "Certificate Invitation"
                      : "Message",
            sender:
                message.sender_credential_email ??
                message.sender_credential_wallet_address ??
                "Unknown",
            date: message.created_at ?? "",
            status: deriveStatus(message),
        }));
    }, [apiMessages]);

    // Apply fuzzy search filter
    const searcher = new FuzzySearch(inboxItems, ["title", "sender"], {
        caseSensitive: false,
    });
    const normalizedQuery = searchQuery?.trim() ?? "";
    const filteredItems =
        normalizedQuery.length === 0 ? inboxItems : searcher.search(normalizedQuery);

    return { inboxItems: filteredItems, isLoading, error };
};
