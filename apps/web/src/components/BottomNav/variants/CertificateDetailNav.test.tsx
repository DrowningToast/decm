import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { CertificateDetailNav } from "./CertificateDetailNav";

const mockOnBack = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

const mockSetIsShareModalOpen = vi.fn();
const mockOnClickShareable = vi.fn();

vi.mock("../stores/certificates", () => ({
    useCertificateDetailNavStore: vi.fn(() => ({
        certificateId: "cert-1",
        isClaimed: true,
        onClickShareable: mockOnClickShareable,
        isShareableLoading: false,
        setIsShareModalOpen: mockSetIsShareModalOpen,
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

import { useCertificateDetailNavStore } from "../stores/certificates";
import { certificateService } from "@/services/services";

describe("CertificateDetailNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.mocked(useCertificateDetailNavStore).mockReturnValue({
            certificateId: "cert-1",
            isClaimed: true,
            onClickShareable: mockOnClickShareable,
            isShareableLoading: false,
            setIsShareModalOpen: mockSetIsShareModalOpen,
        } as ReturnType<typeof useCertificateDetailNavStore>);
    });

    it("renders back button", () => {
        render(<CertificateDetailNav />);
        expect(screen.getByRole("button", { name: /back/i })).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<CertificateDetailNav />);
        await user.click(screen.getByRole("button", { name: /back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });

    it("shows download button when certificate is claimed", () => {
        render(<CertificateDetailNav />);
        expect(screen.getByRole("button", { name: /download as an image/i })).toBeInTheDocument();
    });

    it("hides download button when certificate is not claimed", () => {
        vi.mocked(useCertificateDetailNavStore).mockReturnValue({
            certificateId: "cert-1",
            isClaimed: false,
            setIsShareModalOpen: mockSetIsShareModalOpen,
        } as ReturnType<typeof useCertificateDetailNavStore>);

        render(<CertificateDetailNav />);
        expect(screen.queryByRole("button", { name: /download/i })).not.toBeInTheDocument();
    });

    it("shows share button when claimed", () => {
        render(<CertificateDetailNav />);
        expect(screen.getByRole("button", { name: /^share$/i })).toBeInTheDocument();
    });

    it("hides share button when not claimed", () => {
        vi.mocked(useCertificateDetailNavStore).mockReturnValue({
            certificateId: "cert-1",
            isClaimed: false,
            onClickShareable: mockOnClickShareable,
            isShareableLoading: false,
            setIsShareModalOpen: mockSetIsShareModalOpen,
        } as ReturnType<typeof useCertificateDetailNavStore>);

        render(<CertificateDetailNav />);
        expect(screen.queryByRole("button", { name: /^share$/i })).not.toBeInTheDocument();
    });

    it("calls onClickShareable when share button is clicked and no shareable exists", async () => {
        const user = userEvent.setup();
        render(<CertificateDetailNav />);
        await user.click(screen.getByRole("button", { name: /^share$/i }));
        expect(mockOnClickShareable).toHaveBeenCalledOnce();
        expect(mockSetIsShareModalOpen).not.toHaveBeenCalled();
    });

    it("calls setIsShareModalOpen(true) when share button is clicked and shareable exists", async () => {
        const user = userEvent.setup();
        render(<CertificateDetailNav overrideShowCreateShareButton={false} />);
        await user.click(screen.getByRole("button", { name: /^share$/i }));
        expect(mockSetIsShareModalOpen).toHaveBeenCalledWith(true);
        expect(mockOnClickShareable).not.toHaveBeenCalled();
    });

    it("calls getCertificateImage on download click", async () => {
        vi.mocked(certificateService.getCertificateImage).mockResolvedValue({
            url: "blob:http://localhost/fake",
            blob: new Blob(),
            contentType: "image/png",
        });

        const user = userEvent.setup();
        render(<CertificateDetailNav />);
        await user.click(screen.getByRole("button", { name: /download as an image/i }));

        expect(certificateService.getCertificateImage).toHaveBeenCalledWith("cert-1");
    });
});
