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
    maxLength?: number;
    showCharCount?: boolean;
}

export const StyledFormInput = <T extends FieldValues>({
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
    maxLength,
    showCharCount = false,
}: WrappedInputProps<T>) => {
    const { t } = useTranslation();

    return (
        <Controller
            name={name}
            control={control}
            render={({ field, fieldState: { error } }) => {
                const currentLength = String(field.value ?? "").length;
                const isNearLimit = maxLength && currentLength > maxLength * 0.8;
                const isOverLimit = maxLength && currentLength > maxLength;

                return (
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
                            className={`!border !border-primary min-h-12 ${className}`}
                            min={min}
                            max={max}
                            step={step}
                            maxLength={maxLength}
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
                        {showCharCount && maxLength && (
                            <Typography
                                variant="text"
                                tag="p"
                                className={`text-xs text-right ${
                                    isOverLimit
                                        ? "text-destructive font-medium"
                                        : isNearLimit
                                          ? "text-yellow-600 dark:text-yellow-500"
                                          : "text-muted-foreground"
                                }`}
                            >
                                {currentLength} / {maxLength} {t("common.characters")}
                            </Typography>
                        )}
                    </div>
                );
            }}
        />
    );
};
