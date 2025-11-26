import { ParticipantSettingsForm } from "@/components/forms/ParticipantSettingsForm";
import { type ParticipantSettingsData } from "@/lib/schemas/participantSettingsSchema";

import SectionContainer from "@/components/container/SectionContainer";
import { useUpdateParticipantSetting } from "@/components/forms/ParticipantSettingsForm/useUpdateParticipantSetting";
import { useEventViewModelUsecase } from "@/components/pages/Participant/Events/Detail/useEventViewModelUsecase";
import { useEventRegistrationConfigUsecase } from "@/hooks/events/useEventRegistrationConfigUsecase";
import type { EventRegistrationConfiguration } from "@/services/EventRegistration/EventRegistration";

interface EventParticipantSettingPageProps {
    eventId: string;
}
export const EventParticipantSettingPage = ({ eventId }: EventParticipantSettingPageProps) => {
    const { event, isLoading: isLoadingEvent } = useEventViewModelUsecase({ eventId });
    const { data: eventRegistrationConfig, isLoading: isLoadingEventRegistrationConfig } =
        useEventRegistrationConfigUsecase({ eventId });
    const { updateParticipantSetting, isUpdatingParticipantSetting } =
        useUpdateParticipantSetting(eventId);

    if (isLoadingEvent || isLoadingEventRegistrationConfig || !eventRegistrationConfig) {
        return <div>Loading event...</div>;
    }

    if (!event) {
        return <div>Event not found</div>;
    }

    const onSubmit = async (data: ParticipantSettingsData) => {
        const params: EventRegistrationConfiguration = {
            academicEmail: data.academicEmail,
            academicInstitution: data.academicInstitution,
            address: data.address,
            bio: data.bio,
            email: data.email,
            firstName: data.firstName,
            lastName: data.lastName,
            phoneNumber: data.phoneNumber,
            eventType: data.eventType,
            finalCallForRegistration: data.finalCallRegistrationDate
                ? new Date(data.finalCallRegistrationDate)
                : undefined,
            // TODO: Add these fields
            // password
            // is_booking_request_required: data.isBookingRequired,
            // is_ticket_transferable: data.isTicketTransferable,
        };

        if (data.requireRegistrationPassword) {
            params.registrationPassword = data.registrationPassword ?? undefined;
        }

        await updateParticipantSetting({
            academicEmail: params.academicEmail,
            academicInstitution: params.academicInstitution,
            address: params.address,
            bio: params.bio,
            email: params.email,
            firstName: params.firstName,
            lastName: params.lastName,
            phoneNumber: params.phoneNumber,
            eventType: params.eventType,
            finalCallForRegistration: params.finalCallForRegistration,
            registrationPassword: params.registrationPassword,
        });
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
                        eventType: event.eventType,
                        // Add missing required fields from Event object
                        // isBookingRequired: event.isBookingRequestRequired ?? false,
                        // requireRegistrationPassword: false, // TODO: Get from backend
                        registrationPassword: "", // Password not returned from API for security
                    }}
                    showPreview={true}
                />
            </SectionContainer>
        </div>
    );
};
