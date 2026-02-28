import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { SearchIdentitiesNav } from "./SearchIdentitiesNav";

const mockOnBack = vi.fn();
const mockSetSearchQuery = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

vi.mock("../stores/identities", () => ({
    useSearchIdentitiesNavStore: vi.fn(() => ({
        searchQuery: "",
        setSearchQuery: mockSetSearchQuery,
    })),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => fallback || _key,
    }),
}));

describe("SearchIdentitiesNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders back button and search input", () => {
        render(<SearchIdentitiesNav />);
        expect(screen.getByRole("button", { name: /go back/i })).toBeInTheDocument();
        expect(screen.getByRole("textbox")).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<SearchIdentitiesNav />);
        await user.click(screen.getByRole("button", { name: /go back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });

    it("calls setSearchQuery on input change", async () => {
        const user = userEvent.setup();
        render(<SearchIdentitiesNav />);
        await user.type(screen.getByRole("textbox"), "a");
        expect(mockSetSearchQuery).toHaveBeenCalled();
    });
});
