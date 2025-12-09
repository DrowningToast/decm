import { certificateService } from "@/services/services";
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
            const result = await certificateService.getFontFamilies();
            return result;
        },
        staleTime: Infinity, // Font families don't change often
        gcTime: 1000 * 60 * 30, // Cache for 30 minutes
        retry: 3, // Retry 3 times on failure
    });

    return {
        fontFamilies: fontFamilies || [],
        isLoading,
        error,
    };
}
