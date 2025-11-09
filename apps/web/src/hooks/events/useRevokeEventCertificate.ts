import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

export function useRevokeEventCertificate() {
    const {
        mutate: revokeEventCertificate,
        isPending: isRevoking,
        error: revokeError,
    } = useMutation({
        mutationFn: ({ certificateIds, eventId }: { certificateIds: string[]; eventId: string }) =>
            coreApiClient.v1.revokeEventCertificates(
                { eventId },
                {
                    certificate_ids: certificateIds,
                },
            ),
        onSuccess: () => {
            toast.success("Event certificates revoked successfully");
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
        },
        onError: (error) => {
            toast.error("Failed to revoke event certificates");
            console.error(error);
        },
    });

    return {
        revokeEventCertificate,
        isRevoking,
        revokeError,
    };
}
