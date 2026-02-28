import React from "react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { EventPasswordNav } from "./EventPasswordNav";

const mockOnBack = vi.fn();
const mockSetPassword = vi.fn();
const mockOnSubmitCallback = vi.fn();

vi.mock("../context", () => ({
    useBottomContainerContext: () => ({
        onBack: mockOnBack,
        className: "",
    }),
}));

vi.mock("../stores/event-password", () => ({
    useEventPasswordNavStore: vi.fn(() => ({
        password: "",
        setPassword: mockSetPassword,
        onSubmitCallback: mockOnSubmitCallback,
    })),
}));

vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (_key: string, fallback?: string) => fallback || _key,
    }),
}));

import { useEventPasswordNavStore } from "../stores/event-password";

describe("EventPasswordNav", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders back button, password input, and submit button", () => {
        render(<EventPasswordNav />);
        expect(screen.getByRole("button", { name: /go back/i })).toBeInTheDocument();
        expect(screen.getByPlaceholderText(/passwordPlaceholder/i)).toBeInTheDocument();
        expect(screen.getByRole("button", { name: /submit password/i })).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", async () => {
        const user = userEvent.setup();
        render(<EventPasswordNav />);
        await user.click(screen.getByRole("button", { name: /go back/i }));
        expect(mockOnBack).toHaveBeenCalledOnce();
    });

    it("calls setPassword on input change", async () => {
        const user = userEvent.setup();
        render(<EventPasswordNav />);
        const input = screen.getByPlaceholderText(/passwordPlaceholder/i);
        await user.type(input, "a");
        expect(mockSetPassword).toHaveBeenCalled();
    });

    it("calls onSubmitCallback with password on submit", async () => {
        vi.mocked(useEventPasswordNavStore).mockReturnValue({
            password: "secret123",
            setPassword: mockSetPassword,
            onSubmitCallback: mockOnSubmitCallback,
        } as unknown as ReturnType<typeof useEventPasswordNavStore>);

        const user = userEvent.setup();
        render(<EventPasswordNav />);
        await user.click(screen.getByRole("button", { name: /submit password/i }));
        expect(mockOnSubmitCallback).toHaveBeenCalledWith("secret123");
    });

    it("does not crash when onSubmitCallback is undefined", async () => {
        vi.mocked(useEventPasswordNavStore).mockReturnValue({
            password: "",
            setPassword: mockSetPassword,
            onSubmitCallback: undefined,
        } as unknown as ReturnType<typeof useEventPasswordNavStore>);

        const user = userEvent.setup();
        render(<EventPasswordNav />);
        await user.click(screen.getByRole("button", { name: /submit password/i }));
        // Should not throw
        expect(screen.getByRole("button", { name: /submit password/i })).toBeInTheDocument();
    });
});
