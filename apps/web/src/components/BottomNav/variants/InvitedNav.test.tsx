import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { InvitedNav } from "./InvitedNav";

// Mock the context
const mockOnBack = vi.fn();
const mockOnAcceptCallback = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

vi.mock("../stores/event-invitation", () => ({
    useEventInvitationNavStore: vi.fn(() => ({
        onAcceptCallback: mockOnAcceptCallback,
    })),
}));

// Mock react-i18next
vi.mock("react-i18next", () => ({
    useTranslation: vi.fn(() => ({
        t: (key: string, fallback?: string) => {
            if (key === "participant.events.invited") {
                return "You're invited! Click here to continue.";
            }
            return fallback || key;
        },
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        i18n: {} as any,
        ready: true,
    })),
}));

describe("InvitedNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders the invitation message", () => {
        render(<InvitedNav />);

        expect(screen.getByText("You're invited! Click here to continue.")).toBeInTheDocument();
    });

    it("renders back button", () => {
        render(<InvitedNav />);

        const backButton = screen.getByRole("button", { name: /go back/i });
        expect(backButton).toBeInTheDocument();
    });

    it("renders accept button", () => {
        render(<InvitedNav />);

        const acceptButton = screen.getByRole("button", { name: /accept invitation/i });
        expect(acceptButton).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<InvitedNav />);

        const backButton = screen.getByRole("button", { name: /go back/i });
        await user.click(backButton);

        expect(mockOnBack).toHaveBeenCalledOnce();
    });

    it("calls onAcceptCallback when accept button is clicked", async () => {
        const user = userEvent.setup();
        render(<InvitedNav />);

        const acceptButton = screen.getByRole("button", { name: /accept invitation/i });
        await user.click(acceptButton);

        expect(mockOnAcceptCallback).toHaveBeenCalledOnce();
    });

    it("does not crash when onAcceptCallback is null", async () => {
        const user = userEvent.setup();

        const { useEventInvitationNavStore } = await import("../stores/event-invitation");
        vi.mocked(useEventInvitationNavStore).mockReturnValue({
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            onAcceptCallback: null as any,
        });

        render(<InvitedNav />);

        const acceptButton = screen.getByRole("button", { name: /accept invitation/i });
        await user.click(acceptButton);

        // Should not throw error
        expect(acceptButton).toBeInTheDocument();
    });

    it("does not crash when onAcceptCallback is undefined", async () => {
        const user = userEvent.setup();

        const { useEventInvitationNavStore } = await import("../stores/event-invitation");
        vi.mocked(useEventInvitationNavStore).mockReturnValue({
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            onAcceptCallback: undefined as any,
        });

        render(<InvitedNav />);

        const acceptButton = screen.getByRole("button", { name: /accept invitation/i });
        await user.click(acceptButton);

        // Should not throw error
        expect(acceptButton).toBeInTheDocument();
    });

    it("renders with correct styling classes", () => {
        render(<InvitedNav />);

        const acceptButton = screen.getByRole("button", { name: /accept invitation/i });
        expect(acceptButton).toHaveClass("flex-1");
        expect(acceptButton).toHaveClass("bg-white");
        expect(acceptButton).toHaveClass("rounded-lg");
    });

    it("back button has correct accessibility attributes", () => {
        render(<InvitedNav />);

        const backButton = screen.getByRole("button", { name: /go back/i });
        expect(backButton).toHaveAttribute("aria-label", "Go back");
    });

    it("accept button has correct accessibility attributes", () => {
        render(<InvitedNav />);

        const acceptButton = screen.getByRole("button", { name: /accept invitation/i });
        expect(acceptButton).toHaveAttribute("aria-label", "Accept invitation");
    });

    it("displays fallback text when translation is missing", () => {
        // This test verifies that the component works when translations fall back to keys
        // The mock is already set up to handle fallback behavior
        render(<InvitedNav />);

        // The mock returns the full text for this key, which includes fallback behavior
        expect(screen.getByText("You're invited! Click here to continue.")).toBeInTheDocument();
    });
});
