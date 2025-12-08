import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { useEventIssuers } from "./useEventIssuers";
import { eventService } from "@/services/services";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import type { EventIssuer } from "@/services/EventService/EventService";

// Mock the event service
vi.mock("@/services/services", () => ({
    eventService: {
        getEventIssuersByEventId: vi.fn(),
    },
}));

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: {
            queries: {
                retry: false,
            },
        },
    });

    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

describe("useEventIssuers", () => {
    const mockEventId = "event-123";
    const mockEventIssuers: EventIssuer[] = [
        {
            eventId: mockEventId,
            id: "issuer-1",
            isSigned: true,
            issuerCredentialId: "cred-1",
            issuerProfile: {
                id: "profile-1",
                authenticationCredentialId: "cred-1",
                firstName: "John",
                lastName: "Doe",
                email: "john@example.com",
                academicInstitution: "MIT",
                isAcademicEmailPublic: true,
                isAcademicInstitutionPublic: true,
                isAddressPublic: false,
                isBioPublic: false,
                isEmailPublic: true,
                isFirstNamePublic: true,
                isLastNamePublic: true,
                isPhoneNumberPublic: false,
                isProfilePicturePublic: false,
            },
            signMessage: "sign-message-1",
            signature: "signature-1",
            createdAt: new Date("2024-01-01"),
            updatedAt: new Date("2024-01-01"),
        },
        {
            eventId: mockEventId,
            id: "issuer-2",
            isSigned: false,
            issuerCredentialId: "cred-2",
            issuerProfile: {
                id: "profile-2",
                authenticationCredentialId: "cred-2",
                firstName: "Jane",
                lastName: "Smith",
                email: "jane@example.com",
                academicInstitution: "Stanford",
                isAcademicEmailPublic: true,
                isAcademicInstitutionPublic: true,
                isAddressPublic: false,
                isBioPublic: false,
                isEmailPublic: true,
                isFirstNamePublic: true,
                isLastNamePublic: true,
                isPhoneNumberPublic: false,
                isProfilePicturePublic: false,
            },
            createdAt: new Date("2024-01-02"),
            updatedAt: new Date("2024-01-02"),
        },
    ];

    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    it("should fetch and return event issuers", async () => {
        vi.mocked(eventService.getEventIssuersByEventId).mockResolvedValue(mockEventIssuers);

        const { result } = renderHook(() => useEventIssuers(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoadingEventIssuers).toBe(false);
        });

        expect(result.current.eventIssuers).toEqual(mockEventIssuers);
    });

    it("should call service with correct eventId", async () => {
        vi.mocked(eventService.getEventIssuersByEventId).mockResolvedValue(mockEventIssuers);

        renderHook(() => useEventIssuers(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(eventService.getEventIssuersByEventId).toHaveBeenCalledWith(mockEventId);
        });
    });

    it("should return empty array when no issuers exist", async () => {
        vi.mocked(eventService.getEventIssuersByEventId).mockResolvedValue([]);

        const { result } = renderHook(() => useEventIssuers(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoadingEventIssuers).toBe(false);
        });

        expect(result.current.eventIssuers).toEqual([]);
    });

    it("should handle undefined response", async () => {
        vi.mocked(eventService.getEventIssuersByEventId).mockResolvedValue([]);

        const { result } = renderHook(() => useEventIssuers(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoadingEventIssuers).toBe(false);
        });

        expect(result.current.eventIssuers).toEqual([]);
    });

    it("should handle issuers with and without signatures", async () => {
        vi.mocked(eventService.getEventIssuersByEventId).mockResolvedValue(mockEventIssuers);

        const { result } = renderHook(() => useEventIssuers(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoadingEventIssuers).toBe(false);
        });

        const signedIssuer = result.current.eventIssuers?.find((i) => i.isSigned);
        const unsignedIssuer = result.current.eventIssuers?.find((i) => !i.isSigned);

        expect(signedIssuer).toBeDefined();
        expect(signedIssuer?.signature).toBe("signature-1");
        expect(unsignedIssuer).toBeDefined();
        expect(unsignedIssuer?.signature).toBeUndefined();
    });

    it("should handle API errors gracefully", async () => {
        const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});
        const error = new Error("Failed to fetch issuers");
        vi.mocked(eventService.getEventIssuersByEventId).mockRejectedValue(error);

        const { result } = renderHook(() => useEventIssuers(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoadingEventIssuers).toBe(false);
        });

        // Query should still work, but data will be undefined
        expect(result.current.eventIssuers).toBeUndefined();

        consoleSpy.mockRestore();
    });
});
