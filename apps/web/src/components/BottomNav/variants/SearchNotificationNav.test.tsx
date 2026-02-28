import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { SearchNotificationNav } from "./SearchNotificationNav";

const mockOnBack = vi.fn();
const mockSetSearchQuery = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

vi.mock("../stores/notifications", () => ({
    useSearchNotificationNavStore: vi.fn(() => ({
        searchQuery: "",
        setSearchQuery: mockSetSearchQuery,
    })),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => fallback || _key,
    }),
}));

describe("SearchNotificationNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders back button and search input", () => {
        render(<SearchNotificationNav />);
        expect(screen.getByRole("button", { name: /go back/i })).toBeInTheDocument();
        expect(screen.getByRole("textbox")).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<SearchNotificationNav />);
        await user.click(screen.getByRole("button", { name: /go back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });

    it("calls setSearchQuery on input change", async () => {
        const user = userEvent.setup();
        render(<SearchNotificationNav />);
        await user.type(screen.getByRole("textbox"), "a");
        expect(mockSetSearchQuery).toHaveBeenCalled();
    });
});
