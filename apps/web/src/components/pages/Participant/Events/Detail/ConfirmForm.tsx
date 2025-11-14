import React, { useMemo } from "react";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
    FormField,
    FormItem,
    FormLabel,
    FormControl,
    FormMessage,
    Form,
} from "@/components/ui/form";
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
import type {
    EventRegistrationInvitation,
    RegistrationRequirement,
    RegistrationRequirementStatus,
} from "@/services/EventRegistration/EventRegistration";
import { FieldRequirementBadge } from "./components/FieldRequirementBadge";

interface RegistrationConfirmFormProps {
    eventId: string;
    onSubmit: (data: RegistrationConfirmDataForm) => Promise<void>;
    onCancel: () => void;
    isSubmitting?: boolean;
    // if optional or required, will autofill
    profileData: RegistrationConfirmDataForm;
    // If provide, the field is counted as required and locked
    invitationData?: EventRegistrationInvitation;
}

export type FieldStatus = RegistrationRequirementStatus | "locked";
type Invite = Omit<
    EventRegistrationInvitation,
    | "code"
    | "id"
    | "eventId"
    | "inboxMessageId"
    | "validUntil"
    | "createdAt"
    | "updatedAt"
    | "cancelledAt"
>;

const mapFieldRequirementToStatus = (
    key: keyof RegistrationRequirement,
    requirement: RegistrationRequirement,
    invitation?: Invite,
): FieldStatus => {
    if (!invitation) {
        return requirement[key];
    }
    if (key in invitation && invitation[key as keyof Invite] !== undefined) {
        return "locked";
    }
    return requirement[key];
};

