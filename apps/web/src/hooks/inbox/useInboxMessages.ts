import { useQuery } from "@tanstack/react-query";
import { defaultInboxService, type InboxMessage } from "@/services/InboxService/InboxService";
import { QUERY_KEY } from "@/lib/queryKeys";

interface UseInboxMessagesParams {
    enabled?: boolean;
}

interface UseInboxMessagesResult {
    data: InboxMessage[];
    isLoading: boolean;
    error: Error | null;
    refetch: () => void;
}

/**
 * Custom hook to fetch inbox messages for the authenticated user
 * @param params Configuration options for the query
 * @returns Query result with loading state and inbox messages data
 */
export function useInboxMessages(params: UseInboxMessagesParams = {}): UseInboxMessagesResult {
    const { enabled = true } = params;

    const query = useQuery({
        queryKey: QUERY_KEY.inbox.list(),
        queryFn: () => defaultInboxService.getInboxMessages(),
        enabled,
        staleTime: 2 * 60 * 1000, // 2 minutes
    });

    return {
        data: query.data ?? [],
        isLoading: query.isLoading,
        error: query.error,
        refetch: query.refetch,
    };
}
