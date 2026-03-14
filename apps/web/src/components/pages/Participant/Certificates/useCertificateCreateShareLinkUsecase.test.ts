import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import React from "react";
import { useCertificateCreateShareLinkUsecase } from "./useCertificateCreateShareLinkUsecase";

vi.mock("@/services/services", () => ({
    certificateService: {
        createCertificateShareableLink: vi.fn(),
    },
}));

vi.mock("sonner", () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    },
}));

vi.mock("@/lib/api/queryClient", () => ({
    queryClient: {
        invalidateQueries: vi.fn(),
    },
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string, fallback: string) => fallback ?? key,
    }),
}));

import { certificateService } from "@/services/services";
import { toast } from "sonner";
import { queryClient } from "@/lib/api/queryClient";

describe("useCertificateCreateShareLinkUsecase", () => {
    let testQueryClient: QueryClient;

    beforeEach(() => {
        testQueryClient = new QueryClient({
            defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
        });
        vi.clearAllMocks();
    });

    const wrapper = ({ children }: { children: ReactNode }) =>
        React.createElement(QueryClientProvider, { client: testQueryClient }, children);

    it("returns createCertificateShareLink, isPending=false, error=null initially", () => {
        const { result } = renderHook(() => useCertificateCreateShareLinkUsecase(), { wrapper });

        expect(typeof result.current.createCertificateShareLink).toBe("function");
        expect(result.current.isPending).toBe(false);
        expect(result.current.error).toBeNull();
    });

    it("calls createCertificateShareableLink with certificateId on mutate", async () => {
        vi.mocked(certificateService.createCertificateShareableLink).mockResolvedValue({
            certificateId: "cert-1",
            isPublished: true,
            handle: "abc123",
        });

        const { result } = renderHook(() => useCertificateCreateShareLinkUsecase(), { wrapper });

        await act(async () => {
            await result.current.createCertificateShareLink("cert-1");
        });

        expect(certificateService.createCertificateShareableLink).toHaveBeenCalledWith(
            "cert-1",
            undefined,
        );
    });

    it("on success: invalidates query cache and shows success toast", async () => {
        vi.mocked(certificateService.createCertificateShareableLink).mockResolvedValue({
            certificateId: "cert-1",
            isPublished: true,
            handle: "abc123",
        });
        vi.mocked(queryClient.invalidateQueries).mockResolvedValue(undefined);

        const { result } = renderHook(() => useCertificateCreateShareLinkUsecase(), { wrapper });

        await act(async () => {
            await result.current.createCertificateShareLink("cert-1");
        });

        await waitFor(() => {
            expect(queryClient.invalidateQueries).toHaveBeenCalled();
            expect(toast.success).toHaveBeenCalled();
        });
    });

    it("on error: shows error toast", async () => {
        vi.mocked(certificateService.createCertificateShareableLink).mockRejectedValue(
            new Error("Network error"),
        );

        const { result } = renderHook(() => useCertificateCreateShareLinkUsecase(), { wrapper });

        await act(async () => {
            try {
                await result.current.createCertificateShareLink("cert-1");
            } catch {
                // expected rejection
            }
        });

        await waitFor(() => {
            expect(toast.error).toHaveBeenCalled();
        });
    });
});
