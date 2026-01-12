import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { SearchCertificateNav } from "./SearchCertificateNav";

const mockOnBack = vi.fn();
const mockSetSearchQuery = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

vi.mock("../stores/certificates", () => ({
    useSearchCertificateNavStore: vi.fn(() => ({
        searchQuery: "",
        setSearchQuery: mockSetSearchQuery,
    })),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

describe("SearchCertificateNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders back button and search input", () => {
        render(<SearchCertificateNav />);
        expect(screen.getByRole("button", { name: /go back/i })).toBeInTheDocument();
        expect(screen.getByRole("textbox")).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<SearchCertificateNav />);
        await user.click(screen.getByRole("button", { name: /go back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });

    it("calls setSearchQuery on input change", async () => {
        const user = userEvent.setup();
        render(<SearchCertificateNav />);
        await user.type(screen.getByRole("textbox"), "a");
        expect(mockSetSearchQuery).toHaveBeenCalled();
    });
});
