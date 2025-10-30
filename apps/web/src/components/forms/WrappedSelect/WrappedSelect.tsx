import { Controller, type Control as RHFControl } from "react-hook-form";
import { Typography } from "@/components/typography/typography";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

type Control = RHFControl<Record<string, unknown>>;

interface WrappedSelectProps {
    // Form control
    control: Control;
    name: string;

    // Label and description
    label: string;
    description?: string;
    htmlFor?: string;

    // Select options
    options: Array<{
        value: string;
        label: string;
    }>;

    // State
    disabled?: boolean;
    placeholder?: string;

    // Type casting for the value
    valueAs?: (value: string) => unknown;

    // Layout
    containerClassName?: string;
    labelClassName?: string;
    selectClassName?: string;
    descriptionClassName?: string;
}

export function WrappedSelect({
    control,
    name,
    label,
    description,
    htmlFor,
    options,
    disabled = false,
    placeholder,
    valueAs,
    containerClassName = "",
    labelClassName = "text-sm font-medium",
    selectClassName = "",
    descriptionClassName = "text-xs text-muted-foreground",
}: WrappedSelectProps) {
    return (
        <div className={`space-y-2 ${containerClassName}`}>
            <Label htmlFor={htmlFor}>
                <Typography variant="text" tag="span" className={labelClassName}>
                    {label}
                </Typography>
            </Label>
            <Controller
                name={name}
                control={control}
                render={({ field }) => (
                    <Select
                        value={field.value}
                        onValueChange={(value) => field.onChange(valueAs ? valueAs(value) : value)}
                        disabled={disabled}
                    >
                        <SelectTrigger className={selectClassName}>
                            <SelectValue placeholder={placeholder} />
                        </SelectTrigger>
                        <SelectContent>
                            {options.map((option: { value: string; label: string }) => (
                                <SelectItem key={option.value} value={option.value}>
                                    {option.label}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                )}
            />
            {description && (
                <Typography variant="text" tag="p" className={descriptionClassName}>
                    {description}
                </Typography>
            )}
        </div>
    );
}
