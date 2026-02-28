import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { InboxViewCertificateNav } from "./InboxViewCertificateNav";

const mockOnBack = vi.fn();
const mockOnViewCertificateCallback = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

vi.mock("../stores/inbox", () => ({
    useInboxNavStore: vi.fn(() => ({
        onViewCertificateCallback: mockOnViewCertificateCallback,
    })),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => fallback || _key,
    }),
}));

import { useInboxNavStore } from "../stores/inbox";

describe("InboxViewCertificateNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders back button and view certificate button", () => {
        render(<InboxViewCertificateNav />);
        expect(screen.getByRole("button", { name: /go back/i })).toBeInTheDocument();
        expect(screen.getByRole("button", { name: /view certificate/i })).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<InboxViewCertificateNav />);
        await user.click(screen.getByRole("button", { name: /go back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });

    it("calls onViewCertificateCallback when view button is clicked", async () => {
        const user = userEvent.setup();
        render(<InboxViewCertificateNav />);
        await user.click(screen.getByRole("button", { name: /view certificate/i }));
        expect(mockOnViewCertificateCallback).toHaveBeenCalledOnce();
    });

    it("does not crash when callback is undefined", async () => {
        vi.mocked(useInboxNavStore).mockReturnValue({
            onViewCertificateCallback: undefined,
        } as unknown as ReturnType<typeof useInboxNavStore>);

        const user = userEvent.setup();
        render(<InboxViewCertificateNav />);
        await user.click(screen.getByRole("button", { name: /view certificate/i }));
        expect(screen.getByRole("button", { name: /view certificate/i })).toBeInTheDocument();
    });
});
