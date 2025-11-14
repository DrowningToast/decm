import { authService } from "@/services/services";
import { QUERY_KEY } from "@/lib/queryKeys";
import { useQuery } from "@tanstack/react-query";

export interface UseCheckRolesParams {
    requireHost?: boolean;
    requireIssuer?: boolean;
    enabled?: boolean;
}

export interface UseCheckRolesResult {
    hasRequiredRoles: boolean;
    isHost?: boolean;
    isIssuer?: boolean;
    isLoading: boolean;
    isError: boolean;
    error: Error | null;
}

/**
 * Hook to check if the current user has the required roles.
 * Uses React Query to cache and manage the role check state.
 *
 * @param params - Configuration for role requirements
 * @param params.requireHost - Whether to check for host role
 * @param params.requireIssuer - Whether to check for issuer role
 * @param params.enabled - Whether the query should run (default: true)
 * @returns Object containing role check results and loading state
 *
 * @example
 * ```tsx
 * const { hasRequiredRoles, isLoading } = useCheckRoles({
 *   requireHost: true,
 *   requireIssuer: false
 * });
 *
 * if (isLoading) return <div>Loading...</div>;
 * if (!hasRequiredRoles) return <div>Unauthorized</div>;
 * return <div>Protected Content</div>;
 * ```
 */
export const useCheckRoles = ({
    requireHost = false,
    requireIssuer = false,
    enabled = true,
}: UseCheckRolesParams): UseCheckRolesResult => {
    const { data, isLoading, isError, error } = useQuery({
        queryKey: QUERY_KEY.user.roles(requireHost, requireIssuer),
        queryFn: () =>
            authService.checkRoles({
                requireHost,
                requireIssuer,
            }),
        enabled,
        staleTime: 5 * 60 * 1000, // 5 minutes
        retry: 1,
        retryDelay: 500,
    });

    // Check if user has all required roles
    const hasRequiredHost = !requireHost || data?.is_host === true;
    const hasRequiredIssuer = !requireIssuer || data?.is_issuer === true;
    const hasRequiredRoles = hasRequiredHost && hasRequiredIssuer;

    return {
        hasRequiredRoles,
        isHost: data?.is_host,
        isIssuer: data?.is_issuer,
        isLoading,
        isError,
        error: error as Error | null,
    };
};
