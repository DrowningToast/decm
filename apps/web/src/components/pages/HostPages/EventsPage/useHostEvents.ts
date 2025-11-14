import { useAuth } from "@/context/AuthContext";
import { useQuery } from "@tanstack/react-query";
import { usePaginationState } from "@/hooks/usePaginationState";
import { QUERY_KEY } from "@/lib/queryKeys";
import { eventService } from "@/services/services";

interface UseHostEventsOptions {
    initialPage?: number;
    initialRowsPerPage?: number;
}

export const useHostEvents = (options: UseHostEventsOptions = {}) => {
    const { initialPage = 1, initialRowsPerPage = 10 } = options;
    const { user } = useAuth();

    const { page, rowsPerPage, offset, handlePageChange, handleRowsPerPageChange } =
        usePaginationState(initialPage, initialRowsPerPage);

    const {
        data: events = [],
        isLoading: isLoadingEvents,
        error: isLoadingEventsError,
        refetch,
    } = useQuery({
        queryKey: QUERY_KEY.hostEvents.list(user?.authenticationCredentialId, rowsPerPage, offset),
        queryFn: async () => {
            const result = await eventService.getEventsByOwnerCredentialId(
                user?.authenticationCredentialId ?? "",
            );
            return result;
        },
        enabled: !!user?.authenticationCredentialId,
    });

    // Determine if there are more pages
    const flatEvents = events?.flat() || [];
    const hasMorePages = flatEvents.length === rowsPerPage;
    const hasPreviousPage = page > 1;

    return {
        events,
        isLoadingEvents,
        isLoadingEventsError,
        page,
        rowsPerPage,
        hasPreviousPage,
        hasNextPage: hasMorePages,
        handlePageChange,
        handleRowsPerPageChange,
        refetch,
    };
};
