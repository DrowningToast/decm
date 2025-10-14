import { useTranslation } from "react-i18next";
import { Controller } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";
import { Typography } from "@/components/typography/typography";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";

interface WrappedTextareaProps<T extends FieldValues> {
    name: Path<T>;
    control: Control<T>;
    label: string;
    placeholder?: string;
    required?: boolean;
    disabled?: boolean;
    className?: string;
    rows?: number;
    maxLength?: number;
}

export const WrappedTextarea = <T extends FieldValues>({
    name,
    control,
    label,
    placeholder,
    required = false,
    disabled = false,
    className = "",
    rows,
    maxLength,
}: WrappedTextareaProps<T>) => {
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
                    <Textarea
                        {...field}
                        id={name}
                        placeholder={placeholder}
                        disabled={disabled}
                        aria-invalid={!!error}
                        className={`!border !border-[#D9D9D91A] ${className}`}
                        rows={rows}
                        maxLength={maxLength}
                        value={field.value ?? ""}
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
