import { z } from "zod";

/**
 * Event type enum
 */
export const eventTypeEnum = z.enum(["public", "private"]);
export type EventType = z.infer<typeof eventTypeEnum>;

/**
 * Field requirement enum
 */
export const fieldRequirementEnum = z.enum(["not_required", "required", "optional"]);
export type FieldRequirement = z.infer<typeof fieldRequirementEnum>;

/**
 * Participant settings form schema
 */
export const participantSettingsSchema = z.object({
    // Registration Settings
    eventType: eventTypeEnum,
    isBookingRequired: z.boolean(),
    isTicketTransferable: z.boolean(),

    // Participant Requirements
    firstName: fieldRequirementEnum,
    lastName: fieldRequirementEnum,
    email: fieldRequirementEnum,
    bio: fieldRequirementEnum,
    phoneNumber: fieldRequirementEnum,
    address: fieldRequirementEnum,
    academicInstitution: fieldRequirementEnum,
    academicEmail: fieldRequirementEnum,
});

export type ParticipantSettingsData = z.infer<typeof participantSettingsSchema>;

/**
 * Default values for participant settings
 */
export const defaultParticipantSettings: ParticipantSettingsData = {
    eventType: "public",
    isBookingRequired: false,
    isTicketTransferable: true,
    firstName: "not_required",
    lastName: "not_required",
    email: "not_required",
    bio: "not_required",
    phoneNumber: "not_required",
    address: "not_required",
    academicInstitution: "not_required",
    academicEmail: "not_required",
};
