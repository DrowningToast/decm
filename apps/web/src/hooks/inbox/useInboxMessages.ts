import { useQuery } from "@tanstack/react-query";
import { InboxService } from "@/services/InboxService";
import { QUERY_KEY } from "@/lib/queryKeys";

interface UseInboxMessagesParams {
    limit?: number;
    offset?: number;
    enabled?: boolean;
}

/**
 * Custom hook to fetch inbox messages for the authenticated user
 * @param params Configuration options for the query
 * @returns Query result with loading state and inbox messages data
 */
export function useInboxMessages(params: UseInboxMessagesParams = {}) {
    const { limit = 10, offset = 0, enabled = true } = params;

    return useQuery({
        queryKey: QUERY_KEY.inbox.list(limit, offset),
        queryFn: () =>
            inboxService.getInboxMessages({
                limit,
                offset,
            }),
        enabled,
        staleTime: 2 * 60 * 1000, // 2 minutes
    });
}
