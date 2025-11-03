import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";

export function useEventCertificateConfig(eventId: string) {
    const {
        data: eventCertificateConfig,
        isLoading: isLoadingEventCertificateConfig,
        error: errorEventCertificateConfig,
    } = useQuery({
        queryKey: queryKeys.event.certificate.config(eventId),
        queryFn: () => coreApiClient.v1.getEventCertificateConfig({ eventId: eventId }),
    });

    return {
        data: eventCertificateConfig,
        isLoading: isLoadingEventCertificateConfig,
        error: errorEventCertificateConfig,
    };
}
