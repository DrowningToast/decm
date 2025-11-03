import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";

export function useEventRegistrationConfig(eventId: string) {
    const { data, isLoading, error } = useQuery({
        queryKey: queryKeys.event.registrationConfig(eventId),
        queryFn: () => coreApiClient.v1.getEventRegistrationConfig({ eventId }),
        enabled: !!eventId,
    });

    return { data, isLoading, error };
}
