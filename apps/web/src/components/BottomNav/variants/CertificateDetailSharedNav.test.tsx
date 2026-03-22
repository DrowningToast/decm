import React from "react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, act, fireEvent } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { CertificateDetailSharedNav } from "./CertificateDetailSharedNav";

const { mockToastLoading, mockToastDismiss, mockToastPromise } = vi.hoisted(() => ({
    mockToastLoading: vi.fn().mockReturnValue("toast-id-1"),
    mockToastDismiss: vi.fn(),
    mockToastPromise: vi.fn(),
}));

vi.mock("sonner", () => ({
    toast: {
        loading: mockToastLoading,
        dismiss: mockToastDismiss,
        promise: mockToastPromise,
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

    describe("password dialog", () => {
        it("opens the password dialog when the lock button is clicked", async () => {
            const mockSetIsPasswordDialogOpen = vi.fn();
            vi.mocked(useCertificateDetailsSharedNavStore).mockReturnValue({
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
