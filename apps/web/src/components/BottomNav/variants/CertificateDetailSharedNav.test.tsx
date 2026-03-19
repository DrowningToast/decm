import React from "react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, act, fireEvent } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { CertificateDetailSharedNav } from "./CertificateDetailSharedNav";

const { mockToastLoading, mockToastDismiss, mockToastPromise, mockToastSuccess } = vi.hoisted(
    () => ({
        mockToastLoading: vi.fn().mockReturnValue("toast-id-1"),
        mockToastDismiss: vi.fn(),
        mockToastPromise: vi.fn(),
        mockToastSuccess: vi.fn(),
    }),
);

vi.mock("sonner", () => ({
    toast: {
        loading: mockToastLoading,
        dismiss: mockToastDismiss,
        promise: mockToastPromise,
        success: mockToastSuccess,
    },
}));

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: vi.fn(),
        className: "",
    }),
}));

vi.mock("../stores/certificates", () => ({
    useCertificateDetailNavStore: vi.fn(() => ({
        certificateId: "cert-1",
        isClaimed: true,
        onClickShareable: vi.fn(),
        isShareableLoading: false,
    })),
    useCertificateDetailsSharedNavStore: vi.fn(() => ({
        shareableHandle: "handle-123",
        isPublished: false,
        onChangePublish: vi.fn().mockResolvedValue(undefined),
        isPasswordProtected: false,
        setIsPasswordDialogOpen: vi.fn(),
    })),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => fallback || _key,
    }),
}));

vi.mock("@/services/services", () => ({
    certificateService: {
        getCertificateImage: vi.fn(),
    },
}));

import { useCertificateDetailsSharedNavStore } from "../stores/certificates";

describe("CertificateDetailSharedNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        mockToastLoading.mockReturnValue("toast-id-1");
    });

    describe("publish toggle — debounce & toast lifecycle", () => {
        beforeEach(() => {
            vi.useFakeTimers();
        });

        afterEach(() => {
            vi.useRealTimers();
        });

        it("shows a loading toast when the publish toggle is changed", () => {
            render(<CertificateDetailSharedNav />);
            fireEvent.click(screen.getByRole("switch"));
            expect(mockToastLoading).toHaveBeenCalledOnce();
        });

        it("calls toast.promise after the 3 s debounce delay", () => {
            render(<CertificateDetailSharedNav />);
            fireEvent.click(screen.getByRole("switch"));
            act(() => vi.advanceTimersByTime(3000));
            expect(mockToastPromise).toHaveBeenCalledOnce();
        });

        it("passes the loading toast id to toast.promise so it updates in-place", () => {
            render(<CertificateDetailSharedNav />);
            fireEvent.click(screen.getByRole("switch"));
            act(() => vi.advanceTimersByTime(3000));
            expect(mockToastPromise).toHaveBeenCalledWith(
                expect.anything(),
                expect.objectContaining({ id: "toast-id-1" }),
            );
        });

        it("dismisses the success/error toast when the toggle is changed again after the promise fires", () => {
            render(<CertificateDetailSharedNav />);

            // First toggle — loading toast created, debounce fires
            fireEvent.click(screen.getByRole("switch"));
            act(() => vi.advanceTimersByTime(3000));

            // Second toggle — must dismiss the toast that toast.promise was using
            fireEvent.click(screen.getByRole("switch"));

            expect(mockToastDismiss).toHaveBeenCalledWith("toast-id-1");
        });

        it("dismisses the loading toast when the user toggles again before the debounce fires", () => {
            render(<CertificateDetailSharedNav />);

            fireEvent.click(screen.getByRole("switch"));
            // Second toggle within 3 s — must dismiss the in-progress loading toast
            fireEvent.click(screen.getByRole("switch"));

            expect(mockToastDismiss).toHaveBeenCalledWith("toast-id-1");
        });
    });

    describe("clipboard actions", () => {
        it("copies the share code to clipboard when the copy code button is clicked", async () => {
            const clipboardSpy = vi
                .spyOn(navigator.clipboard, "writeText")
                .mockResolvedValue(undefined);
            render(<CertificateDetailSharedNav />);

            await act(async () => {
                fireEvent.click(screen.getByRole("button", { name: /copy code/i }));
            });

            expect(clipboardSpy).toHaveBeenCalledWith("handle-123");
        });

        it("copies a URL containing the handle when the copy URL button is clicked", async () => {
            const clipboardSpy = vi
                .spyOn(navigator.clipboard, "writeText")
                .mockResolvedValue(undefined);
            render(<CertificateDetailSharedNav />);

            await act(async () => {
                fireEvent.click(screen.getByRole("button", { name: /copy shareable url/i }));
            });

            expect(clipboardSpy).toHaveBeenCalledWith(expect.stringContaining("handle-123"));
        });
    });

    describe("password dialog", () => {
        it("opens the password dialog when the lock button is clicked", async () => {
            const mockSetIsPasswordDialogOpen = vi.fn();
            vi.mocked(useCertificateDetailsSharedNavStore).mockReturnValue({
                shareableHandle: "handle-123",
                isPublished: false,
                onChangePublish: vi.fn().mockResolvedValue(undefined),
                isPasswordProtected: false,
                setIsPasswordDialogOpen: mockSetIsPasswordDialogOpen,
            } as ReturnType<typeof useCertificateDetailsSharedNavStore>);

            const user = userEvent.setup();
            render(<CertificateDetailSharedNav />);

            await user.click(screen.getByRole("button", { name: /edit password/i }));

            expect(mockSetIsPasswordDialogOpen).toHaveBeenCalledWith(true);
        });
    });
});
