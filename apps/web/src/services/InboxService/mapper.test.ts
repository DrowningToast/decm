import { describe, it, expect } from "vitest";
import {
    mapEntityInboxMessageTypeToInboxMessageType,
    mapInboxMessagesViewModelToInboxMessage,
    mapInboxMessagesDetailViewModelToInboxMessageDetail,
    mapEntityInboxMessageToInboxMessage,
} from "./mapper";
import {
    EntityInboxMessageType,
    type InboxInboxMessagesViewModel,
    type InboxmessagesGetInboxMessageResponse,
    type EntityInboxMessage,
} from "@decm/api";

describe("mapEntityInboxMessageTypeToInboxMessageType", () => {
    it("maps General to 'general'", () => {
        expect(
            mapEntityInboxMessageTypeToInboxMessageType(
                EntityInboxMessageType.InboxMessageTypeGeneral,
            ),
        ).toBe("general");
    });

    it("maps EventRegistrationInvitation to 'event_registration_invitation'", () => {
        expect(
            mapEntityInboxMessageTypeToInboxMessageType(
                EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
            ),
        ).toBe("event_registration_invitation");
    });

    it("maps EventCertificateInvitation to 'event_certificate_invitation'", () => {
        expect(
            mapEntityInboxMessageTypeToInboxMessageType(
                EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
            ),
        ).toBe("event_certificate_invitation");
    });

    it("defaults to 'general' for unknown values", () => {
        expect(mapEntityInboxMessageTypeToInboxMessageType(99 as EntityInboxMessageType)).toBe(
            "general",
        );
    });
});

describe("mapInboxMessagesViewModelToInboxMessage", () => {
    const baseViewModel: InboxInboxMessagesViewModel = {
        id: "msg-1",
        message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
        message_content: "Hello",
        is_read: 1,
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-02T00:00:00Z",
    };

    it("maps all base fields correctly", () => {
        const result = mapInboxMessagesViewModelToInboxMessage(baseViewModel);

        expect(result.id).toBe("msg-1");
        expect(result.messageType).toBe("general");
        expect(result.messageContent).toBe("Hello");
        expect(result.isRead).toBe(true);
        expect(result.createdAt).toEqual(new Date("2024-01-01T00:00:00Z"));
        expect(result.updatedAt).toEqual(new Date("2024-01-02T00:00:00Z"));
    });

    it("maps is_read 0 to false", () => {
        const result = mapInboxMessagesViewModelToInboxMessage({ ...baseViewModel, is_read: 0 });
        expect(result.isRead).toBe(false);
    });

    it("handles optional string fields as undefined when null", () => {
        const result = mapInboxMessagesViewModelToInboxMessage(baseViewModel);
        expect(result.receiverEmail).toBeUndefined();
        expect(result.receiverWalletAddress).toBeUndefined();
        expect(result.senderCredentialEmail).toBeUndefined();
        expect(result.senderCredentialWalletAddress).toBeUndefined();
        expect(result.hiddenAt).toBeUndefined();
        expect(result.deletedAt).toBeUndefined();
    });

    it("maps optional date fields when present", () => {
        const result = mapInboxMessagesViewModelToInboxMessage({
            ...baseViewModel,
            hidden_at: "2024-03-01T00:00:00Z",
            deleted_at: "2024-04-01T00:00:00Z",
        });
        expect(result.hiddenAt).toEqual(new Date("2024-03-01T00:00:00Z"));
        expect(result.deletedAt).toEqual(new Date("2024-04-01T00:00:00Z"));
    });

    it("maps certificate-specific fields", () => {
        const result = mapInboxMessagesViewModelToInboxMessage({
            ...baseViewModel,
            event_id: "evt-1",
            certificate_id: "cert-1",
            certificate_title: "My Cert",
            token_id: "tok-1",
            has_participant_joined_event: true,
            event_name: "Test Event",
        });
        expect(result.eventId).toBe("evt-1");
        expect(result.certificateId).toBe("cert-1");
        expect(result.certificateTitle).toBe("My Cert");
        expect(result.tokenId).toBe("tok-1");
        expect(result.hasParticipantJoinedEvent).toBe(true);
        expect(result.eventName).toBe("Test Event");
    });

    it("defaults message_content to empty string when null", () => {
        const vm = {
            ...baseViewModel,
            message_content: null,
        } as unknown as InboxInboxMessagesViewModel;
        const result = mapInboxMessagesViewModelToInboxMessage(vm);
        expect(result.messageContent).toBe("");
    });
});

