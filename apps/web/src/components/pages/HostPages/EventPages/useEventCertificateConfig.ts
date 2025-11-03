import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";

export function useEventCertificateConfig(eventId: string) {
    const {
        data: eventCertificateConfig,
        isLoading: isLoadingEventCertificateConfig,
        error: errorEventCertificateConfig,
    } = useQuery({
        queryKey: QUERY_KEY.event.certificate.config(eventId),
        queryFn: () => coreApiClient.v1.getEventCertificateConfig({ eventId: eventId }),
    });

    return {
        data: eventCertificateConfig,
        isLoading: isLoadingEventCertificateConfig,
        error: errorEventCertificateConfig,
    };
}
