import HostEventDetailsPage from "@/components/pages/HostPages/EventsPage/HostEventDetailsPage";
import { useParams } from "@/router";
import { useEventRegistrationConfig } from "@/hooks/events/useEventRegistrationConfig";
import { useEventCertificateConfig } from "@/components/pages/HostPages/EventPages/useEventCertificateConfig";
import { useEventIssuers } from "@/components/pages/HostPages/EventPages/useEventIssuers";
import { useEventContract } from "@/hooks/events/useEventContracts";
import { useEventInvitedParticipants } from "@/hooks/events/useEventInvitedParticipants";
import { useEventViewModelUsecase } from "@/components/pages/Participant/Events/Detail/useEventViewModelUsecase";

export default function Page() {
    const { eventId } = useParams("/host/events/:eventId");
    const {
        event: eventViewModel,
        isLoading: isLoadingEventViewModel,
        error: errorEventViewModel,
    } = useEventViewModelUsecase({ eventId: eventId! });

    const {
        data: eventRegistrationConfig,
        isLoading: isLoadingEventRegistrationConfig,
        error: isErrorEventRegistrationConfig,
    } = useEventRegistrationConfig(eventId);

    const { data: eventCertificateConfig, isLoading: isLoadingEventCertificateConfig } =
        useEventCertificateConfig(eventId!);

    const { data: eventContract, isLoading: isLoadingEventContract } = useEventContract(eventId!);

    const { eventIssuers, isLoadingEventIssuers } = useEventIssuers(eventId!);
    const { invitations, isLoading: isLoadingEventInvitations } = useEventInvitedParticipants(
        eventId!,
    );

    const isLoading =
        isLoadingEventViewModel ||
        isLoadingEventRegistrationConfig ||
        isLoadingEventCertificateConfig ||
        isLoadingEventIssuers ||
        isLoadingEventContract ||
        isLoadingEventInvitations;

    if (isLoading) {
        return <div>Loading event...</div>;
    }

    if (
        errorEventViewModel ||
        isErrorEventRegistrationConfig ||
        !eventViewModel ||
        !eventRegistrationConfig ||
        !eventContract
    ) {
        return <div>Error loading event</div>;
    }

    return (
        <HostEventDetailsPage
            eventId={eventId}
            event={eventViewModel}
            eventRegistrationConfig={eventRegistrationConfig}
            eventCertificateConfig={eventCertificateConfig}
            eventIssuers={eventIssuers}
            eventContract={eventContract}
            eventInvitations={invitations}
        />
    );
}
