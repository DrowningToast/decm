import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { useEventCertificates } from "./useEventCertificates";
import { coreApiClient } from "@/lib/api/api";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import type { EntityEventCertificate } from "@decm/api";

// Mock the API client
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            getEventCertificates: vi.fn(),
        },
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

describe("useEventCertificates", () => {
    const mockEventId = "event-123";
    const mockCertificates: EntityEventCertificate[] = [
        {
            id: "cert-1",
            event_id: mockEventId,
            event_contract_address: "0x1234567890123456789012345678901234567890",
            receiver_email: "user1@example.com",
            name: "John Doe",
            certificate_title: "Certificate of Completion",
            created_at: "2024-01-01T00:00:00Z",
        },
        {
            id: "cert-2",
            event_id: mockEventId,
            event_contract_address: "0x1234567890123456789012345678901234567890",
            receiver_email: "user2@example.com",
            name: "Jane Doe",
            certificate_title: "Certificate of Achievement",
            created_at: "2024-01-02T00:00:00Z",
        },
    ];

    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    it("should fetch and return event certificates", async () => {
        vi.mocked(coreApiClient.v1.getEventCertificates).mockResolvedValue({
            certificates: mockCertificates,
        });

        const { result } = renderHook(() => useEventCertificates(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificates).toEqual(mockCertificates);
        expect(result.current.error).toBeNull();
    });

    it("should return empty array when no certificates exist", async () => {
        vi.mocked(coreApiClient.v1.getEventCertificates).mockResolvedValue({
            certificates: [],
        });

        const { result } = renderHook(() => useEventCertificates(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificates).toEqual([]);
    });

    it("should handle when response has empty certificates array", async () => {
        vi.mocked(coreApiClient.v1.getEventCertificates).mockResolvedValue({
            certificates: [],
        });

        const { result } = renderHook(() => useEventCertificates(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificates).toEqual([]);
    });

    it("should not fetch when eventId is empty", async () => {
        const { result } = renderHook(() => useEventCertificates(""), { wrapper: createWrapper() });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(coreApiClient.v1.getEventCertificates).not.toHaveBeenCalled();
        expect(result.current.certificates).toEqual([]);
    });

    it("should handle API errors gracefully", async () => {
        const error = new Error("Failed to fetch certificates");
        vi.mocked(coreApiClient.v1.getEventCertificates).mockRejectedValue(error);

        const { result } = renderHook(() => useEventCertificates(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.error).toBeDefined();
        expect(result.current.certificates).toEqual([]);
    });

    it("should call API with correct eventId", async () => {
        vi.mocked(coreApiClient.v1.getEventCertificates).mockResolvedValue({
            certificates: mockCertificates,
        });

        renderHook(() => useEventCertificates(mockEventId), { wrapper: createWrapper() });

        await waitFor(() => {
            expect(coreApiClient.v1.getEventCertificates).toHaveBeenCalledWith({
                eventId: mockEventId,
            });
        });
    });

    it("should provide refetch function", async () => {
        vi.mocked(coreApiClient.v1.getEventCertificates)
            .mockResolvedValueOnce({
                certificates: mockCertificates,
            })
            .mockResolvedValueOnce({
                certificates: [mockCertificates[0]],
            });

        const { result } = renderHook(() => useEventCertificates(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificates).toHaveLength(2);

        // Refetch
        result.current.refetch();

        await waitFor(() => {
            expect(result.current.certificates).toHaveLength(1);
        });
    });

    it("should handle certificates with revoked_at timestamp", async () => {
        const revokedCertificates: EntityEventCertificate[] = [
            ...mockCertificates,
            {
                id: "cert-3",
                event_id: mockEventId,
                event_contract_address: "0x1234567890123456789012345678901234567890",
                receiver_email: "user3@example.com",
                revoked_at: "2024-01-15T00:00:00Z",
                created_at: "2024-01-10T00:00:00Z",
            },
        ];

        vi.mocked(coreApiClient.v1.getEventCertificates).mockResolvedValue({
            certificates: revokedCertificates,
        });

        const { result } = renderHook(() => useEventCertificates(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificates).toHaveLength(3);
        const revokedCert = result.current.certificates.find((c) => c.id === "cert-3");
        expect(revokedCert?.revoked_at).toBeDefined();
    });
});
