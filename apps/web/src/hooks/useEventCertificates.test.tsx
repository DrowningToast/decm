import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { useEventCertificates } from "./useEventCertificates";
import { certificateService } from "@/services/services";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import type { Certificate } from "@/services/CertificateService/mapper";

// Mock the certificate service
vi.mock("@/services/services", () => ({
    certificateService: {
        getEventCertificates: vi.fn(),
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
    const mockCertificates: Certificate[] = [
        {
            id: "cert-1",
            eventId: mockEventId,
            eventContractAddress: "0x1234567890123456789012345678901234567890",
            receiverEmail: "user1@example.com",
            name: "John Doe",
            certificateTitle: "Certificate of Completion",
            createdAt: "2024-01-01T00:00:00Z",
        },
        {
            id: "cert-2",
            eventId: mockEventId,
            eventContractAddress: "0x1234567890123456789012345678901234567890",
            receiverEmail: "user2@example.com",
            name: "Jane Doe",
            certificateTitle: "Certificate of Achievement",
            createdAt: "2024-01-02T00:00:00Z",
        },
    ];

    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    it("should fetch and return event certificates", async () => {
        vi.mocked(certificateService.getEventCertificates).mockResolvedValue(mockCertificates);

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
        vi.mocked(certificateService.getEventCertificates).mockResolvedValue([]);

        const { result } = renderHook(() => useEventCertificates(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificates).toEqual([]);
    });

    it("should handle when response has empty certificates array", async () => {
        vi.mocked(certificateService.getEventCertificates).mockResolvedValue([]);

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

        expect(certificateService.getEventCertificates).not.toHaveBeenCalled();
        expect(result.current.certificates).toEqual([]);
    });

    it("should handle API errors gracefully", async () => {
        const error = new Error("Failed to fetch certificates");
        vi.mocked(certificateService.getEventCertificates).mockRejectedValue(error);

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
        vi.mocked(certificateService.getEventCertificates).mockResolvedValue(mockCertificates);

        renderHook(() => useEventCertificates(mockEventId), { wrapper: createWrapper() });

        await waitFor(() => {
            expect(certificateService.getEventCertificates).toHaveBeenCalledWith(mockEventId);
        });
    });

    it("should provide refetch function", async () => {
        vi.mocked(certificateService.getEventCertificates)
            .mockResolvedValueOnce(mockCertificates)
            .mockResolvedValueOnce([mockCertificates[0]]);

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

    it("should handle certificates with revokedAt timestamp", async () => {
        const revokedCertificates: Certificate[] = [
            ...mockCertificates,
            {
                id: "cert-3",
                eventId: mockEventId,
                eventContractAddress: "0x1234567890123456789012345678901234567890",
                receiverEmail: "user3@example.com",
                revokedAt: "2024-01-15T00:00:00Z",
                createdAt: "2024-01-10T00:00:00Z",
            },
        ];

        vi.mocked(certificateService.getEventCertificates).mockResolvedValue(revokedCertificates);

        const { result } = renderHook(() => useEventCertificates(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.certificates).toHaveLength(3);
        const revokedCert = result.current.certificates.find((c) => c.id === "cert-3");
        expect(revokedCert?.revokedAt).toBeDefined();
    });
});
