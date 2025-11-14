import { useQuery } from "@tanstack/react-query";
import { eventRegistrationService } from "@/services/services";
import { QUERY_KEY } from "@/lib/queryKeys";

export const useEventRegistrationConfiguration = (eventId: string) => {
    const { data, isLoading, error } = useQuery({
        queryKey: QUERY_KEY.event.registrationConfig(eventId),
        queryFn: () => eventRegistrationService.getConfiguration(eventId),
        enabled: !!eventId,
    });

    return { data, isLoading, error };
};
