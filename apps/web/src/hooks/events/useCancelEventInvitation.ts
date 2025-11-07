import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

export function useCancelEventInvitation(eventId: string) {
    const {
        mutate: cancelEventInvitation,
        isPending: isCancelling,
        error: cancelError,
    } = useMutation({
        mutationFn: (eventInvitationId: string) =>
            coreApiClient.v1.cancelEventRegistrationInvitation({
                eventRegistrationInvitationId: eventInvitationId,
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: [QUERY_KEY.event.all] });
            toast.success("Invitation cancelled successfully");
        },
    });

    return {
        cancelEventInvitation,
        isCancelling,
        cancelError,
    };
}
