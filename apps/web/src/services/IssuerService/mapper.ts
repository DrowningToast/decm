import type { EventEventIssuerResponse, IssuerIssuerEventResponse } from "@decm/api";
import type { IssuerEvent } from "./IssuerService";
import { mapProfileViewModel } from "../AuthService/mapper";
import type { EventIssuer } from "../EventService/EventService";

export const mapToEventIssuer = (eventIssuerResponse: EventEventIssuerResponse): EventIssuer => {
    return {
        eventId: eventIssuerResponse.event_id,
        id: eventIssuerResponse.id,
        isSigned: eventIssuerResponse.is_signed === 1,
        issuerCredentialId: eventIssuerResponse.issuer_credential_id,
        issuerProfile: eventIssuerResponse.issuer_profile
            ? mapProfileViewModel(eventIssuerResponse.issuer_profile)
            : undefined,
        signMessage: eventIssuerResponse.sign_message ?? undefined,
        signature: eventIssuerResponse.signature ?? undefined,
    };
};

export const mapToIssuerEvent = (issuerEventResponse: IssuerIssuerEventResponse): IssuerEvent => {
    return {
        eventEndDate: new Date(issuerEventResponse.event_end_date),
        eventId: issuerEventResponse.event_id || "",
        eventLocation: issuerEventResponse.event_location || "",
        eventOwnerCredentialId: issuerEventResponse.event_owner_credential_id || "",
        eventShortDescription: issuerEventResponse.event_short_description || "",
        eventStartDate: new Date(issuerEventResponse.event_start_date),
        eventTitle: issuerEventResponse.event_title || "",
        id: issuerEventResponse.id,
        isSigned: issuerEventResponse.is_signed === 1,
        issuerCredentialId: issuerEventResponse.issuer_credential_id || "",
        signMessage: issuerEventResponse.sign_message || "",
        signature: issuerEventResponse.signature || "",
    };
};
