import { coreApiClient } from "@/lib/api/api";

interface GetEventsListParams {
    includeActiveEvents?: boolean;
    includeInactiveEvents?: boolean;
    includeClosedEvents?: boolean;
    onlyUserJoinedEvents?: boolean;
}

export class EventService {
    private coreApi: typeof coreApiClient;

    constructor() {
        this.coreApi = coreApiClient;
    }

    public async getEventById(eventId: string) {
        const response = await this.coreApi.v1.getEventById({ eventId });
        return response;
    }

    public async getEvents(params: GetEventsListParams) {
        const response = await this.coreApi.v1.getEventsList({
            include_active_events: params.includeActiveEvents,
            include_inactive_events: params.includeInactiveEvents,
            include_closed_events: params.includeClosedEvents,
            only_user_joined_events: params.onlyUserJoinedEvents,
        });

        return response;
    }
}
