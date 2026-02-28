import { describe, it, expect, vi, beforeEach } from "vitest";
import type { CoreApiType } from "@/lib/api/api";
import { EntityEventStatus, EntityEventType } from "@decm/api";
import { EventService } from "./EventService";

// Mock the coreApiClient
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            getEventById: vi.fn(),
            getEventsList: vi.fn(),
            getEventViewmodelById: vi.fn(),
            getEventsByOwnerCredentialsId: vi.fn(),
            getEventIssuersByEventId: vi.fn(),
            publishEventCertificates: vi.fn(),
        },
    },
}));

import { coreApiClient } from "@/lib/api/api";

describe("EventService", () => {
    let eventService: EventService;
    let mockCoreApi: CoreApiType;

    beforeEach(() => {
        vi.clearAllMocks();
        mockCoreApi = coreApiClient as unknown as CoreApiType;
        eventService = new EventService(mockCoreApi);
    });

    describe("getEventById", () => {
        it("should fetch event by ID and map to view model", async () => {
            const mockResponse = {
                id: "event-1",
                title: "Test Event",
                event_status: EntityEventStatus.EventStatusActive,
                event_type: EntityEventType.EventTypePrivate,
                chain_id: 1,
                contact_address: "123 Main St",
                contact_number: "123-456-7890",
                created_at: "2024-01-01T00:00:00Z",
                end_date: "2024-01-02T00:00:00Z",
                start_date: "2024-01-01T00:00:00Z",
                google_map_query: "Test Location",
                icon_presigned_url: "https://example.com/icon.png",
                banner_presigned_url: "https://example.com/banner.png",
                is_public: true,
                is_ticket_transferable: true,
                is_verified: true,
                is_booking_request_required: false,
                location: "Test Location",
                long_description: "Long description",
                max_attendees: 100,
                attendees_count: 50,
                owner_credential_id: "cred-1",
                short_description: "Short description",
                updated_at: "2024-01-01T00:00:00Z",
            };

            vi.mocked(mockCoreApi.v1.getEventById).mockResolvedValue(mockResponse);

            const result = await eventService.getEventById("event-1");

            expect(mockCoreApi.v1.getEventById).toHaveBeenCalledWith({ eventId: "event-1" });
            expect(result.id).toBe("event-1");
            expect(result.title).toBe("Test Event");
            expect(result.eventStatus).toBe("active");
        });
    });

    describe("getEvents", () => {
        it("should fetch events list with default params", async () => {
            const mockResponse = {
                events: [
                    {
                        id: "event-1",
                        title: "Test Event",
                        event_status: EntityEventStatus.EventStatusActive,
                        event_type: EntityEventType.EventTypePrivate,
                        chain_id: 1,
                        event_contract_address: "0x123",
                        banner_storage_key: "banner-key",
                        contact_address: "123 Main St",
                        contact_number: "123-456-7890",
                        google_map_query: "Test Location",
                        icon_storage_key: "icon-key",
                        is_booking_request_required: false,
                        created_at: "2024-01-01T00:00:00Z",
                        end_date: "2024-01-02T00:00:00Z",
                        start_date: "2024-01-01T00:00:00Z",
                        is_public: true,
                        is_ticket_transferable: true,
                        is_verified: true,
                        location: "Test Location",
                        long_description: "Long description",
                        max_attendees: 100,
                        attendees_count: 50,
                        owner_credential_id: "cred-1",
                        short_description: "Short description",
                        updated_at: "2024-01-01T00:00:00Z",
                    },
                ],
            };

            vi.mocked(mockCoreApi.v1.getEventsList).mockResolvedValue(mockResponse);

            const result = await eventService.getEvents({});

            expect(mockCoreApi.v1.getEventsList).toHaveBeenCalledWith({
                include_active_events: true,
                include_inactive_events: undefined,
                include_closed_events: undefined,
                only_user_joined_events: undefined,
            });
            expect(result).toHaveLength(1);
            expect(result[0].id).toBe("event-1");
        });

        it("should fetch events with custom params", async () => {
            const mockResponse = { events: [] };
            vi.mocked(mockCoreApi.v1.getEventsList).mockResolvedValue(mockResponse);

            await eventService.getEvents({
                includeActiveEvents: false,
                includeInactiveEvents: true,
                includeClosedEvents: true,
                onlyUserJoinedEvents: true,
            });

            expect(mockCoreApi.v1.getEventsList).toHaveBeenCalledWith({
                include_active_events: false,
                include_inactive_events: true,
                include_closed_events: true,
                only_user_joined_events: true,
            });
        });

        it("should return empty array when events is undefined", async () => {
            vi.mocked(mockCoreApi.v1.getEventsList).mockResolvedValue({ events: [] });

            const result = await eventService.getEvents({});

            expect(result).toEqual([]);
        });
    });

    describe("getEventViewModelExtended", () => {
        it("should fetch extended event view model", async () => {
            const mockResponse = {
                id: "event-1",
                title: "Test Event",
                event_status: EntityEventStatus.EventStatusActive,
                event_type: EntityEventType.EventTypePrivate,
                chain_id: 1,
                contact_address: "123 Main St",
                contact_number: "123-456-7890",
                google_map_query: "Test Location",
                created_at: "2024-01-01T00:00:00Z",
                end_date: "2024-01-02T00:00:00Z",
                start_date: "2024-01-01T00:00:00Z",
                icon_presigned_url: "https://example.com/icon.png",
                banner_presigned_url: "https://example.com/banner.png",
                is_public: true,
                is_ticket_transferable: true,
                is_verified: true,
                is_booking_request_required: false,
                location: "Test Location",
                long_description: "Long description",
                max_attendees: 100,
                attendees_count: 50,
                owner_credential_id: "cred-1",
                short_description: "Short description",
                updated_at: "2024-01-01T00:00:00Z",
                is_invited: true,
                is_joined: true,
                is_full: false,
                registration_config: {
                    id: "regconfig-1",
                    event_id: "event-1",
                    first_name_requirement_status: 1,
                    last_name_requirement_status: 1,
                    email_requirement_status: 1,
                    bio_requirement_status: 0,
                    phone_number_requirement_status: 0,
                    address_requirement_status: 0,
                    academic_institution_requirement_status: 0,
                    academic_email_requirement_status: 0,
                    final_call_for_registration: undefined,
                    is_identity_verification_required: false,
                    is_registration_password_required: false,
                },
                event_contract: {
                    access_manager_contract_address: "0xaccessmanager",
                    event_contract_address: "0xeventcontract",
                    event_id: "event-1",
                    id: "contract-1",
                },
            };

            vi.mocked(mockCoreApi.v1.getEventViewmodelById).mockResolvedValue(mockResponse);

            const result = await eventService.getEventViewModelExtended("event-1");

            expect(mockCoreApi.v1.getEventViewmodelById).toHaveBeenCalledWith({
                eventId: "event-1",
            });
            expect(result.id).toBe("event-1");
            expect(result.isInvited).toBe(true);
            expect(result.isJoined).toBe(true);
        });
    });

    describe("getEventsByOwnerCredentialId", () => {
        it("should fetch events by owner credential ID", async () => {
            const mockResponse = [
                {
                    id: "event-1",
                    title: "Test Event",
                    event_status: EntityEventStatus.EventStatusActive,
                    event_type: EntityEventType.EventTypePrivate,
                    chain_id: 1,
                    contact_address: "123 Main St",
                    contact_number: "123-456-7890",
                    google_map_query: "Test Location",
                    created_at: "2024-01-01T00:00:00Z",
                    end_date: "2024-01-02T00:00:00Z",
                    start_date: "2024-01-01T00:00:00Z",
                    icon_presigned_url: "https://example.com/icon.png",
                    banner_presigned_url: "https://example.com/banner.png",
                    is_public: true,
                    is_ticket_transferable: true,
                    is_verified: true,
                    is_booking_request_required: false,
                    location: "Test Location",
                    long_description: "Long description",
                    max_attendees: 100,
                    attendees_count: 50,
                    owner_credential_id: "cred-1",
                    short_description: "Short description",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            ];

            vi.mocked(mockCoreApi.v1.getEventsByOwnerCredentialsId).mockResolvedValue(mockResponse);

            const result = await eventService.getEventsByOwnerCredentialId("cred-1");

            expect(mockCoreApi.v1.getEventsByOwnerCredentialsId).toHaveBeenCalledWith({
                ownerCredentialId: "cred-1",
            });
            expect(result).toHaveLength(1);
            expect(result[0].ownerCredentialId).toBe("cred-1");
        });
    });

    describe("getEventIssuersByEventId", () => {
        it("should fetch event issuers by event ID", async () => {
            const mockResponse = [
                {
                    id: "issuer-1",
                    event_id: "event-1",
                    issuer_credential_id: "cred-1",
                    is_signed: 0,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            ];

            vi.mocked(mockCoreApi.v1.getEventIssuersByEventId).mockResolvedValue(mockResponse);

            const result = await eventService.getEventIssuersByEventId("event-1");

            expect(mockCoreApi.v1.getEventIssuersByEventId).toHaveBeenCalledWith({
                eventId: "event-1",
            });
            expect(result).toHaveLength(1);
            expect(result[0].eventId).toBe("event-1");
            expect(result[0].issuerCredentialId).toBe("cred-1");
        });
    });

    describe("publishEventCertificates", () => {
        it("should publish event certificates", async () => {
            const mockResponse = {
                event_id: "event-1",
                inbox_messages_created: 10,
                is_published: true,
                published_count: 10,
            };
            vi.mocked(mockCoreApi.v1.publishEventCertificates).mockResolvedValue(mockResponse);

            const result = await eventService.publishEventCertificates("event-1");

            expect(mockCoreApi.v1.publishEventCertificates).toHaveBeenCalledWith({
                eventId: "event-1",
            });
            expect(result).toEqual(mockResponse);
        });
    });
});
