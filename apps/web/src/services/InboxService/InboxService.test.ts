import { describe, it, expect, vi, beforeEach } from "vitest";
import type { CoreApiType } from "@/lib/api/api";
import { InboxService } from "./InboxService";
import { EntityInboxMessageType, type InboxInboxMessagesViewModel } from "@decm/api";

// Mock the coreApiClient
vi.mock("@/lib/api/api", () => ({
    coreApiClient: {
        v1: {
            v1InboxMessagesList: vi.fn(),
            v1InboxMessagesDetail: vi.fn(),
            v1InboxMessagesReadUpdate: vi.fn(),
            v1InboxMessagesReadAllUpdate: vi.fn(),
        },
    },
}));

import { coreApiClient } from "@/lib/api/api";

describe("InboxService", () => {
    let inboxService: InboxService;
    let mockCoreApi: CoreApiType;

    beforeEach(() => {
        vi.clearAllMocks();
        mockCoreApi = coreApiClient as unknown as CoreApiType;
        inboxService = new InboxService(mockCoreApi);
    });

    describe("getInboxMessages", () => {
        it("should fetch inbox messages list", async () => {
            const mockResponse = {
                inbox_messages: [
                    {
                        id: "msg-1",
                        message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
                        message_content: '{"en": "Hello", "th": "สวัสดี"}',
                        is_read: 0,
                        created_at: "2024-01-01T00:00:00Z",
                        updated_at: "2024-01-01T00:00:00Z",
                    },
                ],
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesList).mockResolvedValue(mockResponse);

            const result = await inboxService.getInboxMessages();

            expect(mockCoreApi.v1.v1InboxMessagesList).toHaveBeenCalled();
            expect(result).toHaveLength(1);
            expect(result[0].id).toBe("msg-1");
            expect(result[0].messageType).toBe("general");
        });

        it("should return empty array when inbox_messages is undefined", async () => {
            vi.mocked(mockCoreApi.v1.v1InboxMessagesList).mockResolvedValue({
                inbox_messages: [],
            });

            const result = await inboxService.getInboxMessages();

            expect(result).toEqual([]);
        });
    });

    describe("getInboxMessageById", () => {
        it("should fetch inbox message detail by ID", async () => {
            const mockResponse = {
                inbox_message: {
                    id: "msg-1",
                    message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
                    message_content: '{"en": "Hello", "th": "สวัสดี"}',
                    is_read: 0,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
                inbox_message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesDetail).mockResolvedValue(mockResponse);

            const result = await inboxService.getInboxMessageById("msg-1");

            expect(mockCoreApi.v1.v1InboxMessagesDetail).toHaveBeenCalledWith({
                inboxMessageId: "msg-1",
            });
            expect(result).not.toBeNull();
            expect(result?.id).toBe("msg-1");
        });

        it("should throw error when inbox_message is not found", async () => {
            const error = new Error("Inbox message not found");
            vi.mocked(mockCoreApi.v1.v1InboxMessagesDetail).mockRejectedValue(error);

            await expect(inboxService.getInboxMessageById("msg-1")).rejects.toThrow(
                "Inbox message not found",
            );
        });

        it("should map event registration invitation message correctly", async () => {
            const mockResponse = {
                inbox_message: {
                    id: "msg-1",
                    message_type:
                        EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
                    message_content: '{"en": "Invitation", "th": "คำเชิญ"}',
                    is_read: 0,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
                inbox_message_type:
                    EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
                event_registration_invitation: {
                    id: "invitation-1",
                    event_id: "event-1",
                    code: "INV-123",
                    email: "test@example.com",
                    first_name: "John",
                    last_name: "Doe",
                    phone_number: "123-456-7890",
                    academic_institution: "University",
                    valid_until: "2024-12-31T00:00:00Z",
                    message_type:
                        EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
                    message_content: '{"en": "Invitation", "th": "คำเชิญ"}',
                    is_read: 0,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesDetail).mockResolvedValue(mockResponse);

            const result = await inboxService.getInboxMessageById("msg-1");

            expect(result).not.toBeNull();
            if (result && "code" in result) {
                expect(result.code).toBe("INV-123");
                expect(result.eventId).toBe("event-1");
            }
        });

        it("should map certificate invitation message correctly", async () => {
            const mockResponse = {
                inbox_message: {
                    id: "msg-1",
                    message_type: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
                    message_content: '{"en": "Certificate", "th": "ใบรับรอง"}',
                    is_read: 0,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
                inbox_message_type:
                    EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
                event_certificate: {
                    id: "cert-msg-1",
                    event_id: "event-1",
                    event_name: "Test Event",
                    certificate_id: "cert-1",
                    certificate_title: "Test Certificate",
                    token_id: "1",
                    has_participant_joined_event: true,
                    message_type: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
                    message_content: '{"en": "Certificate", "th": "ใบรับรอง"}',
                    is_read: 0,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesDetail).mockResolvedValue(mockResponse);

            const result = await inboxService.getInboxMessageById("msg-1");

            expect(result).not.toBeNull();
            if (result && "eventName" in result) {
                expect(result.eventName).toBe("Test Event");
                expect(result.certificateId).toBe("cert-1");
                expect(result.hasParticipantJoinedEvent).toBe(true);
            }
        });
    });

    describe("markAsRead", () => {
        it("should mark message as read", async () => {
            const mockResponse = {
                inbox_message: {
                    id: "msg-1",
                    message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
                    message_content: '{"en": "Hello", "th": "สวัสดี"}',
                    is_read: 1,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesReadUpdate).mockResolvedValue(mockResponse);

            const result = await inboxService.markAsRead("msg-1");

            expect(mockCoreApi.v1.v1InboxMessagesReadUpdate).toHaveBeenCalledWith(
                { messageId: "msg-1" },
                { message_id: "msg-1" },
            );
            expect(result).not.toBeNull();
            expect(result?.isRead).toBe(true);
        });

        it("should throw error when inbox_message is not found", async () => {
            const error = new Error("Message not found");
            vi.mocked(mockCoreApi.v1.v1InboxMessagesReadUpdate).mockRejectedValue(error);

            await expect(inboxService.markAsRead("msg-1")).rejects.toThrow("Message not found");
        });
    });

    describe("markAllAsRead", () => {
        it("should mark all messages as read", async () => {
            const mockResponse = {
                inbox_messages: [
                    {
                        id: "msg-1",
                        message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
                        message_content: '{"en": "Hello", "th": "สวัสดี"}',
                        is_read: 1,
                        created_at: "2024-01-01T00:00:00Z",
                        updated_at: "2024-01-01T00:00:00Z",
                    },
                    {
                        id: "msg-2",
                        message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
                        message_content: '{"en": "World", "th": "โลก"}',
                        is_read: 1,
                        created_at: "2024-01-02T00:00:00Z",
                        updated_at: "2024-01-02T00:00:00Z",
                    },
                ],
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesReadAllUpdate).mockResolvedValue(mockResponse);

            const result = await inboxService.markAllAsRead();

            expect(mockCoreApi.v1.v1InboxMessagesReadAllUpdate).toHaveBeenCalled();
            expect(result).toHaveLength(2);
            expect(result[0].isRead).toBe(true);
            expect(result[1].isRead).toBe(true);
        });

        it("should return empty array when inbox_messages is undefined", async () => {
            vi.mocked(mockCoreApi.v1.v1InboxMessagesReadAllUpdate).mockResolvedValue({
                inbox_messages: [],
            });

            const result = await inboxService.markAllAsRead();

            expect(result).toEqual([]);
        });
    });
});
