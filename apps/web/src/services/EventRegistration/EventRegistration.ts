import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import type { InboxMessage } from "../InboxService/InboxService";
import {
    mapEntityEventRegistrationConfigToEventRegistrationConfiguration,
    mapEntityEventRegistrationInvitationToEventRegistrationInvitation,
    mapEntityInboxMessageToInboxMessage,
} from "./mapper";

export type RegistrationRequirementStatus = "required" | "optional" | "not_required";

export interface RegistrationRequirement {
    firstName: RegistrationRequirementStatus;
    lastName: RegistrationRequirementStatus;
    email: RegistrationRequirementStatus;
    bio: RegistrationRequirementStatus;
    phoneNumber: RegistrationRequirementStatus;
    address: RegistrationRequirementStatus;
    academicInstitution: RegistrationRequirementStatus;
    academicEmail: RegistrationRequirementStatus;
}

export interface EventRegistrationConfiguration extends RegistrationRequirement {
    finalCallForRegistration?: Date;
}

export interface EventRegistrationInvitation {
    academicInstitution?: string;
    cancelledAt?: Date;
    code?: string;
    createdAt: Date;
    email?: string;
    eventId: string;
    firstName?: string;
    id: string;
    inboxMessageId: string;
    lastName?: string;
    phoneNumber?: string;
    updatedAt: Date;
    validUntil?: Date;
}

export class EventRegistrationService {
    private _coreApi: CoreApiType;

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    public async getConfiguration(eventId: string): Promise<EventRegistrationConfiguration> {
        const response = await this._coreApi.v1.getEventRegistrationConfig({ eventId });
        return mapEntityEventRegistrationConfigToEventRegistrationConfiguration(response);
    }

    public async checkPassword(eventId: string, password: string): Promise<boolean> {
        const response = await this._coreApi.v1.checkEventPassword({ eventId }, { password });
        return response.is_valid;
    }

    public async getInvitationByEventId(eventId: string): Promise<EventRegistrationInvitation[]> {
        const response = await this._coreApi.v1.getEventRegistrationInvitationsByEventId({
            eventId,
        });
        return response.map((invitation) =>
            mapEntityEventRegistrationInvitationToEventRegistrationInvitation(invitation),
        );
    }

    public async getInvitationOfUserAndEventId(eventId: string): Promise<{
        registrationInvitation: EventRegistrationInvitation;
        inbox: InboxMessage;
    }> {
        const response = await this._coreApi.v1.getEventRegistrationInvitationByUserAndEvent({
            eventId,
        });
        if (!response.registration_invitation || !response.inbox) {
            throw new Error("Registration invitation or inbox not found");
        }
        return {
            registrationInvitation:
                mapEntityEventRegistrationInvitationToEventRegistrationInvitation(
                    response.registration_invitation,
                ),
            inbox: mapEntityInboxMessageToInboxMessage(response.inbox),
        };
    }
}

export const defaultEventRegistrationService = new EventRegistrationService(coreApiClient);
