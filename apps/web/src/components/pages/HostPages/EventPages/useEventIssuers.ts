import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";

export function useEventIssuers(eventId: string) {
    const { data: eventIssuers, isLoading: isLoadingEventIssuers } = useQuery({
        queryKey: ["event", eventId, "issuers"],
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
