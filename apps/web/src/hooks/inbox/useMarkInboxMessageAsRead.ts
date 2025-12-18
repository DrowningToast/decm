import { useMutation, useQueryClient } from "@tanstack/react-query";
import { defaultInboxService } from "@/services/InboxService/InboxService";
import { QUERY_KEY } from "@/lib/queryKeys";

/**
 * Marks an inbox message as read and refreshes related inbox queries when the mutation succeeds.
 *
 * @returns The mutation object whose `mutate`/`mutateAsync` accepts a message ID string and, on success, invalidates the specific inbox message query and the inbox list query
 */
export function useMarkInboxMessageAsRead() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (messageId: string) => defaultInboxService.markAsRead(messageId),
        onSuccess: (_, messageId) => {
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
