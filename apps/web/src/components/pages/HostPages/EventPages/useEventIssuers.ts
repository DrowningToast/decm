import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";

export function useEventIssuers(eventId: string) {
    const { data: eventIssuers, isLoading: isLoadingEventIssuers } = useQuery({
        queryKey: QUERY_KEY.event.issuers.byEventId(eventId),
        queryFn: () =>
            coreApiClient.v1.getEventIssuersByEventId({
                eventId,
            }),
    });

    return {
        eventIssuers,
        isLoadingEventIssuers,
    };
}
