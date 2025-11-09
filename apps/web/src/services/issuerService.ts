import { coreApiClient } from "@/lib/api/api";
import type { GetIssuerEventsData } from "@decm/api";

export interface IssuerEvent {
    id: string;
    event_id: string;
    event_title: string;
    event_short_description: string;
    event_start_date: string;
    event_end_date: string;
    event_location: string;
    event_owner_credential_id: string;
    issuer_credential_id: string;
    is_signed: number;
    signature: string;
    sign_message: string;
    created_at: string;
    updated_at: string;
}

export interface IssuerEventsResponse {
    data: IssuerEvent[];
    total: number;
    page: number;
    pageSize: number;
}

/**
 * Get events for the authenticated issuer
 * @param limit Number of events to return
 * @param offset Number of events to skip
 * @returns Promise with issuer events data
 */
export const getIssuerEvents = async (
    issuer_credential_id: string,
    limit?: number,
    offset?: number,
): Promise<GetIssuerEventsData> => {
    try {
        const response = await coreApiClient.v1.getIssuerEvents({
            limit: limit || 10,
            offset: offset || 0,
            issuer_credential_id,
        });
        return response;
    } catch (error) {
        console.error("Error fetching issuer events:", error);
        throw error;
    }
};

/**
 * Filter events by signing status
 * @param events Array of events to filter
 * @param isSigned Status to filter by (0 for pending, 1 for signed)
 * @returns Filtered events array
 */
export const filterEventsByStatus = (events: GetIssuerEventsData, isSigned: number) => {
    return events.filter((event) => event.is_signed === isSigned);
};

/**
 * Get pending events for issuer
 * @param limit Number of events to return
 * @param offset Number of events to skip
 * @returns Promise with pending events
 */
export const getPendingEvents = async (
    issuer_credential_id: string,
    limit?: number,
    offset?: number,
): Promise<GetIssuerEventsData> => {
    const allEvents = await getIssuerEvents(issuer_credential_id, limit || 10, offset || 0);
    return filterEventsByStatus(allEvents, 0);
};

/**
 * Get signed events for issuer
 * @param limit Number of events to return
 * @param offset Number of events to skip
 * @returns Promise with signed events
 */
export const getSignedEvents = async (
    issuer_credential_id: string,
    limit?: number,
    offset?: number,
): Promise<GetIssuerEventsData> => {
    const allEvents = await getIssuerEvents(issuer_credential_id, limit, offset);
    return filterEventsByStatus(allEvents, 1);
};