export const RegistrationConfirmForm: React.FC<RegistrationConfirmFormProps> = ({
    eventId,
    onSubmit,
    onCancel,
    isSubmitting = false,
    profileData,
    invitationData,
}) => {
    const { t } = useTranslation();

    const { event } = useEventViewModelUsecase({ eventId });
    const formRequirements = event?.registrationRequirement;

    const formSchema = createRegistrationConfirmDataFormSchema(t, formRequirements);
    const initialData = useMemo<RegistrationConfirmDataForm>(() => {
        return {
            firstName:
                invitationData?.firstName ??
                (formRequirements?.firstName === "not_required"
                    ? undefined
                    : profileData?.firstName),
            lastName:
                invitationData?.lastName ??
                (formRequirements?.lastName === "not_required" ? undefined : profileData?.lastName),
            bio:
                profileData?.bio ??
                (formRequirements?.bio === "not_required" ? undefined : profileData?.bio),
            email:
                invitationData?.email ??
                (formRequirements?.email === "not_required" ? undefined : profileData?.email),
            phoneNumber:
                invitationData?.phoneNumber ??
                (formRequirements?.phoneNumber === "not_required"
                    ? undefined
                    : profileData?.phoneNumber),
            address:
                profileData?.address ??
                (formRequirements?.address === "not_required" ? undefined : profileData?.address),
            academicEmail:
                profileData?.academicEmail ??
                (formRequirements?.academicEmail === "not_required"
                    ? undefined
                    : profileData?.academicEmail),
            academicInstitution:
                invitationData?.academicInstitution ??
                (formRequirements?.academicInstitution === "not_required"
                    ? undefined
                    : profileData?.academicInstitution),
        };
    }, [invitationData, profileData, formRequirements]);

    const form = useForm({
        resolver: zodResolver(formSchema),
        mode: "onChange" as const,
        defaultValues: initialData,
    });

    // Helper function to get field requirement status
    const getFieldStatus = (
        fieldName: keyof RegistrationRequirement | keyof Invite,
    ): FieldStatus => {
        if (!formRequirements) {
            return "not_required";
        }
        return mapFieldRequirementToStatus(fieldName, formRequirements, invitationData);
    };

    // Helper function to check if field should be shown
    const shouldShowField = (fieldName: keyof RegistrationRequirement): boolean => {
        const requirement = getFieldStatus(fieldName);
        return requirement === "required" || requirement === "optional" || requirement === "locked";
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
                                        <FormField
                                            control={form.control}
                                            name="firstName"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <div className="flex items-center gap-2">
                                                        <FormLabel className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70">
                                                            {t("profile.firstName")}
                                                        </FormLabel>
                                                        <FieldRequirementBadge
                                                            status={getFieldStatus("firstName")}
                                                            t={t}
                                                        />
                                                    </div>
                                                    <FormControl>
                                                        <Input
                                                            {...field}
                                                            type="text"
                                                            className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                            readOnly={
                                                                getFieldStatus("firstName") ===
                                                                "locked"
                                                            }
                                                            disabled={
                                                                getFieldStatus("firstName") ===
                                                                "locked"
                                                            }
                                                        />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    )}

                                    {/* Last Name */}
                                    {shouldShowField("lastName") && (
                                        <FormField
                                            control={form.control}
                                            name="lastName"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <div className="flex items-center gap-2">
                                                        <FormLabel className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70">
                                                            {t("profile.lastName")}
                                                        </FormLabel>
                                                        <FieldRequirementBadge
                                                            status={getFieldStatus("lastName")}
                                                            t={t}
                                                        />
                                                    </div>
                                                    <FormControl>
                                                        <Input
                                                            {...field}
                                                            type="text"
                                                            className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                            readOnly={
                                                                getFieldStatus("lastName") ===
                                                                "locked"
                                                            }
                                                            disabled={
                                                                getFieldStatus("lastName") ===
                                                                "locked"
                                                            }
                                                        />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    )}

                                    {/* Bio */}
                                    {shouldShowField("bio") && (
                                        <FormField
                                            control={form.control}
                                            name="bio"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <div className="flex items-center gap-2">
                                                        <FormLabel className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground-alt/70">
                                                            {t("profile.bio")}
                                                        </FormLabel>
                                                        <FieldRequirementBadge
                                                            status={getFieldStatus("bio")}
                                                            t={t}
                                                        />
                                                    </div>
                                                    <FormControl>
                                                        <Textarea
                                                            {...field}
                                                            className="w-full min-h-24 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                            placeholder={t(
                                                                "profile.bioPlaceholder",
                                                            )}
                                                            readOnly={
                                                                getFieldStatus("bio") === "locked"
                                                            }
                                                            disabled={
                                                                getFieldStatus("bio") === "locked"
                                                            }
                                                        />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
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
                                        <FormField
                                            control={form.control}
                                            name="email"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <div className="flex items-center gap-2">
                                                        <FormLabel className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70">
                                                            {t("profile.email")}
                                                        </FormLabel>
                                                        <FieldRequirementBadge
                                                            status={getFieldStatus("email")}
                                                            t={t}
                                                        />
                                                    </div>
                                                    <FormControl>
                                                        <Input
                                                            {...field}
                                                            type="email"
                                                            className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                            readOnly={
                                                                getFieldStatus("email") === "locked"
                                                            }
                                                            disabled={
                                                                getFieldStatus("email") === "locked"
                                                            }
                                                        />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    )}

                                    {/* Phone Number */}
                                    {shouldShowField("phoneNumber") && (
                                        <FormField
                                            control={form.control}
                                            name="phoneNumber"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <div className="flex items-center gap-2">
                                                        <FormLabel className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70">
                                                            {t("profile.phoneNumber")}
                                                        </FormLabel>
                                                        <FieldRequirementBadge
                                                            status={getFieldStatus("phoneNumber")}
                                                            t={t}
                                                        />
                                                    </div>
                                                    <FormControl>
                                                        <Input
                                                            {...field}
                                                            type="tel"
                                                            className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                            readOnly={
                                                                getFieldStatus("phoneNumber") ===
                                                                "locked"
                                                            }
                                                            disabled={
                                                                getFieldStatus("phoneNumber") ===
                                                                "locked"
                                                            }
                                                        />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    )}

                                    {/* Address */}
                                    {shouldShowField("address") && (
                                        <FormField
                                            control={form.control}
                                            name="address"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <div className="flex items-center gap-2">
                                                        <FormLabel className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70">
                                                            {t("profile.address")}
                                                        </FormLabel>
                                                        <FieldRequirementBadge
                                                            status={getFieldStatus("address")}
                                                            t={t}
                                                        />
                                                    </div>
                                                    <FormControl>
                                                        <Textarea
                                                            {...field}
                                                            className="w-full min-h-20 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                            placeholder={t(
                                                                "profile.addressPlaceholder",
                                                            )}
                                                            readOnly={
                                                                getFieldStatus("address") ===
                                                                "locked"
                                                            }
                                                            disabled={
                                                                getFieldStatus("address") ===
                                                                "locked"
                                                            }
                                                        />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
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
                                        <FormField
                                            control={form.control}
                                            name="academicEmail"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <div className="flex items-center gap-2">
                                                        <FormLabel className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70">
                                                            {t("profile.academicEmail")}
                                                        </FormLabel>
                                                        <FieldRequirementBadge
                                                            status={getFieldStatus("academicEmail")}
                                                            t={t}
                                                        />
                                                    </div>
                                                    <FormControl>
                                                        <Input
                                                            {...field}
                                                            type="email"
                                                            className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                            readOnly={
                                                                getFieldStatus("academicEmail") ===
                                                                "locked"
                                                            }
                                                            disabled={
                                                                getFieldStatus("academicEmail") ===
                                                                "locked"
                                                            }
                                                        />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    )}

                                    {/* Academic Institution */}
                                    {shouldShowField("academicInstitution") && (
                                        <FormField
                                            control={form.control}
                                            name="academicInstitution"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <div className="flex items-center gap-2">
                                                        <FormLabel className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground/70">
                                                            {t("profile.academicInstitution")}
                                                        </FormLabel>
                                                        <FieldRequirementBadge
                                                            status={getFieldStatus(
                                                                "academicInstitution",
                                                            )}
                                                            t={t}
                                                        />
                                                    </div>
                                                    <FormControl>
                                                        <Input
                                                            {...field}
                                                            type="text"
                                                            className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-background/50 border-border border-[0.5px] rounded-[12px] text-foreground placeholder:text-muted-foreground"
                                                            readOnly={
                                                                getFieldStatus(
                                                                    "academicInstitution",
                                                                ) === "locked"
                                                            }
                                                            disabled={
                                                                getFieldStatus(
                                                                    "academicInstitution",
                                                                ) === "locked"
                                                            }
                                                        />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
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
