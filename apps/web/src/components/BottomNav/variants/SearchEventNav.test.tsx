import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { SearchEventNav } from "./SearchEventNav";

const mockOnBack = vi.fn();
const mockSetSearchQuery = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

vi.mock("../stores/events", () => ({
    useSearchEventNavStore: vi.fn(() => ({
        searchQuery: "",
        setSearchQuery: mockSetSearchQuery,
    })),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

describe("SearchEventNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders back button and search input", () => {
        render(<SearchEventNav />);
        expect(screen.getByRole("button", { name: /back/i })).toBeInTheDocument();
        expect(screen.getByRole("textbox")).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<SearchEventNav />);
        await user.click(screen.getByRole("button", { name: /back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });

    it("calls setSearchQuery on input change", async () => {
        const user = userEvent.setup();
        render(<SearchEventNav />);
        await user.type(screen.getByRole("textbox"), "a");
        expect(mockSetSearchQuery).toHaveBeenCalled();
    });
});
