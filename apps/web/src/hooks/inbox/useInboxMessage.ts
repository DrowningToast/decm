import { useQuery } from "@tanstack/react-query";
import { InboxService } from "@/services/InboxService";
import { QUERY_KEY } from "@/lib/queryKeys";

interface UseInboxMessageParams {
    messageId: string;
    enabled?: boolean;
}

/**
 * Custom hook to fetch a single inbox message by ID
 * @param params Configuration options for the query
 * @returns Query result with loading state and inbox message data
 */
export function useInboxMessage(params: UseInboxMessageParams) {
    const { messageId, enabled = true } = params;

    const inboxService = new InboxService();

    return useQuery({
        queryKey: QUERY_KEY.inbox.byId(messageId),
        queryFn: () => inboxService.getInboxMessageById(messageId),
        enabled: enabled && !!messageId,
        staleTime: 2 * 60 * 1000, // 2 minutes
    });
}
