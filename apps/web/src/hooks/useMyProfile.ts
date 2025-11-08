import { coreApiClient } from "@/lib/api/api";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useQuery } from "@tanstack/react-query";

export const useMyProfile = () => {
    const query = useQuery({
        queryKey: QUERY_KEY.user.profile,
        queryFn: () => coreApiClient.v1.getMyProfile(),

        retry: false,
        refetchOnWindowFocus: false,
        staleTime: 5 * 60 * 1000,
    });

    return query;
};
