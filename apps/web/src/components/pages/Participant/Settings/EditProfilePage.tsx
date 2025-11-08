import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { FormField, FormItem, FormControl, FormMessage, Form } from "@/components/ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { createEditProfileSchema } from "./editProfileSchema";
import type { EditProfileSchema } from "./editProfileSchema";
import { useTranslation } from "react-i18next";
import { useMyProfile } from "@/hooks/useMyProfile";
import { useUpdateProfile } from "@/hooks/profile/useUpdateProfile";
import { useNavigate } from "@/router";
import { useEffect } from "react";
import { Loader2 } from "lucide-react";
import { PrivateNavbar } from "@/components/layouts/navigations/PrivateNavbar";

export const EditProfilePage: React.FC = () => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { data: profile, isLoading: isLoadingProfile } = useMyProfile();
    const { updateProfile, isLoading: isUpdating } = useUpdateProfile();

    const EditProfileFormSchema = createEditProfileSchema(t);
    const form = useForm<EditProfileSchema>({
        resolver: zodResolver(EditProfileFormSchema),
        mode: "onChange",
        defaultValues: {
            first_name: "",
            is_first_name_public: false,
            last_name: "",
            is_last_name_public: false,
            bio: "",
            is_bio_public: false,
            email: "",
            is_email_public: false,
            phone_number: "",
            is_phone_number_public: false,
            address: "",
            is_address_public: false,
            academic_email: "",
            is_academic_email_public: false,
            academic_institution: "",
            is_academic_institution_public: false,
            profile_picture_url: "",
            is_profile_picture_public: false,
        },
    });

    // Load profile data into form
    useEffect(() => {
        if (profile) {
            form.reset({
                first_name: profile.first_name || "",
                is_first_name_public: profile.is_first_name_public || false,
                last_name: profile.last_name || "",
                is_last_name_public: profile.is_last_name_public || false,
                bio: profile.bio || "",
                is_bio_public: profile.is_bio_public || false,
                email: profile.email || "",
                is_email_public: profile.is_email_public || false,
                phone_number: profile.phone_number || "",
                is_phone_number_public: profile.is_phone_number_public || false,
                address: profile.address || "",
                is_address_public: profile.is_address_public || false,
                academic_email: profile.academic_email || "",
                is_academic_email_public: profile.is_academic_email_public || false,
                academic_institution: profile.academic_institution || "",
                is_academic_institution_public: profile.is_academic_institution_public || false,
                profile_picture_url: profile.profile_picture_url || "",
                is_profile_picture_public: profile.is_profile_picture_public || false,
            });
        }
    }, [profile, form]);

    const onSubmit = async (data: EditProfileSchema) => {
        try {
            await updateProfile({
                first_name: data.first_name || undefined,
                is_first_name_public: data.is_first_name_public,
                last_name: data.last_name || undefined,
                is_last_name_public: data.is_last_name_public,
                bio: data.bio || undefined,
                is_bio_public: data.is_bio_public,
                email: data.email || undefined,
                is_email_public: data.is_email_public,
                phone_number: data.phone_number || undefined,
                is_phone_number_public: data.is_phone_number_public,
                address: data.address || undefined,
                is_address_public: data.is_address_public,
                academic_email: data.academic_email || undefined,
                is_academic_email_public: data.is_academic_email_public,
                academic_institution: data.academic_institution || undefined,
                is_academic_institution_public: data.is_academic_institution_public,
                profile_picture_url: data.profile_picture_url || undefined,
                is_profile_picture_public: data.is_profile_picture_public,
            });
        } catch (error) {
            console.error("Failed to update profile:", error);
        }
    };

    if (isLoadingProfile) {
        return (
            <div className="min-h-screen bg-[#e9dede] flex items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
        );
    }

    return (
        <Form {...form}>
            <div className="min-h-screen bg-[#e9dede] flex flex-col">
                <PrivateNavbar />

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
                                {t("profile.editProfile")}
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                color="background-alt"
                                className="text-base leading-normal [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px]"
                            >
                                {t("profile.editSubtitle")}
                            </Typography>
                        </div>

                        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
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
                                            name="first_name"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background"
                                                            >
                                                                {t("profile.firstName")}
                                                            </Label>
                                                            <Input
                                                                {...field}
                                                                type="text"
                                                                id={field.name}
                                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background/50"
                                                            />
                                                        </div>
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="is_first_name_public"
                                            render={({ field }) => (
                                                <FormItem className="flex flex-row gap-x-1">
                                                    <FormControl>
                                                        <Checkbox
                                                            id="is_first_name_public"
                                                            checked={field.value}
                                                            onCheckedChange={field.onChange}
                                                            className="mt-0.5"
                                                        />
                                                    </FormControl>
                                                    <Label
                                                        htmlFor="is_first_name_public"
                                                        className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50"
                                                    >
                                                        {t("profile.makePublic")}
                                                    </Label>
                                                </FormItem>
                                            )}
                                        />
                                    </div>

                                    {/* Last Name */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="last_name"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background"
                                                            >
                                                                {t("profile.lastName")}
                                                            </Label>
                                                            <Input
                                                                {...field}
                                                                type="text"
                                                                id={field.name}
                                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background/50"
                                                            />
                                                        </div>
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="is_last_name_public"
                                            render={({ field }) => (
                                                <FormItem className="flex flex-row items-start gap-x-1">
                                                    <FormControl>
                                                        <Checkbox
                                                            id="is_last_name_public"
                                                            checked={field.value}
                                                            onCheckedChange={field.onChange}
                                                            className="mt-0.5"
                                                        />
                                                    </FormControl>
                                                    <Label
                                                        htmlFor="is_last_name_public"
                                                        className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50"
                                                    >
                                                        {t("profile.makePublic")}
                                                    </Label>
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
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background"
                                                            >
                                                                {t("profile.bio")}
                                                            </Label>
                                                            <Textarea
                                                                {...field}
                                                                id={field.name}
                                                                className="w-full min-h-24 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background/50"
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
                                        <FormField
                                            control={form.control}
                                            name="is_bio_public"
                                            render={({ field }) => (
                                                <FormItem className="flex flex-row items-start gap-x-1">
                                                    <FormControl>
                                                        <Checkbox
                                                            id="is_bio_public"
                                                            checked={field.value}
                                                            onCheckedChange={field.onChange}
                                                            className="mt-0.5"
                                                        />
                                                    </FormControl>
                                                    <Label
                                                        htmlFor="is_bio_public"
                                                        className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50"
                                                    >
                                                        {t("profile.makePublic")}
                                                    </Label>
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
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background"
                                                            >
                                                                {t("profile.email")}
                                                            </Label>
                                                            <Input
                                                                {...field}
                                                                type="email"
                                                                id={field.name}
                                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background/50"
                                                            />
                                                        </div>
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="is_email_public"
                                            render={({ field }) => (
                                                <FormItem className="flex flex-row items-start gap-x-1">
                                                    <FormControl>
                                                        <Checkbox
                                                            id="is_email_public"
                                                            checked={field.value}
                                                            onCheckedChange={field.onChange}
                                                            className="mt-0.5"
                                                        />
                                                    </FormControl>
                                                    <Label
                                                        htmlFor="is_email_public"
                                                        className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50"
                                                    >
                                                        {t("profile.makePublic")}
                                                    </Label>
                                                </FormItem>
                                            )}
                                        />
                                    </div>

                                    {/* Phone Number */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="phone_number"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background"
                                                            >
                                                                {t("profile.phoneNumber")}
                                                            </Label>
                                                            <Input
                                                                {...field}
                                                                type="tel"
                                                                id={field.name}
                                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background/50"
                                                            />
                                                        </div>
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="is_phone_number_public"
                                            render={({ field }) => (
                                                <FormItem className="flex flex-row items-start gap-x-1">
                                                    <FormControl>
                                                        <Checkbox
                                                            id="is_phone_number_public"
                                                            checked={field.value}
                                                            onCheckedChange={field.onChange}
                                                            className="mt-0.5"
                                                        />
                                                    </FormControl>
                                                    <Label
                                                        htmlFor="is_phone_number_public"
                                                        className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50"
                                                    >
                                                        {t("profile.makePublic")}
                                                    </Label>
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
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background"
                                                            >
                                                                {t("profile.address")}
                                                            </Label>
                                                            <Textarea
                                                                {...field}
                                                                id={field.name}
                                                                className="w-full min-h-20 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background/50"
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
                                        <FormField
                                            control={form.control}
                                            name="is_address_public"
                                            render={({ field }) => (
                                                <FormItem className="flex flex-row items-start gap-x-1">
                                                    <FormControl>
                                                        <Checkbox
                                                            id="is_address_public"
                                                            checked={field.value}
                                                            onCheckedChange={field.onChange}
                                                            className="mt-0.5"
                                                        />
                                                    </FormControl>
                                                    <Label
                                                        htmlFor="is_address_public"
                                                        className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50"
                                                    >
                                                        {t("profile.makePublic")}
                                                    </Label>
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
                                            name="academic_email"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background"
                                                            >
                                                                {t("profile.academicEmail")}
                                                            </Label>
                                                            <Input
                                                                {...field}
                                                                type="email"
                                                                id={field.name}
                                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background/50"
                                                            />
                                                        </div>
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="is_academic_email_public"
                                            render={({ field }) => (
                                                <FormItem className="flex flex-row items-start gap-x-1">
                                                    <FormControl>
                                                        <Checkbox
                                                            id="is_academic_email_public"
                                                            checked={field.value}
                                                            onCheckedChange={field.onChange}
                                                            className="mt-0.5"
                                                        />
                                                    </FormControl>
                                                    <Label
                                                        htmlFor="is_academic_email_public"
                                                        className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50"
                                                    >
                                                        {t("profile.makePublic")}
                                                    </Label>
                                                </FormItem>
                                            )}
                                        />
                                    </div>

                                    {/* Academic Institution */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="academic_institution"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background"
                                                            >
                                                                {t("profile.academicInstitution")}
                                                            </Label>
                                                            <Input
                                                                {...field}
                                                                type="text"
                                                                id={field.name}
                                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background/50"
                                                            />
                                                        </div>
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="is_academic_institution_public"
                                            render={({ field }) => (
                                                <FormItem className="flex flex-row items-start gap-x-1">
                                                    <FormControl>
                                                        <Checkbox
                                                            id="is_academic_institution_public"
                                                            checked={field.value}
                                                            onCheckedChange={field.onChange}
                                                            className="mt-0.5"
                                                        />
                                                    </FormControl>
                                                    <Label
                                                        htmlFor="is_academic_institution_public"
                                                        className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50"
                                                    >
                                                        {t("profile.makePublic")}
                                                    </Label>
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                </div>
                            </div>

                            {/* Profile Picture Section */}
                            <div className="space-y-4">
                                <Typography
                                    variant="header"
                                    tag="h2"
                                    color="primary"
                                    className="text-[24px] leading-[28px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px]"
                                >
                                    {t("profile.profilePicture")}
                                </Typography>

                                <div className="space-y-2.5">
                                    {/* Profile Picture URL */}
                                    <div className="space-y-1.5">
                                        <FormField
                                            control={form.control}
                                            name="profile_picture_url"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <div className="space-y-1">
                                                            <Label
                                                                htmlFor={field.name}
                                                                className="text-base leading-[15px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] tracking-[0.06px] font-normal text-background"
                                                            >
                                                                {t("profile.profilePictureUrl")}
                                                            </Label>
                                                            <Input
                                                                {...field}
                                                                type="url"
                                                                id={field.name}
                                                                placeholder={t(
                                                                    "profile.profilePictureUrlPlaceholder",
                                                                )}
                                                                className="w-full h-12 backdrop-blur-[2px] backdrop-filter bg-[rgba(252,252,252,0.5)] border-[#b8b8b8] border-[0.5px] rounded-[12px] text-background placeholder:text-background/50"
                                                            />
                                                        </div>
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="is_profile_picture_public"
                                            render={({ field }) => (
                                                <FormItem className="flex flex-row items-start gap-x-1">
                                                    <FormControl>
                                                        <Checkbox
                                                            id="is_profile_picture_public"
                                                            checked={field.value}
                                                            onCheckedChange={field.onChange}
                                                            className="mt-0.5"
                                                        />
                                                    </FormControl>
                                                    <Label
                                                        htmlFor="is_profile_picture_public"
                                                        className="text-xs font-medium leading-normal text-background-alt cursor-pointer opacity-50"
                                                    >
                                                        {t("profile.makePublic")}
                                                    </Label>
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                </div>
                            </div>

                            {/* Action Buttons Section */}
                            <div className="space-y-3 pt-6">
                                {/* Save Button */}
                                <Button
                                    type="submit"
                                    variant="primary"
                                    size="xl"
                                    className="w-full"
                                    disabled={isUpdating}
                                >
                                    {isUpdating ? (
                                        <>
                                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                            {t("common.saving")}
                                        </>
                                    ) : (
                                        t("common.save")
                                    )}
                                </Button>

                                {/* Cancel Button */}
                                <Button
                                    type="button"
                                    onClick={() => navigate("/app")}
                                    variant="secondary-light"
                                    size="xl"
                                    className="w-full"
                                    disabled={isUpdating}
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
