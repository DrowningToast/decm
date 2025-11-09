import { useQuery } from "@tanstack/react-query";
import { EventService } from "@/services/EventService";
import { QUERY_KEY } from "@/lib/queryKeys";

interface UseGetEventListParams {
    includeActiveEvents?: boolean;
    includeInactiveEvents?: boolean;
    includeClosedEvents?: boolean;
    onlyUserJoinedEvents?: boolean;
    enabled?: boolean;
}

/**
 * Custom hook to fetch event list for participant view
 * @param params Configuration options for the query
 * @returns Query result with loading state and event data
 */
export function useGetEventList(params: UseGetEventListParams = {}) {
    const {
        includeActiveEvents,
        includeInactiveEvents,
        includeClosedEvents,
        onlyUserJoinedEvents,
        enabled = true,
    } = params;

    const eventService = new EventService();

    return useQuery({
        queryKey: QUERY_KEY.event.list({
            includeActiveEvents,
            includeInactiveEvents,
            includeClosedEvents,
            onlyUserJoinedEvents,
        }),
        queryFn: () =>
            eventService.getEvents({
                includeActiveEvents,
                includeInactiveEvents,
                includeClosedEvents,
                onlyUserJoinedEvents,
            }),
        enabled,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });
}
