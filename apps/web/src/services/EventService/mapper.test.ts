import { describe, it, expect, vi } from "vitest";
import {
    mapEntityEventStatusToEventStatus,
    mapEntityEventTypeToEventType,
    mapEventTypeToEntityEventType,
    mapEntityEventToEventItem,
    mapEventResponseToViewModel,
    mapEventViewModelExtended,
    mapToEventIssuer,
} from "./mapper";
import {
    EntityEventStatus,
    EntityEventType,
    type EntityEvent,
    type EventEventResponse,
    type EventEventViewModel,
    type EventEventIssuerResponse,
} from "@decm/api";

// Mock dependencies used by mapEventViewModelExtended
vi.mock("../EventRegistration/mapper", () => ({
    mapEntityEventRegistrationConfigToEventRegistrationConfiguration: vi.fn(() => ({
        firstName: "required",
    })),
}));

vi.mock("../AuthService/mapper", () => ({
    mapProfileViewModel: vi.fn((profile) => ({ id: profile.id, email: profile.email })),
}));

describe("mapEntityEventStatusToEventStatus", () => {
    it("maps Active to 'active'", () => {
        expect(mapEntityEventStatusToEventStatus(EntityEventStatus.EventStatusActive)).toBe(
            "active",
        );
    });

    it("maps Inactive to 'inactive'", () => {
        expect(mapEntityEventStatusToEventStatus(EntityEventStatus.EventStatusInactive)).toBe(
            "inactive",
        );
    });

    it("maps Closed to 'closed'", () => {
        expect(mapEntityEventStatusToEventStatus(EntityEventStatus.EventStatusClosed)).toBe(
            "closed",
        );
    });

    it("throws for invalid status", () => {
        expect(() => mapEntityEventStatusToEventStatus("unknown" as EntityEventStatus)).toThrow(
            "Invalid entity event status",
        );
    });
});

describe("mapEntityEventTypeToEventType", () => {
    it("maps Private to 'private'", () => {
        expect(mapEntityEventTypeToEventType(EntityEventType.EventTypePrivate)).toBe("private");
    });

    it("maps Invite to 'invite'", () => {
        expect(mapEntityEventTypeToEventType(EntityEventType.EventTypeInvite)).toBe("invite");
    });
});

describe("mapEventTypeToEntityEventType", () => {
    it("maps 'private' to EventTypePrivate", () => {
        expect(mapEventTypeToEntityEventType("private")).toBe(EntityEventType.EventTypePrivate);
    });

    it("maps 'invite' to EventTypeInvite", () => {
        expect(mapEventTypeToEntityEventType("invite")).toBe(EntityEventType.EventTypeInvite);
    });
});

describe("mapEntityEventToEventItem", () => {
    const baseEntity: EntityEvent = {
        attendees_count: 10,
        banner_storage_key: "banner-key",
        chain_id: 1,
        contact_address: "contact@test.com",
        contact_number: "123456",
        created_at: "2024-01-01T00:00:00Z",
        end_date: "2024-02-01T00:00:00Z",
        event_status: EntityEventStatus.EventStatusActive,
        event_type: EntityEventType.EventTypePrivate,
        google_map_query: "test location",
        icon_storage_key: "icon-key",
        id: "evt-1",
        is_booking_request_required: true,
        is_public: false,
        is_ticket_transferable: true,
        is_verified: true,
        location: "Bangkok",
        long_description: "Long desc",
        max_attendees: 100,
        owner_credential_id: "owner-1",
        short_description: "Short desc",
        start_date: "2024-01-15T00:00:00Z",
        title: "Test Event",
        updated_at: "2024-01-10T00:00:00Z",
    };

    it("maps all fields correctly", () => {
        const result = mapEntityEventToEventItem(baseEntity);

        expect(result.id).toBe("evt-1");
        expect(result.title).toBe("Test Event");
        expect(result.chainId).toBe(1);
        expect(result.eventStatus).toBe("active");
        expect(result.eventType).toBe("private");
        expect(result.attendeesCount).toBe(10);
        expect(result.maxAttendees).toBe(100);
        expect(result.isPublic).toBe(false);
        expect(result.isVerified).toBe(true);
        expect(result.location).toBe("Bangkok");
    });

    it("converts date strings to Date objects", () => {
        const result = mapEntityEventToEventItem(baseEntity);

        expect(result.startDate).toEqual(new Date("2024-01-15T00:00:00Z"));
        expect(result.endDate).toEqual(new Date("2024-02-01T00:00:00Z"));
        expect(result.createdAt).toEqual(new Date("2024-01-01T00:00:00Z"));
        expect(result.updatedAt).toEqual(new Date("2024-01-10T00:00:00Z"));
    });
});

