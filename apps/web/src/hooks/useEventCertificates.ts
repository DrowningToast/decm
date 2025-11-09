import { useQuery } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import { QUERY_KEY } from "@/lib/queryKeys";

export interface EventCertificate {
    id?: string;
    event_id?: string;
    receiver_credential_id?: string;
    receiver_email?: string;
    name?: string;
    academic_institution?: string;
    certificate_title?: string;
    certificate_subtitle?: string;
    event_contract_address?: string;
    event_certificate_address?: string;
    certificate_token_id?: string;
    created_at?: string;
    revoked_at?: string;
}

interface UseEventCertificatesReturn {
    certificates: EventCertificate[];
    isLoading: boolean;
    error: Error | null;
    refetch: () => void;
}

export const useEventCertificates = (eventId: string): UseEventCertificatesReturn => {
    const {
        data: response,
        isLoading,
        error,
        refetch,
    } = useQuery({
        queryKey: [QUERY_KEY.event.certificates(eventId)],
        queryFn: async () => {
            const response = await coreApiClient.v1.getEventCertificates({ eventId });
            return response;
        },
        enabled: !!eventId,
    });

    return {
        certificates: response?.certificates || [],
        isLoading,
        error,
        refetch,
    };
};
