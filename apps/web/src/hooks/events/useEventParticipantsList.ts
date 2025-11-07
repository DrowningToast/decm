import { useQuery } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";

interface UseEventParticipantsListParams {
    eventId: string;
    currentPage: number;
    pageSize: number;
    searchQuery?: string;
}

export const useEventParticipantsList = ({
    eventId,
    currentPage,
    pageSize,
    searchQuery,
}: UseEventParticipantsListParams) => {
    const {
        data: participantData,
        isLoading: isParticipantsLoading,
        error: participantsError,
    } = useQuery({
        queryKey: ["eventParticipants", eventId, currentPage, pageSize, searchQuery],
        queryFn: async () => {
            if (!eventId) return [];

            const data = await coreApiClient.v1.getEventRegistrationInvitationsByEventId({
                eventId,
            });
            return data;
        },
        enabled: !!eventId,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });

    // Simple client-side pagination and filtering for now
    let filteredData = participantData || [];

    if (searchQuery) {
        const searchLower = searchQuery.toLowerCase();
        filteredData = filteredData.filter(
            (participant) =>
                participant.first_name?.toLowerCase().includes(searchLower) ||
                participant.last_name?.toLowerCase().includes(searchLower) ||
                participant.email?.toLowerCase().includes(searchLower),
        );
    }

    const startIndex = (currentPage - 1) * pageSize;
    const paginatedData = filteredData.slice(startIndex, startIndex + pageSize);

    return {
        participants: paginatedData,
        totalItems: filteredData.length,
        isLoading: isParticipantsLoading,
        error: participantsError,
    };
};
