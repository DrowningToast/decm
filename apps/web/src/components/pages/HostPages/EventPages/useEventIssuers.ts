import { useQuery } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";
import { eventService } from "@/services/services";

export function useEventIssuers(eventId: string) {
    const { data: eventIssuers, isLoading: isLoadingEventIssuers } = useQuery({
        queryKey: QUERY_KEY.event.issuers.byEventId(eventId),
        queryFn: () => eventService.getEventIssuersByEventId(eventId),
    });

    return {
        eventIssuers,
        isLoadingEventIssuers,
    };
}
