/**
 * Centralized React Query Keys
 *
 * This file contains all query keys used throughout the application.
 * Each key is defined as a factory function that returns an array of strings.
 *
 * @example
 * ```tsx
 * import { queryKeys } from '@/lib/queryKeys';
 *
 * useQuery({
 *   queryKey: queryKeys.user.profile,
 *   queryFn: fetchUserProfile
 * });
 * ```
 */

export const queryKeys = {
    // User & Authentication
    user: {
        profile: ["user", "profile"] as const,
    },

    // Onboarding
    onboard: {
        status: {
            all: ["onboardStatus"] as const,
            wallet: (signSignature: string) => ["onboardStatus", signSignature] as const,
            google: (accessToken: string, expiresIn: number) =>
                ["onboardStatus", accessToken, expiresIn] as const,
        },
        signMessage: ["getSignMessage"] as const,
    },

    // Images
    image: {
        byUrl: (url: string) => ["image", url] as const,
    },

    // Events
    event: {
        all: ["event"] as const,
        byId: (eventId: string) => ["event", eventId] as const,
        registrationConfig: (eventId: string) => ["event-registration-config", eventId] as const,
        issuers: {
            byEventId: (eventId: string) => ["event", eventId, "issuers"] as const,
        },
        certificate: {
            config: (eventId: string) => ["event", eventId, "certificate", "config"] as const,
        },
    },

    // Host Events
    hostEvents: {
        all: ["host-events"] as const,
        list: (userId: string | undefined, rowsPerPage: number, offset: number) =>
            ["host-events", userId, rowsPerPage, offset] as const,
    },

    // Issuers
    issuers: {
        verified: ["issuers"] as const,
    },
} as const;
