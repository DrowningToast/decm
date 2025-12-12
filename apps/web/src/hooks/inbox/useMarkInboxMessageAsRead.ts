import { useMutation, useQueryClient } from "@tanstack/react-query";
import { defaultInboxService } from "@/services/InboxService/InboxService";
import { QUERY_KEY } from "@/lib/queryKeys";

/**
 * Hook to mark an inbox message as read
 * @returns Mutation function to mark message as read
 */
export function useMarkInboxMessageAsRead() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (messageId: string) => defaultInboxService.markAsRead(messageId),
        onSuccess: (data, messageId) => {
            // Invalidate the specific message query to refetch with updated read status
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.inbox.byId(messageId),
            });
            // Also invalidate the inbox list to update the read status in the list
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.inbox.list(),
            });
        },
    });
}
