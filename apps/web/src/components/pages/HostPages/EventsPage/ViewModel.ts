import type { EventType, EventStatus } from "@/services/EventService/EventService";

export const EventTypesViewModel: Record<EventType, string> = {
    private: "Private",
    invite: "Invite",
};

export const EventStatusesViewModel: Record<EventStatus, string> = {
    active: "Active",
    inactive: "Inactive",
    closed: "Closed",
} as const;
