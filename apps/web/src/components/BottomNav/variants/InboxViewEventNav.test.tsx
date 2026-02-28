import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { InboxViewEventNav } from "./InboxViewEventNav";

const mockOnBack = vi.fn();
const mockOnViewEventCallback = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

vi.mock("../stores/inbox", () => ({
    useInboxNavStore: vi.fn(() => ({
        onViewEventCallback: mockOnViewEventCallback,
    })),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => fallback || _key,
    }),
}));

import { useInboxNavStore } from "../stores/inbox";

describe("InboxViewEventNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders back button and view event button", () => {
        render(<InboxViewEventNav />);
        expect(screen.getByRole("button", { name: /go back/i })).toBeInTheDocument();
        expect(screen.getByRole("button", { name: /view event/i })).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<InboxViewEventNav />);
        await user.click(screen.getByRole("button", { name: /go back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });

    it("calls onViewEventCallback when view button is clicked", async () => {
        const user = userEvent.setup();
        render(<InboxViewEventNav />);
        await user.click(screen.getByRole("button", { name: /view event/i }));
        expect(mockOnViewEventCallback).toHaveBeenCalledOnce();
    });

    it("does not crash when callback is undefined", async () => {
        vi.mocked(useInboxNavStore).mockReturnValue({
            onViewEventCallback: undefined,
        } as unknown as ReturnType<typeof useInboxNavStore>);

        const user = userEvent.setup();
        render(<InboxViewEventNav />);
        await user.click(screen.getByRole("button", { name: /view event/i }));
        expect(screen.getByRole("button", { name: /view event/i })).toBeInTheDocument();
    });
});
