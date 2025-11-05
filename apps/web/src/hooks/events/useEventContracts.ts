import { coreApiClient } from "@/lib/api/api";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useQuery } from "@tanstack/react-query";

export function useEventContract(eventId: string) {
    const { data, isLoading, error } = useQuery({
        queryKey: QUERY_KEY.event.contract(eventId),
        queryFn: () => coreApiClient.v1.getEventContractByEventId({ eventId }),
        enabled: !!eventId,
    });

    return { data, isLoading, error };
}
