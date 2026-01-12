import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { CertificateSigningNav } from "./CertificateSigningNav";

const mockOnBack = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => fallback || _key,
    }),
}));

describe("CertificateSigningNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders back button", () => {
        render(<CertificateSigningNav />);
        expect(screen.getByRole("button", { name: /go back/i })).toBeInTheDocument();
    });

    it("renders sign button", () => {
        render(<CertificateSigningNav />);
        expect(
            screen.getByRole("button", { name: /participant\.certificates\.sign/i }),
        ).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<CertificateSigningNav />);
        await user.click(screen.getByRole("button", { name: /go back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });
});
