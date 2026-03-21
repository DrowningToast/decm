import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { CertificateShareModal } from "./CertificateShareModal";

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => (typeof fallback === "string" ? fallback : _key),
    }),
}));

const mockSetIsShareModalOpen = vi.fn();

vi.mock("@/components/BottomNav/stores/certificates", () => ({
    useCertificateDetailNavStore: vi.fn(() => ({
        isShareModalOpen: false,
        setIsShareModalOpen: mockSetIsShareModalOpen,
    })),
}));

import { useCertificateDetailNavStore } from "@/components/BottomNav/stores/certificates";

describe("CertificateShareModal", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("does not render dialog content when closed", () => {
        vi.mocked(useCertificateDetailNavStore).mockReturnValue({
            isShareModalOpen: false,
            setIsShareModalOpen: mockSetIsShareModalOpen,
        } as ReturnType<typeof useCertificateDetailNavStore>);
        render(<CertificateShareModal />);
        expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
    });

    it("renders dialog when open", () => {
        vi.mocked(useCertificateDetailNavStore).mockReturnValue({
            isShareModalOpen: true,
            setIsShareModalOpen: mockSetIsShareModalOpen,
        } as ReturnType<typeof useCertificateDetailNavStore>);
        render(<CertificateShareModal />);
        expect(screen.getByRole("dialog")).toBeInTheDocument();
    });

    it("calls setIsShareModalOpen(false) when dialog is closed", () => {
        vi.mocked(useCertificateDetailNavStore).mockReturnValue({
            isShareModalOpen: true,
            setIsShareModalOpen: mockSetIsShareModalOpen,
        } as ReturnType<typeof useCertificateDetailNavStore>);
        render(<CertificateShareModal />);
        // Close via the Dialog's onOpenChange (e.g. pressing Escape)
        const dialog = screen.getByRole("dialog");
        fireEvent.keyDown(dialog, { key: "Escape" });
        expect(mockSetIsShareModalOpen).toHaveBeenCalledWith(false);
    });
});
