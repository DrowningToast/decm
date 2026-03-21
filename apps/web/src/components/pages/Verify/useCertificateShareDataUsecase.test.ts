import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import React from "react";
import { AxiosError } from "axios";
import { useCertificateShareDataUsecase } from "./useCertificateShareDataUsecase";

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

vi.mock("@/services/services", () => ({
    certificateService: {
        getCertificateShareData: vi.fn(),
    },
}));

vi.mock("sonner", () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
        info: vi.fn(),
        warning: vi.fn(),
    },
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

vi.mock("@/common/Err", () => ({
    ToastFromAxiosError: vi.fn(),
}));

import { certificateService } from "@/services/services";
import { toast } from "sonner";
import { ToastFromAxiosError } from "@/common/Err";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function makeWrapper() {
    const queryClient = new QueryClient({
        defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
    });
    return ({ children }: { children: ReactNode }) =>
        React.createElement(QueryClientProvider, { client: queryClient }, children);
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("useCertificateShareDataUsecase", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("returns null initial state when handle is null", () => {
        const wrapper = makeWrapper();
        const { result } = renderHook(() => useCertificateShareDataUsecase(null), { wrapper });

        expect(result.current.shareStatus).toBeNull();
        expect(result.current.shareData).toBeNull();
        expect(result.current.isLoading).toBe(false);
        expect(result.current.isUnlocking).toBe(false);
        expect(result.current.passwordError).toBe("");
    });

    it("on successful fetch: sets shareStatus to READY and shareData, shows no toast", async () => {
        const mockData = { data: { data: { certificateTitle: "My Cert" } } };
        vi.mocked(certificateService.getCertificateShareData).mockResolvedValue(mockData as never);

        const wrapper = makeWrapper();
        const { result } = renderHook(() => useCertificateShareDataUsecase("abc123"), { wrapper });

        await waitFor(() => {
            expect(result.current.shareStatus).toBe("READY");
        });

        expect(result.current.shareData).toBe(mockData);
        expect(toast.info).not.toHaveBeenCalled();
        expect(toast.error).not.toHaveBeenCalled();
    });

    it("on 403: sets shareStatus to PASSWORD_LOCKED, calls toast.info, does not set isError", async () => {
        vi.mocked(certificateService.getCertificateShareData).mockRejectedValue({ status: 403 });

        const wrapper = makeWrapper();
        const { result } = renderHook(() => useCertificateShareDataUsecase("abc123"), { wrapper });

        await waitFor(() => {
            expect(result.current.shareStatus).toBe("PASSWORD_LOCKED");
        });

        expect(toast.info).toHaveBeenCalledWith("certificateVerify.passwordProtected");
        expect(result.current.isError).toBe(false);
    });

    it("on 404 AxiosError: calls toast.error with notFound key and sets isError", async () => {
        const err = new AxiosError("Not found");
        Object.assign(err, { response: { status: 404 } });
        vi.mocked(certificateService.getCertificateShareData).mockRejectedValue(err);

        const wrapper = makeWrapper();
        const { result } = renderHook(() => useCertificateShareDataUsecase("abc123"), { wrapper });

        await waitFor(() => {
            expect(result.current.isError).toBe(true);
        });

        expect(toast.error).toHaveBeenCalledWith("certificateVerify.notFound");
        expect(ToastFromAxiosError).not.toHaveBeenCalled();
    });

    it("on other AxiosError: calls ToastFromAxiosError and sets isError", async () => {
        const err = new AxiosError("Server error");
        Object.assign(err, { response: { status: 500 } });
        vi.mocked(certificateService.getCertificateShareData).mockRejectedValue(err);

        const wrapper = makeWrapper();
        const { result } = renderHook(() => useCertificateShareDataUsecase("abc123"), { wrapper });

        await waitFor(() => {
            expect(result.current.isError).toBe(true);
        });

        expect(ToastFromAxiosError).toHaveBeenCalled();
        expect(toast.error).not.toHaveBeenCalled();
    });

    describe("handleUnlock", () => {
        it("with empty password: sets passwordError and does not call service for unlock", async () => {
            vi.mocked(certificateService.getCertificateShareData).mockReturnValue(
                new Promise(() => {}),
            ); // initial query hangs

            const wrapper = makeWrapper();
            const { result } = renderHook(() => useCertificateShareDataUsecase("abc123"), {
                wrapper,
            });

            await act(async () => {
                await result.current.handleUnlock("");
            });

            expect(result.current.passwordError).toBe("certificateVerify.passwordRequired");
            // Only the initial query call; unlock must not have fired a second call
            expect(certificateService.getCertificateShareData).toHaveBeenCalledTimes(1);
        });

        it("with wrong password: sets passwordError to incorrectPassword and shows error toast", async () => {
            vi.mocked(certificateService.getCertificateShareData)
                .mockRejectedValueOnce({ status: 403 }) // initial query
                .mockRejectedValueOnce(new Error("Wrong password")); // unlock attempt

            const wrapper = makeWrapper();
            const { result } = renderHook(() => useCertificateShareDataUsecase("abc123"), {
                wrapper,
            });

            await waitFor(() => expect(result.current.shareStatus).toBe("PASSWORD_LOCKED"));

            await act(async () => {
                await result.current.handleUnlock("wrongpass");
            });

            expect(result.current.passwordError).toBe("certificateVerify.incorrectPassword");
            expect(result.current.isUnlocking).toBe(false);
            expect(toast.error).toHaveBeenCalledWith("certificateVerify.incorrectPassword");
        });

        it("with correct password: sets shareStatus to READY and shareData, shows success toast", async () => {
            const mockData = { data: { data: { certificateTitle: "My Cert" } } };
            vi.mocked(certificateService.getCertificateShareData)
                .mockRejectedValueOnce({ status: 403 }) // initial query
                .mockResolvedValueOnce(mockData as never); // unlock attempt

            const wrapper = makeWrapper();
            const { result } = renderHook(() => useCertificateShareDataUsecase("abc123"), {
                wrapper,
            });

            await waitFor(() => expect(result.current.shareStatus).toBe("PASSWORD_LOCKED"));

            await act(async () => {
                await result.current.handleUnlock("correctpass");
            });

            expect(result.current.shareStatus).toBe("READY");
            expect(result.current.shareData).toBe(mockData);
            expect(result.current.passwordError).toBe("");
            expect(toast.success).toHaveBeenCalledWith("certificateVerify.unlockSuccess");
        });
    });

    it("resets shareStatus and shareData when handle changes", async () => {
        const mockData = { data: { data: { certificateTitle: "My Cert" } } };
        vi.mocked(certificateService.getCertificateShareData).mockResolvedValue(mockData as never);

        const wrapper = makeWrapper();
        const { result, rerender } = renderHook(
            ({ handle }: { handle: string }) => useCertificateShareDataUsecase(handle),
            { wrapper, initialProps: { handle: "abc123" } },
        );

        await waitFor(() => expect(result.current.shareStatus).toBe("READY"));

        vi.mocked(certificateService.getCertificateShareData).mockReturnValue(
            new Promise(() => {}),
        ); // new query hangs
        rerender({ handle: "newhandle" });

        await waitFor(() => {
            expect(result.current.shareStatus).toBeNull();
            expect(result.current.shareData).toBeNull();
        });
    });
});
