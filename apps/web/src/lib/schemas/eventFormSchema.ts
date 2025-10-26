import { z } from "zod";

const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ACCEPTED_IMAGE_TYPES = ["image/jpeg", "image/jpg", "image/png", "image/webp"];

/**
 * Event Form Schema with conditional validation
 * Validates event creation/update form data
 * Image fields are required in create mode, optional in edit mode
 */
export const createEventFormSchema = (mode: "create" | "edit" = "create") => {
    const fileValidation = z
        .instanceof(File, { message: "events.validation.eventBannerInvalid" })
        .refine((file) => file.size <= MAX_FILE_SIZE, "events.validation.eventBannerSize")
        .refine(
            (file) => ACCEPTED_IMAGE_TYPES.includes(file.type),
            "events.validation.eventBannerType",
        )
        .nullable();

    return z
        .object({
            name: z
                .string()
                .min(1, "events.validation.nameRequired")
                .min(3, "events.validation.nameMinLength"),
            shortDescription: z
                .string()
                .min(1, "events.validation.shortDescriptionRequired")
                .max(255, "events.validation.shortDescriptionMaxLength"),
            description: z.string().optional(),
            eventBanner:
                mode === "create"
                    ? fileValidation.refine(
                          (file) => file !== null,
                          "events.validation.eventBannerRequired",
                      )
                    : fileValidation.optional(),
            eventIcon:
                mode === "create"
                    ? fileValidation.refine(
                          (file) => file !== null,
                          "events.validation.eventIconRequired",
                      )
                    : fileValidation.optional(),
            startDate: z.date({
                message: "events.validation.startDateRequired",
            }),
            endDate: z.date({
                message: "events.validation.endDateRequired",
            }),
            seatsCount: z
                .number({
                    message: "events.validation.seatsCountRequired",
                })
                .int("events.validation.seatsCountInteger")
                .min(1, "events.validation.seatsCountMin"),
            // Contact Information
            contactNumber: z
                .string()
                .min(1, "events.validation.contactNumberRequired")
                .min(10, "events.validation.contactNumberMinLength"),
            contactAddress: z
                .string()
                .min(1, "events.validation.contactAddressRequired")
                .min(10, "events.validation.contactAddressMinLength"),
            // Venue Information
            location: z
                .string()
                .min(1, "events.validation.locationRequired")
                .min(3, "events.validation.locationMinLength"),
            googleMapQuery: z
                .string()
                .min(1, "events.validation.googleMapQueryRequired")
                .min(3, "events.validation.googleMapQueryMinLength"),
        })
        .refine(
            (data) => {
                // endDate must be on or after startDate
                return data.endDate >= data.startDate;
            },
            {
                message: "events.validation.endDateAfterStart",
                path: ["endDate"],
            },
        );
};

/**
 * Event Form Schema for Create Mode
 * Validates event creation form data
 */
export const eventFormSchema = createEventFormSchema("create");

/**
 * Event Form Schema for Edit Mode
 * Validates event update form data - allows null values for images
 */
export const eventFormEditSchema = createEventFormSchema("edit");

/**
 * TypeScript type inferred from the create schema
 */
export type EventFormData = z.infer<typeof eventFormSchema>;

/**
 * TypeScript type inferred from the edit schema
 */
export type EventFormEditData = z.infer<typeof eventFormEditSchema>;

/**
 * Combined type that accepts both create and edit form data
 * This is the most flexible type that works for both modes
 */
export type EventFormUnion = EventFormData | EventFormEditData;
