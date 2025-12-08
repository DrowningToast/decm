import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import {
    mapInboxMessagesViewModelToInboxMessage,
    mapInboxMessagesDetailViewModelToInboxMessageDetail,
} from "./mapper";
import type { EntityInboxMessageType } from "@decm/api";

export type InboxMessageType =
    | "general"
    | "event_registration_invitation"
    | "event_certificate_invitation";

export interface InboxMessageContent {
    en: string;
    th: string;
}

export interface InboxMessage {
    id: string;
    messageType: InboxMessageType;
    messageContent: string;
    isRead: boolean;
    receiverEmail?: string;
    receiverWalletAddress?: string;
    senderCredentialEmail?: string;
    senderCredentialWalletAddress?: string;
    createdAt: Date;
    updatedAt: Date;
    hiddenAt?: Date;
    deletedAt?: Date;
}

export type InboxMessageDetail =
    | (InboxMessage & {
          entityInboxMessageType: EntityInboxMessageType.InboxMessageTypeGeneral;
      })
    | (InboxMessage & {
          entityInboxMessageType: EntityInboxMessageType.InboxMessageTypeEventRegistrationInvitation;
          eventId?: string;
          code?: string;
          email?: string;
          firstName?: string;
          lastName?: string;
          phoneNumber?: string;
          academicInstitution?: string;
          validUntil?: Date;
          cancelledAt?: Date;
          acceptedAt?: Date;
      })
    | (InboxMessage & {
          entityInboxMessageType: EntityInboxMessageType.InboxMessageTypeEventCertificateInvitation;
          eventName: string;
          certificateId: string;
          certificateTitle?: string;
          tokenId?: string;
          revokedAt?: Date;
          hasParticipantJoinedEvent: boolean;
      });

export class InboxService {
    private _coreApi: CoreApiType;

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    public async getInboxMessages(): Promise<InboxMessage[]> {
        const response = await this._coreApi.v1.v1InboxMessagesList();
        return (response.inbox_messages ?? []).map(mapInboxMessagesViewModelToInboxMessage);
    }

    public async getInboxMessageById(messageId: string): Promise<InboxMessageDetail | null> {
        const response = await this._coreApi.v1.v1InboxMessagesDetail({
            inboxMessageId: messageId,
        });
        if (!response.inbox_message) {
            return null;
        }
        return mapInboxMessagesDetailViewModelToInboxMessageDetail(response);
    }

    public async markAsRead(messageId: string): Promise<InboxMessage | null> {
        const response = await this._coreApi.v1.v1InboxMessagesReadUpdate(
            { messageId },
            { message_id: messageId },
        );
        if (!response.inbox_message) {
            return null;
        }
        return mapInboxMessagesViewModelToInboxMessage(response.inbox_message);
    }

    public async markAllAsRead(): Promise<InboxMessage[]> {
        const response = await this._coreApi.v1.v1InboxMessagesReadAllUpdate();
        return (response.inbox_messages ?? []).map(mapInboxMessagesViewModelToInboxMessage);
    }
}

export const defaultInboxService = new InboxService(coreApiClient);
