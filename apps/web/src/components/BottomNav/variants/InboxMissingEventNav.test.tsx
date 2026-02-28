import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { InboxMissingEventNav } from "./InboxMissingEventNav";

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

describe("InboxMissingEventNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders back button", () => {
        render(<InboxMissingEventNav />);
        expect(screen.getByRole("button", { name: /go back/i })).toBeInTheDocument();
    });

    it("renders missing event message", () => {
        render(<InboxMissingEventNav />);
        expect(screen.getByText("Must be in the event first")).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<InboxMissingEventNav />);
        await user.click(screen.getByRole("button", { name: /go back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });
});
