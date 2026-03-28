import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import {
    mapEntityEventToEventItem,
    mapToEventIssuer,
    mapEventResponseToViewModel,
    mapEventViewModelExtended,
} from "./mapper";
import type { EventRegistrationConfiguration } from "../EventRegistration/EventRegistration";
import type { Profile } from "../AuthService/AuthService";
import type {
    CreateEventPayload,
    UpdateEventPayload,
    CreateEventData,
    UpdateEventData,
    DeleteEventByIdData,
    UpdateEventParams,
    DeleteEventByIdParams,
} from "@decm/api";

type CreateEventCommonFields = Omit<CreateEventPayload, "host_password">;
type EditEventCommonFields = Omit<UpdateEventPayload, "host_password">;

type CreateEventWalletPayload = CreateEventCommonFields & {
    signature: string;
    sign_message: string;
};
type UpdateEventWalletPayload = EditEventCommonFields & { signature: string; sign_message: string };
type DeleteEventWalletPayload = { signature: string; sign_message: string };

type CreateEventFn = (data: CreateEventWalletPayload) => Promise<CreateEventData>;
type UpdateEventFn = (
    params: UpdateEventParams,
    data: UpdateEventWalletPayload,
) => Promise<UpdateEventData>;
type DeleteEventFn = (
    params: DeleteEventByIdParams,
    data: DeleteEventWalletPayload,
) => Promise<DeleteEventByIdData>;

export interface CreateEventWithPasswordInput extends CreateEventCommonFields {
    hostPassword: string;
}

export interface CreateEventWithSignatureInput extends CreateEventCommonFields {
    signature: string;
    signMessage: string;
}

export interface EditEventWithPasswordInput extends EditEventCommonFields {
    hostPassword: string;
}

export interface EditEventWithSignatureInput extends EditEventCommonFields {
    signature: string;
    signMessage: string;
}

export interface WalletSignatureAuth {
    signature: string;
    signMessage: string;
}

interface GetEventsListParams {
    includeActiveEvents?: boolean;
    includeInactiveEvents?: boolean;
    includeClosedEvents?: boolean;
    onlyUserJoinedEvents?: boolean;
}

export type EventStatus = "active" | "inactive" | "closed";
export type EventType = "private" | "invite";

export interface EventItem {
    chainId: number;
    contactAddress: string;
    contactNumber: string;
    createdAt: Date;
    endDate: Date;
    eventStatus: EventStatus;
    eventType: EventType;
    googleMapQuery: string;
    iconStorageKey: string;
    id: string;
    isBookingRequestRequired: boolean;
    isPublic: boolean;
    isTicketTransferable: boolean;
    isVerified: boolean;
    location: string;
    longDescription: string;
    maxAttendees: number;
    attendeesCount: number;
    ownerCredentialId: string;
    shortDescription: string;
    startDate: Date;
    title: string;
    updatedAt: Date;
}

export interface EventViewModel {
    bannerPresignedUrl?: string;
    chainId: number;
    contactAddress: string;
    contactNumber: string;
    createdAt: Date;
    endDate: Date;
    eventStatus: EventStatus;
    eventType: EventType;
    googleMapQuery?: string;
    iconPresignedUrl?: string;
    id: string;
    isBookingRequestRequired?: boolean;
    isPublic: boolean;
    isTicketTransferable: boolean;
    isVerified: boolean;
    location: string;
    longDescription: string;
    maxAttendees: number;
    attendeesCount: number;
    ownerCredentialId: string;
    shortDescription: string;
    startDate: Date;
    title: string;
    updatedAt: Date;
}

export interface EventViewModelExtended extends EventViewModel {
    isInvited?: boolean;
    isJoined?: boolean;
    isFull?: boolean;

    // Joined status fields
    joinedSignature?: string;
    joinedSignMessage?: string;
    joinedAt?: Date;
    joinedIsAccepted?: boolean;
    broadcastedAtEstimatedDeadline?: Date;
    broadcastedAtDeadlineBlock?: number;
    broadcastedAt?: Date;
    signatureExpiredAt?: Date;

