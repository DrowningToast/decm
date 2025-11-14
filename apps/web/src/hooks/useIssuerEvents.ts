import { useQuery } from "@tanstack/react-query";
import { issuerService } from "@/services/services";
import type { IssuerEvent } from "@/services/IssuerService/IssuerService";
import { QUERY_KEY } from "@/lib/queryKeys";

interface UseIssuerEventsOptions {
    limit?: number;
    offset?: number;
    enabled?: boolean;
    issuer_credential_id: string;
}

/**
 * Custom hook to fetch issuer events with React Query
 * @param options Configuration options for the query
 * @returns Query result with loading state and data
 */
export const useIssuerEvents = (options: UseIssuerEventsOptions = { issuer_credential_id: "" }) => {
    const { limit = 10, offset = 0, issuer_credential_id } = options;

    return useQuery<IssuerEvent[]>({
        queryKey: QUERY_KEY.issuers.taskedEvents(issuer_credential_id, limit, offset),
        queryFn: () =>
            issuerService.getTaskedEvents(issuer_credential_id, limit || 10, offset || 0),
        enabled: !!issuer_credential_id,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });
};

/**
 * Custom hook to fetch only pending events
 * @param options Configuration options for the query
 * @returns Query result with loading state and data
 */
export const usePendingEvents = (
    options: UseIssuerEventsOptions = { issuer_credential_id: "" },
) => {
    const { limit = 10, offset = 0, issuer_credential_id } = options;

    return useQuery<IssuerEvent[]>({
        queryKey: QUERY_KEY.issuers.pendingEvents(issuer_credential_id, limit, offset),
        queryFn: () =>
            issuerService
                .getTaskedEvents(issuer_credential_id, limit || 10, offset || 0)
                .then((events) => events.filter((event) => !event.isSigned)),
        enabled: !!issuer_credential_id,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });
};

/**
 * Custom hook to fetch only signed events
 * @param options Configuration options for the query
 * @returns Query result with loading state and data
 */
export const useSignedEvents = (options: UseIssuerEventsOptions = { issuer_credential_id: "" }) => {
    const { limit = 10, offset = 0, issuer_credential_id } = options;

    return useQuery<IssuerEvent[]>({
        queryKey: QUERY_KEY.issuers.signedEvents(issuer_credential_id, limit, offset),
        queryFn: () =>
            issuerService
                .getTaskedEvents(issuer_credential_id, limit || 10, offset || 0)
                .then((events) => events.filter((event) => event.isSigned)),
        enabled: !!issuer_credential_id,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });
};
