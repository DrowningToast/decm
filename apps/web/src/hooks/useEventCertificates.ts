import { useQuery } from "@tanstack/react-query";
import { certificateService } from "@/services/services";
import { QUERY_KEY } from "@/lib/queryKeys";
import type { Certificate } from "@/services/CertificateService/mapper";

interface UseEventCertificatesReturn {
    certificates: Certificate[];
    isLoading: boolean;
    error: Error | null;
    refetch: () => void;
}

export const useEventCertificates = (eventId: string): UseEventCertificatesReturn => {
    const { data, isLoading, error, refetch } = useQuery({
        queryKey: QUERY_KEY.event.certificates(eventId),
        queryFn: () => certificateService.getEventCertificates(eventId),
        enabled: !!eventId,
    });

    return {
        certificates: data || [],
        isLoading,
        error: error as Error | null,
        refetch,
    };
};
