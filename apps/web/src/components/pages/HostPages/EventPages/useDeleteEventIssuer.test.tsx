import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { useDeleteEventIssuer } from "./useDeleteEventIssuer";
import { coreApiClient } from "@/lib/api/api";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";

// Mock the API client
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            deleteEventIssuer: vi.fn(),
        },
    },
}));

// Mock the query client
vi.mock("@/lib/api/queryClient", () => ({
    queryClient: {
        invalidateQueries: vi.fn(),
    },
}));

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: {
            mutations: {
                retry: false,
            },
        },
    });

    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

describe("useDeleteEventIssuer", () => {
    const mockEventId = "event-123";
    const mockIssuerId = "issuer-1";

    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    it("should delete event issuer successfully", async () => {
        const mockResponse = {};
        vi.mocked(coreApiClient.v1.deleteEventIssuer).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useDeleteEventIssuer(), {
            wrapper: createWrapper(),
        });

        let response;
        await waitFor(async () => {
            response = await result.current.deleteEventIssuerAsync({
                eventId: mockEventId,
                issuerId: mockIssuerId,
            });
        });

        expect(response).toEqual(mockResponse);
        expect(coreApiClient.v1.deleteEventIssuer).toHaveBeenCalledWith({
            eventId: mockEventId,
            issuerId: mockIssuerId,
        });
    });

    it("should set isDeletingEventIssuer to true during mutation", async () => {
        vi.mocked(coreApiClient.v1.deleteEventIssuer).mockImplementation(
            () => new Promise((resolve) => setTimeout(() => resolve({}), 100)),
        );

        const { result } = renderHook(() => useDeleteEventIssuer(), {
            wrapper: createWrapper(),
        });

        expect(result.current.isDeletingEventIssuer).toBe(false);

        result.current.deleteEventIssuerAsync({
            eventId: mockEventId,
            issuerId: mockIssuerId,
        });

        await waitFor(() => {
            expect(result.current.isDeletingEventIssuer).toBe(true);
        });

        await waitFor(() => {
            expect(result.current.isDeletingEventIssuer).toBe(false);
        });
    });

    it("should handle API errors", async () => {
        const error = new Error("Failed to delete event issuer");
        vi.mocked(coreApiClient.v1.deleteEventIssuer).mockRejectedValue(error);

        const { result } = renderHook(() => useDeleteEventIssuer(), {
            wrapper: createWrapper(),
        });

        await expect(
            result.current.deleteEventIssuerAsync({
                eventId: mockEventId,
                issuerId: mockIssuerId,
            }),
        ).rejects.toThrow("Failed to delete event issuer");
    });

    it("should handle multiple delete operations", async () => {
        vi.mocked(coreApiClient.v1.deleteEventIssuer).mockResolvedValue({});

        const { result } = renderHook(() => useDeleteEventIssuer(), {
            wrapper: createWrapper(),
        });

        await result.current.deleteEventIssuerAsync({
            eventId: mockEventId,
            issuerId: "issuer-1",
        });

        await result.current.deleteEventIssuerAsync({
            eventId: mockEventId,
            issuerId: "issuer-2",
        });

        expect(coreApiClient.v1.deleteEventIssuer).toHaveBeenCalledTimes(2);
        expect(coreApiClient.v1.deleteEventIssuer).toHaveBeenNthCalledWith(1, {
            eventId: mockEventId,
            issuerId: "issuer-1",
        });
        expect(coreApiClient.v1.deleteEventIssuer).toHaveBeenNthCalledWith(2, {
            eventId: mockEventId,
            issuerId: "issuer-2",
        });
    });

    it("should handle different event IDs", async () => {
        vi.mocked(coreApiClient.v1.deleteEventIssuer).mockResolvedValue({});

        const { result } = renderHook(() => useDeleteEventIssuer(), {
            wrapper: createWrapper(),
        });

        await result.current.deleteEventIssuerAsync({
            eventId: "event-1",
            issuerId: mockIssuerId,
        });

        await result.current.deleteEventIssuerAsync({
            eventId: "event-2",
            issuerId: mockIssuerId,
        });

        expect(coreApiClient.v1.deleteEventIssuer).toHaveBeenCalledTimes(2);
        expect(coreApiClient.v1.deleteEventIssuer).toHaveBeenNthCalledWith(1, {
            eventId: "event-1",
            issuerId: mockIssuerId,
        });
        expect(coreApiClient.v1.deleteEventIssuer).toHaveBeenNthCalledWith(2, {
            eventId: "event-2",
            issuerId: mockIssuerId,
        });
    });

    it("should handle network errors gracefully", async () => {
        const networkError = new Error("Network error");
        vi.mocked(coreApiClient.v1.deleteEventIssuer).mockRejectedValue(networkError);

        const { result } = renderHook(() => useDeleteEventIssuer(), {
            wrapper: createWrapper(),
        });

        await expect(
            result.current.deleteEventIssuerAsync({
                eventId: mockEventId,
                issuerId: mockIssuerId,
            }),
        ).rejects.toThrow("Network error");

        expect(result.current.isDeletingEventIssuer).toBe(false);
    });
});
