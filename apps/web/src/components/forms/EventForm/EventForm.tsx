import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslation } from "react-i18next";
import { ChevronRight, ChevronLeft, Trash } from "lucide-react";
import { createEventFormSchema, type EventFormData } from "@/lib/schemas/eventFormSchema";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { GoogleMapsEmbed } from "@/components/ui/google-maps-embed";
import {
    WrappedInput,
    WrappedTextarea,
    WrappedDateSelect,
    WrappedInputFile,
} from "@/components/forms/wrapped-inputs";
import ConfirmModal from "@/components/ConfirmModal";
import { PasswordPinModal } from "@/components/ui/password-pin-modal";
import type { GetEventContractByEventIdData } from "@decm/api";

interface EventFormProps {
    defaultValues?: Partial<EventFormData>;
    onSubmit: (data: EventFormData, hostPassword: string) => void | Promise<void>;
    onDelete?: (hostPassword: string) => void | Promise<void>;
    isLoading?: boolean;
    mode?: "create" | "edit";
    previewBannerUrl?: string;
    previewIconUrl?: string;
    eventContract?: GetEventContractByEventIdData;
}

export const EventForm = ({
    defaultValues,
    onSubmit,
    onDelete,
    isLoading = false,
    mode = "create",
    previewBannerUrl,
    previewIconUrl,
    eventContract,
}: EventFormProps) => {
    const { t } = useTranslation();
    const [showPasswordModal, setShowPasswordModal] = useState(false);
    const [showDeleteModal, setShowDeleteModal] = useState(false);
    const [currentStep, setCurrentStep] = useState(1);
    const totalSteps = 2;

    // Create schema based on mode
    const schema = createEventFormSchema(mode);

    const { handleSubmit, control, trigger, watch, getValues } = useForm<EventFormData>({
        resolver: zodResolver(schema),
        defaultValues: {
            name: defaultValues?.name || "",
            shortDescription: defaultValues?.shortDescription || "",
            description: defaultValues?.description || "",
            eventBanner: defaultValues?.eventBanner || null,
            eventIcon: defaultValues?.eventIcon || null,
            startDate: defaultValues?.startDate || undefined,
            endDate: defaultValues?.endDate || undefined,
            seatsCount: defaultValues?.seatsCount || undefined,
            contactNumber: defaultValues?.contactNumber || "",
            contactAddress: defaultValues?.contactAddress || "",
            location: defaultValues?.location || "",
            googleMapQuery: defaultValues?.googleMapQuery || "",
        },
    });

    const googleMapQuery = watch("googleMapQuery");

    const handleFormSubmit = async () => {
        setShowPasswordModal(true);
    };

    const onPasswordVerificationSuccess = async (result: {
        type: "pin" | "password";
        value: string;
    }) => {
        const data = getValues();
        const eventBanner = data.eventBanner ?? null;
        const eventIcon = data.eventIcon ?? null;

        await onSubmit(
            {
                ...data,
                eventBanner: eventBanner,
                eventIcon: eventIcon,
            },
            result.value,
        );
    };

    const handleNext = async (e?: React.MouseEvent<HTMLButtonElement>) => {
        // Prevent form submission if this is called from a button click
        if (e) {
            e.preventDefault();
        }

        // Validate current step fields
        const fieldsToValidate: Array<keyof EventFormData> = [];

        if (currentStep === 1) {
            fieldsToValidate.push(
                "name",
                "shortDescription",
                "description",
                "startDate",
                "endDate",
                "seatsCount",
            );
        }

        // Only validate banner and icon in create mode
        if (mode === "create") {
            fieldsToValidate.push("eventBanner", "eventIcon");
        }

        const isValid = await trigger(fieldsToValidate);

        if (isValid && currentStep < totalSteps) {
            setCurrentStep((prev) => prev + 1);
        }
    };

    const handlePrevious = () => {
        if (currentStep > 1) {
            setCurrentStep((prev) => prev - 1);
        }
    };

    const renderStepIndicator = () => (
        <div className="flex items-center justify-center mb-8 space-x-2">
            {Array.from({ length: totalSteps }).map((_, index) => {
                const step = index + 1;
                return (
                    <div key={step} className="flex items-center">
                        <div
                            className={`flex items-center justify-center w-10 h-10 rounded-full border-2 transition-colors ${
                                currentStep === step
                                    ? "border-primary bg-primary text-primary-foreground"
                                    : currentStep > step
                                      ? "border-primary bg-primary/20 text-primary"
                                      : "border-muted-foreground text-muted-foreground"
                            }`}
                        >
                            <Typography variant="text" tag="span" className="font-semibold">
                                {step}
                            </Typography>
                        </div>
                        {step < totalSteps && (
                            <div
                                className={`w-16 h-1 mx-2 ${
                                    currentStep > step ? "bg-primary" : "bg-muted"
                                }`}
                            />
                        )}
                    </div>
                );
            })}
        </div>
    );

    const renderStepTitle = () => {
        const titles = {
            1: t("events.form.stepEventInfo"),
            2: t("events.form.stepVenueInfo"),
        };

        return (
            <div className="mb-6">
                <Typography variant="header" tag="h2" className="text-2xl font-bold">
                    {titles[currentStep as keyof typeof titles]}
                </Typography>
            </div>
        );
    };

    return (
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-6">
            {renderStepIndicator()}
            {renderStepTitle()}

            {/* Step 1: Event Information */}
            {currentStep === 1 && (
                <div className="space-y-6">
                    {/* Event Banner & Icon */}
                    <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
                        <div className="col-span-2 lg:col-span-3">
                            <WrappedInputFile
                                name="eventBanner"
                                control={control}
                                label={t("events.form.eventBanner")}
                                accept="image/jpeg,image/jpg,image/png,image/webp"
                                disabled={isLoading}
                                required={mode === "create"}
                                previewClassName="object-cover"
                                previewUrl={previewBannerUrl}
                            />
                        </div>
                        <div>
                            <WrappedInputFile
                                name="eventIcon"
                                control={control}
                                label={t("events.form.eventIcon")}
                                accept="image/jpeg,image/jpg,image/png,image/webp"
                                disabled={isLoading}
                                required={mode === "create"}
                                previewUrl={previewIconUrl}
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

                    {/* Short Description */}
                    <WrappedInput
                        name="shortDescription"
                        control={control}
                        label={t("events.form.shortDescription")}
                        placeholder={t("events.form.shortDescriptionPlaceholder")}
                        required
                        disabled={isLoading}
                        maxLength={255}
                        showCharCount
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
                        <WrappedDateSelect
                            name="startDate"
                            control={control}
                            label={t("events.form.startDate")}
                            placeholder={t("events.form.startDatePlaceholder")}
                            required
                            disabled={isLoading}
                        />
                        <WrappedDateSelect
                            name="endDate"
                            control={control}
                            label={t("events.form.endDate")}
                            placeholder={t("events.form.endDatePlaceholder")}
                            required
                            disabled={isLoading}
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
                </div>
            )}

            {/* Step 2: Venue Information */}
            {currentStep === 2 && (
                <>
                    <div className="space-y-6">
                        {/* Contact Number */}
                        <WrappedInput
                            name="contactNumber"
                            control={control}
                            label={t("events.form.contactNumber")}
                            placeholder={t("events.form.contactNumberPlaceholder")}
                            type="tel"
                            required
                            disabled={isLoading}
                        />

                        {/* Contact Address */}
                        <WrappedTextarea
                            name="contactAddress"
                            control={control}
                            label={t("events.form.contactAddress")}
                            placeholder={t("events.form.contactAddressPlaceholder")}
                            required
                            disabled={isLoading}
                        />
                    </div>

                    <div className="space-y-6">
                        {/* Location Name */}
                        <WrappedInput
                            name="location"
                            control={control}
                            label={t("events.form.location")}
                            placeholder={t("events.form.locationPlaceholder")}
                            required
                            disabled={isLoading}
                        />

                        {/* Google Maps Query */}
                        <WrappedInput
                            name="googleMapQuery"
                            control={control}
                            label={t("events.form.googleMapQuery")}
                            placeholder={t("events.form.googleMapQueryPlaceholder")}
                            required
                            disabled={isLoading}
                        />

                        {/* Google Maps Embed */}
                        <div className="space-y-2">
                            <Label>
                                <Typography
                                    variant="text"
                                    tag="span"
                                    className="text-sm font-medium"
                                >
                                    {t("events.location")}
                                </Typography>
                            </Label>
                            <GoogleMapsEmbed query={googleMapQuery || ""} height="400px" />
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-xs text-muted-foreground"
                            >
                                {t("events.form.googleMapPlaceholder")}
                            </Typography>
                        </div>
                    </div>
                </>
            )}

            {/* Navigation Buttons */}
            <div className="flex justify-between gap-4 pt-4">
                {/* Previous Button */}
                {currentStep > 1 && (
                    <Button
                        type="button"
                        variant="secondary-dark"
                        size="lg"
                        onClick={handlePrevious}
                        disabled={isLoading}
                        className="min-w-[150px]"
                    >
                        <ChevronLeft className="h-4 w-4 mr-2" />
                        <Typography variant="text" tag="span" className="font-medium">
                            {t("events.form.previous")}
                        </Typography>
                    </Button>
                )}

                {/* Spacer */}
                {currentStep === 1 && <div />}

                {mode === "edit" && (
                    <ConfirmModal
                        title={t("events.form.deleteEvent")}
                        message={t("events.form.deleteEventMessage")}
                        onConfirm={() => setShowDeleteModal(true)}
                        onCancel={() => {}}
                        cancelText={t("events.form.cancel")}
                        confirmText={t("events.form.delete")}
                    >
                        <Button
                            type="button"
                            variant="ghost"
                            size="lg"
                            onClick={() => {}}
                            disabled={isLoading}
                            className="min-w-[150px]"
                        >
                            <Trash className="h-4 w-4 mr-2" />
                            <Typography variant="text" tag="span" className="font-medium">
                                {t("events.form.delete")}
                            </Typography>
                        </Button>
                    </ConfirmModal>
                )}

                {/* Next or Submit Button */}
                {currentStep === 1 ? (
                    <Button
                        type="button"
                        variant="primary"
                        size="lg"
                        onClick={(e) => handleNext(e)}
                        disabled={isLoading}
                        className="min-w-[150px]"
                    >
                        <Typography variant="text" tag="span" className="font-medium">
                            {t("events.form.next")}
                        </Typography>
                        <ChevronRight className="h-4 w-4 ml-2" />
                    </Button>
                ) : (
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
                )}
            </div>

            <PasswordPinModal
                isOpen={showPasswordModal}
                onClose={() => setShowPasswordModal(false)}
                onSuccess={onPasswordVerificationSuccess}
                showSigningDetails
                signingDetails={{
                    details:
                        mode === "create"
                            ? "You're about to create a new event. Please confirm to proceed."
                            : "You're about to update the event. Please confirm to proceed.",
                    contractAddress: eventContract?.event_contract_address ?? undefined,
                    transactionType: mode === "create" ? "Create Event" : "Update Event",
                }}
            />

            <PasswordPinModal
                isOpen={showDeleteModal}
                onClose={() => setShowDeleteModal(false)}
                onSuccess={(result) => onDelete?.(result.value)}
                showSigningDetails
                signingDetails={{
                    details: "You're about to delete the event. Please confirm to proceed.",
                    transactionType: "Delete Event",
                    contractAddress: eventContract?.event_contract_address ?? undefined,
                }}
            />
        </form>
    );
};
