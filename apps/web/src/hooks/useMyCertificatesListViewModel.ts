import { useQuery } from "@tanstack/react-query";
import { certificateService } from "@/services/services";
import type { MyCertificatesViewModel } from "@/services/CertificateService/mapper";

interface UseMyCertificatesListViewModelReturn extends MyCertificatesViewModel {
    isLoading: boolean;
    isError: boolean;
    error: Error | null;
}

export const useMyCertificatesListViewModel = (): UseMyCertificatesListViewModelReturn => {
    const { data, isLoading, isError, error } = useQuery({
        queryKey: ["my-certificates-list-viewmodel"],
        queryFn: () => certificateService.getMyCertificatesList(),
        staleTime: 1000 * 60 * 5, // 5 minutes
    });

    return {
        claimedCertificates: data?.claimedCertificates || [],
        unclaimedCertificates: data?.unclaimedCertificates || [],
        totalClaimed: data?.totalClaimed || 0,
        totalUnclaimed: data?.totalUnclaimed || 0,
        isLoading,
        isError,
        error: error as Error | null,
    };
};
