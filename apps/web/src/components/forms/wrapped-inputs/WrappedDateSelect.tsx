import { useTranslation } from "react-i18next";
import { Controller } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";
import { format } from "date-fns";
import { Calendar as CalendarIcon } from "lucide-react";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Label } from "@/components/ui/label";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { cn } from "@/lib/utils";

interface WrappedDateSelectProps<T extends FieldValues> {
    name: Path<T>;
    control: Control<T>;
    label: string;
    placeholder?: string;
    required?: boolean;
    disabled?: boolean;
    className?: string;
    minDate?: Date;
    maxDate?: Date;
    disablePastDates?: boolean;
}

export const WrappedDateSelect = <T extends FieldValues>({
    name,
    control,
    label,
    placeholder,
    required = false,
    disabled = false,
    className = "",
    minDate,
    maxDate,
    disablePastDates = false,
}: WrappedDateSelectProps<T>) => {
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
                    <Popover>
                        <PopoverTrigger asChild>
                            <Button
                                variant="secondary-dark"
                                id={name}
                                className={cn(
                                    "w-full justify-start text-left font-normal border border-[#D9D9D91A] bg-transparent",
                                    !field.value && "text-muted-foreground",
                                    className,
                                )}
                                disabled={disabled}
                            >
                                <CalendarIcon className="mr-2 h-4 w-4" />
                                {field.value
                                    ? format(field.value, "PPP")
                                    : placeholder || t("common.selectDate")}
                            </Button>
                        </PopoverTrigger>
                        <PopoverContent className="w-auto p-0" align="start">
                            <Calendar
                                mode="single"
                                selected={field.value}
                                onSelect={field.onChange}
                                disabled={(date) => {
                                    if (
                                        disablePastDates &&
                                        date < new Date(new Date().setHours(0, 0, 0, 0))
                                    ) {
                                        return true;
                                    }
                                    if (minDate && date < minDate) {
                                        return true;
                                    }
                                    if (maxDate && date > maxDate) {
                                        return true;
                                    }
                                    return false;
                                }}
                                initialFocus
                            />
                        </PopoverContent>
                    </Popover>
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
