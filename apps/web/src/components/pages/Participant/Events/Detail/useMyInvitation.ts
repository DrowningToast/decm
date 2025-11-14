import { useAuth } from "@/context/AuthContext";
import { QUERY_KEY } from "@/lib/queryKeys";
import { eventRegistrationService } from "@/services/services";
import { useQuery } from "@tanstack/react-query";

export const useMyInvitation = (eventId: string) => {
    const { user } = useAuth();
    const query = useQuery({
        queryKey: [QUERY_KEY.event.invitations.byEventId(eventId)],
        queryFn: () => eventRegistrationService.getInvitationOfUserAndEventId(eventId),
        enabled: !!user && !!eventId,
    });

    return query;
};
