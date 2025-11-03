import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";

export function useEventIssuers(eventId: string) {
    const { data: eventIssuers, isLoading: isLoadingEventIssuers } = useQuery({
        queryKey: queryKeys.event.issuers.byEventId(eventId),
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
