import { useAuth } from "@/context/AuthContext";
import { useQuery } from "@tanstack/react-query";
import { coreApiClient } from "@/lib/api/api";
import { usePaginationState } from "@/hooks/usePaginationState";
import { queryKeys } from "@/lib/queryKeys";

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
        data: events,
        isLoading: isLoadingEvents,
        error: isLoadingEventsError,
        refetch,
    } = useQuery({
        queryKey: queryKeys.hostEvents.list(
            user?.authentication_credential_id,
            rowsPerPage,
            offset,
        ),
        queryFn: async () => {
            const result = await coreApiClient.v1.getEventsByOwnerCredentialsId({
                ownerCredentialId: user?.authentication_credential_id ?? "",
                limit: rowsPerPage,
                offset: offset,
            });

            return result;
        },
        enabled: !!user?.authentication_credential_id,
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
