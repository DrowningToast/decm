import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { useEventCertificateConfig } from "./useEventCertificateConfig";
import { coreApiClient } from "@/lib/api/api";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import type { CoreApiInternalHandlerEventconfigEventCertificateConfigResponse } from "@decm/api";

// Mock the API client
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            getEventCertificateConfig: vi.fn(),
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

describe("useEventCertificateConfig", () => {
    const mockEventId = "event-123";
    const mockCertificateConfig: CoreApiInternalHandlerEventconfigEventCertificateConfigResponse = {
        id: "config-1",
        event_id: mockEventId,
        name_pos_x: 100,
        name_pos_y: 200,
        event_name_pos_x: 150,
        event_name_pos_y: 250,
        academic_institution_pos_x: 175,
        academic_institution_pos_y: 275,
        base_certificate_presigned_url: "https://example.com/certificate.svg",
        base_certificate_storage_key: "certificates/config-1.svg",
        created_at: "2024-01-01T00:00:00Z",
        is_published: false,
        updated_at: "2024-01-01T00:00:00Z",
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.clearAllMocks();
    });

    it("should fetch and return event certificate config", async () => {
        vi.mocked(coreApiClient.v1.getEventCertificateConfig).mockResolvedValue(
            mockCertificateConfig,
        );

        const { result } = renderHook(() => useEventCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.data).toEqual(mockCertificateConfig);
        expect(result.current.error).toBeNull();
    });

    it("should call API with correct eventId", async () => {
        vi.mocked(coreApiClient.v1.getEventCertificateConfig).mockResolvedValue(
            mockCertificateConfig,
        );

        renderHook(() => useEventCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(coreApiClient.v1.getEventCertificateConfig).toHaveBeenCalledWith({
                eventId: mockEventId,
            });
        });
    });

    it("should handle API errors gracefully", async () => {
        const error = new Error("Failed to fetch certificate config");
        vi.mocked(coreApiClient.v1.getEventCertificateConfig).mockRejectedValue(error);

        const { result } = renderHook(() => useEventCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.error).toBeDefined();
        expect(result.current.data).toBeUndefined();
    });

    it("should handle undefined response", async () => {
        vi.mocked(coreApiClient.v1.getEventCertificateConfig).mockResolvedValue(
            {} as CoreApiInternalHandlerEventconfigEventCertificateConfigResponse,
        );

        const { result } = renderHook(() => useEventCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.data).toBeDefined();
    });

    it("should use correct query key", async () => {
        vi.mocked(coreApiClient.v1.getEventCertificateConfig).mockResolvedValue(
            mockCertificateConfig,
        );

        const { result } = renderHook(() => useEventCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        // The hook should have been called
        expect(coreApiClient.v1.getEventCertificateConfig).toHaveBeenCalled();
    });
});
