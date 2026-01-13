import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";

export function useEvent(eventId: string) {
    const {
        data: event,
        isLoading: isLoadingEvent,
        isError: isLoadingEventError,
    } = useQuery({
        queryKey: QUERY_KEY.event.byId(eventId),
        queryFn: async () => coreApiClient.v1.getEventViewmodelById({ eventId }),
    });

    return {
        event,
        isLoadingEvent,
        isLoadingEventError,
    };
}
