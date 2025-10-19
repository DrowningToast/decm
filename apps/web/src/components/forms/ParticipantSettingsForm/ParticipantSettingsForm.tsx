import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslation } from "react-i18next";
import {
    participantSettingsSchema,
    type ParticipantSettingsData,
    type EventType,
    type FieldRequirement,
} from "@/lib/schemas/participantSettingsSchema";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Input } from "@/components/ui/input";
import { WrappedSelect } from "@/components/forms/WrappedSelect";
import { WrappedDateSelect } from "@/components/forms/wrapped-inputs";
import { Controller } from "react-hook-form";
import type { Resolver } from "react-hook-form";
import { RegistrationFormPreview } from "./RegistrationFormPreview";
import { Eye } from "lucide-react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { useState } from "react";

interface ParticipantSettingsFormProps {
    /**
     * Default values for the form (used for editing existing settings)
     */
    defaultValues?: Partial<ParticipantSettingsData>;
    /**
     * Callback function called when form is successfully submitted
     */
    onSubmit: (data: ParticipantSettingsData) => void | Promise<void>;
    /**
     * Whether the form is in loading/submitting state
     */
    isLoading?: boolean;
    /**
     * Whether to show preview button
     */
    showPreview?: boolean;
}

export const ParticipantSettingsForm = ({
    defaultValues,
    onSubmit,
    isLoading = false,
    showPreview = false,
}: ParticipantSettingsFormProps) => {
    const { t } = useTranslation();

    const form = useForm<ParticipantSettingsData>({
        resolver: zodResolver(participantSettingsSchema) as Resolver<ParticipantSettingsData>,
        defaultValues: {
            eventType: defaultValues?.eventType || "public",
            isBookingRequired: defaultValues?.isBookingRequired ?? false,
            isTicketTransferable: defaultValues?.isTicketTransferable ?? true,
            requireRegistrationPassword: defaultValues?.requireRegistrationPassword ?? false,
            registrationPassword: defaultValues?.registrationPassword || "",
            finalCallRegistrationDate: defaultValues?.finalCallRegistrationDate,
            firstName: defaultValues?.firstName || "not_required",
            lastName: defaultValues?.lastName || "not_required",
            email: defaultValues?.email || "not_required",
            bio: defaultValues?.bio || "not_required",
            phoneNumber: defaultValues?.phoneNumber || "not_required",
            address: defaultValues?.address || "not_required",
            academicInstitution: defaultValues?.academicInstitution || "not_required",
            academicEmail: defaultValues?.academicEmail || "not_required",
        },
        mode: "onChange",
    });

    const { control, watch } = form;
    const [isPreviewOpen, setIsPreviewOpen] = useState(false);

    // Watch form values for preview and conditional rendering
    const formValues = watch();
    const eventType = watch("eventType");
    const isBookingRequired = watch("isBookingRequired");

    const requireRegistrationPassword = watch("requireRegistrationPassword");

    const handleFormSubmit = async (data: ParticipantSettingsData) => {
        await onSubmit(data);
    };

    const fieldRequirementOptions: { value: FieldRequirement; label: string }[] = [
        { value: "not_required", label: t("participantSettings.fieldRequirement.notRequired") },
        { value: "optional", label: t("participantSettings.fieldRequirement.optional") },
        { value: "required", label: t("participantSettings.fieldRequirement.required") },
    ];

    return (
        <form
            onSubmit={form.handleSubmit(
                async (data) => await handleFormSubmit(data as ParticipantSettingsData),
            )}
            className="space-y-8 "
        >
            {/* Registration Settings Section */}
            <div className="space-y-6 ">
                <div>
                    <Typography variant="header" tag="h2" className="text-xl font-bold mb-2 ">
                        {t("participantSettings.registrationSettings")}
                    </Typography>
                    <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                        {t("participantSettings.registrationSettingsDescription")}
                    </Typography>
                </div>

                <div className="space-y-4 rounded-lg border p-6">
                    {/* Event Type */}
                    <WrappedSelect
                        control={control}
                        name="eventType"
                        label={t("participantSettings.eventType")}
                        description={t("participantSettings.eventTypeDescription")}
                        htmlFor="eventType"
                        options={[
                            { value: "public", label: t("participantSettings.eventTypePublic") },
                            { value: "private", label: t("participantSettings.eventTypePrivate") },
                            {
                                value: "invite",
                                label: t("participantSettings.eventTypeInviteOnly"),
                            },
                        ]}
                        disabled={isLoading}
                        valueAs={(value) => value as EventType}
                    />

                    {/* Switches Grid */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {/* Booking Required */}
                        <div className="flex items-center justify-between space-x-2">
                            <div className="space-y-1 flex-1">
                                <Label htmlFor="isBookingRequired">
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        className="text-sm font-medium"
                                    >
                                        {t("participantSettings.bookingRequired")}
                                    </Typography>
                                </Label>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-xs text-muted-foreground"
                                >
                                    {t("participantSettings.bookingRequiredDescription")}
                                </Typography>
                            </div>
                            <Controller
                                name="isBookingRequired"
                                control={control}
                                render={({ field }) => (
                                    <Switch
                                        id="isBookingRequired"
                                        checked={field.value}
                                        onCheckedChange={field.onChange}
                                        disabled={isLoading}
                                    />
                                )}
                            />
                        </div>

                        {/* Ticket Transferable */}
                        <div className="flex items-center justify-between space-x-2">
                            <div className="space-y-1 flex-1">
                                <Label htmlFor="isTicketTransferable">
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        className="text-sm font-medium"
                                    >
                                        {t("participantSettings.ticketTransferable")}
                                    </Typography>
                                </Label>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-xs text-muted-foreground"
                                >
                                    {t("participantSettings.ticketTransferableDescription")}
                                </Typography>
                            </div>
                            <Controller
                                name="isTicketTransferable"
                                control={control}
                                render={({ field }) => (
                                    <Switch
                                        id="isTicketTransferable"
                                        checked={field.value}
                                        onCheckedChange={field.onChange}
                                        disabled={isLoading}
                                    />
                                )}
                            />
                        </div>
                    </div>

                    {/* Registration Password Required */}
                    <div className="flex items-center justify-between space-x-2">
                        <div className="space-y-1 flex-1">
                            <Label htmlFor="requireRegistrationPassword">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    className="text-sm font-medium"
                                >
                                    {t("participantSettings.requireRegistrationPassword")}
                                </Typography>
                            </Label>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-xs text-muted-foreground"
                            >
                                {t("participantSettings.requireRegistrationPasswordDescription")}
                            </Typography>
                        </div>
                        <Controller
                            name="requireRegistrationPassword"
                            control={control}
                            render={({ field }) => (
                                <Switch
                                    id="requireRegistrationPassword"
                                    checked={field.value}
                                    onCheckedChange={field.onChange}
                                    disabled={isLoading}
                                />
                            )}
                        />
                    </div>

                    {/* Registration Password (conditional) */}
                    {requireRegistrationPassword && (
                        <div className="space-y-2 pt-2 border-t pb-6">
                            <Label htmlFor="registrationPassword" className="pt-6">
                                <Typography
                                    variant="text"
                                    tag="span"
                                    className="text-sm font-medium"
                                >
                                    {t("participantSettings.registrationPassword")}
                                </Typography>
                            </Label>
                            <Controller
                                name="registrationPassword"
                                control={control}
                                render={({ field, fieldState: { error } }) => (
                                    <div className="space-y-2">
                                        <Input
                                            {...field}
                                            id="registrationPassword"
                                            type="password"
                                            placeholder={t(
                                                "participantSettings.registrationPasswordPlaceholder",
                                            )}
                                            disabled={isLoading}
                                            className={`!border !border-[#D9D9D91A] ${error ? "border-destructive" : ""}`}
                                            value={field.value || ""}
                                            onChange={(e) => field.onChange(e.target.value)}
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
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-xs text-muted-foreground"
                            >
                                {t("participantSettings.registrationPasswordDescription")}
                            </Typography>
                        </div>
                    )}

                    {/* Final Call Registration Date (conditional) */}
                    {(eventType === "public" || eventType === "private") && isBookingRequired && (
                        <div className="space-y-2 pt-2 pb-4">
                            <WrappedDateSelect
                                control={control}
                                name="finalCallRegistrationDate"
                                label={t("participantSettings.finalCallRegistrationDate")}
                                placeholder={t(
                                    "participantSettings.finalCallRegistrationDatePlaceholder",
                                )}
                                disabled={isLoading}
                                disablePastDates
                            />
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-xs text-muted-foreground"
                            >
                                {t("participantSettings.finalCallRegistrationDateDescription")}
                            </Typography>
                        </div>
                    )}
                </div>
            </div>

            {/* Participant Requirements Section */}
            <div className="space-y-6 ">
                <div>
                    <Typography variant="header" tag="h2" className="text-xl font-bold mb-2">
                        {t("participantSettings.participantRequirements")}
                    </Typography>
                    <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                        {t("participantSettings.participantRequirementsDescription")}
                    </Typography>
                </div>

                <div className="space-y-4 rounded-lg border p-6 ">
                    {/* Basic Information */}
                    <div className="space-y-4">
                        <Typography variant="header" tag="h3" className="text-sm font-semibold">
                            {t("participantSettings.basicInformation")}
                        </Typography>

                        {/* First Name */}
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                            <WrappedSelect
                                control={control}
                                name="firstName"
                                label={t("participantSettings.fields.firstName")}
                                options={fieldRequirementOptions}
                                disabled={isLoading}
                                valueAs={(value) => value as FieldRequirement}
                                labelClassName="text-sm font-medium"
                            />
                        </div>

                        {/* Last Name */}
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                            <WrappedSelect
                                control={control}
                                name="lastName"
                                label={t("participantSettings.fields.lastName")}
                                options={fieldRequirementOptions}
                                disabled={isLoading}
                                valueAs={(value) => value as FieldRequirement}
                                labelClassName="text-sm font-medium"
                            />
                        </div>

                        {/* Email */}
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                            <WrappedSelect
                                control={control}
                                name="email"
                                label={t("participantSettings.fields.email")}
                                options={fieldRequirementOptions}
                                disabled={isLoading}
                                valueAs={(value) => value as FieldRequirement}
                                labelClassName="text-sm font-medium"
                            />
                        </div>

                        {/* Phone Number */}
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                            <WrappedSelect
                                control={control}
                                name="phoneNumber"
                                label={t("participantSettings.fields.phoneNumber")}
                                options={fieldRequirementOptions}
                                disabled={isLoading}
                                valueAs={(value) => value as FieldRequirement}
                                labelClassName="text-sm font-medium"
                            />
                        </div>
                    </div>

                    {/* Additional Information */}
                    <div className="pt-4 space-y-4 border-t">
                        <Typography variant="header" tag="h3" className="text-sm font-semibold">
                            {t("participantSettings.additionalInformation")}
                        </Typography>

                        {/* Bio */}
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                            <WrappedSelect
                                control={control}
                                name="bio"
                                label={t("participantSettings.fields.bio")}
                                options={fieldRequirementOptions}
                                disabled={isLoading}
                                valueAs={(value) => value as FieldRequirement}
                                labelClassName="text-sm font-medium"
                            />
                        </div>

                        {/* Address */}
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                            <WrappedSelect
                                control={control}
                                name="address"
                                label={t("participantSettings.fields.address")}
                                options={fieldRequirementOptions}
                                disabled={isLoading}
                                valueAs={(value) => value as FieldRequirement}
                                labelClassName="text-sm font-medium"
                            />
                        </div>
                    </div>

                    {/* Academic Information */}
                    <div className="pt-4 space-y-4 border-t">
                        <Typography variant="header" tag="h3" className="text-sm font-semibold">
                            {t("participantSettings.academicInformation")}
                        </Typography>

                        {/* Academic Institution */}
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                            <WrappedSelect
                                control={control}
                                name="academicInstitution"
                                label={t("participantSettings.fields.academicInstitution")}
                                options={fieldRequirementOptions}
                                disabled={isLoading}
                                valueAs={(value) => value as FieldRequirement}
                                labelClassName="text-sm font-medium"
                            />
                        </div>

                        {/* Academic Email */}
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                            <WrappedSelect
                                control={control}
                                name="academicEmail"
                                label={t("participantSettings.fields.academicEmail")}
                                options={fieldRequirementOptions}
                                disabled={isLoading}
                                valueAs={(value) => value as FieldRequirement}
                                labelClassName="text-sm font-medium"
                            />
                        </div>
                    </div>
                </div>
            </div>

            {/* Submit Button */}
            <div className="flex justify-end gap-4 pt-4">
                {showPreview && (
                    <Button
                        type="button"
                        variant="secondary-light"
                        size="lg"
                        onClick={() => setIsPreviewOpen(true)}
                        disabled={isLoading}
                        className="min-w-[150px] "
                    >
                        <Eye className="h-4 w-4 mr-2" />
                        <Typography variant="text" tag="span" className="font-medium text-black">
                            {t("participantSettings.preview.title")}
                        </Typography>
                    </Button>
                )}
                <Button
                    type="submit"
                    variant="primary"
                    size="lg"
                    disabled={isLoading}
                    className="min-w-[150px]"
                >
                    <Typography variant="text" tag="span" className="font-medium">
                        {isLoading ? t("common.loading") : t("participantSettings.saveSettings")}
                    </Typography>
                </Button>
            </div>

            {/* Preview Dialog */}
            <Dialog open={isPreviewOpen} onOpenChange={setIsPreviewOpen}>
                <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
                    <DialogHeader>
                        <DialogTitle></DialogTitle>
                    </DialogHeader>
                    <div className="mt-4">
                        <RegistrationFormPreview settings={formValues} />
                    </div>
                </DialogContent>
            </Dialog>
        </form>
    );
};
