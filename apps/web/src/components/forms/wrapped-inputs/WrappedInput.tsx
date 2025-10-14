import { useTranslation } from "react-i18next";
import { Controller } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";
import { Typography } from "@/components/typography/typography";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface WrappedInputProps<T extends FieldValues> {
    name: Path<T>;
    control: Control<T>;
    label: string;
    placeholder?: string;
    type?: "text" | "email" | "number" | "password" | "tel" | "url";
    required?: boolean;
    disabled?: boolean;
    className?: string;
    min?: number;
    max?: number;
    step?: number;
}

export const WrappedInput = <T extends FieldValues>({
    name,
    control,
    label,
    placeholder,
    type = "text",
    required = false,
    disabled = false,
    className = "",
    min,
    max,
    step,
}: WrappedInputProps<T>) => {
    const { t } = useTranslation();

    return (
        <Controller
            name={name}
            control={control}
            render={({ field, fieldState: { error } }) => (
                <div className="space-y-2">
                    <Label htmlFor={name}>
                        <Typography variant="text" tag="span" className="text-sm font-medium">
                            {label}
                        </Typography>
                        {required && <span className="text-destructive ml-1">*</span>}
                    </Label>
                    <Input
                        {...field}
                        id={name}
                        type={type}
                        placeholder={placeholder}
                        disabled={disabled}
                        aria-invalid={!!error}
                        className={`!border !border-[#D9D9D91A] ${className}`}
                        min={min}
                        max={max}
                        step={step}
                        value={field.value ?? ""}
                        onChange={(e) => {
                            if (type === "number") {
                                field.onChange(
                                    e.target.value === "" ? undefined : Number(e.target.value),
                                );
                            } else {
                                field.onChange(e.target.value);
                            }
                        }}
                    />
                    {error && (
                        <Typography
                            variant="text"
                            tag="p"
                            className="text-sm text-destructive"
                            role="alert"
                        >
                            {t(error.message as string)}
                        </Typography>
                    )}
                </div>
            )}
        />
    );
};
