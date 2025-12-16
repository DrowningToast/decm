import { describe, it, expect, vi, beforeEach } from "vitest";
import type { CoreApiType } from "@/lib/api/api";
import { EventRegistrationService } from "./EventRegistration";
import { EntityInboxMessageType } from "@decm/api";

// Mock the coreApiClient
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            getEventRegistrationConfig: vi.fn(),
            checkEventPassword: vi.fn(),
            getEventRegistrationInvitationsByEventId: vi.fn(),
            getEventRegistrationInvitationByUserAndEvent: vi.fn(),
            updateEventRegistrationConfig: vi.fn(),
            joinEvent: vi.fn(),
        },
    },
}));

import { coreApiClient } from "@/lib/api/api";

describe("EventRegistrationService", () => {
    let eventRegistrationService: EventRegistrationService;
    let mockCoreApi: CoreApiType;

    beforeEach(() => {
        vi.clearAllMocks();
        mockCoreApi = coreApiClient as unknown as CoreApiType;
        eventRegistrationService = new EventRegistrationService(mockCoreApi);
    });

    describe("getConfiguration", () => {
        it("should fetch event registration configuration", async () => {
            const mockResponse = {
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

            vi.mocked(mockCoreApi.v1.getEventRegistrationConfig).mockResolvedValue(mockResponse);

            const result = await eventRegistrationService.getConfiguration("event-1");

            expect(mockCoreApi.v1.getEventRegistrationConfig).toHaveBeenCalledWith({
                eventId: "event-1",
            });
            expect(result.eventType).toBe("private");
            expect(result.firstName).toBe("required");
            expect(result.phoneNumber).toBe("optional");
        });
    });

    describe("checkPassword", () => {
        it("should check event password and return true when valid", async () => {
            vi.mocked(mockCoreApi.v1.checkEventPassword).mockResolvedValue({ is_valid: true });

            const result = await eventRegistrationService.checkPassword("event-1", "password123");

            expect(mockCoreApi.v1.checkEventPassword).toHaveBeenCalledWith(
                { eventId: "event-1" },
                { password: "password123" },
            );
            expect(result).toBe(true);
        });

        it("should return false when password is invalid", async () => {
            vi.mocked(mockCoreApi.v1.checkEventPassword).mockResolvedValue({ is_valid: false });

            const result = await eventRegistrationService.checkPassword("event-1", "wrong");

            expect(result).toBe(false);
        });
    });

    describe("getInvitationByEventId", () => {
        it("should fetch invitations by event ID", async () => {
            const mockResponse = [
                {
                    id: "inv-1",
                    event_id: "event-1",
                    inbox_message_id: "msg-1",
                    accepted_at: "",
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            ];

            vi.mocked(mockCoreApi.v1.getEventRegistrationInvitationsByEventId).mockResolvedValue(
                mockResponse,
            );

            const result = await eventRegistrationService.getInvitationByEventId("event-1");

            expect(mockCoreApi.v1.getEventRegistrationInvitationsByEventId).toHaveBeenCalledWith({
                eventId: "event-1",
            });
            expect(result).toHaveLength(1);
            expect(result[0].id).toBe("inv-1");
            expect(result[0].eventId).toBe("event-1");
        });
    });

    describe("getInvitationOfUserAndEventId", () => {
        it("should fetch invitation of user and event", async () => {
            const mockResponse = {
                registration_invitation: {
                    id: "inv-1",
                    event_id: "event-1",
                    inbox_message_id: "msg-1",
                    accepted_at: "",
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
                inbox: {
                    id: "msg-1",
                    message_type:
                        EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
                    message_content: "{}",
                    is_read: 0,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            };

            vi.mocked(
                mockCoreApi.v1.getEventRegistrationInvitationByUserAndEvent,
            ).mockResolvedValue(mockResponse);

            const result = await eventRegistrationService.getInvitationOfUserAndEventId("event-1");

            expect(
                mockCoreApi.v1.getEventRegistrationInvitationByUserAndEvent,
            ).toHaveBeenCalledWith({
                eventId: "event-1",
            });
            expect(result.registrationInvitation.id).toBe("inv-1");
            expect(result.inbox.id).toBe("msg-1");
        });

        it("should throw error when invitation or inbox is missing", async () => {
            vi.mocked(
                mockCoreApi.v1.getEventRegistrationInvitationByUserAndEvent,
            ).mockResolvedValue({
                registration_invitation: undefined,
                inbox: undefined,
            });

            await expect(
                eventRegistrationService.getInvitationOfUserAndEventId("event-1"),
            ).rejects.toThrow("Registration invitation or inbox not found");
        });
    });

    describe("updateConfiguration", () => {
        it("should update event registration configuration", async () => {
            const configuration = {
                eventType: "private" as const,
                firstName: "required" as const,
                lastName: "required" as const,
                email: "required" as const,
                bio: "not_required" as const,
                phoneNumber: "optional" as const,
                address: "not_required" as const,
                academicInstitution: "not_required" as const,
                academicEmail: "not_required" as const,
            };

            vi.mocked(mockCoreApi.v1.updateEventRegistrationConfig).mockResolvedValue({
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
            });

            await eventRegistrationService.updateConfiguration("event-1", configuration);

            expect(mockCoreApi.v1.updateEventRegistrationConfig).toHaveBeenCalledWith(
                { eventId: "event-1" },
                {
                    academic_email_requirement_status: 0,
                    academic_institution_requirement_status: 0,
                    address_requirement_status: 0,
                    bio_requirement_status: 0,
                    email_requirement_status: 1,
                    first_name_requirement_status: 1,
                    last_name_requirement_status: 1,
                    phone_number_requirement_status: 2,
                    event_type: "private",
                    is_booking_request_required: false,
                    is_ticket_transferable: false,
                },
            );
        });

        it("should throw error when eventType is missing", async () => {
            const configuration = {
                firstName: "required" as const,
                lastName: "required" as const,
                email: "required" as const,
                bio: "not_required" as const,
                phoneNumber: "optional" as const,
                address: "not_required" as const,
                academicInstitution: "not_required" as const,
                academicEmail: "not_required" as const,
            };

            await expect(
                eventRegistrationService.updateConfiguration(
                    "event-1",
                    configuration as unknown as Parameters<
                        typeof eventRegistrationService.updateConfiguration
                    >[1],
                ),
            ).rejects.toThrow("eventType is required");
        });
    });

    describe("joinEventWithAccountPassword", () => {
        it("should join event with account password", async () => {
            const mockResponse = {
                id: "attendee-1",
                event_id: "event-1",
                attendee_credential_id: "cred-1",
                contract_address: "0x123",
                wallet_address: "0xabc",
                created_at: "2024-01-01T00:00:00Z",
                is_attendee_accepted: true,
                updated_at: "2024-01-01T00:00:00Z",
            };

            vi.mocked(mockCoreApi.v1.joinEvent).mockResolvedValue(mockResponse);

            const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});

            await eventRegistrationService.joinEventWithAccountPassword({
                eventId: "event-1",
                eventPassword: "event-pass",
                accountPassword: "account-pass",
                registrationData: {
                    firstName: "John",
                    lastName: "Doe",
                    email: "john@example.com",
                    bio: undefined,
                    phoneNumber: undefined,
                    address: undefined,
                    academicEmail: undefined,
                    academicInstitution: undefined,
                },
            });

            expect(mockCoreApi.v1.joinEvent).toHaveBeenCalledWith(
                { eventId: "event-1" },
                {
                    account_password: "account-pass",
                    event_password: "event-pass",
                    registration_data: {
                        first_name: "John",
                        last_name: "Doe",
                        email: "john@example.com",
                    },
                },
            );
            expect(consoleSpy).toHaveBeenCalledWith("Join event successful:", mockResponse);

            consoleSpy.mockRestore();
        });
    });

    describe("joinEventWithSignature", () => {
        it("should join event with signature", async () => {
            vi.mocked(mockCoreApi.v1.joinEvent).mockResolvedValue({
                id: "attendee-1",
                event_id: "event-1",
                attendee_credential_id: "cred-1",
                contract_address: "0x123",
                wallet_address: "0xabc",
                created_at: "2024-01-01T00:00:00Z",
                is_attendee_accepted: true,
                updated_at: "2024-01-01T00:00:00Z",
            });

            await eventRegistrationService.joinEventWithSignature({
                eventId: "event-1",
                originalSignMessage: "Sign this message",
                signature: "sig-123",
                registrationData: {
                    firstName: "John",
                    lastName: "Doe",
                    email: "john@example.com",
                    bio: undefined,
                    phoneNumber: undefined,
                    address: undefined,
                    academicEmail: undefined,
                    academicInstitution: undefined,
                },
            });

            expect(mockCoreApi.v1.joinEvent).toHaveBeenCalledWith(
                { eventId: "event-1" },
                {
                    sign_message: "Sign this message",
                    signature: "sig-123",
                    registration_data: {
                        first_name: "John",
                        last_name: "Doe",
                        email: "john@example.com",
                    },
                },
            );
        });
    });
});
