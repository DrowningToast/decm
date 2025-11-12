import { useSearchEventNavStore } from "@/components/BottomNav/stores/events";
import FuzzySearch from "fuzzy-search";
import { useGetEventList } from "@/hooks/events/useGetEventList";
import { useMemo } from "react";

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

type EventFilterType = "all" | "my-events";

interface UseEventsListOptions {
    filterType?: EventFilterType;
}

export const useEventsListUsecase = (options?: UseEventsListOptions) => {
    const { searchQuery } = useSearchEventNavStore();

    // Use the API hook with appropriate filters
    const {
        data: apiEvents,
        isLoading,
        error,
    } = useGetEventList({
        includeActiveEvents: true,
        includeInactiveEvents: true,
        includeClosedEvents: true,
        onlyUserJoinedEvents: options?.filterType === "my-events",
    });

    // Transform API response to match Event interface
    const events = useMemo(() => {
        if (!apiEvents?.events) {
            return [];
        }

        return apiEvents.events.map((event) => ({
            id: event.id,
            name: event.title,
            description: event.short_description,
            eventName: event.title,
            contactAddress: event.contact_number,
            dateTime: event.start_date,
            finalCallDate: event.end_date,
            status: event.event_status,
            accessType: event.is_public ? ("public" as const) : ("password" as const),
        }));
    }, [apiEvents]);

    // Apply fuzzy search filter
    const searcher = new FuzzySearch(events, ["name"], {
        caseSensitive: false,
    });
    const filteredEvents = searchQuery ? searcher.search(searchQuery) : events;

    return { events: filteredEvents, isLoading, error };
};
