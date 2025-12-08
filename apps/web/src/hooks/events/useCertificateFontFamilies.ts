import { coreApiClient } from "@/lib/api/api";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useQuery } from "@tanstack/react-query";

export function useCertificateFontFamilies() {
    const {
        data: fontFamilies,
        isLoading,
        error,
    } = useQuery({
        queryKey: QUERY_KEY.event.certificate.fontFamilies,
        queryFn: async () => {
            const response =
                await coreApiClient.certificateFontFamilies.getEventCertificateFontFamilies();
            return response.font_families || [];
        },
        staleTime: Infinity, // Font families don't change often
        gcTime: 1000 * 60 * 30, // Cache for 30 minutes
    });

    return {
        fontFamilies: fontFamilies || [],
        isLoading,
        error,
    };
}
