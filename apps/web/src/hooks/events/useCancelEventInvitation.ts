import { coreApiClient } from "@/lib/api/api";
import { queryClient } from "@/lib/api/queryClient";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

interface CancelEventInvitationParams {
    eventInvitationId: string;
    eventId: string;
}

export function useCancelEventInvitation() {
    const {
        mutate: cancelEventInvitation,
        isPending: isCancelling,
        error: cancelError,
    } = useMutation({
        mutationFn: ({ eventInvitationId }: CancelEventInvitationParams) =>
            coreApiClient.v1.cancelEventRegistrationInvitation({
                event_registration_invitation_id: eventInvitationId,
            }),
        onSuccess: (_, { eventId }) => {
            // Invalidate the specific event invitations query
            queryClient.invalidateQueries({
                queryKey: QUERY_KEY.event.invitations.byEventId(eventId),
            });
            toast.success("Invitation cancelled successfully");
        },
    });

    return {
        cancelEventInvitation,
        isCancelling,
        cancelError,
    };
}
