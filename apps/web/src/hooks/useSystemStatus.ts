import { useQuery } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";
import { blockchainService } from "@/services/BlockchainService/BlockchainService";

/**
 * Hook to fetch the combined system status viewmodel including gas price and latest maintenance schedule
 */
export const useSystemStatus = () => {
    return useQuery({
        queryKey: QUERY_KEY.blockchain.systemStatus,
        queryFn: () => blockchainService.getSystemStatusViewModel(),
        refetchInterval: 60 * 1000, // Refetch every 1 minute
        staleTime: 60 * 1000, // Consider data fresh for 1 minute
        retry: 2,
    });
};
