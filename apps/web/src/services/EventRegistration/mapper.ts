import type {
    EntityEventRegistrationInvitation,
    EntityInboxMessage,
    EventconfigEventRegistrationConfigResponse,
} from "@decm/api";
import type {
    EventRegistrationConfiguration,
    EventRegistrationInvitation,
    RegistrationRequirement,
    RegistrationRequirementStatus,
} from "./EventRegistration";
import type { InboxMessage } from "../InboxService/InboxService";
import { mapEntityInboxMessageToInboxMessage as mapInboxMessage } from "../InboxService/mapper";

export const mapRegistrationRequirementStatus = (status: number): RegistrationRequirementStatus => {
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

export const mapEntityEventRegistrationConfigRequirementStatus = (
    entityEventRegistrationConfig: EventconfigEventRegistrationConfigResponse,
): RegistrationRequirement => {
    return {
        firstName: mapRegistrationRequirementStatus(
            entityEventRegistrationConfig.first_name_requirement_status,
        ),
        lastName: mapRegistrationRequirementStatus(
            entityEventRegistrationConfig.last_name_requirement_status,
        ),
        email: mapRegistrationRequirementStatus(
            entityEventRegistrationConfig.email_requirement_status,
        ),
        bio: mapRegistrationRequirementStatus(entityEventRegistrationConfig.bio_requirement_status),
        phoneNumber: mapRegistrationRequirementStatus(
            entityEventRegistrationConfig.phone_number_requirement_status,
        ),
        address: mapRegistrationRequirementStatus(
            entityEventRegistrationConfig.address_requirement_status,
        ),
        academicInstitution: mapRegistrationRequirementStatus(
            entityEventRegistrationConfig.academic_institution_requirement_status,
        ),
        academicEmail: mapRegistrationRequirementStatus(
            entityEventRegistrationConfig.academic_email_requirement_status,
        ),
    };
};

export const mapEntityEventRegistrationConfigToEventRegistrationConfiguration = (
    entityEventRegistrationConfig: EventconfigEventRegistrationConfigResponse,
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
