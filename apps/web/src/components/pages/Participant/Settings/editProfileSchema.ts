import { z } from "zod";

export const createEditProfileSchema = (t: (key: string) => string) =>
    z.object({
        // Personal Information
        first_name: z
            .string()
            .min(3, { message: t("validation.firstNameMin3") })
            .max(32, { message: t("validation.firstNameMax32") })
            .optional()
            .or(z.literal("")),
        is_first_name_public: z.boolean().default(false).optional(),

        last_name: z
            .string()
            .min(3, { message: t("validation.lastNameMin3") })
            .max(32, { message: t("validation.lastNameMax32") })
            .optional()
            .or(z.literal("")),
        is_last_name_public: z.boolean().default(false).optional(),

        bio: z
            .string()
            .min(10, { message: t("validation.bioMin10") })
            .max(255, { message: t("validation.bioMax255") })
            .optional()
            .or(z.literal("")),
        is_bio_public: z.boolean().default(false).optional(),

        // Contact Information
        email: z
            .string()
            .email({ message: t("validation.invalidEmail") })
            .max(64, { message: t("validation.emailMax64") })
            .optional()
            .or(z.literal("")),
        is_email_public: z.boolean().default(false).optional(),

        phone_number: z
            .string()
            .max(20, { message: t("validation.phoneNumberMax20") })
            .optional()
            .or(z.literal("")),
        is_phone_number_public: z.boolean().default(false).optional(),

        address: z
            .string()
            .min(10, { message: t("validation.addressMin10") })
            .max(255, { message: t("validation.addressMax255") })
            .optional()
            .or(z.literal("")),
        is_address_public: z.boolean().default(false).optional(),

        // Academic Information
        academic_email: z
            .string()
            .email({ message: t("validation.invalidEmail") })
            .optional()
            .or(z.literal("")),
        is_academic_email_public: z.boolean().default(false).optional(),

        academic_institution: z
            .string()
            .min(3, { message: t("validation.academicInstitutionMin3") })
            .max(255, { message: t("validation.academicInstitutionMax255") })
            .optional()
            .or(z.literal("")),
        is_academic_institution_public: z.boolean().default(false).optional(),

        // Profile Picture
        profile_picture_url: z
            .string()
            .url({ message: t("validation.invalidUrl") })
            .optional()
            .or(z.literal("")),
        is_profile_picture_public: z.boolean().default(false).optional(),
    });

export type EditProfileSchema = z.infer<ReturnType<typeof createEditProfileSchema>>;
