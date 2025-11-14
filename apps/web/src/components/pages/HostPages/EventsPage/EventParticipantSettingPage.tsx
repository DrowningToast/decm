import { ParticipantSettingsForm } from "@/components/forms/ParticipantSettingsForm";
import { type ParticipantSettingsData } from "@/lib/schemas/participantSettingsSchema";

import SectionContainer from "@/components/container/SectionContainer";
import { useUpdateParticipantSetting } from "@/components/forms/ParticipantSettingsForm/useUpdateParticipantSetting";
import type {
    EntityEventType,
    EventconfigEventRegistrationConfigResponse,
    EventconfigUpdateEventRegistrationConfigRequest,
    EventEventResponse,
} from "@decm/api";
import {
    toEventRegistrationConfigStatus,
    toEventRegistrationConfigStatusNumber,
} from "@/lib/events/event.utils";

interface EventParticipantSettingPageProps {
    eventId: string;
    event: EventEventResponse;
    eventRegistrationConfig: EventconfigEventRegistrationConfigResponse;
}
export const EventParticipantSettingPage = ({
    eventId,
    eventRegistrationConfig,
    event,
}: EventParticipantSettingPageProps) => {
    const { updateParticipantSetting, isUpdatingParticipantSetting } =
        useUpdateParticipantSetting(eventId);

    const onSubmit = async (data: ParticipantSettingsData) => {
        const params: EventconfigUpdateEventRegistrationConfigRequest = {
            academic_email_requirement_status: toEventRegistrationConfigStatusNumber(
                data.academicEmail,
            ),
            academic_institution_requirement_status: toEventRegistrationConfigStatusNumber(
                data.academicInstitution,
            ),
            address_requirement_status: toEventRegistrationConfigStatusNumber(data.address),
            bio_requirement_status: toEventRegistrationConfigStatusNumber(data.bio),
            email_requirement_status: toEventRegistrationConfigStatusNumber(data.email),
            first_name_requirement_status: toEventRegistrationConfigStatusNumber(data.firstName),
            last_name_requirement_status: toEventRegistrationConfigStatusNumber(data.lastName),
            phone_number_requirement_status: toEventRegistrationConfigStatusNumber(
                data.phoneNumber,
            ),
            final_call_for_registration: data.finalCallRegistrationDate
                ? new Date(data.finalCallRegistrationDate).toISOString()
                : undefined,
            event_type: data.eventType as EntityEventType,
            is_booking_request_required: data.isBookingRequired,
            is_ticket_transferable: data.isTicketTransferable,
        };

        if (data.requireRegistrationPassword) {
            params.registration_password = data.registrationPassword ?? undefined;
        }

        await updateParticipantSetting(params);
    };

    return (
        <div>
            <SectionContainer>
                <ParticipantSettingsForm
                    onSubmit={onSubmit}
                    isLoading={isUpdatingParticipantSetting}
                    defaultValues={{
                        // Event Config
                        academicEmail: toEventRegistrationConfigStatus(
                            eventRegistrationConfig.academic_email_requirement_status,
                        ),
                        academicInstitution: toEventRegistrationConfigStatus(
                            eventRegistrationConfig.academic_institution_requirement_status,
                        ),
                        address: toEventRegistrationConfigStatus(
                            eventRegistrationConfig.address_requirement_status,
                        ),
                        bio: toEventRegistrationConfigStatus(
                            eventRegistrationConfig.bio_requirement_status,
                        ),
                        email: toEventRegistrationConfigStatus(
                            eventRegistrationConfig.email_requirement_status,
                        ),
                        firstName: toEventRegistrationConfigStatus(
                            eventRegistrationConfig.first_name_requirement_status,
                        ),
                        lastName: toEventRegistrationConfigStatus(
                            eventRegistrationConfig.last_name_requirement_status,
                        ),
                        phoneNumber: toEventRegistrationConfigStatus(
                            eventRegistrationConfig.phone_number_requirement_status,
                        ),
                        finalCallRegistrationDate:
                            eventRegistrationConfig.final_call_for_registration
                                ? new Date(eventRegistrationConfig.final_call_for_registration)
                                : undefined,
                        registrationPassword: eventRegistrationConfig.registration_password
                            ? eventRegistrationConfig.registration_password
                            : undefined,
                        requireRegistrationPassword:
                            !!eventRegistrationConfig.registration_password,

                        // Event
                        isBookingRequired: event.is_booking_request_required,
                        isTicketTransferable: event.is_ticket_transferable,
                    }}
                    showPreview={true}
                />
            </SectionContainer>
        </div>
    );
};
