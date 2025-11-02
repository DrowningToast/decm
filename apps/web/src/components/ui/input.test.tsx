import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Input } from "./input";

describe("Input", () => {
    it("should render input element", () => {
        render(<Input />);
        expect(screen.getByRole("textbox")).toBeInTheDocument();
    });

    it("should accept text input", async () => {
        const user = userEvent.setup();
        render(<Input placeholder="Enter text" />);

        const input = screen.getByPlaceholderText("Enter text");
        await user.type(input, "Hello World");

        expect(input).toHaveValue("Hello World");
    });

    it("should render with different input types", () => {
        const { rerender } = render(<Input type="text" data-testid="input" />);
        expect(screen.getByTestId("input")).toHaveAttribute("type", "text");

        rerender(<Input type="password" data-testid="input" />);
        expect(screen.getByTestId("input")).toHaveAttribute("type", "password");

        rerender(<Input type="email" data-testid="input" />);
        expect(screen.getByTestId("input")).toHaveAttribute("type", "email");

        rerender(<Input type="number" data-testid="input" />);
        expect(screen.getByTestId("input")).toHaveAttribute("type", "number");
    });

    it("should apply placeholder", () => {
        render(<Input placeholder="Enter your name" />);
        expect(screen.getByPlaceholderText("Enter your name")).toBeInTheDocument();
    });

    it("should handle onChange events", async () => {
        const handleChange = vi.fn();
        const user = userEvent.setup();

        render(<Input onChange={handleChange} />);

        const input = screen.getByRole("textbox");
        await user.type(input, "Test");

        expect(handleChange).toHaveBeenCalledTimes(4); // Once for each character
    });

    it("should be disabled when disabled prop is true", async () => {
        const handleChange = vi.fn();
        const user = userEvent.setup();

        render(<Input disabled onChange={handleChange} />);

        const input = screen.getByRole("textbox");
        expect(input).toBeDisabled();

        await user.type(input, "Test");
        expect(handleChange).not.toHaveBeenCalled();
    });

    it("should apply custom className", () => {
        render(<Input className="custom-input" data-testid="input" />);
        expect(screen.getByTestId("input")).toHaveClass("custom-input");
    });

    it("should have data-slot attribute", () => {
        render(<Input data-testid="input" />);
        expect(screen.getByTestId("input")).toHaveAttribute("data-slot", "input");
    });

    it("should apply default styling classes", () => {
        render(<Input data-testid="input" />);
        const input = screen.getByTestId("input");

        expect(input).toHaveClass("rounded-md");
        expect(input).toHaveClass("border");
        expect(input).toHaveClass("px-3");
    });

    it("should show invalid state with aria-invalid", () => {
        render(<Input aria-invalid data-testid="input" />);
        const input = screen.getByTestId("input");

        expect(input).toHaveAttribute("aria-invalid", "true");
    });

    it("should accept defaultValue", () => {
        render(<Input defaultValue="Default text" data-testid="input" />);
        expect(screen.getByTestId("input")).toHaveValue("Default text");
    });

    it("should support controlled input", async () => {
        const user = userEvent.setup();

        const ControlledInput = () => {
            const [value, setValue] = React.useState("");
            return <Input value={value} onChange={(e) => setValue(e.target.value)} />;
        };

        render(<ControlledInput />);

        const input = screen.getByRole("textbox");
        await user.type(input, "Controlled");

        expect(input).toHaveValue("Controlled");
    });

    it("should pass through other HTML attributes", () => {
        render(
            <Input
                data-testid="input"
                maxLength={10}
                minLength={2}
                required
                readOnly
                autoComplete="off"
            />,
        );

        const input = screen.getByTestId("input");

        expect(input).toHaveAttribute("maxLength", "10");
        expect(input).toHaveAttribute("minLength", "2");
        expect(input).toBeRequired();
        expect(input).toHaveAttribute("readOnly");
        expect(input).toHaveAttribute("autoComplete", "off");
    });
});

// Add React import for controlled input test
import React from "react";
