import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslation } from "react-i18next";
import { eventFormSchema, type EventFormData } from "@/lib/schemas/eventFormSchema";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import {
    WrappedInput,
    WrappedTextarea,
    WrappedDateSelect,
    WrappedInputFile,
} from "@/components/forms/wrapped-inputs";

interface EventFormProps {
    /**
     * Default values for the form (used for editing existing events)
     */
    defaultValues?: Partial<EventFormData>;
    /**
     * Callback function called when form is successfully submitted
     */
    onSubmit: (data: EventFormData) => void | Promise<void>;
    /**
     * Whether the form is in loading/submitting state
     */
    isLoading?: boolean;
    /**
     * Form mode: 'create' or 'edit'
     */
    mode?: "create" | "edit";
}

export const EventForm = ({
    defaultValues,
    onSubmit,
    isLoading = false,
    mode = "create",
}: EventFormProps) => {
    const { t } = useTranslation();

    const { handleSubmit, control } = useForm<EventFormData>({
        resolver: zodResolver(eventFormSchema),
        defaultValues: {
            name: defaultValues?.name || "",
            description: defaultValues?.description || "",
            eventBanner: defaultValues?.eventBanner || undefined,
            eventIcon: defaultValues?.eventIcon || undefined,
            startDate: defaultValues?.startDate || undefined,
            endDate: defaultValues?.endDate || undefined,
            seatsCount: defaultValues?.seatsCount || undefined,
        },
    });

    const handleFormSubmit = async (data: EventFormData) => {
        await onSubmit(data);
    };

    return (
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                {/* Event Banner */}
                <div className="lg:col-span-3">
                    <WrappedInputFile
                        name="eventBanner"
                        control={control}
                        label={t("events.form.eventBanner")}
                        accept="image/jpeg,image/jpg,image/png,image/webp"
                        disabled={isLoading}
                        required
                        previewClassName="object-cover"
                    />
                </div>

                {/* Event Icon */}
                <div>
                    <WrappedInputFile
                        name="eventIcon"
                        control={control}
                        label={t("events.form.eventIcon")}
                        accept="image/jpeg,image/jpg,image/png,image/webp"
                        disabled={isLoading}
                        required
                    />
                </div>
            </div>

            {/* Event Name */}
            <WrappedInput
                name="name"
                control={control}
                label={t("events.form.name")}
                placeholder={t("events.form.namePlaceholder")}
                required
                disabled={isLoading}
            />

            {/* Description */}
            <WrappedTextarea
                name="description"
                control={control}
                label={t("events.form.description")}
                placeholder={t("events.form.descriptionPlaceholder")}
                disabled={isLoading}
            />

            {/* Date Fields */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Start Date */}
                <WrappedDateSelect
                    name="startDate"
                    control={control}
                    label={t("events.form.startDate")}
                    placeholder={t("events.form.startDatePlaceholder")}
                    required
                    disabled={isLoading}
                    disablePastDates
                />

                {/* End Date */}
                <WrappedDateSelect
                    name="endDate"
                    control={control}
                    label={t("events.form.endDate")}
                    placeholder={t("events.form.endDatePlaceholder")}
                    required
                    disabled={isLoading}
                    disablePastDates
                />
            </div>

            {/* Seats Count */}
            <WrappedInput
                name="seatsCount"
                control={control}
                label={t("events.form.seatsCount")}
                placeholder={t("events.form.seatsCountPlaceholder")}
                type="number"
                required
                disabled={isLoading}
                min={1}
                step={1}
            />

            {/* Submit Button */}
            <div className="flex justify-end gap-4 pt-4">
                <Button
                    type="submit"
                    variant="primary"
                    size="lg"
                    disabled={isLoading}
                    className="min-w-[150px]"
                >
                    <Typography variant="text" tag="span" className="font-medium">
                        {isLoading
                            ? t("common.loading")
                            : mode === "create"
                              ? t("events.form.submitCreate")
                              : t("events.form.submitUpdate")}
                    </Typography>
                </Button>
            </div>
        </form>
    );
};