    finalCallDate?: string;
    registrationRequirement?: EventRegistrationConfiguration;

    accessManagerContractAddress?: string;
    eventContractAddress?: string;
    ticketContractAddress?: string;
    certificateContractAddress?: string;
}

export interface EventIssuer {
    eventId: string;
    id: string;
    isSigned: boolean;
    issuerCredentialId: string;
    issuerProfile?: Profile;

    createdAt: Date;
    updatedAt: Date;
}

export class EventService {
    private _coreApi: CoreApiType;

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    public async getEventById(eventId: string): Promise<EventViewModelExtended> {
        const response = await this._coreApi.v1.getEventById({ eventId });
        return mapEventResponseToViewModel(response);
    }

    public async getEvents(params: GetEventsListParams): Promise<EventItem[]> {
        const response = await this._coreApi.v1.getEventsList({
            include_active_events: params.includeActiveEvents ?? true,
            include_inactive_events: params.includeInactiveEvents,
            include_closed_events: params.includeClosedEvents,
            only_user_joined_events: params.onlyUserJoinedEvents,
        });

        return response.events?.map(mapEntityEventToEventItem) ?? [];
    }

    public async getEventViewModelExtended(eventId: string): Promise<EventViewModelExtended> {
        const response = await this._coreApi.v1.getEventViewmodelById({ eventId });
        return mapEventViewModelExtended(response);
    }

    public async getEventsByOwnerCredentialId(
        ownerCredentialId: string,
    ): Promise<EventViewModel[]> {
        const response = await coreApiClient.v1.getEventsByOwnerCredentialsId({
            ownerCredentialId: ownerCredentialId,
        });

        return response.map(mapEventResponseToViewModel);
    }

    public async getEventIssuersByEventId(eventId: string): Promise<EventIssuer[]> {
        const response = await coreApiClient.v1.getEventIssuersByEventId({ eventId });
        return response.map(mapToEventIssuer);
    }

    public async publishEventCertificates(eventId: string) {
        return await coreApiClient.v1.publishEventCertificates({ eventId });
    }

    public async createEventWithPassword(input: CreateEventWithPasswordInput): Promise<string> {
        const { hostPassword, ...rest } = input;
        const response = await this._coreApi.v1.createEvent({
            ...rest,
            host_password: hostPassword,
        });
        return response.id ?? "";
    }

    public async createEventWithSignature(input: CreateEventWithSignatureInput): Promise<string> {
        const { signature, signMessage, ...rest } = input;
        const createEventWithWallet = this._coreApi.v1.createEvent as unknown as CreateEventFn;
        const response = await createEventWithWallet({
            ...rest,
            signature,
            sign_message: signMessage,
        });
        return response.id ?? "";
    }

    public async editEventWithPassword(
        eventId: string,
        input: EditEventWithPasswordInput,
    ): Promise<void> {
        const { hostPassword, ...rest } = input;
        await this._coreApi.v1.updateEvent({ eventId }, { ...rest, host_password: hostPassword });
    }

    public async editEventWithSignature(
        eventId: string,
        input: EditEventWithSignatureInput,
    ): Promise<void> {
        const { signature, signMessage, ...rest } = input;
        const updateEventWithWallet = this._coreApi.v1.updateEvent as unknown as UpdateEventFn;
        await updateEventWithWallet({ eventId }, { ...rest, signature, sign_message: signMessage });
    }

    public async deleteEventWithPassword(eventId: string, hostPassword: string): Promise<void> {
        await this._coreApi.v1.deleteEventById({ eventId }, { host_password: hostPassword });
    }

    public async deleteEventWithSignature(
        eventId: string,
        auth: WalletSignatureAuth,
    ): Promise<void> {
        const deleteEventWithWallet = this._coreApi.v1.deleteEventById as unknown as DeleteEventFn;
        await deleteEventWithWallet(
            { eventId },
            { signature: auth.signature, sign_message: auth.signMessage },
        );
    }
}

export const defaultEventService = new EventService(coreApiClient);
