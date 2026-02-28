import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { InvitationRequiredNav } from "./InvitationRequiredNav";

const mockOnBack = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

describe("InvitationRequiredNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders back button", () => {
        render(<InvitationRequiredNav />);
        expect(screen.getByRole("button", { name: /go back/i })).toBeInTheDocument();
    });

    it("renders invitation required message", () => {
        render(<InvitationRequiredNav />);
        expect(screen.getByText("participant.events.invitationRequired")).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<InvitationRequiredNav />);
        await user.click(screen.getByRole("button", { name: /go back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });
});
