import { coreApiClient } from "@/lib/api/api";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useQuery } from "@tanstack/react-query";

export function useEventInvitedParticipants(eventId: string) {
    const {
        data: invitations,
        isLoading,
        error,
    } = useQuery({
        queryKey: QUERY_KEY.event.invitations.byEventId(eventId),
        queryFn: () => coreApiClient.v1.getEventRegistrationInvitationsByEventId({ eventId }),
        enabled: !!eventId,
    });

    return {
        invitations,
        isLoading,
        error,
    };
}
