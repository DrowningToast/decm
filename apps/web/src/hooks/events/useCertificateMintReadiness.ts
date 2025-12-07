import { useQuery } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import type { CoreApiInternalHandlerEventconfigCertificateMintReadinessResponse } from "@decm/api";

export const useCertificateMintReadiness = (eventId: string | undefined) => {
    return useQuery<CoreApiInternalHandlerEventconfigCertificateMintReadinessResponse>({
        queryKey: ["certificate-mint-readiness", eventId],
        queryFn: async () => {
            if (!eventId) {
                throw new Error("Event ID is required");
            }
            return await coreApiClient.v1.checkCertificateMintReadiness({ eventId });
        },
        enabled: !!eventId,
        staleTime: 30000, // 30 seconds
        retry: 1,
    });
};
