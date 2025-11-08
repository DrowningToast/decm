import { describe, it, expect } from "vitest";
import { createEventFormSchema, eventFormSchema, eventFormEditSchema } from "./eventFormSchema";

describe("eventFormSchema", () => {
    const validImageFile = new File(["dummy content"], "test.png", { type: "image/png" });
    const oversizedFile = new File([new Array(6 * 1024 * 1024).join("a")], "large.png", {
        type: "image/png",
    });
    const invalidTypeFile = new File(["dummy content"], "test.pdf", { type: "application/pdf" });

    describe("Create Mode", () => {
        const schema = createEventFormSchema("create");

        it("validates a complete valid event form", () => {
            const validData = {
                name: "Test Event",
                shortDescription: "Short description of the event",
                description: "Long description",
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: new Date("2025-01-01"),
                endDate: new Date("2025-01-02"),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Venue Location",
            };

            const result = schema.safeParse(validData);
            expect(result.success).toBe(true);
        });

        it("requires event name", () => {
            const invalidData = {
                name: "",
                shortDescription: "Description",
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0].path).toContain("name");
            }
        });

        it("requires event name to be at least 3 characters", () => {
            const invalidData = {
                name: "AB",
                shortDescription: "Description",
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe("events.validation.nameMinLength");
            }
        });

        it("requires short description", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "",
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
        });

        it("limits short description to 255 characters", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "a".repeat(256),
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe(
                    "events.validation.shortDescriptionMaxLength",
                );
            }
        });

        it("requires event banner in create mode", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: null,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe(
                    "events.validation.eventBannerRequired",
                );
            }
        });

        it("requires event icon in create mode", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: validImageFile,
                eventIcon: null,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe("events.validation.eventIconRequired");
            }
        });

        it("validates image file size", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: oversizedFile,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe("events.validation.eventBannerSize");
            }
        });

        it("validates image file type", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: invalidTypeFile,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe("events.validation.eventBannerType");
            }
        });

        it("requires end date to be on or after start date", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: new Date("2025-01-02"),
                endDate: new Date("2025-01-01"),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe("events.validation.endDateAfterStart");
            }
        });

        it("allows end date to be same as start date", () => {
            const date = new Date("2025-01-01");
            const validData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: date,
                endDate: date,
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(validData);
            expect(result.success).toBe(true);
        });

        it("requires seats count to be at least 1", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 0,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe("events.validation.seatsCountMin");
            }
        });

        it("requires seats count to be an integer", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 10.5,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe("events.validation.seatsCountInteger");
            }
        });

        it("requires contact number with minimum length", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "123",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe(
                    "events.validation.contactNumberMinLength",
                );
            }
        });

        it("requires location with minimum length", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: validImageFile,
                eventIcon: validImageFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "AB",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe("events.validation.locationMinLength");
            }
        });
    });

    describe("Edit Mode", () => {
        const schema = createEventFormSchema("edit");

        it("makes image fields optional in edit mode", () => {
            const validData = {
                name: "Test Event",
                shortDescription: "Description",
                startDate: new Date("2025-01-01"),
                endDate: new Date("2025-01-02"),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(validData);
            expect(result.success).toBe(true);
        });

        it("still validates image files when provided", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: oversizedFile,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = schema.safeParse(invalidData);
            expect(result.success).toBe(false);
        });
    });

    describe("Exported Schemas", () => {
        it("eventFormSchema is in create mode", () => {
            const invalidData = {
                name: "Test Event",
                shortDescription: "Description",
                eventBanner: null,
                eventIcon: null,
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = eventFormSchema.safeParse(invalidData);
            expect(result.success).toBe(false);
        });

        it("eventFormEditSchema is in edit mode", () => {
            const validData = {
                name: "Test Event",
                shortDescription: "Description",
                startDate: new Date(),
                endDate: new Date(),
                seatsCount: 100,
                contactNumber: "1234567890",
                contactAddress: "123 Test Street",
                location: "Test Venue",
                googleMapQuery: "Test Query",
            };

            const result = eventFormEditSchema.safeParse(validData);
            expect(result.success).toBe(true);
        });
    });
});
