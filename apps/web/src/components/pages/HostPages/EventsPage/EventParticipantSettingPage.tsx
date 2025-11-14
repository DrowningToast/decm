import { ParticipantSettingsForm } from "@/components/forms/ParticipantSettingsForm";
import { type ParticipantSettingsData } from "@/lib/schemas/participantSettingsSchema";

import SectionContainer from "@/components/container/SectionContainer";
import { useUpdateParticipantSetting } from "@/components/forms/ParticipantSettingsForm/useUpdateParticipantSetting";
import type { EntityEventType, EventconfigUpdateEventRegistrationConfigRequest } from "@decm/api";
import { toEventRegistrationConfigStatusNumber } from "@/lib/events/event.utils";
import { useEventViewModelUsecase } from "@/components/pages/Participant/Events/Detail/useEventViewModelUsecase";

interface EventParticipantSettingPageProps {
    eventId: string;
}
export const EventParticipantSettingPage = ({ eventId }: EventParticipantSettingPageProps) => {
    const { event, isLoading: isLoadingEvent } = useEventViewModelUsecase({ eventId });
    const { updateParticipantSetting, isUpdatingParticipantSetting } =
        useUpdateParticipantSetting(eventId);

    if (isLoadingEvent) {
        return <div>Loading event...</div>;
    }

    if (!event) {
        return <div>Event not found</div>;
    }

    const eventRegistrationConfig = event.registrationRequirement;
    if (!eventRegistrationConfig) {
        return <div>Event registration config not found</div>;
    }

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

    console.log("eventRegistrationConfig", eventRegistrationConfig);

    return (
        <div>
            <SectionContainer>
                <ParticipantSettingsForm
                    onSubmit={onSubmit}
                    isLoading={isUpdatingParticipantSetting}
                    defaultValues={{
                        // Event Config
                        academicEmail: eventRegistrationConfig.academicEmail,
                        academicInstitution: eventRegistrationConfig.academicInstitution,
                        address: eventRegistrationConfig.address,
                        bio: eventRegistrationConfig.bio,
                        email: eventRegistrationConfig.email,
                        firstName: eventRegistrationConfig.firstName,
                        lastName: eventRegistrationConfig.lastName,
                        phoneNumber: eventRegistrationConfig.phoneNumber,
                        finalCallRegistrationDate: eventRegistrationConfig.finalCallForRegistration
                            ? new Date(eventRegistrationConfig.finalCallForRegistration)
                            : undefined,
                        // registrationPassword: event.registrationPassword,
                        requireRegistrationPassword: event.registrationPassword !== undefined,

                        // Event
                        isBookingRequired: event.isBookingRequestRequired,
                        isTicketTransferable: event.isTicketTransferable,
                    }}
                    showPreview={true}
                />
            </SectionContainer>
        </div>
    );
};
