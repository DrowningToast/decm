import { describe, it, expect } from "vitest";
import { EntityEventType } from "@decm/api";
import {
    mapToRegistrationRequirementStatus,
    mapRegistrationRequirementStatusToNumber,
    mapEventTypeToEntityEventType,
    mapEntityEventRegistrationConfigRequirementStatus,
    mapEntityEventRegistrationConfigToEventRegistrationConfiguration,
    mapEntityEventRegistrationInvitationToEventRegistrationInvitation,
    mapRegistrationToJoinEventParticipant,
} from "./mapper";

const baseConfig = {
    id: "config-1",
    event_id: "event-1",
    first_name_requirement_status: 1,
    last_name_requirement_status: 1,
    email_requirement_status: 1,
    bio_requirement_status: 0,
    phone_number_requirement_status: 2,
    address_requirement_status: 0,
    academic_institution_requirement_status: 0,
    academic_email_requirement_status: 0,
    final_call_for_registration: undefined,
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z",
};

describe("mapToRegistrationRequirementStatus", () => {
    it("maps 0 to not_required", () => {
        expect(mapToRegistrationRequirementStatus(0)).toBe("not_required");
    });

    it("maps 1 to required", () => {
        expect(mapToRegistrationRequirementStatus(1)).toBe("required");
    });

    it("maps 2 to optional", () => {
        expect(mapToRegistrationRequirementStatus(2)).toBe("optional");
    });

    it("throws on invalid status", () => {
        expect(() => mapToRegistrationRequirementStatus(99)).toThrow(
            "Invalid registration requirement status: 99",
        );
    });
});

describe("mapRegistrationRequirementStatusToNumber", () => {
    it("maps not_required to 0", () => {
        expect(mapRegistrationRequirementStatusToNumber("not_required")).toBe(0);
    });

    it("maps required to 1", () => {
        expect(mapRegistrationRequirementStatusToNumber("required")).toBe(1);
    });

    it("maps optional to 2", () => {
        expect(mapRegistrationRequirementStatusToNumber("optional")).toBe(2);
    });

    it("throws on invalid status", () => {
        expect(() => mapRegistrationRequirementStatusToNumber("invalid" as never)).toThrow(
            "Invalid registration requirement status: invalid",
        );
    });
});

describe("mapEventTypeToEntityEventType", () => {
    it("maps private to EntityEventType.EventTypePrivate", () => {
        expect(mapEventTypeToEntityEventType("private")).toBe(EntityEventType.EventTypePrivate);
    });

    it("maps invite to EntityEventType.EventTypeInvite", () => {
        expect(mapEventTypeToEntityEventType("invite")).toBe(EntityEventType.EventTypeInvite);
    });

    it("throws on invalid event type", () => {
        expect(() => mapEventTypeToEntityEventType("public" as never)).toThrow(
            "Invalid event type: public",
        );
    });
});

describe("mapEntityEventRegistrationConfigRequirementStatus", () => {
    it("maps all requirement statuses correctly", () => {
        const result = mapEntityEventRegistrationConfigRequirementStatus(baseConfig);
        expect(result.firstName).toBe("required");
        expect(result.lastName).toBe("required");
        expect(result.email).toBe("required");
        expect(result.bio).toBe("not_required");
        expect(result.phoneNumber).toBe("optional");
        expect(result.address).toBe("not_required");
        expect(result.academicInstitution).toBe("not_required");
        expect(result.academicEmail).toBe("not_required");
    });
});

describe("mapEntityEventRegistrationConfigToEventRegistrationConfiguration", () => {
    it("maps private event type correctly", () => {
        const config = { ...baseConfig, event_type: EntityEventType.EventTypePrivate };
        const result = mapEntityEventRegistrationConfigToEventRegistrationConfiguration(config);
        expect(result.eventType).toBe("private");
    });

    it("maps invite event type correctly", () => {
        const config = { ...baseConfig, event_type: EntityEventType.EventTypeInvite };
        const result = mapEntityEventRegistrationConfigToEventRegistrationConfiguration(config);
        expect(result.eventType).toBe("invite");
    });

    it("maps event type 0 (numeric private) correctly", () => {
        const config = { ...baseConfig, event_type: 0 as never };
        const result = mapEntityEventRegistrationConfigToEventRegistrationConfiguration(config);
        expect(result.eventType).toBe("private");
    });

    it("maps event type 1 (numeric invite) correctly", () => {
        const config = { ...baseConfig, event_type: 1 as never };
        const result = mapEntityEventRegistrationConfigToEventRegistrationConfiguration(config);
        expect(result.eventType).toBe("invite");
    });

    it("maps finalCallForRegistration date when present", () => {
        const config = {
            ...baseConfig,
            event_type: EntityEventType.EventTypePrivate,
            final_call_for_registration: "2025-06-01T00:00:00Z",
        };
        const result = mapEntityEventRegistrationConfigToEventRegistrationConfiguration(config);
        expect(result.finalCallForRegistration).toBeInstanceOf(Date);
        expect(result.finalCallForRegistration?.toISOString()).toBe("2025-06-01T00:00:00.000Z");
    });

    it("leaves finalCallForRegistration undefined when absent", () => {
        const result = mapEntityEventRegistrationConfigToEventRegistrationConfiguration(baseConfig);
        expect(result.finalCallForRegistration).toBeUndefined();
    });

    it("leaves eventType undefined when event_type field is absent", () => {
        const result = mapEntityEventRegistrationConfigToEventRegistrationConfiguration(baseConfig);
        expect(result.eventType).toBeUndefined();
    });

    it("includes all requirement fields", () => {
        const result = mapEntityEventRegistrationConfigToEventRegistrationConfiguration(baseConfig);
        expect(result.firstName).toBe("required");
        expect(result.email).toBe("required");
        expect(result.bio).toBe("not_required");
        expect(result.phoneNumber).toBe("optional");
    });
});

