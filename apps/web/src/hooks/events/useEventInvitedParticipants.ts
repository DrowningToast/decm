import { QUERY_KEY } from "@/lib/queryKeys";
import { eventRegistrationService } from "@/services/services";
import { useQuery } from "@tanstack/react-query";

export function useEventInvitedParticipants(eventId: string) {
    const {
        data: invitations,
        isLoading,
        error,
    } = useQuery({
        queryKey: QUERY_KEY.event.invitations.byEventId(eventId),
        queryFn: () => eventRegistrationService.getInvitationByEventId(eventId),
        enabled: !!eventId,
    });
    console.log(invitations);
    return {
        invitations,
        isLoading,
        error,
    };
}
