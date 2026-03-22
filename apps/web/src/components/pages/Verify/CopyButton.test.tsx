import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { CopyButton } from "./CopyButton";

describe("CopyButton", () => {
    it("renders a copy button", () => {
        render(<CopyButton value="test-value" />);
        expect(screen.getByRole("button", { name: /copy/i })).toBeInTheDocument();
    });

    it("copies value to clipboard on click", async () => {
        render(<CopyButton value="hello-world" />);
        fireEvent.click(screen.getByRole("button", { name: /copy/i }));
        await waitFor(() => {
            expect(navigator.clipboard.writeText).toHaveBeenCalledWith("hello-world");
        });
    });

    it("shows check icon after copying", async () => {
        vi.mocked(navigator.clipboard.writeText).mockResolvedValue(undefined);
        render(<CopyButton value="hello-world" />);
        fireEvent.click(screen.getByRole("button", { name: /copy/i }));
        await waitFor(() => {
            expect(document.querySelector('[data-testid="mock-check"]')).toBeInTheDocument();
        });
    });
});
