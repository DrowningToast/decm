import { useQuery } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import { QUERY_KEY } from "@/lib/queryKeys";

export function useCertificateFontFamilies() {
    const { data, isLoading, error } = useQuery({
        queryKey: QUERY_KEY.event.certificate.fontFamilies,
        queryFn: async () => {
            const response = await coreApiClient.v1.getEventCertificateFontFamilies();
            return response;
        },
    });

    return {
        fontFamilies: data?.font_families || [],
        isLoading,
        error,
    };
}
