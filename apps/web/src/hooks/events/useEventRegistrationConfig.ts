import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";

export function useEventRegistrationConfig(eventId: string) {
    const { data, isLoading, error } = useQuery({
        queryKey: ["event-registration-config", eventId],
        queryFn: () => coreApiClient.v1.getEventRegistrationConfig({ eventId }),
        enabled: !!eventId,
    });

    return { data, isLoading, error };
}
