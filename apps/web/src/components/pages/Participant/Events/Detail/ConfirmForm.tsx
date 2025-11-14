import React from "react";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { FormField, FormItem, FormControl, FormMessage, Form } from "@/components/ui/form";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogDescription,
} from "@/components/ui/alert-dialog";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { createRegistrationConfirmDataFormSchema } from "./RegistrationConfirmDataFormSchema";
import type { RegistrationConfirmDataForm } from "./RegistrationConfirmDataFormSchema";
import { useTranslation } from "react-i18next";
import { Loader2, Info } from "lucide-react";
import { useEventViewModelUsecase } from "./useEventViewModelUsecase";
import type { RegistrationRequirement } from "@/services/EventRegistration/EventRegistration";

interface RegistrationConfirmFormProps {
    eventId: string;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onSubmit: (data: RegistrationConfirmDataForm) => Promise<void>;
    onCancel: () => void;
    isSubmitting?: boolean;
    initialData?: RegistrationConfirmDataForm;
}

type FieldRequirement = "required" | "optional" | "not_needed";

// Helper component for field requirement indicator
const FieldRequirementBadge: React.FC<{
    requirement: FieldRequirement;
    t: (key: string) => string;
}> = ({ requirement, t }) => {
    if (requirement === "required") {
        return (
            <span className="ml-2 inline-flex items-center rounded-md bg-destructive/10 px-2 py-0.5 text-xs font-medium text-destructive ring-1 ring-inset ring-destructive/20">
                {t("common.required")}
            </span>
        );
    }
    if (requirement === "optional") {
        return (
            <span className="ml-2 inline-flex items-center rounded-md bg-secondary/50 px-2 py-0.5 text-xs font-medium text-muted-foreground ring-1 ring-inset ring-border">
                {t("common.optional")}
            </span>
        );
    }
    return null;
};

