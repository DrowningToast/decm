import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { useMutation } from "@tanstack/react-query";
import { queryKeys } from "@/lib/queryKeys";

export function useDeleteEventIssuer() {
    const { mutateAsync: deleteEventIssuerAsync, isPending: isDeletingEventIssuer } = useMutation({
        mutationKey: ["deleteEventIssuer"],
        mutationFn: async ({ eventId, issuerId }: { eventId: string; issuerId: string }) =>
            await coreApiClient.v1.deleteEventIssuer({
                eventId,
                issuerId,
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.event.all });
        },
    });

    return {
        deleteEventIssuerAsync,
        isDeletingEventIssuer,
    };
}
