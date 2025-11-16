import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

export function useRevokeAllEventCertificates() {
    const {
        mutate: revokeAllEventCertificates,
        isPending: isRevokingAll,
        error: revokeAllError,
    } = useMutation({
        mutationFn: ({ eventId }: { eventId: string }) =>
            coreApiClient.v1.revokeAllEventCertificates({ eventId }),
        onSuccess: () => {
            toast.success("All event certificates revoked successfully");
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.event.all });
        },
        onError: (error) => {
            toast.error("Failed to revoke all event certificates");
            console.error(error);
        },
    });

    return {
        revokeAllEventCertificates,
        isRevokingAll,
        revokeAllError,
    };
}