export const RegistrationConfirmForm: React.FC<RegistrationConfirmFormProps> = ({
    eventId,
    onSubmit,
    onCancel,
    isSubmitting = false,
    initialData,
}) => {
    const { t } = useTranslation();

    const { event } = useEventViewModelUsecase({ eventId });
    const formRequirements = event?.registrationRequirement;

    const formSchema = createRegistrationConfirmDataFormSchema(t, formRequirements);

    const form = useForm({
        resolver: zodResolver(formSchema),
        mode: "onChange" as const,
        defaultValues: initialData,
    });

    // Helper function to get field requirement status
    const getFieldRequirement = (fieldName: keyof RegistrationRequirement): FieldRequirement => {
        const value = formRequirements?.[fieldName];
        if (value === "required") return "required";
        if (value === "optional") return "optional";
        return "not_needed";
    };

    // Helper function to check if field should be shown
    const shouldShowField = (fieldName: keyof RegistrationRequirement): boolean => {
        const requirement = getFieldRequirement(fieldName);
        return requirement === "required" || requirement === "optional";
    };

    // Helper function to check if any field in a section should be shown
    const shouldShowSection = (fieldNames: (keyof RegistrationRequirement)[]): boolean => {
        return fieldNames.some((fieldName) => shouldShowField(fieldName));
    };

    const handleSubmit = async (data: RegistrationConfirmDataForm) => {
        await onSubmit(data);
    };

    const handleCancel = () => {
        onCancel();
    };

    if (!formRequirements) {
        return null;
    }

    return (
        <AlertDialog open={true} onOpenChange={() => handleCancel()}>
            <AlertDialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
                <AlertDialogHeader>
                    <AlertDialogTitle>
                        <Typography
                            variant="header"
                            tag="span"
                            color="primary"
                            className="text-2xl"
                        >
                            {t("events.registration.piiForm.title")}
                        </Typography>
                    </AlertDialogTitle>
                    <AlertDialogDescription>
                        <Typography
                            variant="text"
                            tag="span"
                            color="muted-foreground"
                            className="text-sm"
                        >
                            {t("events.registration.piiForm.subtitle")}
                        </Typography>
                    </AlertDialogDescription>
                </AlertDialogHeader>

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
                        {/* Info Alert */}
                        <Alert variant="info" className="bg-background/50 border-muted">
                            <Info className="h-4 w-4" />
                            <AlertTitle>{t("events.registration.piiForm.privacyTitle")}</AlertTitle>
                            <AlertDescription>
                                {t("events.registration.piiForm.privacyDescription")}
                            </AlertDescription>
                        </Alert>

                        {/* Personal Information Section */}
                        {shouldShowSection(["firstName", "lastName", "bio"]) && (
                            <div className="flex flex-col gap-y-3">
                                <Typography
                                    variant="header"
                                    tag="h2"
                                    color="primary"
                                    className="text-xl"
                                >
                                    {t("profile.personalInfo")}
                                </Typography>

                                <div className="space-y-2.5">
                                    {/* First Name */}
                                    {shouldShowField("firstName") && (
                                        <div className="space-y-1.5">
                                            <FormField
                                                control={form.control}
                                                name="firstName"
                                                render={({ field }) => (
                                                    <FormItem>
                                                        <FormControl>
                                                            <div className="flex flex-col gap-y-2">
                                                                <div className="flex items-center">
                                                                    <Label
                                                                        htmlFor={field.name}
                                                                        className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70"
                                                                    >
                                                                        {t("profile.firstName")}
                                                                    </Label>
                                                                    <FieldRequirementBadge
                                                                        requirement={getFieldRequirement(
                                                                            "firstName",
                                                                        )}
                                                                        t={t}
                                                                    />
                                                                </div>
                                                                <Input
                                                                    {...field}
                                                                    type="text"
                                                                    id={field.name}
                                                                    className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                                />
                                                            </div>
                                                        </FormControl>
                                                        <FormMessage />
                                                    </FormItem>
                                                )}
                                            />
                                        </div>
                                    )}

                                    {/* Last Name */}
                                    {shouldShowField("lastName") && (
                                        <div className="space-y-1.5">
                                            <FormField
                                                control={form.control}
                                                name="lastName"
                                                render={({ field }) => (
                                                    <FormItem>
                                                        <FormControl>
                                                            <div className="flex flex-col gap-y-2">
                                                                <div className="flex items-center">
                                                                    <Label
                                                                        htmlFor={field.name}
                                                                        className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70"
                                                                    >
                                                                        {t("profile.lastName")}
                                                                    </Label>
                                                                    <FieldRequirementBadge
                                                                        requirement={getFieldRequirement(
                                                                            "lastName",
                                                                        )}
                                                                        t={t}
                                                                    />
                                                                </div>
                                                                <Input
                                                                    {...field}
                                                                    type="text"
                                                                    id={field.name}
                                                                    className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                                />
                                                            </div>
                                                        </FormControl>
                                                        <FormMessage />
                                                    </FormItem>
                                                )}
                                            />
                                        </div>
                                    )}

                                    {/* Bio */}
                                    {shouldShowField("bio") && (
                                        <div className="space-y-1.5">
                                            <FormField
                                                control={form.control}
                                                name="bio"
                                                render={({ field }) => (
                                                    <FormItem>
                                                        <FormControl>
                                                            <div className="flex flex-col gap-y-2">
                                                                <div className="flex items-center">
                                                                    <Label
                                                                        htmlFor={field.name}
                                                                        className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground-alt/70"
                                                                    >
                                                                        {t("profile.bio")}
                                                                    </Label>
                                                                    <FieldRequirementBadge
                                                                        requirement={getFieldRequirement(
                                                                            "bio",
                                                                        )}
                                                                        t={t}
                                                                    />
                                                                </div>
                                                                <Textarea
                                                                    {...field}
                                                                    id={field.name}
                                                                    className="w-full min-h-24 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                                    placeholder={t(
                                                                        "profile.bioPlaceholder",
                                                                    )}
                                                                />
                                                            </div>
                                                        </FormControl>
                                                        <FormMessage />
                                                    </FormItem>
                                                )}
                                            />
                                        </div>
                                    )}
                                </div>
                            </div>
                        )}

                        {/* Contact Information Section */}
                        {shouldShowSection(["email", "phoneNumber", "address"]) && (
                            <div className="flex flex-col gap-y-3">
                                <Typography
                                    variant="header"
                                    tag="h2"
                                    color="primary"
                                    className="text-xl"
                                >
                                    {t("profile.contactInfo")}
                                </Typography>

                                <div className="space-y-2.5">
                                    {/* Email */}
                                    {shouldShowField("email") && (
                                        <div className="space-y-1.5">
                                            <FormField
                                                control={form.control}
                                                name="email"
                                                render={({ field }) => (
                                                    <FormItem>
                                                        <FormControl>
                                                            <div className="flex flex-col gap-y-2">
                                                                <div className="flex items-center">
                                                                    <Label
                                                                        htmlFor={field.name}
                                                                        className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70"
                                                                    >
                                                                        {t("profile.email")}
                                                                    </Label>
                                                                    <FieldRequirementBadge
                                                                        requirement={getFieldRequirement(
                                                                            "email",
                                                                        )}
                                                                        t={t}
                                                                    />
                                                                </div>
                                                                <Input
                                                                    {...field}
                                                                    type="email"
                                                                    id={field.name}
                                                                    className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                                />
                                                            </div>
                                                        </FormControl>
                                                        <FormMessage />
                                                    </FormItem>
                                                )}
                                            />
                                        </div>
                                    )}

                                    {/* Phone Number */}
                                    {shouldShowField("phoneNumber") && (
                                        <div className="space-y-1.5">
                                            <FormField
                                                control={form.control}
                                                name="phoneNumber"
                                                render={({ field }) => (
                                                    <FormItem>
                                                        <FormControl>
                                                            <div className="flex flex-col gap-y-2">
                                                                <div className="flex items-center">
                                                                    <Label
                                                                        htmlFor={field.name}
                                                                        className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70"
                                                                    >
                                                                        {t("profile.phoneNumber")}
                                                                    </Label>
                                                                    <FieldRequirementBadge
                                                                        requirement={getFieldRequirement(
                                                                            "phoneNumber",
                                                                        )}
                                                                        t={t}
                                                                    />
                                                                </div>
                                                                <Input
                                                                    {...field}
                                                                    type="tel"
                                                                    id={field.name}
                                                                    className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                                />
                                                            </div>
                                                        </FormControl>
                                                        <FormMessage />
                                                    </FormItem>
                                                )}
                                            />
                                        </div>
                                    )}

                                    {/* Address */}
                                    {shouldShowField("address") && (
                                        <div className="space-y-1.5">
                                            <FormField
                                                control={form.control}
                                                name="address"
                                                render={({ field }) => (
                                                    <FormItem>
                                                        <FormControl>
                                                            <div className="flex flex-col gap-y-2">
                                                                <div className="flex items-center">
                                                                    <Label
                                                                        htmlFor={field.name}
                                                                        className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70"
                                                                    >
                                                                        {t("profile.address")}
                                                                    </Label>
                                                                    <FieldRequirementBadge
                                                                        requirement={getFieldRequirement(
                                                                            "address",
                                                                        )}
                                                                        t={t}
                                                                    />
                                                                </div>
                                                                <Textarea
                                                                    {...field}
                                                                    id={field.name}
                                                                    className="w-full min-h-20 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                                    placeholder={t(
                                                                        "profile.addressPlaceholder",
                                                                    )}
                                                                />
                                                            </div>
                                                        </FormControl>
                                                        <FormMessage />
                                                    </FormItem>
                                                )}
                                            />
                                        </div>
                                    )}
                                </div>
                            </div>
                        )}

                        {/* Academic Information Section */}
                        {shouldShowSection(["academicEmail", "academicInstitution"]) && (
                            <div className="flex flex-col gap-y-3">
                                <Typography
                                    variant="header"
                                    tag="h2"
                                    color="primary"
                                    className="text-xl"
                                >
                                    {t("profile.academicInfo")}
                                </Typography>

                                <div className="space-y-2.5">
                                    {/* Academic Email */}
                                    {shouldShowField("academicEmail") && (
                                        <div className="space-y-1.5">
                                            <FormField
                                                control={form.control}
                                                name="academicEmail"
                                                render={({ field }) => (
                                                    <FormItem>
                                                        <FormControl>
                                                            <div className="flex flex-col gap-y-2">
                                                                <div className="flex items-center">
                                                                    <Label
                                                                        htmlFor={field.name}
                                                                        className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70"
                                                                    >
                                                                        {t("profile.academicEmail")}
                                                                    </Label>
                                                                    <FieldRequirementBadge
                                                                        requirement={getFieldRequirement(
                                                                            "academicEmail",
                                                                        )}
                                                                        t={t}
                                                                    />
                                                                </div>
                                                                <Input
                                                                    {...field}
                                                                    type="email"
                                                                    id={field.name}
                                                                    className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                                />
                                                            </div>
                                                        </FormControl>
                                                        <FormMessage />
                                                    </FormItem>
                                                )}
                                            />
                                        </div>
                                    )}

                                    {/* Academic Institution */}
                                    {shouldShowField("academicInstitution") && (
                                        <div className="space-y-1.5">
                                            <FormField
                                                control={form.control}
                                                name="academicInstitution"
                                                render={({ field }) => (
                                                    <FormItem>
                                                        <FormControl>
                                                            <div className="flex flex-col gap-y-2">
                                                                <div className="flex items-center">
                                                                    <Label
                                                                        htmlFor={field.name}
                                                                        className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70"
                                                                    >
                                                                        {t(
                                                                            "profile.academicInstitution",
                                                                        )}
                                                                    </Label>
                                                                    <FieldRequirementBadge
                                                                        requirement={getFieldRequirement(
                                                                            "academicInstitution",
                                                                        )}
                                                                        t={t}
                                                                    />
                                                                </div>
                                                                <Input
                                                                    {...field}
                                                                    type="text"
                                                                    id={field.name}
                                                                    className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                                />
                                                            </div>
                                                        </FormControl>
                                                        <FormMessage />
                                                    </FormItem>
                                                )}
                                            />
                                        </div>
                                    )}
                                </div>
                            </div>
                        )}

                        {/* Action Buttons Section */}
                        <div className="flex gap-3 pt-4">
                            {/* Cancel Button */}
                            <Button
                                type="button"
                                onClick={handleCancel}
                                variant="secondary-light"
                                size="lg"
                                className="flex-1"
                                disabled={isSubmitting}
                            >
                                {t("common.cancel")}
                            </Button>

                            {/* Submit Button */}
                            <Button
                                type="submit"
                                variant="primary"
                                size="lg"
                                className="flex-1"
                                disabled={isSubmitting}
                            >
                                {isSubmitting ? (
                                    <>
                                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                        {t("common.submitting")}
                                    </>
                                ) : (
                                    t("common.submit")
                                )}
                            </Button>
                        </div>
                    </form>
                </Form>
            </AlertDialogContent>
        </AlertDialog>
    );
};
