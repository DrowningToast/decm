import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";

export function useEvent(eventId: string) {
    const {
        data: event,
        isLoading: isLoadingEvent,
        isError: isLoadingEventError,
    } = useQuery({
        queryKey: ["event", eventId],
        queryFn: async () => coreApiClient.v1.getEventById({ eventId }),
    });

    return {
        event,
        isLoadingEvent,
        isLoadingEventError,
    };
}
