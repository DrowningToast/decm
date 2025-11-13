import { eventRegistrationService } from "@/services/services";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEY } from "@/lib/queryKeys";

export const useEventInvitationByUserAndEvent = (eventId: string, userId: string | undefined) => {
    const {
        data: invitation,
        isLoading,
        error,
    } = useQuery({
        queryKey: [QUERY_KEY.event.invitations.ofUserAndEventId(eventId, userId ?? "")],
        queryFn: () => eventRegistrationService.getInvitationOfUserAndEventId(eventId),
        enabled: !!eventId && !!userId,
    });

    return {
        invitation,
        isLoading,
        error,
    };
};
