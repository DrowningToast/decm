import type { RegistrationRequirement } from "@/services/EventRegistration/EventRegistration";
import { z } from "zod";

export const createRegistrationConfirmDataFormSchema = (
    t: (key: string) => string,
    registrationRequirement: Partial<RegistrationRequirement> | undefined,
) =>
    z.object({
        // Personal Information
        firstName: z
            .string()
            .min(3, { message: t("validation.firstNameMin3") })
            .max(32, { message: t("validation.firstNameMax32") })
            .optional()
            .or(z.literal(undefined))
            .transform((v) => (v === "" ? undefined : v))
            .refine(
                (value) => {
                    if (registrationRequirement?.firstName === "required") {
                        return value !== undefined && value !== "";
                    }
                    return true;
                },
                {
                    message: t("validation.firstNameRequired"),
                },
            ),
        lastName: z
            .string()
            .min(3, { message: t("validation.lastNameMin3") })
            .max(32, { message: t("validation.lastNameMax32") })
            .optional()
            .or(z.literal(undefined))
            .transform((v) => (v === "" ? undefined : v))
            .refine(
                (value) => {
                    if (registrationRequirement?.lastName === "required") {
                        return value !== undefined && value !== "";
                    }
                    return true;
                },
                {
                    message: t("validation.lastNameRequired"),
                },
            ),

        bio: z
            .string()
            .min(10, { message: t("validation.bioMin10") })
            .max(255, { message: t("validation.bioMax255") })
            .optional()
            .or(z.literal(undefined))
            .transform((v) => (v === "" ? undefined : v))
            .refine(
                (value) => {
                    if (registrationRequirement?.bio === "required") {
                        return value !== undefined && value !== "";
                    }
                    return true;
                },
                {
                    message: t("validation.bioRequired"),
                },
            ),

        // Contact Information
        email: z
            .string()
            .email({ message: t("validation.invalidEmail") })
            .max(64, { message: t("validation.emailMax64") })
            .optional()
            .or(z.literal(undefined))
            .transform((v) => (v === "" ? undefined : v))
            .refine(
                (value) => {
                    if (registrationRequirement?.email === "required") {
                        return value !== undefined && value !== "";
                    }
                    return true;
                },
                {
                    message: t("validation.emailRequired"),
                },
            ),

        phoneNumber: z
            .string()
            .max(15, { message: t("validation.phoneNumberMax15") })
            .optional()
            .or(z.literal(undefined))
            .transform((v) => (v === "" ? undefined : v))
            .refine(
                (value) => {
                    if (registrationRequirement?.phoneNumber === "required") {
                        return value !== undefined && value !== "";
                    }
                    return true;
                },
                {
                    message: t("validation.phoneNumberRequired"),
                },
            ),

        address: z
            .string()
            .min(10, { message: t("validation.addressMin10") })
            .max(255, { message: t("validation.addressMax255") })
            .optional()
            .or(z.literal(undefined))
            .transform((v) => (v === "" ? undefined : v))
            .refine(
                (value) => {
                    if (registrationRequirement?.address === "required") {
                        return value !== undefined && value !== "";
                    }
                    return true;
                },
                {
                    message: t("validation.addressRequired"),
                },
            ),
        // Academic Information
        academicEmail: z
            .string()
            .email({ message: t("validation.invalidEmail") })
            .optional()
            .or(z.literal(undefined))
            .transform((v) => (v === "" ? undefined : v))
            .refine(
                (value) => {
                    if (registrationRequirement?.academicEmail === "required") {
                        return value !== undefined && value !== "";
                    }
                    return true;
                },
                {
                    message: t("validation.academicEmailRequired"),
                },
            ),
        academicInstitution: z
            .string()
            .min(3, { message: t("validation.academicInstitutionMin3") })
            .max(255, { message: t("validation.academicInstitutionMax255") })
            .optional()
            .or(z.literal(undefined))
            .transform((v) => (v === "" ? undefined : v))
            .refine(
                (value) => {
                    if (registrationRequirement?.academicInstitution === "required") {
                        return value !== undefined && value !== "";
                    }
                    return true;
                },
                {
                    message: t("validation.academicInstitutionRequired"),
                },
            ),
    });

export type RegistrationConfirmDataForm = z.infer<
    ReturnType<typeof createRegistrationConfirmDataFormSchema>
>;
