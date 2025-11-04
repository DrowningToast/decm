import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";

export function useEventRegistrationConfig(eventId: string) {
    const { data, isLoading, error } = useQuery({
        queryKey: QUERY_KEY.event.registrationConfig(eventId),
        queryFn: () => coreApiClient.v1.getEventRegistrationConfig({ eventId }),
        enabled: !!eventId,
    });

    return { data, isLoading, error };
}
