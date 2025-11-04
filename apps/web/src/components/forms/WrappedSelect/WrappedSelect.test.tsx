import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { useForm } from "react-hook-form";
import { WrappedSelect } from "./WrappedSelect";

// Mock components
vi.mock("@/components/typography/typography", () => ({
    Typography: ({ children, ...props }: React.ComponentProps<"span">) => (
        <span {...props}>{children}</span>
    ),
}));

vi.mock("@/components/ui/label", () => ({
    Label: ({ children, ...props }: React.ComponentProps<"label">) => (
        <label {...props}>{children}</label>
    ),
}));

vi.mock("@/components/ui/select", () => ({
    Select: ({
        children,
        value,
        disabled,
    }: React.PropsWithChildren<{
        value?: string;
        onValueChange?: (value: string) => void;
        disabled?: boolean;
    }>) => (
        <div data-testid="select" data-value={value} data-disabled={disabled}>
            {children}
        </div>
    ),
    SelectTrigger: ({ children, ...props }: React.ComponentProps<"button">) => (
        <button data-testid="select-trigger" {...props}>
            {children}
        </button>
    ),
    SelectValue: ({ placeholder }: { placeholder?: string }) => (
        <span data-testid="select-value">{placeholder || "Select..."}</span>
    ),
    SelectContent: ({ children }: React.PropsWithChildren) => (
        <div data-testid="select-content">{children}</div>
    ),
    SelectItem: ({
        children,
        value,
        onClick,
    }: React.PropsWithChildren<{ value?: string; onClick?: () => void }>) => (
        <div data-testid={`select-item-${value}`} onClick={onClick}>
            {children}
        </div>
    ),
}));

function TestSelectForm({ defaultValue = "" }: { defaultValue?: string }) {
    const { control } = useForm({
        defaultValues: {
            category: defaultValue,
        },
    });

    const options = [
        { value: "option1", label: "Option 1" },
        { value: "option2", label: "Option 2" },
        { value: "option3", label: "Option 3" },
    ];

    return (
        <WrappedSelect
            control={control}
            name="category"
            label="Select Category"
            options={options}
            placeholder="Choose an option"
        />
    );
}

