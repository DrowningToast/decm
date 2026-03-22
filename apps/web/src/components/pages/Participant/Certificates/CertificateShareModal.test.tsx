import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { CertificateShareModal } from "./CertificateShareModal";

// ─── Mocks ────────────────────────────────────────────────────────────────────

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => (typeof fallback === "string" ? fallback : _key),
    }),
}));

vi.mock("react-router-dom", () => ({
    useNavigate: () => vi.fn(),
    Link: ({
        children,
        to,
        ...props
    }: React.PropsWithChildren<{ to: string; [key: string]: unknown }>) =>
        React.createElement("a", { href: to, ...props }, children),
}));

vi.mock("qrcode.react", () => ({
    QRCodeCanvas: ({ value }: { value: string }) => (
        <canvas data-testid="qr-code" data-value={value} />
    ),
}));

const mockSetIsShareModalOpen = vi.fn();
const mockNavDetailStore = {
    isShareModalOpen: false,
    setIsShareModalOpen: mockSetIsShareModalOpen,
};

vi.mock("@/components/BottomNav/stores/certificates", () => ({
    useCertificateDetailNavStore: vi.fn(() => mockNavDetailStore),
    useCertificateDetailsSharedNavStore: vi.fn(() => mockSharedNavStore),
}));

const mockSharedNavStore: {
    shareableHandle: string | null;
    isPublished: boolean;
    isPasswordProtected: boolean;
} = {
    shareableHandle: "abc123",
    isPublished: true,
    isPasswordProtected: false,
};

import {
    useCertificateDetailNavStore,
    useCertificateDetailsSharedNavStore,
} from "@/components/BottomNav/stores/certificates";

// ─── Helpers ──────────────────────────────────────────────────────────────────

function renderOpenModal(
    navOverrides: Partial<typeof mockNavDetailStore> = {},
    sharedOverrides: Partial<typeof mockSharedNavStore> = {},
) {
    vi.mocked(useCertificateDetailNavStore).mockReturnValue({
        ...mockNavDetailStore,
        isShareModalOpen: true,
        ...navOverrides,
    } as ReturnType<typeof useCertificateDetailNavStore>);
    vi.mocked(useCertificateDetailsSharedNavStore).mockReturnValue({
        ...mockSharedNavStore,
        ...sharedOverrides,
    } as ReturnType<typeof useCertificateDetailsSharedNavStore>);
    return render(<CertificateShareModal />);
}

// ─── Tests ────────────────────────────────────────────────────────────────────

