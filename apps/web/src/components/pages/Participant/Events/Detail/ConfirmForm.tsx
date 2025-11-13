import React from "react";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { FormField, FormItem, FormControl, FormMessage, Form } from "@/components/ui/form";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { createRegistrationConfirmDataFormSchema } from "./RegistrationConfirmDataFormSchema";
import type { RegistrationConfirmDataForm } from "./RegistrationConfirmDataFormSchema";
import { useTranslation } from "react-i18next";
import { Loader2, Info } from "lucide-react";
import { useEventViewModelUsecase } from "./useEventViewModelUsecase";

interface RegistrationConfirmFormProps {
    eventId: string;
    onSubmit: (data: RegistrationConfirmDataForm) => Promise<void>;
    onCancel: () => void;
    isSubmitting?: boolean;
    initialData?: RegistrationConfirmDataForm;
}

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

    const handleSubmit = async (data: RegistrationConfirmDataForm) => {
        await onSubmit(data);
    };

    if (!formRequirements) {
        return <div>No registration requirements found</div>;
    }

    return (
        <Form {...form}>
            <div className="min-h-screen bg-muted flex flex-col">
                {/* Background decorative image */}
                <div className="absolute inset-0 pointer-events-none opacity-30 overflow-hidden">
                    <div className="absolute top-[53%] left-1/2 -translate-x-1/2 w-[90%] md:w-[45%] h-auto">
                        <div className="w-full aspect-[3/4] bg-gradient-to-b from-transparent via-muted/20 to-transparent rounded-full blur-3xl" />
                    </div>
                </div>

                {/* Main content */}
                <div className="relative z-10 flex flex-col items-center px-6 py-4 md:py-16">
                    <div className="w-full max-w-[720px] space-y-8 md:space-y-10">
                        {/* Header Section */}
                        <div className="space-y-1.5">
                            <Typography
                                variant="header"
                                tag="h1"
                                color="primary"
                                className="text-[36px] leading-[40px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] tracking-[0.06px]"
                            >
                                {t("events.registration.piiForm.title")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background-alt"
                                className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px]"
                            >
                                {t("events.registration.piiForm.subtitle")}
                            </Typography>
                        </div>

                        {/* Info Alert */}
                        <Alert
                            variant="info"
                            className="backdrop-blur-[2px] backdrop-filter bg-background/50 border-border"
                        >
                            <Info className="h-4 w-4" />
                            <AlertTitle>{t("events.registration.piiForm.privacyTitle")}</AlertTitle>
                            <AlertDescription>
                                {t("events.registration.piiForm.privacyDescription")}
                            </AlertDescription>
                        </Alert>

                        <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-8">
                            {/* Personal Information Section */}
                            <div className="space-y-4">
                                <Typography
                                    variant="header"
                                    tag="h2"
                                    color="primary"
                                    className="text-[24px] leading-[28px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                >
                                    {t("profile.personalInfo")}
                                </Typography>

                                <div className="space-y-2.5">
                                    {/* First Name */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="firstName"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground"
                                                            >
                                                                {t("profile.firstName")}
                                                            </Label>
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

                                    {/* Last Name */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="lastName"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground"
                                                            >
                                                                {t("profile.lastName")}
                                                            </Label>
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

                                    {/* Bio */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="bio"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground"
                                                            >
                                                                {t("profile.bio")}
                                                            </Label>
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
                                </div>
                            </div>

                            {/* Contact Information Section */}
                            <div className="space-y-4">
                                <Typography
                                    variant="header"
                                    tag="h2"
                                    color="primary"
                                    className="text-[24px] leading-[28px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                >
                                    {t("profile.contactInfo")}
                                </Typography>

                                <div className="space-y-2.5">
                                    {/* Email */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="email"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground"
                                                            >
                                                                {t("profile.email")}
                                                            </Label>
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

                                    {/* Phone Number */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="phoneNumber"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground"
                                                            >
                                                                {t("profile.phoneNumber")}
                                                            </Label>
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

                                    {/* Address */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="address"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground"
                                                            >
                                                                {t("profile.address")}
                                                            </Label>
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
                                </div>
                            </div>

                            {/* Academic Information Section */}
                            <div className="space-y-4">
                                <Typography
                                    variant="header"
                                    tag="h2"
                                    color="primary"
                                    className="text-[24px] leading-[28px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                >
                                    {t("profile.academicInfo")}
                                </Typography>

                                <div className="space-y-2.5">
                                    {/* Academic Email */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="academicEmail"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground"
                                                            >
                                                                {t("profile.academicEmail")}
                                                            </Label>
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

                                    {/* Academic Institution */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="academicInstitution"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-foreground"
                                                            >
                                                                {t("profile.academicInstitution")}
                                                            </Label>
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
                                </div>
                            </div>

                            {/* Action Buttons Section */}
                            <div className="space-y-3 pt-6">
                                {/* Submit Button */}
                                <Button
                                    type="submit"
                                    variant="primary"
                                    size="xl"
                                    className="w-full"
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

                                {/* Cancel Button */}
                                <Button
                                    type="button"
                                    onClick={onCancel}
                                    variant="secondary-light"
                                    size="xl"
                                    className="w-full"
                                    disabled={isSubmitting}
                                >
                                    {t("common.cancel")}
                                </Button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </Form>
    );
};
