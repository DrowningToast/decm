import { useSearchEventNavStore } from "@/components/BottomNav/stores/events";
import FuzzySearch from "fuzzy-search";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useQuery } from "@tanstack/react-query";
import { eventService } from "@/services/services";

type EventFilterType = "all" | "my-events";

const defaultParam = {
    includeActiveEvents: true,
    includeInactiveEvents: true,
    includeClosedEvents: false,
    onlyUserJoinedEvents: false,
};

interface UseEventsListOptions {
    filterType?: EventFilterType;
}

export const useEventsListUsecase = (options?: UseEventsListOptions) => {
    const { searchQuery } = useSearchEventNavStore();

    const {
        data: events,
        isLoading,
        error,
    } = useQuery({
        queryKey: QUERY_KEY.event.list({
            includeActiveEvents: defaultParam.includeActiveEvents,
            includeInactiveEvents: defaultParam.includeInactiveEvents,
            includeClosedEvents: defaultParam.includeClosedEvents,
            onlyUserJoinedEvents: options?.filterType === "my-events",
        }),
        queryFn: () =>
            eventService.getEvents({
                includeActiveEvents: defaultParam.includeActiveEvents,
                includeInactiveEvents: defaultParam.includeInactiveEvents,
                includeClosedEvents: defaultParam.includeClosedEvents,
                onlyUserJoinedEvents: options?.filterType === "my-events",
            }),
        enabled: true,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });

    // Apply fuzzy search filter
    const searcher = new FuzzySearch(events ?? [], ["name"], {
        caseSensitive: false,
    });
    const filteredEvents = searchQuery ? searcher.search(searchQuery) : events;

    return { events: filteredEvents, isLoading, error };
};
