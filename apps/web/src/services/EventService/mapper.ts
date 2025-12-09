import {
    type EventEventViewModel,
    EntityEventStatus,
    type EntityEvent,
    type EventEventResponse,
    EntityEventType,
    type EventEventIssuerResponse,
} from "@decm/api";
import type {
    EventViewModelExtended,
    EventStatus,
    EventItem,
    EventType,
    EventViewModel,
    EventIssuer,
} from "./EventService";
import { mapEntityEventRegistrationConfigToEventRegistrationConfiguration } from "../EventRegistration/mapper";
import { mapProfileViewModel } from "../AuthService/mapper";

export const mapEntityEventStatusToEventStatus = (
    entityEventStatus: EntityEventStatus,
): EventStatus => {
    switch (entityEventStatus) {
        case EntityEventStatus.EventStatusActive:
            return "active";
        case EntityEventStatus.EventStatusInactive:
            return "inactive";
        case EntityEventStatus.EventStatusClosed:
            return "closed";
        default:
            console.error(`Invalid entity event status: ${entityEventStatus}`);
            throw new Error(`Invalid entity event status: ${entityEventStatus}`);
    }
};

export const mapEntityEventTypeToEventType = (entityEventType: EntityEventType): EventType => {
    switch (entityEventType) {
        case EntityEventType.EventTypePrivate:
            return "private";
        case EntityEventType.EventTypeInvite:
            return "invite";
    }
};

export const mapEventTypeToEntityEventType = (eventType: EventType): EntityEventType => {
    switch (eventType) {
        case "private":
            return EntityEventType.EventTypePrivate;
        case "invite":
            return EntityEventType.EventTypeInvite;
    }
};

export const mapEntityEventToEventItem = (entityEvent: EntityEvent): EventItem => {
    return {
        chainId: entityEvent.chain_id,
        contactAddress: entityEvent.contact_address,
        contactNumber: entityEvent.contact_number,
        isVerified: entityEvent.is_verified,
        eventStatus: mapEntityEventStatusToEventStatus(entityEvent.event_status),
        eventType: mapEntityEventTypeToEventType(entityEvent.event_type),
        googleMapQuery: entityEvent.google_map_query,
        iconStorageKey: entityEvent.icon_storage_key,
        id: entityEvent.id,
        isBookingRequestRequired: entityEvent.is_booking_request_required,
        isPublic: entityEvent.is_public,
        isTicketTransferable: entityEvent.is_ticket_transferable,
        location: entityEvent.location,
        longDescription: entityEvent.long_description,
        maxAttendees: entityEvent.max_attendees,
        attendeesCount: entityEvent.attendees_count,
        ownerCredentialId: entityEvent.owner_credential_id,
        shortDescription: entityEvent.short_description,
        title: entityEvent.title,
        startDate: new Date(entityEvent.start_date),
        endDate: new Date(entityEvent.end_date),
        createdAt: new Date(entityEvent.created_at),
        updatedAt: new Date(entityEvent.updated_at),
    };
};

export const mapEventResponseToViewModel = (eventResponse: EventEventResponse): EventViewModel => {
    return {
        bannerPresignedUrl: eventResponse.banner_presigned_url,
        chainId: eventResponse.chain_id,
        contactAddress: eventResponse.contact_address,
        contactNumber: eventResponse.contact_number,
        createdAt: new Date(eventResponse.created_at),
        endDate: new Date(eventResponse.end_date),
        eventStatus: mapEntityEventStatusToEventStatus(eventResponse.event_status),
        eventType: mapEntityEventTypeToEventType(eventResponse.event_type),
        iconPresignedUrl: eventResponse.icon_presigned_url,
        id: eventResponse.id,
        isBookingRequestRequired: eventResponse.is_booking_request_required,
        isPublic: eventResponse.is_public,
        isTicketTransferable: eventResponse.is_ticket_transferable,
        isVerified: eventResponse.is_verified,
        location: eventResponse.location,
        longDescription: eventResponse.long_description,
        maxAttendees: eventResponse.max_attendees,
        attendeesCount: eventResponse.attendees_count,
        ownerCredentialId: eventResponse.owner_credential_id,
        shortDescription: eventResponse.short_description,
        startDate: new Date(eventResponse.start_date),
        title: eventResponse.title,
        updatedAt: new Date(eventResponse.updated_at),
        googleMapQuery: eventResponse.google_map_query,
    };
};

export const mapEventViewModelExtended = (
    eventViewModel: EventEventViewModel,
): EventViewModelExtended => {
    return {
        ...mapEventResponseToViewModel(eventViewModel),
        isInvited: eventViewModel.is_invited,
        isJoined: eventViewModel.is_joined,
        isFull: eventViewModel.is_full,
        finalCallDate: eventViewModel.registration_config?.final_call_for_registration,
        registrationRequirement: mapEntityEventRegistrationConfigToEventRegistrationConfiguration(
            eventViewModel.registration_config,
        ),
        accessManagerContractAddress:
            eventViewModel.event_contract?.access_manager_contract_address,
        eventContractAddress: eventViewModel.event_contract?.event_contract_address,
        ticketContractAddress: eventViewModel.event_contract?.ticket_contract_address,
        certificateContractAddress: eventViewModel.event_contract?.certificate_contract_address,
    };
};

export const mapToEventIssuer = (eventIssuerResponse: EventEventIssuerResponse): EventIssuer => {
    return {
        eventId: eventIssuerResponse.event_id,
        id: eventIssuerResponse.id,
        isSigned: eventIssuerResponse.is_signed === 1,
        issuerCredentialId: eventIssuerResponse.issuer_credential_id,
        issuerProfile: eventIssuerResponse.issuer_profile
            ? mapProfileViewModel(eventIssuerResponse.issuer_profile)
            : undefined,

        createdAt: new Date(eventIssuerResponse.created_at),
        updatedAt: new Date(eventIssuerResponse.updated_at),
    };
};
