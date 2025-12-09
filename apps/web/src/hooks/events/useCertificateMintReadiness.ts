import { useQuery } from "@tanstack/react-query";
import { certificateService } from "@/services/services";
import type { CertificateMintReadiness } from "@/services/CertificateService/CertificateService";

export const useCertificateMintReadiness = (eventId: string | undefined) => {
    return useQuery<CertificateMintReadiness>({
        queryKey: ["certificate-mint-readiness", eventId],
        queryFn: async () => {
            if (!eventId) {
                throw new Error("Event ID is required");
            }
            return await certificateService.checkCertificateMintReadiness(eventId);
        },
        enabled: !!eventId,
        staleTime: 30000, // 30 seconds
        retry: 1,
    });
};
