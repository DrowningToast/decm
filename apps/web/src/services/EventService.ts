import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import {
    EntityEventStatus,
    EntityEventType,
    type EventEventResponse,
    type EventEventViewModel,
} from "@decm/api";

interface GetEventsListParams {
    includeActiveEvents?: boolean;
    includeInactiveEvents?: boolean;
    includeClosedEvents?: boolean;
    onlyUserJoinedEvents?: boolean;
}

export type RegistrationRequirementStatus = "required" | "optional" | "not_required";

interface RegistrationRequirement {
    firstName: RegistrationRequirementStatus;
    lastName: RegistrationRequirementStatus;
    email: RegistrationRequirementStatus;
    bio: RegistrationRequirementStatus;
    phoneNumber: RegistrationRequirementStatus;
    address: RegistrationRequirementStatus;
    academicInstitution: RegistrationRequirementStatus;
    academicEmail: RegistrationRequirementStatus;
}

export interface Event {
    banner_presigned_url?: string;
    chainId?: number;
    contactAddress?: string;
    contactNumber?: string;
    createdAt?: string;
    endDate?: string;
    eventStatus?: EntityEventStatus;
    eventType?: EntityEventType;
    googleMapQuery?: string;
    icon_presigned_url?: string;
    id?: string;
    isBookingRequestRequired?: boolean;
    isPublic?: boolean;
    isTicketTransferable?: boolean;
    isVerified?: boolean;
    location?: string;
    longDescription?: string;
    maxAttendees?: number;
    ownerCredentialId?: string;
    shortDescription?: string;
    startDate?: string;
    title?: string;
    updatedAt?: string;
}

export interface EventViewModel extends Event {
    isInvited?: boolean;
    isJoined?: boolean;

    finalCallDate?: string;
    registrationRequirement?: RegistrationRequirement;

    accessManagerContractAddress?: string;
    eventContractAddress?: string;
    ticketContractAddress?: string;
    certificateContractAddress?: string;
}

export const EventStatus = EntityEventStatus;
export const EventType = EntityEventType;

export class EventService {
    private _coreApi: CoreApiType;

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    private transformRegistrationRequirementStatus(status: number): RegistrationRequirementStatus {
        switch (status) {
            case 0:
                return "not_required";
            case 1:
                return "required";
            case 2:
                return "optional";
        }
        throw new Error(`Invalid registration requirement status: ${status}`);
    }

    /**
     * Transform API response from snake_case to camelCase
     */
    private transformEntityEvent(entityEvent: EventEventResponse): Event {
        return {
            banner_presigned_url: entityEvent.banner_presigned_url,
            chainId: entityEvent.chain_id,
            contactAddress: entityEvent.contact_address,
            contactNumber: entityEvent.contact_number,
            createdAt: entityEvent.created_at as string | undefined,
            endDate: entityEvent.end_date as string | undefined,
            eventStatus: entityEvent.event_status as EntityEventStatus | undefined,
            eventType: entityEvent.event_type as EntityEventType | undefined,
            googleMapQuery: entityEvent.google_map_query as string | undefined,
            icon_presigned_url: entityEvent.icon_presigned_url as string | undefined,
            id: entityEvent.id as string | undefined,
            isBookingRequestRequired: entityEvent.is_booking_request_required as
                | boolean
                | undefined,
            isPublic: entityEvent.is_public as boolean | undefined,
            isTicketTransferable: entityEvent.is_ticket_transferable as boolean | undefined,
            isVerified: entityEvent.is_verified as boolean | undefined,
            location: entityEvent.location as string | undefined,
            longDescription: entityEvent.long_description as string | undefined,
            maxAttendees: entityEvent.max_attendees as number | undefined,
            ownerCredentialId: entityEvent.owner_credential_id as string | undefined,
            shortDescription: entityEvent.short_description as string | undefined,
            startDate: entityEvent.start_date as string | undefined,
            title: entityEvent.title as string | undefined,
            updatedAt: entityEvent.updated_at as string | undefined,
        };
    }

    private transformEntityEventViewModel(
        entityEventViewModel: EventEventViewModel,
    ): EventViewModel {
        return {
            ...this.transformEntityEvent(entityEventViewModel),
            isInvited: entityEventViewModel.is_invited as boolean | undefined,
            isJoined: entityEventViewModel.is_joined as boolean | undefined,

            registrationRequirement: {
                firstName: this.transformRegistrationRequirementStatus(
                    entityEventViewModel.registration_config
                        ?.first_name_requirement_status as number,
                ),
                lastName: this.transformRegistrationRequirementStatus(
                    entityEventViewModel.registration_config
                        ?.last_name_requirement_status as number,
                ),
                email: this.transformRegistrationRequirementStatus(
                    entityEventViewModel.registration_config?.email_requirement_status as number,
                ),
                bio: this.transformRegistrationRequirementStatus(
                    entityEventViewModel.registration_config?.bio_requirement_status as number,
                ),
                phoneNumber: this.transformRegistrationRequirementStatus(
                    entityEventViewModel.registration_config
                        ?.phone_number_requirement_status as number,
                ),
                address: this.transformRegistrationRequirementStatus(
                    entityEventViewModel.registration_config?.address_requirement_status as number,
                ),
                academicInstitution: this.transformRegistrationRequirementStatus(
                    entityEventViewModel.registration_config
                        ?.academic_institution_requirement_status as number,
                ),
                academicEmail: this.transformRegistrationRequirementStatus(
                    entityEventViewModel.registration_config
                        ?.academic_email_requirement_status as number,
                ),
            },

            finalCallDate: entityEventViewModel.registration_config?.final_call_for_registration,
            eventContractAddress: entityEventViewModel.event_contract?.event_contract_address,
            ticketContractAddress: entityEventViewModel.event_contract?.ticket_contract_address,
            certificateContractAddress:
                entityEventViewModel.event_contract?.certificate_contract_address,
        };
    }

    public async getEventById(eventId: string): Promise<Event> {
        const response = await this._coreApi.v1.getEventById({ eventId });
        return this.transformEntityEvent(response);
    }

    public async getEvents(params: GetEventsListParams): Promise<Event[]> {
        const response = await this._coreApi.v1.getEventsList({
            include_active_events: params.includeActiveEvents ?? true,
            include_inactive_events: params.includeInactiveEvents,
            include_closed_events: params.includeClosedEvents,
            only_user_joined_events: params.onlyUserJoinedEvents,
        });

        return response.events?.map((event) => this.transformEntityEvent(event)) ?? [];
    }

    public async getEventViewModel(eventId: string): Promise<EventViewModel> {
        const response = await this._coreApi.v1.getEventViewmodelById({ eventId });
        return this.transformEntityEventViewModel(response);
    }
}

export const defaultEventService = new EventService(coreApiClient);
