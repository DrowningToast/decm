import { useSearchNotificationNavStore } from "@/components/BottomNav/stores/notifications";
import { useQuery } from "@tanstack/react-query";
import FuzzySearch from "fuzzy-search";

export interface InboxItem {
    id: string;
    title: string;
    sender: string;
    date: string;
    status: "pending" | "available" | "expired" | "action-required";
}

const MOCK_INBOX_ITEMS: InboxItem[] = [
    {
        id: "1",
        title: "Event Invitation",
        sender: "ToBeIT69",
        date: "24 Sep 2025",
        status: "pending",
    },
    {
        id: "2",
        title: "Event Invitation",
        sender: "ToBeIT69",
        date: "24 Sep 2025",
        status: "available",
    },
    {
        id: "3",
        title: "Event Invitation",
        sender: "ToBeIT69",
        date: "24 Sep 2025",
        status: "expired",
    },
    {
        id: "4",
        title: "New certificate",
        sender: "ToBeIT69",
        date: "24 Sep 2025",
        status: "action-required",
    },
    {
        id: "5",
        title: "New certificate",
        sender: "ToBeIT69",
        date: "24 Sep 2025",
        status: "available",
    },
];

export const useInboxListUsecase = () => {
    // const api = useApi();
    const { searchQuery } = useSearchNotificationNavStore();

    const {
        data: inboxItems = [],
        isLoading,
        error,
    } = useQuery({
        queryKey: ["participant-inbox"],
        queryFn: async () => {
            try {
                // TODO: Replace with actual API call once endpoint is available
                // const response = await api.getParticipantInbox();
                return MOCK_INBOX_ITEMS;
            } catch (error) {
                console.error("Failed to fetch inbox items:", error);
                return [] as InboxItem[];
            }
        },
    });

    const searcher = new FuzzySearch(inboxItems, ["title", "sender"], {
        caseSensitive: false,
    });
    const normalizedQuery = searchQuery?.trim() ?? "";
    const filteredItems =
        normalizedQuery.length === 0 ? inboxItems : searcher.search(normalizedQuery);

    return { inboxItems: filteredItems, isLoading, error };
};