describe("mapEventResponseToViewModel", () => {
    const baseResponse: EventEventResponse = {
        attendees_count: 5,
        banner_presigned_url: "https://example.com/banner.png",
        chain_id: 1,
        contact_address: "contact@test.com",
        contact_number: "123456",
        created_at: "2024-01-01T00:00:00Z",
        end_date: "2024-02-01T00:00:00Z",
        event_status: EntityEventStatus.EventStatusActive,
        event_type: EntityEventType.EventTypeInvite,
        google_map_query: "test query",
        icon_presigned_url: "https://example.com/icon.png",
        id: "evt-1",
        is_booking_request_required: false,
        is_public: true,
        is_ticket_transferable: false,
        is_verified: true,
        location: "Bangkok",
        long_description: "Long desc",
        max_attendees: 50,
        owner_credential_id: "owner-1",
        short_description: "Short desc",
        start_date: "2024-01-15T00:00:00Z",
        title: "Test Event",
        updated_at: "2024-01-10T00:00:00Z",
    };

    it("maps all fields correctly", () => {
        const result = mapEventResponseToViewModel(baseResponse);

        expect(result.id).toBe("evt-1");
        expect(result.title).toBe("Test Event");
        expect(result.eventStatus).toBe("active");
        expect(result.eventType).toBe("invite");
        expect(result.bannerPresignedUrl).toBe("https://example.com/banner.png");
        expect(result.iconPresignedUrl).toBe("https://example.com/icon.png");
        expect(result.googleMapQuery).toBe("test query");
    });

    it("converts date strings to Date objects", () => {
        const result = mapEventResponseToViewModel(baseResponse);

        expect(result.startDate).toEqual(new Date("2024-01-15T00:00:00Z"));
        expect(result.endDate).toEqual(new Date("2024-02-01T00:00:00Z"));
        expect(result.createdAt).toEqual(new Date("2024-01-01T00:00:00Z"));
        expect(result.updatedAt).toEqual(new Date("2024-01-10T00:00:00Z"));
    });
});