describe("CertificateShareModal", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.mocked(useCertificateDetailNavStore).mockReturnValue({
            ...mockNavDetailStore,
        } as ReturnType<typeof useCertificateDetailNavStore>);
        vi.mocked(useCertificateDetailsSharedNavStore).mockReturnValue({
            ...mockSharedNavStore,
        } as ReturnType<typeof useCertificateDetailsSharedNavStore>);
    });

    // ── Open / close ──────────────────────────────────────────────────────────

    it("does not render dialog content when closed", () => {
        render(<CertificateShareModal />);
        expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
    });

    it("renders dialog when open", () => {
        renderOpenModal();
        expect(screen.getByRole("dialog")).toBeInTheDocument();
    });

    it("calls setIsShareModalOpen(false) when Escape is pressed", () => {
        renderOpenModal();
        fireEvent.keyDown(screen.getByRole("dialog"), { key: "Escape" });
        expect(mockSetIsShareModalOpen).toHaveBeenCalledWith(false);
    });

    // ── Header ────────────────────────────────────────────────────────────────

    it("shows the modal title", () => {
        renderOpenModal();
        expect(screen.getByRole("heading", { name: /share certificate/i })).toBeInTheDocument();
    });

    it("shows a description", () => {
        renderOpenModal();
        // description text is present somewhere in the dialog
        expect(screen.getByRole("dialog")).toHaveTextContent(/.+/);
    });

    it("has a close button", () => {
        renderOpenModal();
        expect(screen.getByRole("button", { name: /close/i })).toBeInTheDocument();
    });

    // ── Status badges ─────────────────────────────────────────────────────────

    it("shows 'Hidden' warning badge when not published", () => {
        renderOpenModal({}, { isPublished: false });
        expect(screen.getByTestId("badge-hidden")).toBeInTheDocument();
    });

    it("does not show hidden warning badge when published", () => {
        renderOpenModal({}, { isPublished: true });
        expect(screen.queryByTestId("badge-hidden")).not.toBeInTheDocument();
    });

    it("shows 'Published' badge when published", () => {
        renderOpenModal({}, { isPublished: true });
        expect(screen.getByTestId("badge-published")).toBeInTheDocument();
    });

    it("shows password-protected badge when password protected", () => {
        renderOpenModal({}, { isPasswordProtected: true });
        expect(screen.getByTestId("badge-password-protected")).toBeInTheDocument();
    });

    it("shows no-password badge when not password protected", () => {
        renderOpenModal({}, { isPasswordProtected: false });
        expect(screen.getByTestId("badge-no-password")).toBeInTheDocument();
    });

    // ── QR code ───────────────────────────────────────────────────────────────

    it("renders the QR code with the share URL", () => {
        renderOpenModal();
        const qr = screen.getByTestId("qr-code");
        expect(qr).toBeInTheDocument();
        expect(qr.getAttribute("data-value")).toContain("abc123");
    });

    it("shows a download QR button", () => {
        renderOpenModal();
        expect(screen.getByRole("button", { name: /download qr/i })).toBeInTheDocument();
    });

    it("shows a copy QR button", () => {
        renderOpenModal();
        expect(screen.getByRole("button", { name: /copy qr/i })).toBeInTheDocument();
    });

    // ── Share URL ─────────────────────────────────────────────────────────────

    it("displays the full share URL containing the handle", () => {
        renderOpenModal();
        expect(screen.getByTestId("share-url")).toHaveTextContent("abc123");
    });

    it("has a copy button for the share URL", () => {
        renderOpenModal();
        expect(screen.getByRole("button", { name: /copy url/i })).toBeInTheDocument();
    });

    // ── Share code ────────────────────────────────────────────────────────────

    it("displays the share handle/code", () => {
        renderOpenModal();
        expect(screen.getByTestId("share-code")).toHaveTextContent("abc123");
    });

    it("has a copy button for the share code", () => {
        renderOpenModal();
        expect(screen.getByRole("button", { name: /copy code/i })).toBeInTheDocument();
    });

    // ── Action button ─────────────────────────────────────────────────────────

    it("has a button that links to the verify page with handle prefilled", () => {
        renderOpenModal();
        const link = screen.getByRole("link", { name: /open verify page/i });
        expect(link).toBeInTheDocument();
        expect(link.getAttribute("href")).toContain("abc123");
    });

    // ── Null handle ───────────────────────────────────────────────────────────

    it("does not render share content when handle is null", () => {
        renderOpenModal({}, { shareableHandle: null });
        expect(screen.queryByTestId("share-url")).not.toBeInTheDocument();
        expect(screen.queryByTestId("share-code")).not.toBeInTheDocument();
        expect(screen.queryByTestId("qr-code")).not.toBeInTheDocument();
    });

    // ── Copy interactions ─────────────────────────────────────────────────────

    it("copies the URL to clipboard when copy URL button is clicked", async () => {
        const writeText = vi.spyOn(navigator.clipboard, "writeText").mockResolvedValue(undefined);
        renderOpenModal();
        fireEvent.click(screen.getByTestId("copy-url-btn"));
        await Promise.resolve(); // flush async handler
        expect(writeText).toHaveBeenCalledWith(expect.stringContaining("abc123"));
    });

    it("copies the handle to clipboard when copy code button is clicked", async () => {
        const writeText = vi.spyOn(navigator.clipboard, "writeText").mockResolvedValue(undefined);
        renderOpenModal();
        fireEvent.click(screen.getByTestId("copy-code-btn"));
        await Promise.resolve(); // flush async handler
        expect(writeText).toHaveBeenCalledWith("abc123");
    });
});
