import { useQuery } from "@tanstack/react-query";
import { certificateService } from "@/services/services";
import type {
    MyCertificatesViewModel,
    ClaimedCertificate,
    UnclaimedCertificate,
} from "@/services/CertificateService/mapper";
import { QUERY_KEY } from "@/lib/queryKeys";

interface UseMyCertificatesListViewModelReturn {
    claimedCertificates: ClaimedCertificate[];
    unclaimedCertificates: UnclaimedCertificate[];
    totalClaimed: number;
    totalUnclaimed: number;
    isLoading: boolean;
    isError: boolean;
    error: Error | null;
}

export const useMyCertificatesListViewModel = (): UseMyCertificatesListViewModelReturn => {
    const { data, isLoading, isError, error } = useQuery<MyCertificatesViewModel>({
        queryKey: QUERY_KEY.certificate.myCertificatesListViewModel,
        queryFn: () => certificateService.getMyCertificatesList(),
        staleTime: 1000 * 60 * 5, // 5 minutes
        refetchInterval: (query) => {
            const unclaimedCerts = query.state.data?.unclaimedCertificates ?? [];
            // Poll only while at least one cert is queued and its deadline hasn't passed yet.
            const now = new Date();
            const hasActivelyQueued = unclaimedCerts.some(
                (c) =>
                    c.estimatedDeadline && !c.broadcastedAt && new Date(c.estimatedDeadline) > now,
            );
            return hasActivelyQueued ? 15_000 : false;
        },
    });

    return {
        claimedCertificates: data?.claimedCertificates ?? [],
        unclaimedCertificates: data?.unclaimedCertificates ?? [],
        totalClaimed: data?.totalClaimed ?? 0,
        totalUnclaimed: data?.totalUnclaimed ?? 0,
        isLoading,
        isError,
        error: error as Error | null,
    };
};