describe("mapEventViewModelExtended", () => {
    const baseViewModel = {
        attendees_count: 5,
        banner_presigned_url: "https://example.com/banner.png",
        chain_id: 1,
        contact_address: "contact@test.com",
        contact_number: "123456",
        created_at: "2024-01-01T00:00:00Z",
        end_date: "2024-02-01T00:00:00Z",
        event_status: EntityEventStatus.EventStatusActive,
        event_type: EntityEventType.EventTypePrivate,
        google_map_query: "test",
        icon_presigned_url: "https://example.com/icon.png",
        id: "evt-1",
        is_booking_request_required: false,
        is_public: true,
        is_ticket_transferable: false,
        is_verified: true,
        location: "Bangkok",
        long_description: "Long",
        max_attendees: 50,
        owner_credential_id: "owner-1",
        short_description: "Short",
        start_date: "2024-01-15T00:00:00Z",
        title: "Test Event",
        updated_at: "2024-01-10T00:00:00Z",
        is_invited: true,
        is_joined: true,
        is_full: false,
        joined_signature: "sig-123",
        joined_sign_message: "msg",
        registration_config: {
            id: "reg-1",
            event_id: "evt-1",
            first_name_requirement_status: 1,
            last_name_requirement_status: 0,
            email_requirement_status: 1,
            bio_requirement_status: 0,
            phone_number_requirement_status: 0,
            address_requirement_status: 0,
            academic_institution_requirement_status: 0,
            academic_email_requirement_status: 0,
            is_identity_verification_required: false,
            is_registration_password_required: false,
        },
        event_contract: {
            id: "contract-1",
            event_id: "evt-1",
            access_manager_contract_address: "0xAM",
            event_contract_address: "0xEC",
            ticket_contract_address: "0xTC",
            certificate_contract_address: "0xCC",
        },
    } as unknown as EventEventViewModel;

    it("maps extended fields", () => {
        const result = mapEventViewModelExtended(baseViewModel);

        expect(result.isInvited).toBe(true);
        expect(result.isJoined).toBe(true);
        expect(result.isFull).toBe(false);
        expect(result.joinedSignature).toBe("sig-123");
        expect(result.joinedSignMessage).toBe("msg");
    });

    it("maps contract addresses from event_contract", () => {
        const result = mapEventViewModelExtended(baseViewModel);

        expect(result.accessManagerContractAddress).toBe("0xAM");
        expect(result.eventContractAddress).toBe("0xEC");
        expect(result.ticketContractAddress).toBe("0xTC");
        expect(result.certificateContractAddress).toBe("0xCC");
    });

    it("maps optional date fields when present", () => {
        const result = mapEventViewModelExtended({
            ...baseViewModel,
            joined_at: "2024-01-20T00:00:00Z",
            broadcasted_at: "2024-01-25T00:00:00Z",
            broadcasted_at_estimated_deadline: "2024-02-01T00:00:00Z",
        } as unknown as EventEventViewModel);

        expect(result.joinedAt).toEqual(new Date("2024-01-20T00:00:00Z"));
        expect(result.broadcastedAt).toEqual(new Date("2024-01-25T00:00:00Z"));
        expect(result.broadcastedAtEstimatedDeadline).toEqual(new Date("2024-02-01T00:00:00Z"));
    });

    it("leaves optional date fields undefined when not present", () => {
        const result = mapEventViewModelExtended(baseViewModel);

        expect(result.joinedAt).toBeUndefined();
        expect(result.broadcastedAt).toBeUndefined();
        expect(result.broadcastedAtEstimatedDeadline).toBeUndefined();
    });

    it("maps signature_expired_at to signatureExpiredAt as a Date when present", () => {
        const result = mapEventViewModelExtended({
            ...baseViewModel,
            signature_expired_at: "2024-03-01T10:00:00Z",
        } as unknown as EventEventViewModel);

        expect(result.signatureExpiredAt).toEqual(new Date("2024-03-01T10:00:00Z"));
    });

    it("leaves signatureExpiredAt undefined when signature_expired_at is absent", () => {
        const result = mapEventViewModelExtended(baseViewModel);

        expect(result.signatureExpiredAt).toBeUndefined();
    });
});

describe("mapToEventIssuer", () => {
    const baseIssuer: EventEventIssuerResponse = {
        event_id: "evt-1",
        id: "issuer-1",
        is_signed: 1,
        issuer_credential_id: "cred-1",
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-10T00:00:00Z",
    };

    it("maps all fields correctly", () => {
        const result = mapToEventIssuer(baseIssuer);

        expect(result.eventId).toBe("evt-1");
        expect(result.id).toBe("issuer-1");
        expect(result.issuerCredentialId).toBe("cred-1");
        expect(result.createdAt).toEqual(new Date("2024-01-01T00:00:00Z"));
        expect(result.updatedAt).toEqual(new Date("2024-01-10T00:00:00Z"));
    });

    it("maps is_signed 1 to true", () => {
        const result = mapToEventIssuer(baseIssuer);
        expect(result.isSigned).toBe(true);
    });

    it("maps is_signed 0 to false", () => {
        const result = mapToEventIssuer({ ...baseIssuer, is_signed: 0 });
        expect(result.isSigned).toBe(false);
    });

    it("maps issuer_profile when present", () => {
        const result = mapToEventIssuer({
            ...baseIssuer,
            issuer_profile: {
                id: "profile-1",
                authentication_credential_id: "auth-1",
                email: "issuer@test.com",
                is_profile_picture_public: false,
                is_academic_email_public: false,
                is_academic_institution_public: false,
                is_address_public: false,
                is_bio_public: false,
                is_email_public: false,
                is_first_name_public: false,
                is_last_name_public: false,
                is_phone_number_public: false,
            },
        });
        expect(result.issuerProfile).toBeDefined();
        expect(result.issuerProfile).toEqual({ id: "profile-1", email: "issuer@test.com" });
    });

    it("leaves issuerProfile undefined when not present", () => {
        const result = mapToEventIssuer(baseIssuer);
        expect(result.issuerProfile).toBeUndefined();
    });
});
