import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import type { InboxMessage } from "../InboxService/InboxService";
import {
    mapEntityEventRegistrationConfigToEventRegistrationConfiguration,
    mapEntityEventRegistrationInvitationToEventRegistrationInvitation,
    mapEntityInboxMessageToInboxMessage,
    mapRegistrationRequirementStatusToNumber,
    mapRegistrationToJoinEventParticipant,
} from "./mapper";
import type { RegistrationConfirmDataForm } from "@/components/pages/Participant/Events/Detail/RegistrationConfirmDataFormSchema";

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
    registrationPassword?: string;
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

    public async updateConfiguration(
        eventId: string,
        configuration: EventRegistrationConfiguration,
    ): Promise<void> {
        await this._coreApi.v1.updateEventRegistrationConfig(
            { eventId },
            {
                academic_email_requirement_status: mapRegistrationRequirementStatusToNumber(
                    configuration.academicEmail,
                ),
                academic_institution_requirement_status: mapRegistrationRequirementStatusToNumber(
                    configuration.academicInstitution,
                ),
                address_requirement_status: mapRegistrationRequirementStatusToNumber(
                    configuration.address,
                ),
                bio_requirement_status: mapRegistrationRequirementStatusToNumber(configuration.bio),
                email_requirement_status: mapRegistrationRequirementStatusToNumber(
                    configuration.email,
                ),
                first_name_requirement_status: mapRegistrationRequirementStatusToNumber(
                    configuration.firstName,
                ),
                last_name_requirement_status: mapRegistrationRequirementStatusToNumber(
                    configuration.lastName,
                ),
                phone_number_requirement_status: mapRegistrationRequirementStatusToNumber(
                    configuration.phoneNumber,
                ),
                // TODO
                is_booking_request_required: false,
                is_ticket_transferable: false,
            },
        );
        return;
    }

    public async joinEventWithAccountPassword({
        eventId,
        eventPassword,
        accountPassword,
        registrationData,
    }: {
        eventId: string;
        eventPassword?: string;
        accountPassword: string;
        registrationData: RegistrationConfirmDataForm;
    }): Promise<void> {
        await this._coreApi.v1.joinEvent(
            { eventId },
            {
                account_password: accountPassword,
                event_password: eventPassword,
                registration_data: mapRegistrationToJoinEventParticipant(registrationData),
            },
        );
        return;
    }

    public async fuckJoinEvent({
        eventId,
        pinCode,
        registrationData,
    }: {
        eventId: string;
        pinCode: string;
        registrationData: RegistrationConfirmDataForm;
    }): Promise<void> {
        await this._coreApi.v1.fuckJoinEvent(
            { eventId },
            {
                pin_code: pinCode,
                academic_email: registrationData.academicEmail,
                academic_institution: registrationData.academicInstitution,
                address: registrationData.address,
                bio: registrationData.bio,
                email: registrationData.email,
                first_name: registrationData.firstName,
                last_name: registrationData.lastName,
                phone_number: registrationData.phoneNumber,
            },
        );
    }

    public async joinEventWithSignature({
        eventId,
        originalSignMessage,
        signature,
        registrationData,
    }: {
        eventId: string;
        originalSignMessage: string;
        signature: string;
        registrationData: RegistrationConfirmDataForm;
    }): Promise<void> {
        await this._coreApi.v1.joinEvent(
            { eventId },
            {
                sign_message: originalSignMessage,
                signature: signature,
                registration_data: mapRegistrationToJoinEventParticipant(registrationData),
            },
        );
        return;
    }
}

export const defaultEventRegistrationService = new EventRegistrationService(coreApiClient);
