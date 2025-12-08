import { useMutation } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import type { EventUpdateEventCertificateTextConfigRequest } from "@decm/api";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";

export const useUpdateCertificateTextConfig = (eventId: string) => {
    const mutation = useMutation({
        mutationFn: async (textConfig: EventUpdateEventCertificateTextConfigRequest) => {
            return await coreApiClient.eventId.updateEventCertificateTextConfig(
                { eventId },
                textConfig,
            );
        },
        onSuccess: () => {
            // Invalidate certificate config query to refetch updated data
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.certificate.config(eventId),
            });
        },
    });

    return {
        updateCertificateTextConfig: mutation.mutateAsync,
        isUpdatingCertificateTextConfig: mutation.isPending,
    };
};
