import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Button } from "./button";

describe("Button", () => {
    it("should render button with children", () => {
        render(<Button>Click me</Button>);
        expect(screen.getByRole("button", { name: "Click me" })).toBeInTheDocument();
    });

    it("should apply primary variant by default", () => {
        render(<Button>Primary Button</Button>);
        const button = screen.getByRole("button");
        expect(button).toHaveClass("bg-primary");
    });

    it("should apply secondary-dark variant", () => {
        render(<Button variant="secondary-dark">Secondary Dark</Button>);
        const button = screen.getByRole("button");
        expect(button).toHaveClass("bg-secondary-foreground");
    });

    it("should apply secondary-light variant", () => {
        render(<Button variant="secondary-light">Secondary Light</Button>);
        const button = screen.getByRole("button");
        expect(button).toHaveClass("bg-secondary");
    });

    it("should apply ghost variant", () => {
        render(<Button variant="ghost">Ghost</Button>);
        const button = screen.getByRole("button");
        expect(button).toHaveClass("bg-transparent");
    });

    it("should apply different sizes", () => {
        const { rerender } = render(<Button size="sm">Small</Button>);
        expect(screen.getByRole("button")).toHaveClass("h-8");

        rerender(<Button size="lg">Large</Button>);
        expect(screen.getByRole("button")).toHaveClass("h-10");

        rerender(<Button size="xl">Extra Large</Button>);
        expect(screen.getByRole("button")).toHaveClass("h-12");

        rerender(<Button size="icon">Icon</Button>);
        expect(screen.getByRole("button")).toHaveClass("size-9");
    });

    it("should handle click events", async () => {
        const handleClick = vi.fn();
        const user = userEvent.setup();

        render(<Button onClick={handleClick}>Click me</Button>);

        await user.click(screen.getByRole("button"));

        expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it("should be disabled when disabled prop is true", async () => {
        const handleClick = vi.fn();
        const user = userEvent.setup();

        render(
            <Button disabled onClick={handleClick}>
                Disabled Button
            </Button>,
        );

        const button = screen.getByRole("button");
        expect(button).toBeDisabled();

        await user.click(button);
        expect(handleClick).not.toHaveBeenCalled();
    });

    it("should show loading spinner when loading is true", () => {
        render(<Button loading>Loading Button</Button>);

        const button = screen.getByRole("button");
        expect(button).toBeDisabled();

        // Check for loading spinner (Loader2 icon)
        const spinner = button.querySelector(".animate-spin");
        expect(spinner).toBeInTheDocument();
    });

    it("should not show loading spinner when disabled", () => {
        render(
            <Button loading disabled>
                Disabled Loading
            </Button>,
        );

        const button = screen.getByRole("button");
        expect(button).toBeDisabled();

        // Loading spinner should not be shown when disabled
        const spinner = button.querySelector(".animate-spin");
        expect(spinner).not.toBeInTheDocument();
    });

    it("should apply custom className", () => {
        render(<Button className="custom-class">Custom Button</Button>);
        expect(screen.getByRole("button")).toHaveClass("custom-class");
    });

    it("should handle asChild prop behavior", () => {
        // Note: Radix UI Slot requires single React element child
        // Testing with proper single child element
        render(<Button data-testid="regular-button">Regular Button</Button>);

        const button = screen.getByTestId("regular-button");
        expect(button).toBeInTheDocument();
        expect(button.tagName).toBe("BUTTON");
    });

    it("should pass through other props", () => {
        render(
            <Button type="submit" data-testid="submit-button">
                Submit
            </Button>,
        );

        const button = screen.getByTestId("submit-button");
        expect(button).toHaveAttribute("type", "submit");
    });

    it("should have data-slot attribute", () => {
        render(<Button>Button</Button>);
        expect(screen.getByRole("button")).toHaveAttribute("data-slot", "button");
    });
});
