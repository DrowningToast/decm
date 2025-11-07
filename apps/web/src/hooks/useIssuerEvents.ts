import { useQuery } from "@tanstack/react-query";
import { getIssuerEvents } from "@/services/issuerService";
import type { GetIssuerEventsData } from "@decm/api";

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

    return useQuery<GetIssuerEventsData>({
        queryKey: ["issuer-events", limit, offset],
        queryFn: () => getIssuerEvents(issuer_credential_id, limit || 10, offset || 0),
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

    return useQuery<GetIssuerEventsData>({
        queryKey: ["issuer-events-pending", limit, offset],
        queryFn: () =>
            getIssuerEvents(issuer_credential_id, limit || 10, offset || 0).then((events) =>
                events.filter((event) => event.is_signed === 0),
            ),
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

    return useQuery<GetIssuerEventsData>({
        queryKey: ["issuer-events-signed", limit, offset],
        queryFn: () =>
            getIssuerEvents(issuer_credential_id, limit || 10, offset || 0).then((events) =>
                events.filter((event) => event.is_signed === 1),
            ),
        enabled: !!issuer_credential_id,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });
};
