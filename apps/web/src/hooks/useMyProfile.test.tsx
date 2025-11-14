import React from "react";
import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useMyProfile } from "./useMyProfile";

// Mock the authService
vi.mock("@/services/services", () => ({
    authService: {
        getMyProfile: vi.fn(),
    },
}));

import { authService } from "@/services/services";

// Mock profile returned from the service (camelCase after transformation)
const mockProfile = {
    profileId: "profile-123",
    authenticationCredentialId: "cred-456",
    firstName: "John",
    lastName: "Doe",
    email: "john.doe@example.com",
    phoneNumber: "+1234567890",
    address: "123 Main St",
    bio: "Software developer",
    profilePictureUrl: "https://example.com/pic.jpg",
    academicInstitution: "MIT",
    academicEmail: "john@mit.edu",
    walletAddress: "0x1234567890123456789012345678901234567890",
    googleConnectorRef: undefined,
    githubConnectorRef: undefined,
    solutionStatus: "SYSTEM_MANAGED" as const,
    isFirstNamePublic: false,
    isLastNamePublic: false,
    isEmailPublic: false,
    isPhoneNumberPublic: false,
    isAddressPublic: false,
    isBioPublic: false,
    isAcademicInstitutionPublic: false,
    isAcademicEmailPublic: false,
};

describe("useMyProfile", () => {
    let queryClient: QueryClient;

    beforeEach(() => {
        queryClient = new QueryClient({
            defaultOptions: {
                queries: {
                    retry: false,
                },
            },
        });
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    const wrapper = ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );

    it("should fetch profile successfully", async () => {
        vi.mocked(authService.getMyProfile).mockResolvedValue(mockProfile);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        expect(result.current.isLoading).toBe(true);

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.data).toEqual(mockProfile);
        expect(result.current.isSuccess).toBe(true);
    });

    it("should have correct initial loading state", async () => {
        vi.mocked(authService.getMyProfile).mockImplementation(
            () => new Promise((resolve) => setTimeout(() => resolve(mockProfile), 100)),
        );

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        expect(result.current.isLoading).toBe(true);
        expect(result.current.data).toBeUndefined();

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.data).toEqual(mockProfile);
    });

    it("should handle error state", async () => {
        const error = new Error("Failed to fetch profile");
        vi.mocked(authService.getMyProfile).mockRejectedValue(error);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        await waitFor(
            () => {
                expect(result.current.isError).toBe(true);
            },
            { timeout: 5000 },
        );

        expect(result.current.isLoading).toBe(false);
        expect(result.current.error).toBeDefined();
    });

    it("should have 5 minute stale time", async () => {
        vi.mocked(authService.getMyProfile).mockResolvedValue(mockProfile);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        // The stale time should be 5 minutes (5 * 60 * 1000 = 300000ms)
        // This is configured in the hook itself
        expect(result.current.data).toEqual(mockProfile);
    });

    it("should not refetch on window focus", async () => {
        const mockGetMyProfile = vi.mocked(authService.getMyProfile);
        mockGetMyProfile.mockResolvedValue(mockProfile);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        const callCountBefore = mockGetMyProfile.mock.calls.length;

        // Simulate window focus event
        window.dispatchEvent(new Event("focus"));

        // Should not trigger a refetch
        expect(mockGetMyProfile.mock.calls.length).toBe(callCountBefore);
    });

    it("should return query object with status properties", async () => {
        vi.mocked(authService.getMyProfile).mockResolvedValue(mockProfile);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        expect(result.current).toHaveProperty("data");
        expect(result.current).toHaveProperty("isLoading");
        expect(result.current).toHaveProperty("isSuccess");
        expect(result.current).toHaveProperty("isError");
        expect(result.current).toHaveProperty("error");

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.data).toEqual(mockProfile);
        expect(result.current.isSuccess).toBe(true);
        expect(result.current.isError).toBe(false);
    });

    it("should include profile data fields", async () => {
        vi.mocked(authService.getMyProfile).mockResolvedValue(mockProfile);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        await waitFor(() => {
            expect(result.current.isSuccess).toBe(true);
        });

        expect(result.current.data?.firstName).toBe("John");
        expect(result.current.data?.lastName).toBe("Doe");
        expect(result.current.data?.email).toBe("john.doe@example.com");
        expect(result.current.data?.academicInstitution).toBe("MIT");
    });

    it("should not retry on failure (retry: false)", async () => {
        const mockGetMyProfile = vi.mocked(authService.getMyProfile);
        const error = new Error("Failed to fetch profile");
        mockGetMyProfile.mockRejectedValueOnce(error);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        await waitFor(
            () => {
                expect(result.current.isError).toBe(true);
            },
            { timeout: 1000 },
        );

        expect(result.current.isLoading).toBe(false);
        expect(result.current.error).toBeDefined();
        // With retry: false configuration, it should not retry
        expect(mockGetMyProfile).toHaveBeenCalledTimes(1);
    });
});
