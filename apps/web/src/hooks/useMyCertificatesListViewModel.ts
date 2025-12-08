import { useQuery } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import type { EventGetMyCertificatesListViewModelResponse } from "@decm/api";

interface UseMyCertificatesListViewModelReturn {
    data: EventGetMyCertificatesListViewModelResponse | undefined;
    claimedCertificates: EventGetMyCertificatesListViewModelResponse["claimed_certificates"];
    unclaimedCertificates: EventGetMyCertificatesListViewModelResponse["unclaimed_certificates"];
    totalClaimed: number;
    totalUnclaimed: number;
    isLoading: boolean;
    isError: boolean;
    error: Error | null;
}

export const useMyCertificatesListViewModel = (): UseMyCertificatesListViewModelReturn => {
    const { data, isLoading, isError, error } = useQuery({
        queryKey: ["my-certificates-list-viewmodel"],
        queryFn: async () => {
            const response = await coreApiClient.v1.getMyCertificatesListViewmodel();
            return response;
        },
        staleTime: 1000 * 60 * 5, // 5 minutes
    });

    return {
        data,
        claimedCertificates: data?.claimed_certificates || [],
        unclaimedCertificates: data?.unclaimed_certificates || [],
        totalClaimed: data?.total_claimed || 0,
        totalUnclaimed: data?.total_unclaimed || 0,
        isLoading,
        isError,
        error: error as Error | null,
    };
};
