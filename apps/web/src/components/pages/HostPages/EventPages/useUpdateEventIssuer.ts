import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import type { UpdateEventIssuerPayload } from "@decm/api";
import { useMutation } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";

export function useUpdateEventIssuer(eventID: string) {
    const { mutateAsync: updateEventIssuer, isPending: isUpdatingEventIssuer } = useMutation({
        mutationKey: ["updateEventIssuer"],
        mutationFn: async (data: UpdateEventIssuerPayload) =>
            await coreApiClient.v1.updateEventIssuer(
                {
                    eventId: eventID,
                },
                data,
            ),
        onSuccess: () => {
            // Invalidate event issuers query
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.issuers.byEventId(eventID),
            });
            // Invalidate all event queries to refresh event details page
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
        },
    });

    return {
        updateEventIssuer,
        isUpdatingEventIssuer,
    };
}
