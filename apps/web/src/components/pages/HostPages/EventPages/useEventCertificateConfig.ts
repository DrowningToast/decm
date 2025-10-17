import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";

export function useEventCertificateConfig(eventId: string) {
    const {
        data: eventCertificateConfig,
        isLoading: isLoadingEventCertificateConfig,
        error: errorEventCertificateConfig,
    } = useQuery({
        queryKey: ["event", eventId, "certificate", "config"],
        queryFn: () => coreApiClient.v1.getEventCertificateConfig({ eventId: eventId }),
    });

    return {
        data: eventCertificateConfig,
        isLoading: isLoadingEventCertificateConfig,
        error: errorEventCertificateConfig,
    };
}