describe("WrappedSelect Component", () => {
    it("should render without crashing", () => {
        render(<TestSelectForm />);
        expect(screen.getByText("Select Category")).toBeInTheDocument();
    });

    it("should render label", () => {
        render(<TestSelectForm />);
        expect(screen.getByText("Select Category")).toBeInTheDocument();
    });

    it("should render select element", () => {
        render(<TestSelectForm />);
        expect(screen.getByTestId("select")).toBeInTheDocument();
    });

    it("should render placeholder text", () => {
        render(<TestSelectForm />);
        expect(screen.getByTestId("select-value")).toHaveTextContent("Choose an option");
    });

    it("should render all options", () => {
        render(<TestSelectForm />);

        expect(screen.getByText("Option 1")).toBeInTheDocument();
        expect(screen.getByText("Option 2")).toBeInTheDocument();
        expect(screen.getByText("Option 3")).toBeInTheDocument();
    });

    it("should handle disabled state", () => {
        function DisabledTest() {
            const { control } = useForm({
                defaultValues: { field: "" },
            });
            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Disabled Select"
                    options={[
                        { value: "opt1", label: "Option 1" },
                        { value: "opt2", label: "Option 2" },
                    ]}
                    disabled={true}
                />
            );
        }

        render(<DisabledTest />);

        const select = screen.getByTestId("select");
        expect(select).toHaveAttribute("data-disabled", "true");
    });

    it("should apply custom containerClassName", () => {
        function ContainerTest() {
            const { control } = useForm();
            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Custom Container"
                    options={[{ value: "opt1", label: "Option 1" }]}
                    containerClassName="custom-container"
                />
            );
        }

        const { container } = render(<ContainerTest />);

        const wrapper = container.querySelector(".custom-container");
        expect(wrapper).toBeInTheDocument();
    });

    it("should apply custom labelClassName", () => {
        function LabelClassTest() {
            const { control } = useForm();
            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Custom Label"
                    options={[{ value: "opt1", label: "Option 1" }]}
                    labelClassName="custom-label-class"
                />
            );
        }

        render(<LabelClassTest />);

        const label = screen.getByText("Custom Label");
        expect(label).toHaveClass("custom-label-class");
    });

    it("should apply custom selectClassName", () => {
        function SelectClassTest() {
            const { control } = useForm();
            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Custom Select"
                    options={[{ value: "opt1", label: "Option 1" }]}
                    selectClassName="custom-select"
                />
            );
        }

        const { container } = render(<SelectClassTest />);

        // The select trigger would have the custom class
        expect(container.querySelector(".custom-select")).toBeInTheDocument();
    });

    it("should apply custom descriptionClassName to description", () => {
        function DescClassTest() {
            const { control } = useForm();
            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Field"
                    options={[{ value: "opt1", label: "Option 1" }]}
                    description="This is a description"
                    descriptionClassName="custom-desc-class"
                />
            );
        }

        render(<DescClassTest />);

        const description = screen.getByText("This is a description");
        expect(description).toHaveClass("custom-desc-class");
    });

    it("should render description when provided", () => {
        function DescTest() {
            const { control } = useForm();
            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Field"
                    options={[{ value: "opt1", label: "Option 1" }]}
                    description="This is helpful text"
                />
            );
        }

        render(<DescTest />);

        expect(screen.getByText("This is helpful text")).toBeInTheDocument();
    });

    it("should not render description when not provided", () => {
        render(<TestSelectForm />);

        // Should not have any description
        expect(screen.queryByText(/This is/)).not.toBeInTheDocument();
    });

    it("should handle empty options array", () => {
        function EmptyOptionsTest() {
            const { control } = useForm();
            return (
                <WrappedSelect control={control} name="field" label="Empty Options" options={[]} />
            );
        }

        render(<EmptyOptionsTest />);

        expect(screen.getByText("Empty Options")).toBeInTheDocument();
    });

    it("should handle default value", () => {
        function DefaultValueTest() {
            const { control } = useForm({
                defaultValues: { field: "default-value" },
            });
            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Select"
                    options={[
                        { value: "default-value", label: "Default Option" },
                        { value: "other", label: "Other Option" },
                    ]}
                />
            );
        }

        const { container } = render(<DefaultValueTest />);

        const select = container.querySelector('[data-testid="select"]');
        expect(select).toHaveAttribute("data-value", "default-value");
    });

    it("should support valueAs transformer", () => {
        function ValueAsTest() {
            const { control } = useForm({
                defaultValues: { field: "1" },
            });

            const valueAsNumber = (value: string) => parseInt(value, 10);

            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Number Select"
                    options={[
                        { value: "1", label: "One" },
                        { value: "2", label: "Two" },
                        { value: "3", label: "Three" },
                    ]}
                    valueAs={valueAsNumber}
                />
            );
        }

        render(<ValueAsTest />);

        expect(screen.getByText("Number Select")).toBeInTheDocument();
    });

    it("should handle htmlFor attribute", () => {
        function HtmlForTest() {
            const { control } = useForm();
            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Field"
                    options={[{ value: "opt1", label: "Option 1" }]}
                    htmlFor="custom-id"
                />
            );
        }

        const { container } = render(<HtmlForTest />);

        const label = container.querySelector("label");
        expect(label).toHaveAttribute("htmlFor", "custom-id");
    });

    it("should render multiple options correctly", () => {
        function ManyOptionsTest() {
            const { control } = useForm();

            const manyOptions = Array.from({ length: 10 }, (_, i) => ({
                value: `opt${i}`,
                label: `Option ${i}`,
            }));

            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Many Options"
                    options={manyOptions}
                />
            );
        }

        render(<ManyOptionsTest />);

        expect(screen.getByText("Option 0")).toBeInTheDocument();
        expect(screen.getByText("Option 9")).toBeInTheDocument();
    });

    it("should apply default labelClassName when not overridden", () => {
        function DefaultLabelTest() {
            const { control } = useForm();
            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Default Label"
                    options={[{ value: "opt1", label: "Option 1" }]}
                />
            );
        }

        render(<DefaultLabelTest />);

        const label = screen.getByText("Default Label");
        expect(label).toHaveClass("text-sm");
        expect(label).toHaveClass("font-medium");
    });

    it("should apply default descriptionClassName when not overridden", () => {
        function DefaultDescTest() {
            const { control } = useForm();
            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Field"
                    options={[{ value: "opt1", label: "Option 1" }]}
                    description="Default description"
                />
            );
        }

        render(<DefaultDescTest />);

        const description = screen.getByText("Default description");
        expect(description).toHaveClass("text-xs");
        expect(description).toHaveClass("text-muted-foreground");
    });

    it("should handle options with special characters", () => {
        function SpecialCharsTest() {
            const { control } = useForm();

            const optionsWithSpecialChars = [
                { value: "opt&special", label: "Option & Special" },
                { value: "opt<tag>", label: "Option <tag>" },
            ];

            return (
                <WrappedSelect
                    control={control}
                    name="field"
                    label="Special Characters"
                    options={optionsWithSpecialChars}
                />
            );
        }

        render(<SpecialCharsTest />);

        expect(screen.getByText("Option & Special")).toBeInTheDocument();
        expect(screen.getByText("Option <tag>")).toBeInTheDocument();
    });

    it("should work with controlled component pattern", () => {
        function ControlledTest() {
            const { control } = useForm({
                defaultValues: { category: "option2" },
            });

            return (
                <WrappedSelect
                    control={control}
                    name="category"
                    label="Category"
                    options={[
                        { value: "option1", label: "Option 1" },
                        { value: "option2", label: "Option 2" },
                    ]}
                />
            );
        }

        const { container } = render(<ControlledTest />);

        const select = container.querySelector('[data-testid="select"]');
        expect(select).toHaveAttribute("data-value", "option2");
    });
});
