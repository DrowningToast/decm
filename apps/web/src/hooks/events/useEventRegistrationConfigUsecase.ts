import { QUERY_KEY } from "@/lib/queryKeys";
import { eventRegistrationService } from "@/services/services";
import { useQuery } from "@tanstack/react-query";
export const useEventRegistrationConfigUsecase = ({ eventId }: { eventId: string }) => {
    const { data, isLoading, error } = useQuery({
        queryKey: QUERY_KEY.event.registrationConfig(eventId),
        queryFn: () => eventRegistrationService.getConfiguration(eventId),
        enabled: !!eventId,
    });

    return { data, isLoading, error };
};
