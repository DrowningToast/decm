import { useQuery } from "@tanstack/react-query";
import { defaultInboxService, type InboxMessageDetail } from "@/services/InboxService/InboxService";
import { QUERY_KEY } from "@/lib/queryKeys";

interface UseInboxMessageParams {
    messageId: string;
    enabled?: boolean;
}

interface UseInboxMessageResult {
    data: InboxMessageDetail | null;
    isLoading: boolean;
    error: Error | null;
    refetch: () => void;
}

/**
 * Custom hook to fetch a single inbox message by ID
 * @param params Configuration options for the query
 * @returns Query result with loading state and inbox message data
 */
export function useInboxMessage(params: UseInboxMessageParams): UseInboxMessageResult {
    const { messageId, enabled = true } = params;

    const query = useQuery({
        queryKey: QUERY_KEY.inbox.byId(messageId),
        queryFn: () => defaultInboxService.getInboxMessageById(messageId),
        enabled: enabled && !!messageId,
        staleTime: 2 * 60 * 1000, // 2 minutes
    });

    return {
        data: query.data ?? null,
        isLoading: query.isLoading,
        error: query.error,
        refetch: query.refetch,
    };
}
