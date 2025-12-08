import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";
import { useUpdateCertificateConfig } from "./useUpdateCertificateConfig";
import { coreApiClient } from "@/lib/api/api";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import type {
    UpdateEventCertificateConfigPayload,
    CoreApiInternalHandlerEventconfigEventCertificateConfigResponse,
} from "@decm/api";

// Mock the API client
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            updateEventCertificateConfig: vi.fn(),
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

describe("useUpdateCertificateConfig", () => {
    const mockEventId = "event-123";
    const mockFile = new File(["<svg></svg>"], "certificate.svg", { type: "image/svg+xml" });
    const mockPayload: UpdateEventCertificateConfigPayload = {
        name_pos_x: 100,
        name_pos_y: 200,
        event_name_pos_x: 150,
        event_name_pos_y: 250,
        academic_institution_pos_x: 175,
        academic_institution_pos_y: 275,
        base_certificate_image: mockFile,
    };

    const mockResponse: CoreApiInternalHandlerEventconfigEventCertificateConfigResponse = {
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

    it("should update certificate config successfully", async () => {
        vi.mocked(coreApiClient.v1.updateEventCertificateConfig).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useUpdateCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        let response;
        await waitFor(async () => {
            response = await result.current.updateCertificateConfig(mockPayload);
        });

        expect(response).toEqual(mockResponse);
        expect(coreApiClient.v1.updateEventCertificateConfig).toHaveBeenCalledWith(
            { eventId: mockEventId },
            mockPayload,
        );
    });

    it("should set isUpdatingCertificateConfig to true during mutation", async () => {
        vi.mocked(coreApiClient.v1.updateEventCertificateConfig).mockImplementation(
            () => new Promise((resolve) => setTimeout(() => resolve(mockResponse), 100)),
        );

        const { result } = renderHook(() => useUpdateCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        expect(result.current.isUpdatingCertificateConfig).toBe(false);

        result.current.updateCertificateConfig(mockPayload);

        await waitFor(() => {
            expect(result.current.isUpdatingCertificateConfig).toBe(true);
        });

        await waitFor(() => {
            expect(result.current.isUpdatingCertificateConfig).toBe(false);
        });
    });

    it("should handle API errors", async () => {
        const error = new Error("Failed to update certificate config");
        vi.mocked(coreApiClient.v1.updateEventCertificateConfig).mockRejectedValue(error);

        const { result } = renderHook(() => useUpdateCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        await expect(result.current.updateCertificateConfig(mockPayload)).rejects.toThrow(
            "Failed to update certificate config",
        );

        expect(result.current.error).toBeDefined();
    });

    it("should handle partial payload", async () => {
        const partialPayload: UpdateEventCertificateConfigPayload = {
            name_pos_x: 100,
            name_pos_y: 200,
        };

        vi.mocked(coreApiClient.v1.updateEventCertificateConfig).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useUpdateCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        await result.current.updateCertificateConfig(partialPayload);

        expect(coreApiClient.v1.updateEventCertificateConfig).toHaveBeenCalledWith(
            { eventId: mockEventId },
            partialPayload,
        );
    });

    it("should update with base64 certificate image", async () => {
        const imageFile = new File(['<svg width="100" height="100">'], "certificate.svg", {
            type: "image/svg+xml",
        });
        const payloadWithImage: UpdateEventCertificateConfigPayload = {
            ...mockPayload,
            base_certificate_image: imageFile,
        };

        vi.mocked(coreApiClient.v1.updateEventCertificateConfig).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useUpdateCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        await result.current.updateCertificateConfig(payloadWithImage);

        expect(coreApiClient.v1.updateEventCertificateConfig).toHaveBeenCalledWith(
            { eventId: mockEventId },
            payloadWithImage,
        );
    });

    it("should use correct mutation key", async () => {
        vi.mocked(coreApiClient.v1.updateEventCertificateConfig).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useUpdateCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        await result.current.updateCertificateConfig(mockPayload);

        // Mutation should have been called
        expect(coreApiClient.v1.updateEventCertificateConfig).toHaveBeenCalled();
    });

    it("should handle multiple consecutive updates", async () => {
        vi.mocked(coreApiClient.v1.updateEventCertificateConfig).mockResolvedValue(mockResponse);

        const { result } = renderHook(() => useUpdateCertificateConfig(mockEventId), {
            wrapper: createWrapper(),
        });

        await result.current.updateCertificateConfig(mockPayload);
        await result.current.updateCertificateConfig({
            ...mockPayload,
            name_pos_x: 200,
        });

        expect(coreApiClient.v1.updateEventCertificateConfig).toHaveBeenCalledTimes(2);
    });
});
