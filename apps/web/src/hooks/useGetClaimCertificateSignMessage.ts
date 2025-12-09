import { useQuery } from "@tanstack/react-query";
import { certificateService } from "@/services/services";
import { QUERY_KEY } from "@/lib/queryKeys";

interface UseGetClaimCertificateSignMessageOptions {
    certificateId: string;
    enabled?: boolean;
}

/**
 * Hook to fetch the sign message for claiming a certificate via wallet signature
 * This message needs to be signed by the user's wallet
 */
export function useGetClaimCertificateSignMessage({
    certificateId,
    enabled = true,
}: UseGetClaimCertificateSignMessageOptions) {
    const query = useQuery({
        queryKey: QUERY_KEY.certificate.claimSignMessage(certificateId),
        queryFn: () => certificateService.getClaimCertificateSignMessage(certificateId),
        enabled: enabled && !!certificateId,
        // Cache for 5 minutes
        staleTime: 5 * 60 * 1000,
    });

    return {
        signMessage: query.data?.signMessage,
        isLoading: query.isLoading,
        error: query.error,
        refetch: query.refetch,
    };
}
