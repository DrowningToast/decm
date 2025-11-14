import { useQuery } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";
import { issuerService } from "@/services/services";

interface UseSearchIssuerOptions {
    search?: string;
    limit?: number;
    offset?: number;
}

export const useSearchIssuer = (options: UseSearchIssuerOptions) => {
    const { search = "", limit = 25, offset = 0 } = options;

    const query = useQuery({
        queryKey: QUERY_KEY.issuers.search(search, limit, offset),
        queryFn: () => issuerService.searchPublicIssuers(search, limit, offset),
    });

    return query;
};
