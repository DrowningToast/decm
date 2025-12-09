import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { useCertificateImage } from "./useCertificateImage";

// Mock the certificateService
vi.mock("@/services/services", () => ({
    certificateService: {
        getCertificateImage: vi.fn(),
    },
}));

import { certificateService } from "@/services/services";

describe("useCertificateImage", () => {
    const mockCreateObjectURL = vi.spyOn(URL, "createObjectURL");
    const mockRevokeObjectURL = vi.spyOn(URL, "revokeObjectURL");

    beforeEach(() => {
        vi.clearAllMocks();
        mockCreateObjectURL.mockReturnValue("blob:mock-url");
    });

    afterEach(() => {
        mockCreateObjectURL.mockRestore();
        mockRevokeObjectURL.mockRestore();
    });

    it("should return null when certificateId is not provided", () => {
        const { result } = renderHook(() => useCertificateImage({ certificateId: undefined }));

        expect(result.current.imageUrl).toBeNull();
        expect(result.current.imageData).toBeNull();
        expect(result.current.isLoading).toBe(false);
        expect(certificateService.getCertificateImage).not.toHaveBeenCalled();
    });

    it("should return null when enabled is false", () => {
        const { result } = renderHook(() =>
            useCertificateImage({ certificateId: "cert-123", enabled: false }),
        );

        expect(result.current.imageUrl).toBeNull();
        expect(result.current.imageData).toBeNull();
        expect(certificateService.getCertificateImage).not.toHaveBeenCalled();
    });

    it("should fetch certificate image when certificateId is provided", async () => {
        const mockBlob = new Blob(["image data"], { type: "image/png" });
        const mockImageData = {
            url: "blob:mock-url",
            blob: mockBlob,
            contentType: "image/png",
        };

        vi.mocked(certificateService.getCertificateImage).mockResolvedValue(mockImageData);

        const { result } = renderHook(() => useCertificateImage({ certificateId: "cert-123" }));

        expect(result.current.isLoading).toBe(true);

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(certificateService.getCertificateImage).toHaveBeenCalledWith("cert-123");
        expect(result.current.imageUrl).toBe("blob:mock-url");
        expect(result.current.imageData).toEqual(mockImageData);
        expect(result.current.error).toBeNull();
    });

    it("should handle fetch errors", async () => {
        const error = new Error("Failed to fetch image");
        vi.mocked(certificateService.getCertificateImage).mockRejectedValue(error);

        const { result } = renderHook(() => useCertificateImage({ certificateId: "cert-123" }));

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(result.current.error).toBe(error);
        expect(result.current.imageUrl).toBeNull();
        expect(result.current.imageData).toBeNull();
    });

    it("should refetch when certificateId changes", async () => {
        const mockBlob = new Blob(["image data"], { type: "image/png" });
        const mockImageData = {
            url: "blob:mock-url",
            blob: mockBlob,
            contentType: "image/png",
        };

        vi.mocked(certificateService.getCertificateImage).mockResolvedValue(mockImageData);

        const { result, rerender } = renderHook(
            ({ certificateId }) => useCertificateImage({ certificateId }),
            {
                initialProps: { certificateId: "cert-123" },
            },
        );

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(certificateService.getCertificateImage).toHaveBeenCalledWith("cert-123");
        vi.clearAllMocks();

        rerender({ certificateId: "cert-456" });

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        expect(certificateService.getCertificateImage).toHaveBeenCalledWith("cert-456");
    });

    it("should not fetch when enabled changes to false", async () => {
        const mockBlob = new Blob(["image data"], { type: "image/png" });
        const mockImageData = {
            url: "blob:mock-url",
            blob: mockBlob,
            contentType: "image/png",
        };

        vi.mocked(certificateService.getCertificateImage).mockResolvedValue(mockImageData);

        const { result, rerender } = renderHook(
            ({ enabled }) => useCertificateImage({ certificateId: "cert-123", enabled }),
            {
                initialProps: { enabled: true },
            },
        );

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        vi.clearAllMocks();

        rerender({ enabled: false });

        expect(certificateService.getCertificateImage).not.toHaveBeenCalled();
        expect(result.current.imageUrl).toBeNull();
        expect(result.current.imageData).toBeNull();
    });

    it("should handle abort errors gracefully", async () => {
        const abortError = new Error("Aborted");
        abortError.name = "AbortError";
        vi.mocked(certificateService.getCertificateImage).mockRejectedValue(abortError);

        const { result } = renderHook(() => useCertificateImage({ certificateId: "cert-123" }));

        await waitFor(() => {
            expect(result.current.isLoading).toBe(false);
        });

        // AbortError should not set the error state
        expect(result.current.error).toBeNull();
    });
});
