import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { useForm } from "react-hook-form";
import { WrappedInput } from "./WrappedInput";

// Mock components
vi.mock("@/components/typography/typography", () => ({
    Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
}));

vi.mock("@/components/ui/input", () => ({
    Input: ({ ...props }: any) => <input data-testid="input" {...props} />,
}));

vi.mock("@/components/ui/label", () => ({
    Label: ({ children, ...props }: any) => <label {...props}>{children}</label>,
}));

// Mock translation
vi.mock("react-i18next", () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

function TestForm() {
    const { control } = useForm({
        defaultValues: {
            testField: "",
            emailField: "",
            numberField: 0,
        },
    });

    return (
        <form>
            <WrappedInput
                name="testField"
                control={control}
                label="Test Label"
                placeholder="Enter text"
            />
        </form>
    );
}

describe("WrappedInput Component", () => {
    it("should render label", () => {
        render(<TestForm />);
        expect(screen.getByText("Test Label")).toBeInTheDocument();
    });

    it("should render input field", () => {
        render(<TestForm />);
        expect(screen.getByTestId("input")).toBeInTheDocument();
    });

    it("should apply placeholder", () => {
        render(<TestForm />);
        const input = screen.getByTestId("input") as HTMLInputElement;
        expect(input).toHaveAttribute("placeholder", "Enter text");
    });

    it("should render required indicator when required is true", () => {
        function RequiredTest() {
            const { control } = useForm();
            return (
                <WrappedInput
                    name="field"
                    control={control}
                    label="Required Field"
                    required={true}
                />
            );
        }

        render(<RequiredTest />);

        // Should render asterisk for required field
        const requiredMarker = screen.getByText("*");
        expect(requiredMarker).toBeInTheDocument();
    });

    it("should not render required indicator when required is false", () => {
        function OptionalTest() {
            const { control } = useForm();
            return (
                <WrappedInput
                    name="field"
                    control={control}
                    label="Optional Field"
                    required={false}
                />
            );
        }

        render(<OptionalTest />);

        // Should not have asterisk
        const asterisk = screen.queryByText("*");
        expect(asterisk).not.toBeInTheDocument();
    });

    it("should handle different input types", () => {
        function TypeTest() {
            const { control } = useForm();
            return <WrappedInput name="field" control={control} label="Email" type="email" />;
        }

        const { rerender } = render(<TypeTest />);

        let input = screen.getByTestId("input") as HTMLInputElement;
        expect(input.type).toBe("email");

        function PasswordTest() {
            const { control } = useForm();
            return <WrappedInput name="field" control={control} label="Password" type="password" />;
        }

        rerender(<PasswordTest />);

        input = screen.getByTestId("input") as HTMLInputElement;
        expect(input.type).toBe("password");
    });

    it("should handle disabled state", () => {
        function DisabledTest() {
            const { control } = useForm();
            return (
                <WrappedInput
                    name="field"
                    control={control}
                    label="Disabled Input"
                    disabled={true}
                />
            );
        }

        render(<DisabledTest />);

        const input = screen.getByTestId("input");
        expect(input).toBeDisabled();
    });

    it("should apply custom className", () => {
        function ClassNameTest() {
            const { control } = useForm();
            return (
                <WrappedInput
                    name="field"
                    control={control}
                    label="Custom Class"
                    className="custom-class"
                />
            );
        }

        render(<ClassNameTest />);

        const input = screen.getByTestId("input");
        expect(input).toHaveClass("custom-class");
    });

    it("should handle maxLength attribute", () => {
        function MaxLengthTest() {
            const { control } = useForm();
            return (
                <WrappedInput
                    name="field"
                    control={control}
                    label="Limited"
                    maxLength={10}
                    showCharCount={true}
                />
            );
        }

        render(<MaxLengthTest />);

        const input = screen.getByTestId("input");
        expect(input).toHaveAttribute("maxLength", "10");
    });

    it("should display character count when showCharCount is true", () => {
        function CharCountTest() {
            const { control } = useForm({
                defaultValues: { field: "hello" },
            });
            return (
                <WrappedInput
                    name="field"
                    control={control}
                    label="Text"
                    maxLength={20}
                    showCharCount={true}
                />
            );
        }

        render(<CharCountTest />);

        // Character count should be displayed
        const charCountText = screen.getByText(/\/\s*20/);
        expect(charCountText).toBeInTheDocument();
    });

    it("should not display character count when showCharCount is false", () => {
        function NoCharCountTest() {
            const { control } = useForm();
            return (
                <WrappedInput
                    name="field"
                    control={control}
                    label="Text"
                    maxLength={20}
                    showCharCount={false}
                />
            );
        }

        render(<NoCharCountTest />);

        // Should not have character count
        const charCount = screen.queryByText(/common\.characters/);
        expect(charCount).not.toBeInTheDocument();
    });

    it("should handle min and max for number inputs", () => {
        function NumberRangeTest() {
            const { control } = useForm();
            return (
                <WrappedInput
                    name="field"
                    control={control}
                    label="Number"
                    type="number"
                    min={0}
                    max={100}
                />
            );
        }

        render(<NumberRangeTest />);

        const input = screen.getByTestId("input");
        expect(input).toHaveAttribute("min", "0");
        expect(input).toHaveAttribute("max", "100");
    });

    it("should handle step for number inputs", () => {
        function StepTest() {
            const { control } = useForm();
            return (
                <WrappedInput
                    name="field"
                    control={control}
                    label="Step Number"
                    type="number"
                    step={0.5}
                />
            );
        }

        render(<StepTest />);

        const input = screen.getByTestId("input");
        expect(input).toHaveAttribute("step", "0.5");
    });

    it("should set input id to field name", () => {
        function IdTest() {
            const { control } = useForm();
            return <WrappedInput name="myField" control={control} label="My Field" />;
        }

        render(<IdTest />);

        const input = screen.getByTestId("input");
        expect(input).toHaveAttribute("id", "myField");
    });

    it("should handle text input default value", () => {
        function DefaultValueTest() {
            const { control } = useForm({
                defaultValues: { field: "default text" },
            });
            return <WrappedInput name="field" control={control} label="Text" />;
        }

        render(<DefaultValueTest />);

        const input = screen.getByTestId("input") as HTMLInputElement;
        expect(input.value).toBe("default text");
    });

    it("should handle empty value gracefully", () => {
        function EmptyValueTest() {
            const { control } = useForm({
                defaultValues: { field: null },
            });
            return <WrappedInput name="field" control={control} label="Text" />;
        }

        render(<EmptyValueTest />);

        const input = screen.getByTestId("input") as HTMLInputElement;
        expect(input.value).toBe("");
    });

    it("should handle undefined value", () => {
        function UndefinedValueTest() {
            const { control } = useForm({
                defaultValues: { field: undefined },
            });
            return <WrappedInput name="field" control={control} label="Text" />;
        }

        render(<UndefinedValueTest />);

        const input = screen.getByTestId("input") as HTMLInputElement;
        expect(input.value).toBe("");
    });

    it("should support all text input types", () => {
        const types: Array<"text" | "email" | "number" | "password" | "tel" | "url"> = [
            "text",
            "email",
            "number",
            "password",
            "tel",
            "url",
        ];

        types.forEach((type) => {
            function TypeVariantTest() {
                const { control } = useForm();
                return (
                    <WrappedInput
                        name="field"
                        control={control}
                        label={`${type} input`}
                        type={type}
                    />
                );
            }

            const { unmount } = render(<TypeVariantTest />);

            const input = screen.getByTestId("input") as HTMLInputElement;
            expect(input.type).toBe(type);

            unmount();
        });
    });

    it("should have aria-invalid attribute for error state", () => {
        function ErrorStateTest() {
            const { control } = useForm();
            return <WrappedInput name="field" control={control} label="Field" />;
        }

        render(<ErrorStateTest />);

        const input = screen.getByTestId("input");
        // aria-invalid should be false by default
        expect(input).toHaveAttribute("aria-invalid", "false");
    });

    it("should be able to input text", () => {
        function InputTest() {
            const { control } = useForm({
                defaultValues: { field: "" },
            });
            return <WrappedInput name="field" control={control} label="Text Input" />;
        }

        render(<InputTest />);

        const input = screen.getByTestId("input") as HTMLInputElement;
        fireEvent.change(input, { target: { value: "test input" } });

        expect(input.value).toBe("test input");
    });
});
