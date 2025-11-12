import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import type { EntityEventStatus, EntityEventType } from "@decm/api";

interface GetEventsListParams {
    includeActiveEvents?: boolean;
    includeInactiveEvents?: boolean;
    includeClosedEvents?: boolean;
    onlyUserJoinedEvents?: boolean;
}

export interface Event {
    banner_storage_key?: string;
    chain_id?: number;
    contact_address?: string;
    contact_number?: string;
    created_at?: string;
    end_date?: string;
    event_status?: EntityEventStatus;
    event_type?: EntityEventType;
    google_map_query?: string;
    icon_storage_key?: string;
    id?: string;
    is_booking_request_required?: boolean;
    is_public?: boolean;
    is_ticket_transferable?: boolean;
    is_verified?: boolean;
    location?: string;
    long_description?: string;
    max_attendees?: number;
    owner_credential_id?: string;
    short_description?: string;
    start_date?: string;
    title?: string;
    updated_at?: string;
}

export class EventService {
    private _coreApi: CoreApiType;

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    public async getEventById(eventId: string) {
        const response = await this._coreApi.v1.getEventById({ eventId });
        return response;
    }

    public async getEvents(params: GetEventsListParams): Promise<Event[]> {
        const response = await this._coreApi.v1.getEventsList({
            include_active_events: params.includeActiveEvents,
            include_inactive_events: params.includeInactiveEvents,
            include_closed_events: params.includeClosedEvents,
            only_user_joined_events: params.onlyUserJoinedEvents,
        });

        return response;
    }
}

export const defaultEventService = new EventService(coreApiClient);
