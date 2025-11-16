import type {
    CoreApiInternalHandlerEventRegistrationJoinEventPayload,
    EntityEventRegistrationInvitation,
    EntityInboxMessage,
    EventconfigEventRegistrationConfigViewModel,
} from "@decm/api";
import type {
    EventRegistrationConfiguration,
    EventRegistrationInvitation,
    RegistrationRequirement,
    RegistrationRequirementStatus,
} from "./EventRegistration";
import type { InboxMessage } from "../InboxService/InboxService";
import { mapEntityInboxMessageToInboxMessage as mapInboxMessage } from "../InboxService/mapper";
import type { RegistrationConfirmDataForm } from "@/components/pages/Participant/Events/Detail/RegistrationConfirmDataFormSchema";

export type EventRegistrationConfigResponse = Omit<
    EventconfigEventRegistrationConfigViewModel,
    "created_at" | "updated_at"
>;

export const mapToRegistrationRequirementStatus = (
    status: number,
): RegistrationRequirementStatus => {
    switch (status) {
        case 0:
            return "not_required";
        case 1:
            return "required";
        case 2:
            return "optional";
    }
    throw new Error(`Invalid registration requirement status: ${status}`);
};

export const mapRegistrationRequirementStatusToNumber = (
    status: RegistrationRequirementStatus,
): number => {
    switch (status) {
        case "not_required":
            return 0;
        case "required":
            return 1;
        case "optional":
            return 2;
    }
    throw new Error(`Invalid registration requirement status: ${status}`);
};

export const mapEntityEventRegistrationConfigRequirementStatus = (
    entityEventRegistrationConfig: EventRegistrationConfigResponse,
): RegistrationRequirement => {
    return {
        firstName: mapToRegistrationRequirementStatus(
            entityEventRegistrationConfig.first_name_requirement_status,
        ),
        lastName: mapToRegistrationRequirementStatus(
            entityEventRegistrationConfig.last_name_requirement_status,
        ),
        email: mapToRegistrationRequirementStatus(
            entityEventRegistrationConfig.email_requirement_status,
        ),
        bio: mapToRegistrationRequirementStatus(
            entityEventRegistrationConfig.bio_requirement_status,
        ),
        phoneNumber: mapToRegistrationRequirementStatus(
            entityEventRegistrationConfig.phone_number_requirement_status,
        ),
        address: mapToRegistrationRequirementStatus(
            entityEventRegistrationConfig.address_requirement_status,
        ),
        academicInstitution: mapToRegistrationRequirementStatus(
            entityEventRegistrationConfig.academic_institution_requirement_status,
        ),
        academicEmail: mapToRegistrationRequirementStatus(
            entityEventRegistrationConfig.academic_email_requirement_status,
        ),
    };
};

export const mapEntityEventRegistrationConfigToEventRegistrationConfiguration = (
    entityEventRegistrationConfig: EventRegistrationConfigResponse,
): EventRegistrationConfiguration => {
    return {
        finalCallForRegistration: entityEventRegistrationConfig.final_call_for_registration
            ? new Date(entityEventRegistrationConfig.final_call_for_registration)
            : undefined,
        ...mapEntityEventRegistrationConfigRequirementStatus(entityEventRegistrationConfig),
    };
};

export const mapEntityEventRegistrationInvitationToEventRegistrationInvitation = (
    entityEventRegistrationInvitation: EntityEventRegistrationInvitation,
): EventRegistrationInvitation => {
    return {
        id: entityEventRegistrationInvitation.id,
        eventId: entityEventRegistrationInvitation.event_id,
        email: entityEventRegistrationInvitation.email,
        firstName: entityEventRegistrationInvitation.first_name,
        lastName: entityEventRegistrationInvitation.last_name,
        phoneNumber: entityEventRegistrationInvitation.phone_number,
        academicInstitution: entityEventRegistrationInvitation.academic_institution,
        validUntil: entityEventRegistrationInvitation.valid_until
            ? new Date(entityEventRegistrationInvitation.valid_until)
            : undefined,
        createdAt: new Date(entityEventRegistrationInvitation.created_at),
        updatedAt: new Date(entityEventRegistrationInvitation.updated_at),
        cancelledAt: entityEventRegistrationInvitation.cancelled_at
            ? new Date(entityEventRegistrationInvitation.cancelled_at)
            : undefined,
        inboxMessageId: entityEventRegistrationInvitation.inbox_message_id,
    };
};

// Re-export the inbox message mapper for convenience
export const mapEntityInboxMessageToInboxMessage = (
    entityInboxMessage: EntityInboxMessage,
): InboxMessage => {
    return mapInboxMessage(entityInboxMessage);
};

export const mapRegistrationToJoinEventParticipant = (
    participantData: RegistrationConfirmDataForm,
): CoreApiInternalHandlerEventRegistrationJoinEventPayload => {
    return {
        first_name: participantData.firstName,
        last_name: participantData.lastName,
        email: participantData.email,
        phone_number: participantData.phoneNumber,
        academic_institution: participantData.academicInstitution,
        academic_email: participantData.academicEmail,
        address: participantData.address,
        bio: participantData.bio,
    };
};
