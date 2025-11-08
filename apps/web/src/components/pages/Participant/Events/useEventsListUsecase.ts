import { useSearchEventNavStore } from "@/components/BottomNav/stores/events";
import { useQuery } from "@tanstack/react-query";
import FuzzySearch from "fuzzy-search";

export interface Event {
    id: string;
    name: string;
    description: string;
    eventName: string;
    contactAddress?: string;
    dateTime?: string;
    finalCallDate?: string;
    status: "accepting" | "closed";
    accessType: "public" | "password" | "invite-only";
    seatsAvailable?: number;
    totalSeats?: number;
    participationStatus?: "applied" | "accepted" | "rejected" | "shortlisted";
    isShortlisted?: boolean;
    requiresPassword?: boolean;
    image?: string;
}

const mockEvents: Event[] = [
    {
        id: "1",
        name: "ToBelT69",
        description: "Description 1",
        eventName: "ToBelT69",
        contactAddress: "0x1234...abcd",
        dateTime: "2024-09-24",
        finalCallDate: "2024-09-24",
        status: "accepting",
        accessType: "password",
        requiresPassword: true,
    },
    {
        id: "2",
        name: "ToBelT69",
        description: "Description 2",
        eventName: "ToBelT69",
        contactAddress: "0x9876...4321",
        dateTime: "2024-09-25",
        finalCallDate: "2024-09-25",
        status: "accepting",
        accessType: "invite-only",
    },
    {
        id: "3",
        name: "ToBelT69",
        description: "Description 3",
        eventName: "ToBelT69",
        contactAddress: "0xa1b2...c3d4",
        dateTime: "2024-09-20",
        finalCallDate: "2024-09-20",
        status: "closed",
        accessType: "public",
    },
];

type EventFilterType = "all" | "my-events";

interface UseEventsListOptions {
    filterType?: EventFilterType;
}

export const useEventsListUsecase = (options?: UseEventsListOptions) => {
    // const api = useApi();
    const { searchQuery } = useSearchEventNavStore();
    const {
        data: events = [],
        isLoading,
        error,
    } = useQuery({
        queryKey: ["participant-events"],
        queryFn: async () => {
            try {
                // TODO: Replace with actual API call once endpoint is available
                // const response = await api.getParticipantEvents();
                return mockEvents;
            } catch (error) {
                console.error("Failed to fetch events:", error);
                return [] as Event[];
            }
        },
    });

    const searcher = new FuzzySearch(events, ["name"], {
        caseSensitive: false,
    });
    let filteredEvents = searcher.search(searchQuery);

    // Apply filter based on filter type
    if (options?.filterType === "my-events") {
        filteredEvents = filteredEvents.filter((event) => {
            // Filter events that the user has joined (has participation status)
            return event.participationStatus !== undefined;
        });
    }

    return { events: filteredEvents, isLoading, error };
};
