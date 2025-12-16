import { describe, it, expect, vi, beforeEach } from "vitest";
import type { CoreApiType } from "@/lib/api/api";
import { InboxService } from "./InboxService";
import { EntityInboxMessageType } from "@decm/api";

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


// =============================================================================
// Additional tests for InboxService - Enhanced Coverage
// =============================================================================

describe("InboxService - Enhanced Coverage", () => {
    let inboxService: InboxService;
    let mockCoreApi: CoreApiType;

    beforeEach(() => {
        vi.clearAllMocks();
        mockCoreApi = coreApiClient as unknown as CoreApiType;
        inboxService = new InboxService(mockCoreApi);
    });

    describe("getInboxMessages - Array handling", () => {
        it("should return empty array when inbox_messages is undefined", async () => {
            // Arrange
            vi.mocked(mockCoreApi.v1.v1InboxMessagesList).mockResolvedValue({
                inbox_messages: undefined,
            });

            // Act
            const result = await inboxService.getInboxMessages();

            // Assert
            expect(result).toEqual([]);
            expect(Array.isArray(result)).toBe(true);
        });

        it("should return empty array when inbox_messages is undefined", async () => {
            // Arrange
            vi.mocked(mockCoreApi.v1.v1InboxMessagesList).mockResolvedValue({
                inbox_messages: undefined,
            });

            // Act
            const result = await inboxService.getInboxMessages();

            // Assert
            expect(result).toEqual([]);
        });

        it("should correctly map multiple inbox messages", async () => {
            // Arrange
            const mockMessages = [
                {
                    id: "msg-1",
                    sender_credential_email: "sender1@example.com",
                    receiver_credential_id: "receiver-1",
                    message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
                    message_content: "Message 1",
                    is_read: 0,
                    created_at: "2024-01-01T00:00:00Z",
                    updated_at: "2024-01-01T00:00:00Z",
                },
                {
                    id: "msg-2",
                    sender_credential_email: "sender2@example.com",
                    receiver_credential_id: "receiver-2",
                    message_type: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
                    message_content: "Message 2",
                    is_read: 1,
                    created_at: "2024-01-02T00:00:00Z",
                    updated_at: "2024-01-02T00:00:00Z",
                    event_id: "event-123",
                },
                {
                    id: "msg-3",
                    sender_credential_wallet_address: "0xsender3",
                    receiver_email: "receiver3@example.com",
                    message_type: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
                    message_content: "Message 3",
                    is_read: 0,
                    created_at: "2024-01-03T00:00:00Z",
                    updated_at: "2024-01-03T00:00:00Z",
                    certificate_id: "cert-456",
                },
            ];

            vi.mocked(mockCoreApi.v1.v1InboxMessagesList).mockResolvedValue({
                inbox_messages: mockMessages,
            });

            // Act
            const result = await inboxService.getInboxMessages();

            // Assert
            expect(result).toHaveLength(3);
            expect(result[0].id).toBe("msg-1");
            expect(result[0].isRead).toBe(false);
            expect(result[1].id).toBe("msg-2");
            expect(result[1].isRead).toBe(true);
            expect(result[1].eventId).toBe("event-123");
            expect(result[2].id).toBe("msg-3");
            expect(result[2].certificateId).toBe("cert-456");
        });

        it("should handle messages with all optional fields as null", async () => {
            // Arrange
            const mockMessage = {
                id: "msg-minimal",
                sender_credential_id: "sender-id",
                receiver_credential_id: "receiver-id",
                message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
                message_content: null,
                fallback_message_content: null,
                is_read: 0,
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
                sender_credential_email: null,
                sender_credential_wallet_address: null,
                receiver_email: null,
                receiver_wallet_address: null,
                event_id: null,
                certificate_id: null,
                accepted_at: null,
                cancelled_at: null,
                revoked_at: null,
                token_id: null,
                valid_until: null,
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesList).mockResolvedValue({
                inbox_messages: [mockMessage],
            });

            // Act
            const result = await inboxService.getInboxMessages();

            // Assert
            expect(result).toHaveLength(1);
            expect(result[0].id).toBe("msg-minimal");
            expect(result[0].messageContent).toBeUndefined();
            expect(result[0].eventId).toBeUndefined();
            expect(result[0].certificateId).toBeUndefined();
        });
    });

    describe("markAsRead - Enhanced scenarios", () => {
        it("should call API with correct message ID", async () => {
            // Arrange
            const messageId = "msg-to-mark-read";
            vi.mocked(mockCoreApi.v1.v1InboxMessagesReadUpdate).mockResolvedValue({});

            // Act
            await inboxService.markAsRead(messageId);

            // Assert
            expect(mockCoreApi.v1.v1InboxMessagesReadUpdate).toHaveBeenCalledWith(messageId);
            expect(mockCoreApi.v1.v1InboxMessagesReadUpdate).toHaveBeenCalledTimes(1);
        });

        it("should handle API error when marking as read", async () => {
            // Arrange
            const messageId = "msg-error";
            const error = new Error("API Error");
            vi.mocked(mockCoreApi.v1.v1InboxMessagesReadUpdate).mockRejectedValue(error);

            // Act & Assert
            await expect(inboxService.markAsRead(messageId)).rejects.toThrow("API Error");
        });

        it("should handle network error when marking as read", async () => {
            // Arrange
            const messageId = "msg-network-error";
            const networkError = new Error("Network timeout");
            vi.mocked(mockCoreApi.v1.v1InboxMessagesReadUpdate).mockRejectedValue(networkError);

            // Act & Assert
            await expect(inboxService.markAsRead(messageId)).rejects.toThrow("Network timeout");
        });
    });

    describe("getInboxMessageById - Edge cases", () => {
        it("should handle message with revoked_at timestamp", async () => {
            // Arrange
            const messageId = "msg-revoked";
            const revokedAt = "2024-12-01T12:00:00Z";
            const mockMessage = {
                id: messageId,
                message_type: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
                is_read: 1,
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
                revoked_at: revokedAt,
                certificate_id: "cert-123",
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesDetail).mockResolvedValue(mockMessage);

            // Act
            const result = await inboxService.getInboxMessageById(messageId);

            // Assert
            expect(result).not.toBeNull();
            expect(result?.revokedAt).toEqual(new Date(revokedAt));
        });

        it("should handle message with cancelled_at timestamp", async () => {
            // Arrange
            const messageId = "msg-cancelled";
            const cancelledAt = "2024-11-15T10:30:00Z";
            const mockMessage = {
                id: messageId,
                message_type: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
                is_read: 1,
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
                cancelled_at: cancelledAt,
                event_id: "event-123",
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesDetail).mockResolvedValue(mockMessage);

            // Act
            const result = await inboxService.getInboxMessageById(messageId);

            // Assert
            expect(result).not.toBeNull();
            expect(result?.cancelledAt).toEqual(new Date(cancelledAt));
        });

        it("should handle message with valid_until timestamp", async () => {
            // Arrange
            const messageId = "msg-expiring";
            const validUntil = "2025-01-01T00:00:00Z";
            const mockMessage = {
                id: messageId,
                message_type: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
                is_read: 0,
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
                valid_until: validUntil,
                event_id: "event-456",
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesDetail).mockResolvedValue(mockMessage);

            // Act
            const result = await inboxService.getInboxMessageById(messageId);

            // Assert
            expect(result).not.toBeNull();
            expect(result?.validUntil).toEqual(new Date(validUntil));
        });

        it("should handle message with accepted_at timestamp", async () => {
            // Arrange
            const messageId = "msg-accepted";
            const acceptedAt = "2024-10-01T15:45:00Z";
            const mockMessage = {
                id: messageId,
                message_type: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
                is_read: 1,
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
                accepted_at: acceptedAt,
                event_id: "event-789",
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesDetail).mockResolvedValue(mockMessage);

            // Act
            const result = await inboxService.getInboxMessageById(messageId);

            // Assert
            expect(result).not.toBeNull();
            expect(result?.acceptedAt).toEqual(new Date(acceptedAt));
        });

        it("should handle message with token_id for claimed certificate", async () => {
            // Arrange
            const messageId = "msg-claimed";
            const tokenId = "12345";
            const mockMessage = {
                id: messageId,
                message_type: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
                is_read: 1,
                created_at: "2024-01-01T00:00:00Z",
                updated_at: "2024-01-01T00:00:00Z",
                certificate_id: "cert-999",
                token_id: tokenId,
            };

            vi.mocked(mockCoreApi.v1.v1InboxMessagesDetail).mockResolvedValue(mockMessage);

            // Act
            const result = await inboxService.getInboxMessageById(messageId);

            // Assert
            expect(result).not.toBeNull();
            expect(result?.tokenId).toBe(tokenId);
        });
    });
});