describe("mapInboxMessagesDetailViewModelToInboxMessageDetail", () => {
    const baseInboxMessage: InboxInboxMessagesViewModel = {
        id: "msg-1",
        message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
        message_content: "General message",
        is_read: 0,
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-02T00:00:00Z",
    };

    it("maps general inbox message type", () => {
        const response: InboxmessagesGetInboxMessageResponse = {
            inbox_message: baseInboxMessage,
            inbox_message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
        };
        const result = mapInboxMessagesDetailViewModelToInboxMessageDetail(response);

        expect(result.entityInboxMessageType).toBe(EntityInboxMessageType.InboxMessageTypeGeneral);
        expect(result.id).toBe("msg-1");
        expect(result.isRead).toBe(false);
    });

    it("maps registration invitation type with invitation fields", () => {
        const response: InboxmessagesGetInboxMessageResponse = {
            inbox_message: {
                ...baseInboxMessage,
                message_type: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
            },
            inbox_message_type: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
            event_registration_invitation: {
                id: "inv-1",
                event_id: "evt-1",
                code: "ABC123",
                email: "test@example.com",
                first_name: "John",
                message_content: "Invitation",
                message_type: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
                is_read: 0,
                created_at: "2024-01-01T00:00:00Z",
            },
        };
        const result = mapInboxMessagesDetailViewModelToInboxMessageDetail(response);

        expect(result.entityInboxMessageType).toBe(
            EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
        );
        expect(result).toHaveProperty("eventId", "evt-1");
        expect(result).toHaveProperty("code", "ABC123");
        expect(result).toHaveProperty("email", "test@example.com");
        expect(result).toHaveProperty("firstName", "John");
    });

    it("maps certificate invitation type with certificate fields", () => {
        const response: InboxmessagesGetInboxMessageResponse = {
            inbox_message: {
                ...baseInboxMessage,
                message_type: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
            },
            inbox_message_type: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
            event_certificate: {
                id: "cert-vm-1",
                event_id: "evt-1",
                event_name: "Test Event",
                certificate_id: "cert-1",
                certificate_title: "Completion Certificate",
                token_id: "tok-1",
                has_participant_joined_event: true,
                message_content: "Certificate",
                message_type: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
                is_read: 0,
                created_at: "2024-01-01T00:00:00Z",
            },
        };
        const result = mapInboxMessagesDetailViewModelToInboxMessageDetail(response);

        expect(result.entityInboxMessageType).toBe(
            EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
        );
        expect(result).toHaveProperty("eventName", "Test Event");
        expect(result).toHaveProperty("certificateId", "cert-1");
        expect(result).toHaveProperty("certificateTitle", "Completion Certificate");
        expect(result).toHaveProperty("tokenId", "tok-1");
        expect(result).toHaveProperty("hasParticipantJoinedEvent", true);
    });

    it("throws when certificate invitation type is missing event_certificate", () => {
        const response: InboxmessagesGetInboxMessageResponse = {
            inbox_message: {
                ...baseInboxMessage,
                message_type: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
            },
            inbox_message_type: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation,
        };
        expect(() => mapInboxMessagesDetailViewModelToInboxMessageDetail(response)).toThrow(
            "inbox message type is event certificate invitation but event certificate is not found",
        );
    });

    it("maps optional date fields in registration invitation", () => {
        const response: InboxmessagesGetInboxMessageResponse = {
            inbox_message: {
                ...baseInboxMessage,
                message_type: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
            },
            inbox_message_type: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
            event_registration_invitation: {
                id: "inv-1",
                event_id: "evt-1",
                message_content: "Invitation",
                message_type: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation,
                is_read: 0,
                created_at: "2024-01-01T00:00:00Z",
                valid_until: "2024-12-31T00:00:00Z",
                cancelled_at: "2024-06-01T00:00:00Z",
                accepted_at: "2024-02-01T00:00:00Z",
                broadcasted_at: "2024-01-15T00:00:00Z",
            },
        };
        const result = mapInboxMessagesDetailViewModelToInboxMessageDetail(response);

        expect(result).toHaveProperty("validUntil", new Date("2024-12-31T00:00:00Z"));
        expect(result).toHaveProperty("cancelledAt", new Date("2024-06-01T00:00:00Z"));
        expect(result).toHaveProperty("acceptedAt", new Date("2024-02-01T00:00:00Z"));
        expect(result).toHaveProperty("broadcastedAt", new Date("2024-01-15T00:00:00Z"));
    });
});

describe("mapEntityInboxMessageToInboxMessage", () => {
    const baseEntity: EntityInboxMessage = {
        id: "msg-1",
        message_type: EntityInboxMessageType.InboxMessageTypeGeneral,
        message_content: "Hello",
        is_read: 1,
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-02T00:00:00Z",
    };

    it("maps all fields correctly", () => {
        const result = mapEntityInboxMessageToInboxMessage(baseEntity);

        expect(result.id).toBe("msg-1");
        expect(result.messageType).toBe("general");
        expect(result.messageContent).toBe("Hello");
        expect(result.isRead).toBe(true);
        expect(result.createdAt).toEqual(new Date("2024-01-01T00:00:00Z"));
        expect(result.updatedAt).toEqual(new Date("2024-01-02T00:00:00Z"));
        expect(result.senderCredentialWalletAddress).toBeUndefined();
    });

    it("maps is_read 0 to false", () => {
        const result = mapEntityInboxMessageToInboxMessage({ ...baseEntity, is_read: 0 });
        expect(result.isRead).toBe(false);
    });

    it("maps sender_credential_id to senderCredentialEmail", () => {
        const result = mapEntityInboxMessageToInboxMessage({
            ...baseEntity,
            sender_credential_id: "cred-123",
        });
        expect(result.senderCredentialEmail).toBe("cred-123");
    });

    it("maps optional date fields", () => {
        const result = mapEntityInboxMessageToInboxMessage({
            ...baseEntity,
            hidden_at: "2024-03-01T00:00:00Z",
            deleted_at: "2024-04-01T00:00:00Z",
        });
        expect(result.hiddenAt).toEqual(new Date("2024-03-01T00:00:00Z"));
        expect(result.deletedAt).toEqual(new Date("2024-04-01T00:00:00Z"));
    });

    it("leaves optional date fields undefined when not present", () => {
        const result = mapEntityInboxMessageToInboxMessage(baseEntity);
        expect(result.hiddenAt).toBeUndefined();
        expect(result.deletedAt).toBeUndefined();
    });
});
