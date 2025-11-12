import { describe, it, expect } from "vitest";
import {
    participantSettingsSchema,
    eventTypeEnum,
    fieldRequirementEnum,
    defaultParticipantSettings,
} from "./participantSettingsSchema";

describe("participantSettingsSchema", () => {
    const baseValidData = {
        eventType: "private" as const,
        isBookingRequired: false,
        isTicketTransferable: true,
        requireRegistrationPassword: false,
        firstName: "not_required" as const,
        lastName: "not_required" as const,
        email: "not_required" as const,
        bio: "not_required" as const,
        phoneNumber: "not_required" as const,
        address: "not_required" as const,
        academicInstitution: "not_required" as const,
        academicEmail: "not_required" as const,
    };

    describe("Basic Validation", () => {
        it("validates a complete valid participant settings object", () => {
            const result = participantSettingsSchema.safeParse(baseValidData);
            expect(result.success).toBe(true);
        });

        it("accepts all valid event types", () => {
            const eventTypes = ["private", "invite"] as const;

            eventTypes.forEach((eventType) => {
                const data = {
                    ...baseValidData,
                    eventType,
                };
                const result = participantSettingsSchema.safeParse(data);
                expect(result.success).toBe(true);
            });
        });

        it("accepts all valid field requirement values", () => {
            const requirements = ["not_required", "required", "optional"] as const;

            requirements.forEach((requirement) => {
                const data = {
                    ...baseValidData,
                    firstName: requirement,
                    lastName: requirement,
                    email: requirement,
                };
                const result = participantSettingsSchema.safeParse(data);
                expect(result.success).toBe(true);
            });
        });
    });

    describe("Registration Password Validation", () => {
        it("requires password when requireRegistrationPassword is true", () => {
            const invalidData = {
                ...baseValidData,
                requireRegistrationPassword: true,
                registrationPassword: "",
            };

            const result = participantSettingsSchema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe(
                    "participantSettings.validation.registrationPasswordRequired",
                );
                expect(result.error.issues[0]?.path).toContain("registrationPassword");
            }
        });

        it("accepts valid password when requireRegistrationPassword is true", () => {
            const validData = {
                ...baseValidData,
                requireRegistrationPassword: true,
                registrationPassword: "password123",
            };

            const result = participantSettingsSchema.safeParse(validData);
            expect(result.success).toBe(true);
        });

        it("requires password to be at least 8 characters", () => {
            const invalidData = {
                ...baseValidData,
                requireRegistrationPassword: true,
                registrationPassword: "pass123",
            };

            const result = participantSettingsSchema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe(
                    "participantSettings.validation.registrationPasswordMinLength",
                );
            }
        });

        it("validates password length even when requireRegistrationPassword is false but password is provided", () => {
            const invalidData = {
                ...baseValidData,
                requireRegistrationPassword: false,
                registrationPassword: "short",
            };

            const result = participantSettingsSchema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe(
                    "participantSettings.validation.registrationPasswordMinLength",
                );
            }
        });

        it("allows empty password when requireRegistrationPassword is false", () => {
            const validData = {
                ...baseValidData,
                requireRegistrationPassword: false,
                registrationPassword: "",
            };

            const result = participantSettingsSchema.safeParse(validData);
            expect(result.success).toBe(true);
        });

        it("allows undefined password when requireRegistrationPassword is false", () => {
            const validData = {
                ...baseValidData,
                requireRegistrationPassword: false,
            };

            const result = participantSettingsSchema.safeParse(validData);
            expect(result.success).toBe(true);
        });

        it("rejects password with only whitespace when required", () => {
            const invalidData = {
                ...baseValidData,
                requireRegistrationPassword: true,
                registrationPassword: "        ",
            };

            const result = participantSettingsSchema.safeParse(invalidData);
            expect(result.success).toBe(false);
        });
    });

    describe("Final Call Registration Date Validation", () => {
        it("requires final call date to be in the future", () => {
            const pastDate = new Date("2020-01-01");
            const invalidData = {
                ...baseValidData,
                finalCallRegistrationDate: pastDate,
            };

            const result = participantSettingsSchema.safeParse(invalidData);
            expect(result.success).toBe(false);
            if (!result.success) {
                expect(result.error.issues[0]?.message).toBe(
                    "participantSettings.validation.finalCallRegistrationDateFuture",
                );
                expect(result.error.issues[0]?.path).toContain("finalCallRegistrationDate");
            }
        });

        it("accepts future final call date", () => {
            const futureDate = new Date();
            futureDate.setFullYear(futureDate.getFullYear() + 1);

            const validData = {
                ...baseValidData,
                finalCallRegistrationDate: futureDate,
            };

            const result = participantSettingsSchema.safeParse(validData);
            expect(result.success).toBe(true);
        });

        it("allows undefined final call date", () => {
            const validData = {
                ...baseValidData,
            };

            const result = participantSettingsSchema.safeParse(validData);
            expect(result.success).toBe(true);
        });
    });

    describe("Boolean Fields", () => {
        it("accepts true for isBookingRequired", () => {
            const validData = {
                ...baseValidData,
                isBookingRequired: true,
            };

            const result = participantSettingsSchema.safeParse(validData);
            expect(result.success).toBe(true);
        });

        it("accepts false for isTicketTransferable", () => {
            const validData = {
                ...baseValidData,
                isTicketTransferable: false,
            };

            const result = participantSettingsSchema.safeParse(validData);
            expect(result.success).toBe(true);
        });

        it("defaults requireRegistrationPassword to false", () => {
            const dataWithoutPasswordFlag = {
                eventType: "private" as const,
                isBookingRequired: false,
                isTicketTransferable: true,
                firstName: "not_required" as const,
                lastName: "not_required" as const,
                email: "not_required" as const,
                bio: "not_required" as const,
                phoneNumber: "not_required" as const,
                address: "not_required" as const,
                academicInstitution: "not_required" as const,
                academicEmail: "not_required" as const,
            };

            const result = participantSettingsSchema.safeParse(dataWithoutPasswordFlag);
            expect(result.success).toBe(true);
            if (result.success) {
                expect(result.data.requireRegistrationPassword).toBe(false);
            }
        });
    });

    describe("Complex Scenarios", () => {
        it("validates private event with password and all required fields", () => {
            const futureDate = new Date();
            futureDate.setDate(futureDate.getDate() + 7);

            const validData = {
                eventType: "private" as const,
                isBookingRequired: true,
                isTicketTransferable: false,
                requireRegistrationPassword: true,
                registrationPassword: "securepassword123",
                finalCallRegistrationDate: futureDate,
                firstName: "required" as const,
                lastName: "required" as const,
                email: "required" as const,
                bio: "optional" as const,
                phoneNumber: "required" as const,
                address: "optional" as const,
                academicInstitution: "required" as const,
                academicEmail: "required" as const,
            };

            const result = participantSettingsSchema.safeParse(validData);
            expect(result.success).toBe(true);
        });

        it("validates invite-only event with minimal requirements", () => {
            const validData = {
                eventType: "invite" as const,
                isBookingRequired: false,
                isTicketTransferable: true,
                requireRegistrationPassword: false,
                firstName: "not_required" as const,
                lastName: "not_required" as const,
                email: "not_required" as const,
                bio: "not_required" as const,
                phoneNumber: "not_required" as const,
                address: "not_required" as const,
                academicInstitution: "not_required" as const,
                academicEmail: "not_required" as const,
            };

            const result = participantSettingsSchema.safeParse(validData);
            expect(result.success).toBe(true);
        });
    });

    describe("Enum Validation", () => {
        it("eventTypeEnum accepts valid values", () => {
            expect(eventTypeEnum.safeParse("private").success).toBe(true);
            expect(eventTypeEnum.safeParse("invite").success).toBe(true);
        });

        it("eventTypeEnum rejects invalid values", () => {
            expect(eventTypeEnum.safeParse("invalid").success).toBe(false);
            expect(eventTypeEnum.safeParse("").success).toBe(false);
            expect(eventTypeEnum.safeParse("public").success).toBe(false);
        });

        it("fieldRequirementEnum accepts valid values", () => {
            expect(fieldRequirementEnum.safeParse("not_required").success).toBe(true);
            expect(fieldRequirementEnum.safeParse("required").success).toBe(true);
            expect(fieldRequirementEnum.safeParse("optional").success).toBe(true);
        });

        it("fieldRequirementEnum rejects invalid values", () => {
            expect(fieldRequirementEnum.safeParse("invalid").success).toBe(false);
            expect(fieldRequirementEnum.safeParse("mandatory").success).toBe(false);
        });
    });

    describe("Default Values", () => {
        it("defaultParticipantSettings has correct structure", () => {
            expect(defaultParticipantSettings).toHaveProperty("eventType");
            expect(defaultParticipantSettings).toHaveProperty("isBookingRequired");
            expect(defaultParticipantSettings).toHaveProperty("isTicketTransferable");
            expect(defaultParticipantSettings).toHaveProperty("requireRegistrationPassword");
            expect(defaultParticipantSettings).toHaveProperty("firstName");
            expect(defaultParticipantSettings).toHaveProperty("lastName");
            expect(defaultParticipantSettings).toHaveProperty("email");
        });

        it("defaultParticipantSettings passes validation", () => {
            const result = participantSettingsSchema.safeParse(defaultParticipantSettings);
            expect(result.success).toBe(true);
        });

        it("defaultParticipantSettings has expected default values", () => {
            expect(defaultParticipantSettings.eventType).toBe("private");
            expect(defaultParticipantSettings.isBookingRequired).toBe(false);
            expect(defaultParticipantSettings.isTicketTransferable).toBe(true);
            expect(defaultParticipantSettings.requireRegistrationPassword).toBe(false);
            expect(defaultParticipantSettings.firstName).toBe("not_required");
        });
    });
});
