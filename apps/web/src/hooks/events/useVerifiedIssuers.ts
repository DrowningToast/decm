import { coreApiClient } from "@/lib/api/api";
import { useQuery } from "@tanstack/react-query";

export function useVerifiedIssuers() {
    const {
        data: verifiedIssuers,
        isLoading: isLoadingVerifiedIssuers,
        error: errorVerifiedIssuers,
    } = useQuery({
        queryKey: ["verifiedIssuers"],
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
