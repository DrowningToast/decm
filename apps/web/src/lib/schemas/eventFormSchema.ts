import { z } from "zod";

const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ACCEPTED_IMAGE_TYPES = ["image/jpeg", "image/jpg", "image/png", "image/webp"];

/**
 * Event Form Schema
 * Validates event creation/update form data
 */
export const eventFormSchema = z
  .object({
    name: z
      .string()
      .min(1, "events.validation.nameRequired")
      .min(3, "events.validation.nameMinLength"),
    description: z.string().optional(),
    eventBanner: z
      .instanceof(File, { message: "events.validation.eventBannerRequired" })
      .refine((file) => file.size <= MAX_FILE_SIZE, "events.validation.eventBannerSize")
      .refine(
        (file) => ACCEPTED_IMAGE_TYPES.includes(file.type),
        "events.validation.eventBannerType"
      )
      .optional(),
    eventIcon: z
      .instanceof(File, { message: "events.validation.eventIconRequired" })
      .refine((file) => file.size <= MAX_FILE_SIZE, "events.validation.eventIconSize")
      .refine((file) => ACCEPTED_IMAGE_TYPES.includes(file.type), "events.validation.eventIconType")
      .optional(),
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
  })
  .refine(
    (data) => {
      return data.endDate > data.startDate;
    },
    {
      message: "events.validation.endDateAfterStart",
      path: ["endDate"],
    }
  );

/**
 * TypeScript type inferred from the schema
 */
export type EventFormData = z.infer<typeof eventFormSchema>;
