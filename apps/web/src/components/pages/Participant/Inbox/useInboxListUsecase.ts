import { useSearchNotificationNavStore } from "@/components/BottomNav/stores/notifications";
import { useInboxMessages } from "@/hooks/inbox/useInboxMessages";
import FuzzySearch from "fuzzy-search";
import { useMemo } from "react";

export interface InboxItem {
    id: string;
    title: string;
    sender: string;
    date: string;
    status: "pending" | "available" | "expired" | "action-required";
}

export const useInboxListUsecase = () => {
    const { searchQuery } = useSearchNotificationNavStore();

    // Fetch inbox messages from API
    const {
        data: apiMessages = [],
        isLoading,
        error,
    } = useInboxMessages({
        limit: 100,
        offset: 0,
    });

    // Transform API response to match InboxItem interface
    const inboxItems = useMemo(() => {
        return apiMessages.map((message) => ({
            id: message.id,
            title: message.title,
            sender: message.sender,
            date: message.date,
            status: message.status,
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
