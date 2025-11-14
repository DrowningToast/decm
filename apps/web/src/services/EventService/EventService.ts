import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import {
    mapEntityEventToEventItem,
    mapEventResponseToViewModel,
    mapEventViewModelExtended,
} from "./mapper";
import type { EventRegistrationConfiguration } from "../EventRegistration/EventRegistration";

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
    iconPresignedUrl: string;
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

    finalCallDate?: string;
    registrationRequirement?: EventRegistrationConfiguration;

    accessManagerContractAddress?: string;
    eventContractAddress?: string;
    ticketContractAddress?: string;
    certificateContractAddress?: string;
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

    public async getEventViewModel(eventId: string): Promise<EventViewModelExtended> {
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
}

export const defaultEventService = new EventService(coreApiClient);
