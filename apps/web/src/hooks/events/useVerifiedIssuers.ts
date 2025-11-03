import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";

export function useVerifiedIssuers() {
    const {
        data: verifiedIssuers,
        isLoading: isLoadingVerifiedIssuers,
        error: errorVerifiedIssuers,
    } = useQuery({
        queryKey: queryKeys.issuers.verified,
        queryFn: () =>
            coreApiClient.v1.getVerifiedIssuers({
                limit: 25,
                offset: 0,
            }),
    });

    return {
        verifiedIssuers,
        isLoadingVerifiedIssuers,
        errorVerifiedIssuers,
    };
}
