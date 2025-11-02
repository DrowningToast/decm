import {
    Controller,
    type Control as ReactHookFormControl,
    type FieldValues,
    type Path,
} from "react-hook-form";
import { Typography } from "@/components/typography/typography";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

interface WrappedSelectProps<T extends FieldValues = FieldValues> {
    // Form control
    control: ReactHookFormControl<T>;
    name: Path<T>;

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
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    valueAs?: (value: string) => any;

    // Layout
    containerClassName?: string;
    labelClassName?: string;
    selectClassName?: string;
    descriptionClassName?: string;
}

export function WrappedSelect<T extends FieldValues = FieldValues>({
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
}: WrappedSelectProps<T>) {
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
                            {options.map((option) => (
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
