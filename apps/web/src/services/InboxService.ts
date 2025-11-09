import { coreApiClient } from "@/lib/api/api";

export interface InboxMessage {
    id: string;
    title: string;
    sender: string;
    date: string;
    status: "pending" | "available" | "expired" | "action-required";
    contentType: "event-invitation" | "certificate";
    eventId?: string;
    certificateId?: string;
    description?: string;
    isUserInEvent?: boolean;
}

export interface GetInboxMessagesParams {
    limit?: number;
    offset?: number;
}

export class InboxService {
    private coreApi: typeof coreApiClient;

    constructor() {
        this.coreApi = coreApiClient;
    }

    public async getInboxMessages(params?: GetInboxMessagesParams): Promise<InboxMessage[]> {
        try {
            const response = await this.coreApi.v1.getInboxMessages({
                limit: params?.limit,
                offset: params?.offset,
            });
            return response.data;
        } catch (error) {
            console.error("Failed to fetch inbox messages:", error);
            throw error;
        }
    }

    public async getInboxMessageById(messageId: string): Promise<InboxMessage | null> {
        try {
            const response = await this.coreApi.v1.getInboxMessageById({ messageId });
            return response.data;
        } catch (error) {
            console.error(`Failed to fetch inbox message ${messageId}:`, error);
            throw error;
        }
    }

    public async markAsRead(messageId: string): Promise<void> {
        try {
            await this.coreApi.v1.markInboxMessageAsRead({ messageId });
        } catch (error) {
            console.error(`Failed to mark inbox message ${messageId} as read:`, error);
            throw error;
        }
    }

    public async deleteInboxMessage(messageId: string): Promise<void> {
        try {
            await this.coreApi.v1.deleteInboxMessage({ messageId });
        } catch (error) {
            console.error(`Failed to delete inbox message ${messageId}:`, error);
            throw error;
        }
    }
}