describe("mapEntityEventRegistrationInvitationToEventRegistrationInvitation", () => {
    const baseInvitation = {
        id: "inv-1",
        event_id: "event-1",
        email: "user@example.com",
        first_name: "John",
        last_name: "Doe",
        phone_number: "0812345678",
        academic_institution: "MIT",
        valid_until: "2025-12-31T00:00:00Z",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-02T00:00:00Z",
        cancelled_at: undefined,
        accepted_at: undefined,
        inbox_message_id: "msg-1",
    };

    it("maps all fields correctly", () => {
        const result =
            mapEntityEventRegistrationInvitationToEventRegistrationInvitation(baseInvitation);
        expect(result.id).toBe("inv-1");
        expect(result.eventId).toBe("event-1");
        expect(result.email).toBe("user@example.com");
        expect(result.firstName).toBe("John");
        expect(result.lastName).toBe("Doe");
        expect(result.phoneNumber).toBe("0812345678");
        expect(result.academicInstitution).toBe("MIT");
        expect(result.inboxMessageId).toBe("msg-1");
    });

    it("maps validUntil date when present", () => {
        const result =
            mapEntityEventRegistrationInvitationToEventRegistrationInvitation(baseInvitation);
        expect(result.validUntil).toBeInstanceOf(Date);
        expect(result.validUntil?.toISOString()).toBe("2025-12-31T00:00:00.000Z");
    });

    it("leaves validUntil undefined when absent", () => {
        const inv = { ...baseInvitation, valid_until: undefined };
        const result = mapEntityEventRegistrationInvitationToEventRegistrationInvitation(inv);
        expect(result.validUntil).toBeUndefined();
    });

    it("maps cancelledAt date when present", () => {
        const inv = { ...baseInvitation, cancelled_at: "2024-06-01T00:00:00Z" };
        const result = mapEntityEventRegistrationInvitationToEventRegistrationInvitation(inv);
        expect(result.cancelledAt).toBeInstanceOf(Date);
        expect(result.cancelledAt?.toISOString()).toBe("2024-06-01T00:00:00.000Z");
    });

    it("leaves cancelledAt undefined when absent", () => {
        const result =
            mapEntityEventRegistrationInvitationToEventRegistrationInvitation(baseInvitation);
        expect(result.cancelledAt).toBeUndefined();
    });

    it("maps acceptedAt date when present", () => {
        const inv = { ...baseInvitation, accepted_at: "2024-07-01T00:00:00Z" };
        const result = mapEntityEventRegistrationInvitationToEventRegistrationInvitation(inv);
        expect(result.acceptedAt).toBeInstanceOf(Date);
    });

    it("leaves acceptedAt undefined when absent", () => {
        const result =
            mapEntityEventRegistrationInvitationToEventRegistrationInvitation(baseInvitation);
        expect(result.acceptedAt).toBeUndefined();
    });

    it("maps createdAt and updatedAt", () => {
        const result =
            mapEntityEventRegistrationInvitationToEventRegistrationInvitation(baseInvitation);
        expect(result.createdAt).toBeInstanceOf(Date);
        expect(result.updatedAt).toBeInstanceOf(Date);
    });
});

describe("mapRegistrationToJoinEventParticipant", () => {
    it("maps all registration form fields to payload", () => {
        const form = {
            firstName: "John",
            lastName: "Doe",
            email: "john@example.com",
            phoneNumber: "0812345678",
            address: "Bangkok",
            bio: "Engineer",
            academicInstitution: "MIT",
            academicEmail: "john@mit.edu",
        };
        const result = mapRegistrationToJoinEventParticipant(form);
        expect(result.first_name).toBe("John");
        expect(result.last_name).toBe("Doe");
        expect(result.email).toBe("john@example.com");
        expect(result.phone_number).toBe("0812345678");
        expect(result.address).toBe("Bangkok");
        expect(result.bio).toBe("Engineer");
        expect(result.academic_institution).toBe("MIT");
        expect(result.academic_email).toBe("john@mit.edu");
    });

    it("passes through undefined optional fields", () => {
        const form = {
            firstName: "Jane",
            lastName: "Smith",
            email: undefined,
            phoneNumber: undefined,
            address: undefined,
            bio: undefined,
            academicInstitution: undefined,
            academicEmail: undefined,
        };
        const result = mapRegistrationToJoinEventParticipant(form);
        expect(result.first_name).toBe("Jane");
        expect(result.last_name).toBe("Smith");
        expect(result.email).toBeUndefined();
        expect(result.phone_number).toBeUndefined();
    });
});
