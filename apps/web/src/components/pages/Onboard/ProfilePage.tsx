import { useContext } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { FormField, FormItem, FormControl, FormMessage } from "@/components/ui/form";
import { OAuthOnboardContext } from "./OAuth/OAuthOnboardContext";
import { OnboardPageContext } from "@/pages/onboard/[method]";


export const ProfilePage: React.FC = () => {
    const { setStep } = useContext(OnboardPageContext);
    const { form } = useContext(OAuthOnboardContext);
    const { t } = useTranslation();

    const onConfirm = () => {
        setStep(3)
    }

    const onBack = () => {
        setStep(1)
    }

    return (
        <div className="min-h-screen bg-[#e9dede] flex flex-col items-center px-6 py-16 md:py-24 relative overflow-hidden">
            {/* Background decorative image - positioned absolutely */}
            <div className="absolute inset-0 pointer-events-none opacity-30 overflow-hidden">
                <div className="absolute top-[53%] left-1/2 -translate-x-1/2 w-[90%] md:w-[45%] h-auto">
                    {/* Lady Justice placeholder - using background color for now */}
                    <div className="w-full aspect-[3/4] bg-gradient-to-b from-transparent via-muted/20 to-transparent rounded-full blur-3xl" />
                </div>
            </div>

            {/* Main content - positioned relative to stay above background */}
            <div className="relative z-10 w-full max-w-[420px] space-y-8 md:space-y-9">
                {/* Header Section */}
                <div className="space-y-1.5">
                    <Typography
                        variant="header"
                        tag="h1"
                        color="primary"
                        className="text-[36px] leading-[40px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] tracking-[0.06px]"
                    >
                        {t("onboard.profile.title")}
                    </Typography>
                    <Typography
                        variant="text"
                        tag="p"
                        color="background-alt"
                        className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px]"
                    >
                        {t("onboard.profile.subtitle")}
                    </Typography>
                </div>

                {/* Form Fields Section */}
                <div className="flex gap-y-2.5 flex-col">
                    {/* First Name */}
                    <div className="space-y-1.5">
                        <FormField
                            control={form.control}
                            name="firstName"
                            render={({ field }) => (
                                <FormItem>
                                    <FormControl>
                                        <div className="space-y-1">
                                            <Label className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background">
                                                {t("onboard.profile.firstName")}
                                            </Label>
                                            <Input
                                                {...field}
                                                type="text"
                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background"
                                            />
                                        </div>
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="isFirstNamePublic"
                            render={({ field }) => (
                                <FormItem className="flex flex-row gap-x-1 items-center">
                                    <FormControl>
                                        <Checkbox
                                            checked={field.value}
                                            onCheckedChange={field.onChange}
                                            className="mt-0.5"
                                        />
                                    </FormControl>
                                    <Label className="text-xs        font-medium leading-normal text-background-alt cursor-pointer opacity-50">
                                        {t("onboard.profile.makePublic")}
                                    </Label>
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
                                            <Label className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background">
                                                {t("onboard.profile.lastName")}
                                            </Label>
                                            <Input
                                                {...field}
                                                type="text"
                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background"
                                            />
                                        </div>
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="isLastNamePublic"
                            render={({ field }) => (
                                <FormItem className="flex flex-row items-start gap-x-1 items-center">
                                    <FormControl>
                                        <Checkbox
                                            checked={field.value}
                                            onCheckedChange={field.onChange}
                                            className="mt-0.5"
                                        />
                                    </FormControl>
                                    <Label className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50">
                                        {t("onboard.profile.makePublic")}
                                    </Label>
                                </FormItem>
                            )}
                        />
                    </div>

                    {/* Contact Email */}
                    <div className="space-y-1.5">
                        <FormField
                            control={form.control}
                            name="email"
                            render={({ field }) => (
                                <FormItem>
                                    <FormControl>
                                        <div className="space-y-1">
                                            <Label className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background">
                                                {t("onboard.profile.email")}
                                            </Label>
                                            <Input
                                                {...field}
                                                type="email"
                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-foreground placeholder:text-background"
                                            />
                                        </div>
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="isEmailPublic"
                            render={({ field }) => (
                                <FormItem className="flex flex-row items-start gap-x-1 items-center">
                                    <FormControl>
                                        <Checkbox
                                            checked={field.value}
                                            onCheckedChange={field.onChange}
                                            className="mt-0.5"
                                        />
                                    </FormControl>
                                    <Label className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50">
                                        {t("onboard.profile.makePublic")}
                                    </Label>
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
                                            <Label className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background">
                                                {t("onboard.profile.phoneNumber")}
                                            </Label>
                                            <Input
                                                {...field}
                                                type="tel"
                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background"
                                            />
                                        </div>
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="isPhoneNumberPublic"
                            render={({ field }) => (
                                <FormItem className="flex flex-row items-start gap-x-1 items-center">
                                    <FormControl>
                                        <Checkbox
                                            checked={field.value}
                                            onCheckedChange={field.onChange}
                                            className="mt-0.5"
                                        />
                                    </FormControl>
                                    <Label className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50">
                                        {t("onboard.profile.makePublic")}
                                    </Label>
                                </FormItem>
                            )}
                        />
                    </div>
                </div>

                {/* Action Buttons Section */}
                <div className="space-y-3 pt-6 md:pt-8">
                    {/* Confirm Button */}
                    <Button
                        type="button"
                        onClick={onConfirm}
                        variant="primary"
                        size="xl"
                        className="w-full"
                    >
                        {t("common.confirm")}
                    </Button>

                    {/* Back Button - Only visible on desktop */}
                    <Button
                        type="button"
                        onClick={onBack}
                        variant="secondary-light"
                        size="xl"
                        className="w-full hidden md:flex"
                    >
                        {t("common.back")}
                    </Button>
                </div>
            </div>
        </div>
    );
};