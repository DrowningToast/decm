import { z } from "zod";

/**
 * Event type enum
 */
export const eventTypeEnum = z.enum(["private", "invite"]);
export type EventType = z.infer<typeof eventTypeEnum>;

/**
 * Field requirement enum
 */
export const fieldRequirementEnum = z.enum(["not_required", "required", "optional"]);
export type FieldRequirement = z.infer<typeof fieldRequirementEnum>;

/**
 * Participant settings form schema
 */
export const participantSettingsSchema = z
    .object({
        // Registration Settings
        eventType: eventTypeEnum,
        isBookingRequired: z.boolean(),
        isTicketTransferable: z.boolean(),
        requireRegistrationPassword: z.boolean().default(false),
        registrationPassword: z.string().optional(),
        finalCallRegistrationDate: z.date().optional(),

        // Participant Requirements
        firstName: fieldRequirementEnum,
        lastName: fieldRequirementEnum,
        email: fieldRequirementEnum,
        bio: fieldRequirementEnum,
        phoneNumber: fieldRequirementEnum,
        address: fieldRequirementEnum,
        academicInstitution: fieldRequirementEnum,
        academicEmail: fieldRequirementEnum,
    })
    .refine(
        (data) => {
            // If registration password is required, password must be provided
            if (data.requireRegistrationPassword) {
                return data.registrationPassword && data.registrationPassword.trim().length > 0;
            }
            return true;
        },
        {
            message: "participantSettings.validation.registrationPasswordRequired",
            path: ["registrationPassword"],
        },
    )
    .refine(
        (data) => {
            // If registration password is provided, it must meet minimum length
            if (data.registrationPassword && data.registrationPassword.trim().length > 0) {
                return data.registrationPassword.length >= 8;
            }
            return true;
        },
        {
            message: "participantSettings.validation.registrationPasswordMinLength",
            path: ["registrationPassword"],
        },
    )
    .refine(
        (data) => {
            // If final call registration date is provided, it must be in the future
            if (data.finalCallRegistrationDate) {
                return data.finalCallRegistrationDate > new Date();
            }
            return true;
        },
        {
            message: "participantSettings.validation.finalCallRegistrationDateFuture",
            path: ["finalCallRegistrationDate"],
        },
    );

export type ParticipantSettingsData = z.infer<typeof participantSettingsSchema>;

/**
 * Default values for participant settings
 */
export const defaultParticipantSettings: ParticipantSettingsData = {
    eventType: "private",
    isBookingRequired: false,
    isTicketTransferable: true,
    requireRegistrationPassword: false,
    registrationPassword: "",
    firstName: "not_required",
    lastName: "not_required",
    email: "not_required",
    bio: "not_required",
    phoneNumber: "not_required",
    address: "not_required",
    academicInstitution: "not_required",
    academicEmail: "not_required",
};
