import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { useUpdateEventIssuer } from "./useUpdateEventIssuer";
import { coreApiClient } from "@/lib/api/api";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import type { UpdateEventIssuerPayload, EventEventIssuerResponse } from "@decm/api";

// Mock the API client
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            updateEventIssuer: vi.fn(),
        },
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

describe("useUpdateEventIssuer", () => {
    const mockEventId = "event-123";
    const mockPayload: UpdateEventIssuerPayload = [
        {
            event_id: mockEventId,
            issuer_credential_id: "cred-1",
        },
        {
            event_id: mockEventId,
            issuer_credential_id: "cred-2",
        },
    ];

    const mockResponse: EventEventIssuerResponse = {
        id: "issuer-1",
        event_id: mockEventId,
        issuer_credential_id: "cred-1",
        is_signed: 0,
        sign_message: "",
        signature: "",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z",
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    it("should update event issuers successfully", async () => {
        vi.mocked(coreApiClient.v1.updateEventIssuer).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useUpdateEventIssuer(mockEventId), {
            wrapper: createWrapper(),
        });

        let response;
        await waitFor(async () => {
            response = await result.current.updateEventIssuer(mockPayload);
        });

        expect(response).toEqual(mockResponse);
        expect(coreApiClient.v1.updateEventIssuer).toHaveBeenCalledWith(
            { eventId: mockEventId },
            mockPayload,
        );
    });

    it("should set isUpdatingEventIssuer to true during mutation", async () => {
        vi.mocked(coreApiClient.v1.updateEventIssuer).mockImplementation(
            () => new Promise((resolve) => setTimeout(() => resolve(mockResponse), 100)),
        );

        const { result } = renderHook(() => useUpdateEventIssuer(mockEventId), {
            wrapper: createWrapper(),
        });

        expect(result.current.isUpdatingEventIssuer).toBe(false);

        result.current.updateEventIssuer(mockPayload);

        await waitFor(() => {
            expect(result.current.isUpdatingEventIssuer).toBe(true);
        });

        await waitFor(() => {
            expect(result.current.isUpdatingEventIssuer).toBe(false);
        });
    });

    it("should handle API errors", async () => {
        const error = new Error("Failed to update event issuer");
        vi.mocked(coreApiClient.v1.updateEventIssuer).mockRejectedValue(error);

        const { result } = renderHook(() => useUpdateEventIssuer(mockEventId), {
            wrapper: createWrapper(),
        });

        await expect(result.current.updateEventIssuer(mockPayload)).rejects.toThrow(
            "Failed to update event issuer",
        );
    });

    it("should handle single issuer update", async () => {
        const singleIssuerPayload: UpdateEventIssuerPayload = [
            {
                event_id: mockEventId,
                issuer_credential_id: "cred-1",
            },
        ];

        vi.mocked(coreApiClient.v1.updateEventIssuer).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useUpdateEventIssuer(mockEventId), {
            wrapper: createWrapper(),
        });

        await result.current.updateEventIssuer(singleIssuerPayload);

        expect(coreApiClient.v1.updateEventIssuer).toHaveBeenCalledWith(
            { eventId: mockEventId },
            singleIssuerPayload,
        );
    });

    it("should handle empty issuer array", async () => {
        const emptyPayload: UpdateEventIssuerPayload = [];

        vi.mocked(coreApiClient.v1.updateEventIssuer).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useUpdateEventIssuer(mockEventId), {
            wrapper: createWrapper(),
        });

        await result.current.updateEventIssuer(emptyPayload);

        expect(coreApiClient.v1.updateEventIssuer).toHaveBeenCalledWith(
            { eventId: mockEventId },
            emptyPayload,
        );
    });

    it("should handle multiple issuers", async () => {
        const multipleIssuersPayload: UpdateEventIssuerPayload = [
            { event_id: mockEventId, issuer_credential_id: "cred-1" },
            { event_id: mockEventId, issuer_credential_id: "cred-2" },
            { event_id: mockEventId, issuer_credential_id: "cred-3" },
            { event_id: mockEventId, issuer_credential_id: "cred-4" },
        ];

        vi.mocked(coreApiClient.v1.updateEventIssuer).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useUpdateEventIssuer(mockEventId), {
            wrapper: createWrapper(),
        });

        await result.current.updateEventIssuer(multipleIssuersPayload);

        expect(coreApiClient.v1.updateEventIssuer).toHaveBeenCalledWith(
            { eventId: mockEventId },
            multipleIssuersPayload,
        );
    });

    it("should handle multiple consecutive updates", async () => {
        vi.mocked(coreApiClient.v1.updateEventIssuer).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useUpdateEventIssuer(mockEventId), {
            wrapper: createWrapper(),
        });

        await result.current.updateEventIssuer(mockPayload);
        await result.current.updateEventIssuer([
            {
                event_id: mockEventId,
                issuer_credential_id: "cred-3",
            },
        ]);

        expect(coreApiClient.v1.updateEventIssuer).toHaveBeenCalledTimes(2);
    });
});
