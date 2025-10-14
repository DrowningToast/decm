import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Controller } from "react-hook-form";
import { RegistrationFormPreview } from "@/components/forms/ParticipantSettingsForm";
import {
    participantSettingsSchema,
    type ParticipantSettingsData,
    type EventType,
    type FieldRequirement,
    defaultParticipantSettings,
} from "@/lib/schemas/participantSettingsSchema";
import { Eye } from "lucide-react";
import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";

export const EventParticipantSettingPage = () => {
    const { t } = useTranslation();
    const { eventId } = useParams<{ eventId: string }>();
    const [isPreviewOpen, setIsPreviewOpen] = useState(false);
    const [isLoading, setIsLoading] = useState(false);

    const { handleSubmit, control, formState, watch } = useForm<ParticipantSettingsData>({
        resolver: zodResolver(participantSettingsSchema),
        defaultValues: defaultParticipantSettings,
        mode: "onChange",
    });

    // Watch form values for preview
    const formValues = watch();

    const onSubmit = async (data: ParticipantSettingsData) => {
        try {
            setIsLoading(true);

            // TODO: Implement API call to save participant settings
            console.log("Event ID:", eventId);
            console.log("Participant Settings:", data);

            // Simulate API delay
            await new Promise((resolve) => setTimeout(resolve, 1000));

            // Success feedback (can be replaced with toast notification)
            alert(t("participantSettings.saveSuccess"));
        } catch (error) {
            console.error("Error saving participant settings:", error);
            alert(t("participantSettings.saveError"));
        } finally {
            setIsLoading(false);
        }
    };

    const fieldRequirementOptions: { value: FieldRequirement; label: string }[] = [
        { value: "not_required", label: t("participantSettings.fieldRequirement.notRequired") },
        { value: "optional", label: t("participantSettings.fieldRequirement.optional") },
        { value: "required", label: t("participantSettings.fieldRequirement.required") },
    ];

    return (
        <PageContainer
            title={t("participantSettings.pageTitle")}
            description={t("participantSettings.pageDescription")}
        >
            <SectionContainer>
                <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
                    {/* Registration Settings Section */}
                    <div className="space-y-6">
                        <div>
                            <Typography
                                variant="header"
                                tag="h2"
                                className="text-xl font-bold mb-2"
                            >
                                {t("participantSettings.registrationSettings")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-sm text-muted-foreground"
                            >
                                {t("participantSettings.registrationSettingsDescription")}
                            </Typography>
                        </div>

                        <div className="space-y-4 rounded-lg border p-6">
                            {/* Event Type */}
                            <div className="space-y-2">
                                <Label htmlFor="eventType">
                                    <Typography
                                        variant="text"
                                        tag="span"
                                        className="text-sm font-medium"
                                    >
                                        {t("participantSettings.eventType")}
                                    </Typography>
                                </Label>
                                <Controller
                                    name="eventType"
                                    control={control}
                                    render={({ field }) => (
                                        <Select
                                            value={field.value}
                                            onValueChange={(value) =>
                                                field.onChange(value as EventType)
                                            }
                                            disabled={isLoading}
                                        >
                                            <SelectTrigger>
                                                <SelectValue />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="public">
                                                    {t("participantSettings.eventTypePublic")}
                                                </SelectItem>
                                                <SelectItem value="private">
                                                    {t("participantSettings.eventTypePrivate")}
                                                </SelectItem>
                                            </SelectContent>
                                        </Select>
                                    )}
                                />
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-xs text-muted-foreground"
                                >
                                    {t("participantSettings.eventTypeDescription")}
                                </Typography>
                            </div>

                            {/* Booking Required */}
                            <div className="flex items-center justify-between space-x-2 py-2">
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
                            <div className="flex items-center justify-between space-x-2 py-2">
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
                    </div>

                    {/* Participant Requirements Section */}
                    <div className="space-y-6">
                        <div>
                            <Typography
                                variant="header"
                                tag="h2"
                                className="text-xl font-bold mb-2"
                            >
                                {t("participantSettings.participantRequirements")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-sm text-muted-foreground"
                            >
                                {t("participantSettings.participantRequirementsDescription")}
                            </Typography>
                        </div>

                        <div className="space-y-4 rounded-lg border p-6">
                            {/* Basic Information */}
                            <div className="space-y-4">
                                <Typography
                                    variant="header"
                                    tag="h3"
                                    className="text-sm font-semibold"
                                >
                                    {t("participantSettings.basicInformation")}
                                </Typography>

                                {/* First Name */}
                                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                                    <Label>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-sm font-medium"
                                        >
                                            {t("participantSettings.fields.firstName")}
                                        </Typography>
                                    </Label>
                                    <div className="md:col-span-2">
                                        <Controller
                                            name="firstName"
                                            control={control}
                                            render={({ field }) => (
                                                <Select
                                                    value={field.value}
                                                    onValueChange={(value) =>
                                                        field.onChange(value as FieldRequirement)
                                                    }
                                                    disabled={isLoading}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {fieldRequirementOptions.map((option) => (
                                                            <SelectItem
                                                                key={option.value}
                                                                value={option.value}
                                                            >
                                                                {option.label}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            )}
                                        />
                                    </div>
                                </div>

                                {/* Last Name */}
                                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                                    <Label>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-sm font-medium"
                                        >
                                            {t("participantSettings.fields.lastName")}
                                        </Typography>
                                    </Label>
                                    <div className="md:col-span-2">
                                        <Controller
                                            name="lastName"
                                            control={control}
                                            render={({ field }) => (
                                                <Select
                                                    value={field.value}
                                                    onValueChange={(value) =>
                                                        field.onChange(value as FieldRequirement)
                                                    }
                                                    disabled={isLoading}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {fieldRequirementOptions.map((option) => (
                                                            <SelectItem
                                                                key={option.value}
                                                                value={option.value}
                                                            >
                                                                {option.label}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            )}
                                        />
                                    </div>
                                </div>

                                {/* Email */}
                                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                                    <Label>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-sm font-medium"
                                        >
                                            {t("participantSettings.fields.email")}
                                        </Typography>
                                    </Label>
                                    <div className="md:col-span-2">
                                        <Controller
                                            name="email"
                                            control={control}
                                            render={({ field }) => (
                                                <Select
                                                    value={field.value}
                                                    onValueChange={(value) =>
                                                        field.onChange(value as FieldRequirement)
                                                    }
                                                    disabled={isLoading}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {fieldRequirementOptions.map((option) => (
                                                            <SelectItem
                                                                key={option.value}
                                                                value={option.value}
                                                            >
                                                                {option.label}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            )}
                                        />
                                    </div>
                                </div>

                                {/* Phone Number */}
                                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                                    <Label>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-sm font-medium"
                                        >
                                            {t("participantSettings.fields.phoneNumber")}
                                        </Typography>
                                    </Label>
                                    <div className="md:col-span-2">
                                        <Controller
                                            name="phoneNumber"
                                            control={control}
                                            render={({ field }) => (
                                                <Select
                                                    value={field.value}
                                                    onValueChange={(value) =>
                                                        field.onChange(value as FieldRequirement)
                                                    }
                                                    disabled={isLoading}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {fieldRequirementOptions.map((option) => (
                                                            <SelectItem
                                                                key={option.value}
                                                                value={option.value}
                                                            >
                                                                {option.label}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            )}
                                        />
                                    </div>
                                </div>
                            </div>

                            {/* Additional Information */}
                            <div className="pt-4 space-y-4 border-t">
                                <Typography
                                    variant="header"
                                    tag="h3"
                                    className="text-sm font-semibold"
                                >
                                    {t("participantSettings.additionalInformation")}
                                </Typography>

                                {/* Bio */}
                                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                                    <Label>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-sm font-medium"
                                        >
                                            {t("participantSettings.fields.bio")}
                                        </Typography>
                                    </Label>
                                    <div className="md:col-span-2">
                                        <Controller
                                            name="bio"
                                            control={control}
                                            render={({ field }) => (
                                                <Select
                                                    value={field.value}
                                                    onValueChange={(value) =>
                                                        field.onChange(value as FieldRequirement)
                                                    }
                                                    disabled={isLoading}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {fieldRequirementOptions.map((option) => (
                                                            <SelectItem
                                                                key={option.value}
                                                                value={option.value}
                                                            >
                                                                {option.label}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            )}
                                        />
                                    </div>
                                </div>

                                {/* Address */}
                                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                                    <Label>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-sm font-medium"
                                        >
                                            {t("participantSettings.fields.address")}
                                        </Typography>
                                    </Label>
                                    <div className="md:col-span-2">
                                        <Controller
                                            name="address"
                                            control={control}
                                            render={({ field }) => (
                                                <Select
                                                    value={field.value}
                                                    onValueChange={(value) =>
                                                        field.onChange(value as FieldRequirement)
                                                    }
                                                    disabled={isLoading}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {fieldRequirementOptions.map((option) => (
                                                            <SelectItem
                                                                key={option.value}
                                                                value={option.value}
                                                            >
                                                                {option.label}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            )}
                                        />
                                    </div>
                                </div>
                            </div>

                            {/* Academic Information */}
                            <div className="pt-4 space-y-4 border-t">
                                <Typography
                                    variant="header"
                                    tag="h3"
                                    className="text-sm font-semibold"
                                >
                                    {t("participantSettings.academicInformation")}
                                </Typography>

                                {/* Academic Institution */}
                                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                                    <Label>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-sm font-medium"
                                        >
                                            {t("participantSettings.fields.academicInstitution")}
                                        </Typography>
                                    </Label>
                                    <div className="md:col-span-2">
                                        <Controller
                                            name="academicInstitution"
                                            control={control}
                                            render={({ field }) => (
                                                <Select
                                                    value={field.value}
                                                    onValueChange={(value) =>
                                                        field.onChange(value as FieldRequirement)
                                                    }
                                                    disabled={isLoading}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {fieldRequirementOptions.map((option) => (
                                                            <SelectItem
                                                                key={option.value}
                                                                value={option.value}
                                                            >
                                                                {option.label}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            )}
                                        />
                                    </div>
                                </div>

                                {/* Academic Email */}
                                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
                                    <Label>
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            className="text-sm font-medium"
                                        >
                                            {t("participantSettings.fields.academicEmail")}
                                        </Typography>
                                    </Label>
                                    <div className="md:col-span-2">
                                        <Controller
                                            name="academicEmail"
                                            control={control}
                                            render={({ field }) => (
                                                <Select
                                                    value={field.value}
                                                    onValueChange={(value) =>
                                                        field.onChange(value as FieldRequirement)
                                                    }
                                                    disabled={isLoading}
                                                >
                                                    <SelectTrigger>
                                                        <SelectValue />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {fieldRequirementOptions.map((option) => (
                                                            <SelectItem
                                                                key={option.value}
                                                                value={option.value}
                                                            >
                                                                {option.label}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            )}
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Action Buttons */}
                    <div className="flex justify-end gap-4 pt-4">
                        <Button
                            type="button"
                            variant="secondary-dark"
                            size="lg"
                            onClick={() => setIsPreviewOpen(true)}
                            disabled={isLoading}
                            className="min-w-[150px]"
                        >
                            <Eye className="h-4 w-4 mr-2" />
                            <Typography variant="text" tag="span" className="font-medium">
                                {t("participantSettings.preview.title")}
                            </Typography>
                        </Button>
                        <Button
                            type="submit"
                            variant="primary"
                            size="lg"
                            disabled={isLoading || !formState.isDirty}
                            className="min-w-[150px]"
                        >
                            <Typography variant="text" tag="span" className="font-medium">
                                {isLoading
                                    ? t("common.loading")
                                    : t("participantSettings.saveSettings")}
                            </Typography>
                        </Button>
                    </div>
                </form>
            </SectionContainer>

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
        </PageContainer>
    );
};
