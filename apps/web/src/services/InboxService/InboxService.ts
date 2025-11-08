import { coreApiClient } from "@/lib/api/api";
import type {
    InboxInboxMessagesViewModel,
    InboxInboxMessagesEventRegistrationInvitationViewModel,
} from "@decm/api";
import z from "zod";

export const InboxMessageContent = z.object({
    en: z.string(),
    th: z.string(),
});

export type InboxMessageContent = z.infer<typeof InboxMessageContent>;

export type InboxMessageType =
    | "general"
    | "event_registration_invitation"
    | "event_certificate_invitation";

export interface InboxMessage {
    id: string;
    messageType: InboxMessageType;
    messageContent: InboxMessageContent;
    isRead: boolean;
    fallbackMessageContent?: string;
    receiverCredentialId?: string;
    receiverEmail?: string;
    receiverWalletAddress?: string;
    senderCredentialId?: string;
    createdAt: string;
    updatedAt: string;
    hiddenAt?: string;
    deletedAt?: string;
}

export class InboxService {
    private coreApi: typeof coreApiClient;

    constructor() {
        this.coreApi = coreApiClient;
    }

    public async getInboxMessages(): Promise<InboxInboxMessagesViewModel[]> {
        try {
            const response = await this.coreApi.inboxMessagesList();
            return response.inbox_messages ?? [];
        } catch (error) {
            console.error("Failed to fetch inbox messages:", error);
            throw error;
        }
    }

    public async getInboxMessageById(
        messageId: string,
    ): Promise<InboxInboxMessagesEventRegistrationInvitationViewModel | null> {
        try {
            const response = await this.coreApi.inboxMessageId.inboxMessagesDetail({
                inboxMessageId: messageId,
            });
            return response.inbox_message ?? null;
        } catch (error) {
            console.error(`Failed to fetch inbox message ${messageId}:`, error);
            throw error;
        }
    }

    public async markAsRead(messageId: string): Promise<InboxInboxMessagesViewModel | null> {
        try {
            const response = await this.coreApi.read.readUpdate(
                { messageId },
                { message_id: messageId },
            );
            return response.inbox_message ?? null;
        } catch (error) {
            console.error(`Failed to mark inbox message ${messageId} as read:`, error);
            throw error;
        }
    }

    public async markAllAsRead(): Promise<InboxInboxMessagesViewModel[]> {
        try {
            const response = await this.coreApi.readAll.readAllUpdate();
            return response.inbox_messages ?? [];
        } catch (error) {
            console.error("Failed to mark all inbox messages as read:", error);
            throw error;
        }
    }
}
