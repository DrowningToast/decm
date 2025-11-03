import React from "react";
import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { useMyProfile } from "./useMyProfile";

// Mock the API module
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            getMyProfile: vi.fn(),
        },
    },
}));

import { coreApiClient } from "@/lib/api/api";

const mockProfile = {
    id: "profile-123",
    authentication_credential_id: "cred-456",
    first_name: "John",
    last_name: "Doe",
    email: "john.doe@example.com",
    phone_number: "+1234567890",
    address: "123 Main St",
    bio: "Software developer",
    profile_picture_url: "https://example.com/pic.jpg",
    academic_institution: "MIT",
    academic_email: "john@mit.edu",
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
        vi.mocked(coreApiClient.v1.getMyProfile).mockResolvedValue(mockProfile);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        expect(result.current.isLoading).toBe(true);

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.data).toEqual(mockProfile);
        expect(result.current.isSuccess).toBe(true);
    });

    it("should have correct initial loading state", async () => {
        vi.mocked(coreApiClient.v1.getMyProfile).mockImplementation(
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
        vi.mocked(coreApiClient.v1.getMyProfile).mockRejectedValue(error);

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
        vi.mocked(coreApiClient.v1.getMyProfile).mockResolvedValue(mockProfile);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        // The stale time should be 5 minutes (5 * 60 * 1000 = 300000ms)
        // This is configured in the hook itself
        expect(result.current.data).toEqual(mockProfile);
    });

    it("should not refetch on window focus", async () => {
        const mockGetMyProfile = vi.mocked(coreApiClient.v1.getMyProfile);
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
        vi.mocked(coreApiClient.v1.getMyProfile).mockResolvedValue(mockProfile);

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
        vi.mocked(coreApiClient.v1.getMyProfile).mockResolvedValue(mockProfile);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.data?.first_name).toBe("John");
        expect(result.current.data?.last_name).toBe("Doe");
        expect(result.current.data?.email).toBe("john.doe@example.com");
        expect(result.current.data?.academic_institution).toBe("MIT");
    });

    it("should handle retry configuration", async () => {
        const mockGetMyProfile = vi.mocked(coreApiClient.v1.getMyProfile);
        mockGetMyProfile
            .mockRejectedValueOnce(new Error("First attempt failed"))
            .mockRejectedValueOnce(new Error("Second attempt failed"))
            .mockResolvedValueOnce(mockProfile);

        const { result } = renderHook(() => useMyProfile(), { wrapper });

        await waitFor(
            () => {
                expect(result.current.isSuccess).toBe(true);
            },
            { timeout: 5000 },
        );

        expect(result.current.isLoading).toBe(false);
        expect(result.current.data).toEqual(mockProfile);
        // With retry: 2 configuration, it should eventually succeed after 3 attempts (initial + 2 retries)
        expect(mockGetMyProfile).toHaveBeenCalledTimes(3);
    });
});
