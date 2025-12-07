/**
 * Centralized React Query Keys
 *
 * This file contains all query keys used throughout the application.
 * Each key is defined as a factory function that returns an array of strings.
 *
 * @example
 * ```tsx
 * import { QUERY_KEY } from '@/lib/queryKeys';
 *
 * useQuery({
 *   queryKey: QUERY_KEY.user.profile,
 *   queryFn: fetchUserProfile
 * });
 * ```
 */

export const QUERY_KEY = {
    // User & Authentication
    user: {
        profile: ["user", "profile"] as const,
        roles: (requireHost?: boolean, requireIssuer?: boolean) =>
            ["user", "profile", "roles", requireHost, requireIssuer] as const,
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
        list: (params?: {
            includeActiveEvents?: boolean;
            includeInactiveEvents?: boolean;
            includeClosedEvents?: boolean;
            onlyUserJoinedEvents?: boolean;
        }) => ["event", "list", params] as const,
        byId: (eventId: string) => ["event", eventId] as const,
        registrationConfig: (eventId: string) => ["event-registration-config", eventId] as const,
        contract: (eventId: string) => ["event", eventId, "contract"] as const,
        issuers: {
            byEventId: (eventId: string) => ["event", eventId, "issuers"] as const,
        },
        certificate: {
            config: (eventId: string) => ["event", eventId, "certificate", "config"] as const,
        },
        certificates: (eventId: string) => ["event", eventId, "certificates"] as const,
        invitations: {
            byEventId: (eventId: string) => ["event", eventId, "invitations"] as const,
            ofUserAndEventId: (eventId: string, userId: string) =>
                ["event", eventId, "invitations", "ofUserAndEventId", userId] as const,
        },
        attendees: {
            byEventId: (eventId: string) => ["event", eventId, "attendees"] as const,
        },
        viewmodel: (eventId: string, userId: string | undefined) =>
            ["event", eventId, "viewmodel", userId ?? "me"] as const,
    },

    // Host Events
    hostEvents: {
        all: ["host-events"] as const,
        list: (userId: string | undefined, rowsPerPage: number, offset: number) =>
            ["host-events", userId, rowsPerPage, offset] as const,
    },

    // Issuers
    issuers: {
        search: (search?: string, limit?: number, offset?: number) =>
            ["issuers", "search", search ?? "", limit, offset] as const,
        taskedEvents: (issuerCredentialId: string, limit?: number, offset?: number) =>
            ["issuers", "tasked", issuerCredentialId, limit, offset] as const,
        signedEvents: (issuerCredentialId: string, limit?: number, offset?: number) =>
            ["issuers", "signed", issuerCredentialId, limit, offset] as const,
        pendingEvents: (issuerCredentialId: string, limit?: number, offset?: number) =>
            ["issuers", "pending", issuerCredentialId, limit, offset] as const,
    },

    // Certificates
    certificate: {
        all: ["certificate"] as const,
        list: (limit: number, offset: number, status?: "completed" | "pending") =>
            ["certificate", "list", limit, offset, status] as const,
        byId: (certificateId: string) => ["certificate", certificateId] as const,
        byEventId: (eventId: string) => ["certificate", "event", eventId] as const,
    },

    // Inbox
    inbox: {
        all: ["inbox"] as const,
        list: () => ["inbox", "list"] as const,
        byId: (messageId: string) => ["inbox", messageId] as const,
    },
} as const;
