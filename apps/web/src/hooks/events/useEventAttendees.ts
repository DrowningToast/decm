import { useQuery } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import { QUERY_KEY } from "@/lib/queryKeys";

interface UseEventAttendeesParams {
    eventId: string;
}

export const useEventAttendees = ({ eventId }: UseEventAttendeesParams) => {
    const {
        data: attendees,
        isLoading,
        error,
    } = useQuery({
        queryKey: QUERY_KEY.event.attendees.byEventId(eventId),
        queryFn: async () => {
            if (!eventId) return [];

            const data = await coreApiClient.v1.getEventParticipants({
                eventId,
            });
            return data;
        },
        enabled: !!eventId,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });

    return {
        attendees: attendees || [],
        isLoading,
        error,
    };
};
